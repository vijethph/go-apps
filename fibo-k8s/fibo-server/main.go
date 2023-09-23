package main

import (
	"context"
	"fmt"
	"os"
	"os/signal"
	"syscall"
	"time"

	"github.com/gofiber/fiber/v2"
	"github.com/gofiber/fiber/v2/middleware/cors"
	"github.com/jackc/pgx/v5"
	"github.com/jackc/pgx/v5/pgxpool"
	"github.com/redis/go-redis/v9"
	"github.com/rs/zerolog/log"
)

var ctx = context.Background()
var dbPoolClient *pgxpool.Pool
var redisClient *redis.Client
var publisherClient *redis.Client

func init() {
	var commonErr error
	redisClient = redis.NewClient(&redis.Options{
		Addr:            fmt.Sprintf("%s:%s", redisHost, redisPort),
		MaxRetries:      3,
		MinRetryBackoff: 100 * time.Millisecond,
		Password:        "", // no password set
		DB:              0,  // use default DB
	})

	publisherClient = redis.NewClient(&redis.Options{
		Addr:            fmt.Sprintf("%s:%s", redisHost, redisPort),
		MaxRetries:      3,
		MinRetryBackoff: 100 * time.Millisecond,
		Password:        "", // no password set
		DB:              0,  // use default DB
	})

	time.Sleep(5 * time.Second) // wait for postgres container to start

	dbPoolClient, commonErr = pgxpool.New(ctx, fmt.Sprintf("postgres://%s:%s@%s:%s/%s", pgUser, pgPassword, pgHost, pgPort, pgDatabase))
	if commonErr != nil {
		log.Error().Err(commonErr).Msg("Unable to create connection pool")
		os.Exit(1)
	}

	result, err := dbPoolClient.Exec(ctx, "CREATE TABLE IF NOT EXISTS values (number INT)")
	if err != nil {
		log.Error().Err(err).Msg("Unable to create connection pool")
		os.Exit(1)
	}

	if result.RowsAffected() != 0 {
		log.Warn().Msg("CREATE TABLE failed\n")
	}

}

func main() {

	app := fiber.New()
	app.Use(cors.New())

	app.Get("/", func(c *fiber.Ctx) error {
		return c.SendString("Hi!")
	})

	app.Get("/values/all", func(c *fiber.Ctx) error {
		rows, _ := dbPoolClient.Query(context.Background(), "select * from values")
		numbers, err := pgx.CollectRows(rows, pgx.RowTo[int64])
		if err != nil {
			return err
		}
		// fmt.Println(numbers)
		return c.JSON(numbers)
	})

	app.Get("/values/current", func(c *fiber.Ctx) error {
		val, err := redisClient.HGetAll(ctx, "values").Result()
		if err != nil {
			return c.Status(fiber.StatusInternalServerError).JSON(fiber.Map{
				"message": "Redis get failed",
			})
		}

		return c.JSON(val)
	})

	app.Post("/values", func(c *fiber.Ctx) error {
		type RequestPayload struct {
			Index int64 `json:"index"`
		}

		requestPayload := new(RequestPayload)
		if err := c.BodyParser(requestPayload); err != nil {
			log.Error().Err(err).Msg("Invalid request payload")
			return c.Status(fiber.StatusBadRequest).JSON(fiber.Map{
				"message": "Invalid request payload",
			})
		}

		if requestPayload.Index > 40 {
			return c.Status(fiber.StatusBadRequest).JSON(fiber.Map{
				"message": "Index too high",
			})
		}

		redisClient.HSet(ctx, "values", requestPayload.Index, "Nothing yet!")
		err := publisherClient.Publish(ctx, "insert", requestPayload.Index).Err()
		if err != nil {
			log.Error().Err(err).Msg("Failed to publish message to Redis")
		}

		result, err := dbPoolClient.Exec(context.Background(), "INSERT INTO values(number) VALUES($1)", requestPayload.Index)
		if err != nil {
			return c.Status(fiber.StatusInternalServerError).JSON(fiber.Map{
				"message": "Insert to DB failed",
			})
		}

		if result.RowsAffected() != 1 {
			return c.Status(fiber.StatusInternalServerError).JSON(fiber.Map{
				"message": "Insert to DB failed",
			})
		}

		return c.JSON(fiber.Map{
			"working": true,
		})

	})

	// Start the server in a new go routine
	go func() {
		err := app.Listen(":5000")
		if err != nil {
			log.Error().Err(err).Msg("Error starting server")
		}
	}()

	// Wait for a signal to gracefully shut down the server
	// creates a channel with buffer size 1
	quit := make(chan os.Signal, 1)
	signal.Notify(quit, syscall.SIGINT, syscall.SIGTERM)
	<-quit

	// Shut down the server gracefully
	_, cancel := context.WithTimeout(context.Background(), 5*time.Second)
	defer cancel()

	err := app.Shutdown()
	if err != nil {
		log.Error().Err(err).Msg("Error shutting down server")
	}

	// Close the database connection pool
	dbPoolClient.Close()

}

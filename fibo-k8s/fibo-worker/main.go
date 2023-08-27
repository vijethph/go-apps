package main

import (
	"context"
	"fmt"
	"strconv"
	"time"

	"github.com/redis/go-redis/v9"
	"github.com/rs/zerolog"
	"github.com/rs/zerolog/log"
)

var ctx = context.Background()

func main() {
	zerolog.TimeFieldFormat = zerolog.TimeFormatUnix
	// Create a Redis client
	redisClient := redis.NewClient(&redis.Options{
		Addr:            fmt.Sprintf("%s:%s", redisHost, redisPort),
		MaxRetries:      3,
		MinRetryBackoff: 100 * time.Millisecond,
		Password:        "", // no password set
		DB:              0,  // use default DB
	})

	// Create a pubsub client
	pubsub := redisClient.Subscribe(ctx, "insert")

	// Wait for confirmation that subscription is created before publishing anything
	// _, err := pubsub.Receive(ctx)
	// if err != nil {
	// 	log.Fatal().Err(err).Msg("Failed to subscribe to Redis channel")
	// }

	// Listen for messages
	ch := pubsub.Channel()
	for msg := range ch {
		index, err := strconv.Atoi(msg.Payload)
		if err != nil {
			log.Error().Err(err).Msg("Failed to convert message payload to integer")
			continue
		}

		// Calculate the Fibonacci number
		result := fib(index)

		// Store the result in Redis
		err = redisClient.HSet(ctx, "values", msg.Payload, result).Err()
		if err != nil {
			log.Error().Err(err).Msg("Failed to store result in Redis")
			continue
		}

		log.Info().Str("index", msg.Payload).Int("result", result).Msg("Calculated Fibonacci number")
	}

	defer pubsub.Close()
	defer redisClient.Close()
}

func fib(index int) int {
	if index < 2 {
		return 1
	}
	return fib(index-1) + fib(index-2)
}

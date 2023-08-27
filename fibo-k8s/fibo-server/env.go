package main

import (
	"os"

	_ "github.com/joho/godotenv/autoload"
)

var (
	redisHost  = os.Getenv("REDIS_HOST")
	redisPort  = os.Getenv("REDIS_PORT")
	pgUser     = os.Getenv("PG_USER")
	pgHost     = os.Getenv("PG_HOST")
	pgPort     = os.Getenv("PG_PORT")
	pgDatabase = os.Getenv("PG_DATABASE")
	pgPassword = os.Getenv("PG_PASSWORD")
)

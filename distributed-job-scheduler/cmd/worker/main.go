package main

import (
	"context"
	"log"
	"os"

	"github.com/jackc/pgx/v5/pgxpool"
	"github.com/joho/godotenv"

	"github.com/zahid031/distributed-job-scheduler/redisclient"
	"github.com/zahid031/distributed-job-scheduler/worker"
)

func main() {
	if err := godotenv.Load(); err != nil {
		log.Println("no .env file found, relying on real environment variables")
	}

	ctx := context.Background()

	connString := os.Getenv("DATABASE_URL")
	if connString == "" {
		log.Fatal("DATABASE_URL environment variable is not set")
	}

	pool, err := pgxpool.New(ctx, connString)
	if err != nil {
		log.Fatal(err)
	}
	defer pool.Close()

	if err := pool.Ping(ctx); err != nil {
		log.Fatal(err)
	}
	log.Println("worker connected to db")

	redisClient := redisclient.NewClient(ctx)
	defer redisClient.Close()

	w := worker.New(pool, redisClient)
	w.Run(ctx)
}
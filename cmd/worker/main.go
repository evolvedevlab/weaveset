package main

import (
	"context"
	"log"
	"time"

	"github.com/evolvedevlab/weaveset/data"
	"github.com/evolvedevlab/weaveset/scraper"
	"github.com/evolvedevlab/weaveset/util"
	"github.com/google/uuid"
	"github.com/joho/godotenv"
	"github.com/redis/go-redis/v9"
)

var (
	client *redis.Client
	queue  data.Queuer
)

func main() {
	godotenv.Load()
	var (
		hostname  = util.GetEnv("HOSTNAME")
		redisAddr = util.GetEnv("REDIS_ADDR", "127.0.0.1:6379")
		redisPass = util.GetEnv("REDIS_PASSWORD")
	)
	if len(hostname) == 0 {
		log.Fatal("HOSTNAME variable not provided")
	}

	client = redis.NewClient(&redis.Options{
		Addr:       redisAddr,
		Password:   redisPass,
		DB:         0,
		ClientName: "worker",
	})

	ctx := context.Background()
	if err := client.Ping(ctx).Err(); err != nil {
		log.Fatal("redis ping error:", err)
	}

	queue = data.NewRedisQueue(hostname, "jobs", "workers", client)
	queue.Enqueue(ctx, &data.Job{
		ID:        uuid.New().String(),
		URL:       "https://www.goodreads.com/list/show/399714",
		CreatedAt: time.Now(),
	})

	log.Println("Consume loop started...")
	queue.Consume(ctx, scraper.NewHandler(nil))
}

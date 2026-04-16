package main

import (
	"context"
	"log"
	"net/http"
	"strings"

	"github.com/evolvedevlab/weaveset/util"
	"github.com/joho/godotenv"
	"github.com/redis/go-redis/v9"
)

func main() {
	godotenv.Load()
	var (
		listenAddr = util.GetEnv("LISTEN_ADDR", ":3000")

		redisAddr = util.GetEnv("REDIS_ADDR", "127.0.0.1:6379")
		redisPass = util.GetEnv("REDIS_PASSWORD")
	)

	rc := redis.NewClient(&redis.Options{
		Addr:       redisAddr,
		Password:   redisPass,
		DB:         0,
		ClientName: "apiserver",
	})

	ctx := context.Background()
	if err := rc.Ping(ctx).Err(); err != nil {
		log.Fatal("redis ping:", err)
	}

	err := rc.XGroupCreateMkStream(ctx, "jobs", "workers", "0").Err()
	if err != nil && !strings.Contains(err.Error(), "BUSYGROUP") {
		log.Fatal(err)
	}

	http.HandleFunc("/health", handleGetHealth)

	log.Printf("started at %s\n", listenAddr)
	log.Fatal(http.ListenAndServe(listenAddr, nil))
}

func handleGetHealth(w http.ResponseWriter, r *http.Request) {
	w.WriteHeader(http.StatusOK)
	w.Write([]byte("OK!"))
}

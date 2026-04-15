package main

import (
	"log"
	"net/http"

	"github.com/evolvedevlab/weaveset/util"
	"github.com/joho/godotenv"
)

func main() {
	godotenv.Load()
	var (
		listenAddr = util.GetEnv("LISTEN_ADDR", ":3000")
		// redisPass  = util.GetEnv("REDIS_PASSWORD", ":6379")
	)

	http.HandleFunc("/health", handleGetHealth)

	log.Printf("started at %s\n", listenAddr)
	log.Fatal(http.ListenAndServe(listenAddr, nil))
}

func handleGetHealth(w http.ResponseWriter, r *http.Request) {
	w.WriteHeader(http.StatusOK)
	w.Write([]byte("OK!"))
}

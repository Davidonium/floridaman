package main

import (
	"fmt"
	"log"
	"net/http"
	"os"
	"time"

	"github.com/go-redis/redis/v7"
	"github.com/joho/godotenv"

	"github.com/davidonium/floridaman"
)

func main() {

	godotenv.Load()

	logger := log.New(os.Stdout, "", log.LstdFlags)

	client := redis.NewClient(&redis.Options{
		Addr:        GetenvDefault("REDIS_ADDR", "127.0.0.1:6379"),
		Password:    "",
		DB:          0,
		MaxRetries:  3,
		IdleTimeout: 4 * time.Second,
	})

	articleReader := floridaman.NewRedisArticleReader(client)

	mux := http.NewServeMux()
	mux.HandleFunc("/health", floridaman.NewHealthHandler(client))
	mux.HandleFunc("/random", floridaman.NewRandomHandler(logger, articleReader))
	mux.HandleFunc("/random-slack", floridaman.NewSlackRandomHandler(logger, articleReader))
	mux.HandleFunc("/redirect-slack", floridaman.NewSlackOAuthRedirectHandler())

	port := GetenvDefault("APP_PORT", "8081")
	srv := &http.Server{
		Addr:         fmt.Sprintf(":%s", port),
		Handler:      mux,
		ErrorLog:     log.New(os.Stdout, "", log.LstdFlags),
		ReadTimeout:  5 * time.Second,
		WriteTimeout: 10 * time.Second,
		IdleTimeout:  15 * time.Second,
	}

	logger.Printf("Floridaman api listening on port: %s", port)

	err := srv.ListenAndServe()

	logger.Fatalln(err)
}

// GetenvDefault gets the `key` environment variable or returns the default value
func GetenvDefault(key string, d string) string {
	e, ok := os.LookupEnv(key)

	if !ok {
		e = d
	}

	return e
}

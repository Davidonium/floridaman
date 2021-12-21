package main

import (
	"fmt"
	"log"
	"net/http"
	"os"
	"time"

	"github.com/davidonium/floridaman"
	"github.com/go-redis/redis/v8"
	"github.com/joho/godotenv"
)

func ApiServerListen() {
	_, ok := os.LookupEnv("APP_PORT")
	if !ok {
		err := godotenv.Load()
		if err != nil {
			log.Fatalln("Failed to load dotenv environment variables", err)
		}
	}

	logger := log.New(os.Stdout, "", log.LstdFlags)

	client := redis.NewClient(&redis.Options{
		Addr:        GetEnvDefault("REDIS_ADDR", "127.0.0.1:6379"),
		Password:    "",
		DB:          0,
		MaxRetries:  3,
		IdleTimeout: 4 * time.Second,
	})

	articleReader := floridaman.NewRedisArticleReader(client)

	ah := floridaman.NewAPIHandler(logger)

	mux := http.NewServeMux()
	mux.Handle("/health", ah.ToHandler(floridaman.NewHealthHandler(client)))
	mux.Handle("/random", ah.ToHandler(floridaman.NewRandomHandler(articleReader)))
	mux.Handle("/random-slack", ah.ToHandler(floridaman.NewSlackRandomHandler(logger, articleReader, os.Getenv("SLACK_SIGNING_SECRET"))))
	mux.Handle("/redirect-slack", ah.ToHandler(floridaman.NewSlackOAuthRedirectHandler()))

	port := GetEnvDefault("APP_PORT", "8081")
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

// GetEnvDefault gets the `key` environment variable or returns the default value.
func GetEnvDefault(key, d string) string {
	e, ok := os.LookupEnv(key)

	if !ok {
		e = d
	}

	return e
}

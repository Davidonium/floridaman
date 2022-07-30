package main

import (
	"fmt"
	"log"
	"net/http"
	"os"
	"time"

	"github.com/go-redis/redis/v8"

	"github.com/davidonium/floridaman/internal/server"
	"github.com/davidonium/floridaman/internal/storage"
)

func HTTPServerListen(logger *log.Logger) {
	redisClient := redis.NewClient(
		&redis.Options{
			Addr:        GetEnvDefault("REDIS_ADDR", "127.0.0.1:6379"),
			Password:    "",
			DB:          0,
			MaxRetries:  3,
			IdleTimeout: 4 * time.Second,
		},
	)

	articleReader := storage.NewRedisArticleReader(redisClient)

	serv := server.NewServer(
		logger,
		redisClient,
		articleReader,
	)

	port := GetEnvDefault("APP_PORT", "8081")
	srv := &http.Server{
		Addr:              fmt.Sprintf(":%s", port),
		Handler:           serv,
		ErrorLog:          log.New(os.Stdout, "", log.LstdFlags),
		ReadTimeout:       15 * time.Second,
		ReadHeaderTimeout: 10 * time.Second,
		WriteTimeout:      10 * time.Second,
		IdleTimeout:       30 * time.Second,
	}

	logger.Printf("floridaman api listening on port: %s", port)

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

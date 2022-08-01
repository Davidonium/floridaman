package main

import (
	"fmt"
	"log"
	"net/http"
	"os"
	"time"

	"github.com/go-redis/redis/v8"

	"github.com/davidonium/floridaman/internal/util"

	"github.com/davidonium/floridaman/internal/server"
	"github.com/davidonium/floridaman/internal/storage"
)

func HTTPServerListen(logger *log.Logger) {
	redisClient := redis.NewClient(
		&redis.Options{
			Addr:        util.GetEnvDefault("REDIS_ADDR", "127.0.0.1:6379"),
			Password:    "",
			DB:          0,
			MaxRetries:  3,
			IdleTimeout: 4 * time.Second,
		},
	)

	articleStorage := storage.NewRedisArticleStorage(redisClient)

	serv := server.NewServer(
		logger,
		redisClient,
		articleStorage,
	)

	port := util.GetEnvDefault("APP_PORT", "8081")
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

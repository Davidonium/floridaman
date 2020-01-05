package main

import (
	"encoding/json"
	"io"
	"log"
	"net/http"
	"os"
	"time"

	"github.com/go-redis/redis/v7"
)

func main() {

	logger := log.New(os.Stderr, "", log.LstdFlags)

	client := redis.NewClient(&redis.Options{
		Addr:     os.Getenv("REDIS_ADDR"),
		Password: "",
		DB:       0,
	})

	mux := http.NewServeMux()

	mux.HandleFunc("/random", func(w http.ResponseWriter, r *http.Request) {

		key, err := client.RandomKey().Result()

		if err != nil {
			WriteInternalServerError(w)
			logger.Println(err)
			return
		}

		fda, err := client.Get(key).Result()

		if err != nil {
			WriteInternalServerError(w)
			logger.Println(err)
			return
		}

		_, err = io.WriteString(w, fda)

		if err != nil {
			WriteInternalServerError(w)
			logger.Println(err)
			return
		}
	})

	srv := &http.Server{
		Addr:         ":" + os.Getenv("APP_PORT"),
		Handler:      mux,
		ErrorLog:     log.New(os.Stderr, "http: ", log.LstdFlags),
		ReadTimeout:  5 * time.Second,
		WriteTimeout: 10 * time.Second,
		IdleTimeout:  15 * time.Second,
	}

	err := srv.ListenAndServe()

	log.Fatalln(err)
}

func WriteInternalServerError(w http.ResponseWriter) {
	w.WriteHeader(http.StatusInternalServerError)
	json.NewEncoder(w).Encode(struct {
		Message string `json:"message"`
	}{
		Message: "Internal server error",
	})
}

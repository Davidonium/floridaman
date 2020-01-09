package main

import (
	"encoding/json"
	"fmt"
	"io"
	"log"
	"net/http"
	"os"
	"time"

	"github.com/go-redis/redis/v7"
	"github.com/joho/godotenv"
)

func main() {

	godotenv.Load()

	logger := log.New(os.Stdout, "", log.LstdFlags)

	client := redis.NewClient(&redis.Options{
		Addr:     GetenvDefault("REDIS_ADDR", "127.0.0.1:8081"),
		Password: "",
		DB:       0,
	})

	mux := http.NewServeMux()

	mux.HandleFunc("/random", func(w http.ResponseWriter, r *http.Request) {

		w.Header().Set("Content-Type", "application/json")

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
		Addr:         fmt.Sprintf(":%s", GetenvDefault("APP_PORT", "8081")),
		Handler:      mux,
		ErrorLog:     log.New(os.Stderr, "http: ", log.LstdFlags),
		ReadTimeout:  5 * time.Second,
		WriteTimeout: 10 * time.Second,
		IdleTimeout:  15 * time.Second,
	}

	err := srv.ListenAndServe()

	logger.Fatalln(err)
}

func GetenvDefault(key string, d string) string {
	e, ok := os.LookupEnv(key)

	if !ok {
		e = d
	}

	return e
}

func WriteInternalServerError(w http.ResponseWriter) {
	w.WriteHeader(http.StatusInternalServerError)
	w.Header().Set("Content-Type", "application/json")
	json.NewEncoder(w).Encode(struct {
		Message string `json:"message"`
	}{
		Message: "Internal server error",
	})
}

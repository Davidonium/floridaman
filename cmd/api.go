package main

import (
	"encoding/json"
	"fmt"
	"log"
	"net/http"
	"os"
	"time"

	"github.com/go-redis/redis/v7"
	"github.com/joho/godotenv"
	"gitlab.com/davidonium/floridaman"
)

type ErrorResponse struct {
	Message string `json:"message"`
}

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

	mux := http.NewServeMux()

	mux.HandleFunc("/health", func(w http.ResponseWriter, r *http.Request) {
		w.Header().Set("Content-Type", "application/json")

		res, _ := client.Ping().Result()

		if res == "PONG" {
			res = "UP"
		} else {
			res = "DOWN"
		}

		json.NewEncoder(w).Encode(struct {
			Status string `json:"status"`
			Redis  string `json:"redis"`
		}{
			Status: "UP",
			Redis:  res,
		})
	})

	mux.HandleFunc("/random", func(w http.ResponseWriter, r *http.Request) {
		w.Header().Set("Content-Type", "application/json")

		fda, err := floridaman.ReadRandomArticle(client)

		if err != nil {
			WriteInternalServerError(w)
			logger.Printf("%v\n", err)
			return
		}

		_, err = w.Write([]byte(fda))

		if err != nil {
			WriteInternalServerError(w)
			logger.Printf("%v\n", err)
			return
		}
	})

	mux.HandleFunc("/random-slack", func(w http.ResponseWriter, r *http.Request) {
		w.Header().Set("Content-Type", "application/json")

		ok := floridaman.ValidateSlackRequest(r, logger)

		if !ok {
			logger.Printf("invalid slack request %v\n", r)
			w.WriteHeader(http.StatusBadRequest)
			json.NewEncoder(w).Encode(ErrorResponse{Message: "Invalid slack request"})
			return
		}

		fda, err := floridaman.ReadRandomArticle(client)

		if err != nil {
			WriteInternalServerError(w)
			logger.Printf("%v\n", err)
			return
		}

		article := &floridaman.Article{}
		json.Unmarshal([]byte(fda), article)

		response := floridaman.SlackResponse{
			Text: fmt.Sprintf("%s (%s)", article.Title, article.Link),
		}
		json.NewEncoder(w).Encode(response)
	})

	srv := &http.Server{
		Addr:         fmt.Sprintf(":%s", GetenvDefault("APP_PORT", "8081")),
		Handler:      mux,
		ErrorLog:     log.New(os.Stdout, "", log.LstdFlags),
		ReadTimeout:  5 * time.Second,
		WriteTimeout: 10 * time.Second,
		IdleTimeout:  15 * time.Second,
	}

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

func WriteInternalServerError(w http.ResponseWriter) {
	w.WriteHeader(http.StatusInternalServerError)
	w.Header().Set("Content-Type", "application/json")
	json.NewEncoder(w).Encode(ErrorResponse{Message: "Internal server error"})
}

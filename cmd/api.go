package main

import (
	"bytes"
	"crypto/hmac"
	"crypto/sha256"
	"encoding/hex"
	"encoding/json"
	"fmt"
	"io"
	"io/ioutil"
	"log"
	"net/http"
	"os"
	"strconv"
	"time"

	"github.com/go-redis/redis/v7"
	"gitlab.com/davidonium/floridaman"
)

func main() {

	//godotenv.Load()

	logger := log.New(os.Stdout, "", log.LstdFlags)

	logger.Printf("ss %s", os.Getenv("SLACK_SIGNING_SECRET"))
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

		fda, err := floridaman.ReadRandomArticleFromRedis(client)

		if err != nil {
			WriteInternalServerError(w)
			logger.Printf("%v\n", err)
			return
		}

		_, err = io.WriteString(w, fda)

		if err != nil {
			WriteInternalServerError(w)
			logger.Printf("%v\n", err)
			return
		}
	})

	mux.HandleFunc("/random-slack", func(w http.ResponseWriter, r *http.Request) {
		w.Header().Set("Content-Type", "application/json")

		ok := ValidateSlackRequest(r, logger)

		if !ok {
			logger.Printf("invalid slack request %v\n", r)
			w.WriteHeader(http.StatusBadRequest)
			json.NewEncoder(w).Encode(struct {
				Message string `json:"message"`
			}{
				Message: "Invalid slack request",
			})
			return
		}

		fda, err := floridaman.ReadRandomArticleFromRedis(client)

		if err != nil {
			WriteInternalServerError(w)
			logger.Printf("%v\n", err)
			return
		}

		article := &floridaman.Article{}
		json.Unmarshal([]byte(fda), article)

		json.NewEncoder(w).Encode(struct {
			Text string `json:"text"`
		}{
			Text: article.Title,
		})
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
	json.NewEncoder(w).Encode(struct {
		Message string `json:"message"`
	}{
		Message: "Internal server error",
	})
}

func ValidateSlackRequest(r *http.Request, logger *log.Logger) bool {
	s := r.Header.Get("X-Slack-Signature")
	t := r.Header.Get("X-Slack-Request-Timestamp")

	ts, err := strconv.ParseInt(t, 10, 64)

	if err != nil {
		return false
	}

	tsu := time.Unix(ts, 0)

	if time.Now().Sub(tsu) > 5*time.Minute {
		logger.Println("timestamp difference is greater than 5 minutes")
		return false
	}

	body, err := ioutil.ReadAll(r.Body)

	defer r.Body.Close()

	if err != nil {
		return false
	}

	r.Body = ioutil.NopCloser(bytes.NewBuffer(body))

	msg := fmt.Sprintf("v0:%s:%s", t, body)

	ss := os.Getenv("SLACK_SIGNING_SECRET")
	logger.Printf("before hashing ss: %s | %s\n", ss, msg)

	sig := SlackHashHMAC([]byte(msg), []byte(ss))

	ok := hmac.Equal([]byte(sig), []byte(s))
	if !ok {
		logger.Printf("error validating hmac signature from slack: %s, generated %s\n", s, sig)
		return false
	}

	return true
}

func SlackHashHMAC(msg, key []byte) string {
	hm := hmac.New(sha256.New, key)
	hm.Write(msg)
	finalHash := hm.Sum(nil)
	return fmt.Sprintf("v0=%s", hex.EncodeToString(finalHash))
}

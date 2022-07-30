package server

import (
	"net/http"

	"github.com/go-redis/redis/v8"
)

type healthResponse struct {
	Status string `json:"status"`
	Redis  string `json:"redis"`
}

func (s *Server) healthHandler(client *redis.Client) APIHandlerFunc {
	return func(w http.ResponseWriter, r *http.Request) error {
		res, _ := client.Ping(r.Context()).Result()

		if res == "PONG" {
			res = "UP"
		} else {
			res = "DOWN"
		}

		response := &healthResponse{
			Status: "UP",
			Redis:  res,
		}
		return s.writeJSON(w, response)
	}
}

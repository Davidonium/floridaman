package floridaman

import (
	"encoding/json"
	"net/http"

	"github.com/go-redis/redis/v7"
)

type healthResponse struct {
	Status string `json:"status"`
	Redis  string `json:"redis"`
}

func NewHealthHandler(client *redis.Client) ApiHandlerFunc {
	return func(w http.ResponseWriter, r *http.Request) error {
		res, _ := client.Ping().Result()

		if res == "PONG" {
			res = "UP"
		} else {
			res = "DOWN"
		}

		response := &healthResponse{
			Status: "UP",
			Redis:  res,
		}
		return json.NewEncoder(w).Encode(response)
	}
}

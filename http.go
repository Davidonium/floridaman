package floridaman

import (
	"encoding/json"
	"errors"
	"log"
	"net/http"
)

type ErrorResponse struct {
	Message string `json:"message"`
}

func WriteInternalServerError(w http.ResponseWriter) {
	w.WriteHeader(http.StatusInternalServerError)
	json.NewEncoder(w).Encode(ErrorResponse{Message: "Internal server error"})
}

type ApiHandlerFunc func(http.ResponseWriter, *http.Request) error

type ApiHandler struct {
	logger *log.Logger
}

func (eh ApiHandler) ToHandler(handler ApiHandlerFunc) http.Handler {
	return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		w.Header().Set("Content-Type", "application/json")

		err := handler(w, r)

		if err != nil {
			if errors.Is(err, ErrInvalidSlackRequest) {
				w.WriteHeader(http.StatusBadRequest)
				json.NewEncoder(w).Encode(ErrorResponse{Message: err.Error()})
				return
			}
			eh.logger.Printf("unhandled error in request: %v\n", err)
			WriteInternalServerError(w)
			return
		}
	})
}

func NewApiHandler(logger *log.Logger) ApiHandler {
	return ApiHandler{logger: logger}
}

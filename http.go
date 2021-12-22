package floridaman

import (
	"encoding/json"
	"errors"
	"fmt"
	"log"
	"net/http"
)

type ErrorResponse struct {
	Message string `json:"message"`
}

func WriteInternalServerError(w http.ResponseWriter) {
	w.WriteHeader(http.StatusInternalServerError)
	err := json.NewEncoder(w).Encode(ErrorResponse{Message: "Internal server error"})
	if err != nil {
		fmt.Println("could not write response")	
	}
}

type APIHandlerFunc func(http.ResponseWriter, *http.Request) error

type APIHandler struct {
	logger *log.Logger
}

func (eh APIHandler) ToHandler(handler APIHandlerFunc) http.Handler {
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

func NewAPIHandler(logger *log.Logger) APIHandler {
	return APIHandler{logger: logger}
}

package server

import (
	"encoding/json"
	"errors"
	"fmt"
	"log"
	"net/http"
	"os"

	"github.com/davidonium/floridaman/internal/floridaman"
	"github.com/go-redis/redis/v8"
)

type ErrorResponse struct {
	Message string `json:"message"`
}

type APIHandlerFunc func(http.ResponseWriter, *http.Request) error

type Server struct {
	logger *log.Logger
	mux    *http.ServeMux

	redisClient   *redis.Client
	articleReader floridaman.ArticleReader
}

func NewServer(logger *log.Logger, redisClient *redis.Client, articleReader floridaman.ArticleReader) *Server {
	srv := &Server{
		logger:        logger,
		mux:           http.NewServeMux(),
		redisClient:   redisClient,
		articleReader: articleReader,
	}
	srv.routes()

	return srv
}

func (s *Server) ServeHTTP(rw http.ResponseWriter, r *http.Request) {
	s.mux.ServeHTTP(rw, r)
}

func (s *Server) routes() {
	s.mux.Handle("/health", s.handleAPI(s.healthHandler(s.redisClient)))
	s.mux.Handle("/random", s.handleAPI(s.randomArticleHandler(s.articleReader)))
	s.mux.Handle(
		"/slack/random",
		s.handleAPI(s.slackRandomArticleHandler(s.logger, s.articleReader, os.Getenv("SLACK_SIGNING_SECRET"))),
	)
	s.mux.Handle("/slack/redirect", s.handleAPI(s.oauthSlackRedirectHandler()))
}

func (s *Server) handleAPI(handler APIHandlerFunc) http.HandlerFunc {
	return func(w http.ResponseWriter, r *http.Request) {
		w.Header().Set("Content-Type", "application/json")

		err := handler(w, r)
		if err != nil {
			if errors.Is(err, ErrInvalidSlackRequest) {
				w.WriteHeader(http.StatusBadRequest)
				b, _ := json.Marshal(ErrorResponse{Message: err.Error()})
				_, _ = w.Write(b)
				return
			}
			s.logger.Printf("unhandled error in request: %v\n", err)
			s.writeInternalError(w)
			return
		}
	}
}

func (s *Server) writeJSON(w http.ResponseWriter, data any) error {
	b, err := json.Marshal(data)

	if err != nil {
		s.logger.Printf("could not serialize response to json")
		return err
	}

	w.WriteHeader(http.StatusOK)
	_, _ = w.Write(b)

	return nil
}

func (s *Server) writeInternalError(w http.ResponseWriter) {
	w.WriteHeader(http.StatusInternalServerError)

	b, _ := json.Marshal(ErrorResponse{Message: "Internal server error"})
	_, err := w.Write(b)
	if err != nil {
		fmt.Println("could not write response")
	}
}

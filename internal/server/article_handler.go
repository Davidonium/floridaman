package server

import (
	"net/http"

	"github.com/davidonium/floridaman/internal/floridaman"
)

func (s *Server) randomArticleHandler(ar floridaman.ArticleReader) APIHandlerFunc {
	return func(w http.ResponseWriter, r *http.Request) error {
		fda, err := ar.RawRandom(r.Context())
		if err != nil {
			return err
		}

		_, err = w.Write([]byte(fda))

		if err != nil {
			return err
		}

		return nil
	}
}

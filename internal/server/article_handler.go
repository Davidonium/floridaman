package server

import (
	"net/http"

	"github.com/davidonium/floridaman/internal/floridaman"
)

func (s *Server) randomArticleHandler(articleStorage floridaman.ArticleStorage) APIHandlerFunc {
	return func(w http.ResponseWriter, r *http.Request) error {
		rawArticle, err := articleStorage.RawRandom(r.Context())
		if err != nil {
			return err
		}

		return s.writeBytes(w, rawArticle)
	}
}

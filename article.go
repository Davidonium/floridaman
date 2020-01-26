package floridaman

import (
	"log"
	"net/http"

	"github.com/go-redis/redis/v7"
	"github.com/turnage/graw/reddit"
)

type Article struct {
	Title  string `json:"title"`
	Link   string `json:"link"`
	Source string `json:"source"`
}

func NewArticleFromReddit(post *reddit.Post) Article {
	return Article{
		Title:  post.Title,
		Link:   post.URL,
		Source: "reddit",
	}
}

func NewRandomHandler(logger *log.Logger, client *redis.Client) http.HandlerFunc {
	return func(w http.ResponseWriter, r *http.Request) {
		w.Header().Set("Content-Type", "application/json")

		fda, err := ReadRandomArticle(client)

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
	}
}

func ReadRandomArticle(client *redis.Client) (string, error) {
	key, err := client.RandomKey().Result()

	if err != nil {
		return "", err
	}

	return client.Get(key).Result()
}

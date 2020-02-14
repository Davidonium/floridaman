package floridaman

import (
	"encoding/json"
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

type ArticleReader interface {
	Random() (Article, error)
	RawRandom() (string, error)
}

type RedisArticleReader struct {
	client *redis.Client
}

func NewRedisArticleReader(client *redis.Client) *RedisArticleReader {
	return &RedisArticleReader{client: client}
}

func (r *RedisArticleReader) Random() (Article, error) {
	ra, err := r.RawRandom()

	if err != nil {
		return Article{}, err
	}

	a := Article{}
	err = json.Unmarshal([]byte(ra), &a)

	if err != nil {
		return Article{}, err
	}

	return a, nil

}

func (r *RedisArticleReader) RawRandom() (string, error) {
	key, err := r.client.RandomKey().Result()

	if err != nil {
		return "", err
	}

	return r.client.Get(key).Result()
}

func NewArticleFromReddit(post *reddit.Post) Article {
	return Article{
		Title:  post.Title,
		Link:   post.URL,
		Source: "reddit",
	}
}

func NewRandomHandler(logger *log.Logger, ar ArticleReader) http.HandlerFunc {
	return func(w http.ResponseWriter, r *http.Request) {
		w.Header().Set("Content-Type", "application/json")

		fda, err := ar.RawRandom()

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

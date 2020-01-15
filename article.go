package floridaman

import (
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

func ReadRandomArticle(client *redis.Client) (string, error) {
	key, err := client.RandomKey().Result()

	if err != nil {
		return "", err
	}

	return client.Get(key).Result()
}

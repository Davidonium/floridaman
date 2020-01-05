package floridaman

import (
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

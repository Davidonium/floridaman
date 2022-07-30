package floridaman

import "context"

type Article struct {
	Title  string `json:"title"`
	Link   string `json:"link"`
	Source string `json:"source"`
}

type ArticleReader interface {
	Random(context.Context) (Article, error)
	RawRandom(context.Context) (string, error)
}

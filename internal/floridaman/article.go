package floridaman

import "context"

type Article struct {
	Title  string `json:"title"`
	Link   string `json:"link"`
	Source string `json:"source"`
}

type ArticleStorage interface {
	Random(context.Context) (Article, error)
	RawRandom(context.Context) ([]byte, error)
	Save(context.Context, Article)
	ExistsByTitle(context.Context, string) (bool, error)
}

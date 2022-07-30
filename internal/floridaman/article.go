package floridaman

type Article struct {
	Title  string `json:"title"`
	Link   string `json:"link"`
	Source string `json:"source"`
}

type ArticleReader interface {
	Random() (Article, error)
	RawRandom() (string, error)
}

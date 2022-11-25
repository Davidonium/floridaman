package storage

import (
	"context"
	"encoding/json"
	"fmt"

	"github.com/go-redis/redis/v8"

	"github.com/davidonium/floridaman/internal/util"

	"github.com/davidonium/floridaman/internal/floridaman"
)

type RedisArticleStorage struct {
	client *redis.Client
}

func NewRedisArticleStorage(client *redis.Client) *RedisArticleStorage {
	return &RedisArticleStorage{client: client}
}

func (r *RedisArticleStorage) Random(ctx context.Context) (floridaman.Article, error) {
	rawArticle, err := r.RawRandom(ctx)
	if err != nil {
		return floridaman.Article{}, err
	}

	a := floridaman.Article{}
	err = json.Unmarshal(rawArticle, &a)
	if err != nil {
		return floridaman.Article{}, err
	}

	return a, nil
}

func (r *RedisArticleStorage) RawRandom(ctx context.Context) ([]byte, error) {
	key, err := r.client.RandomKey(ctx).Result()
	if err != nil {
		return nil, err
	}

	return r.client.Get(ctx, key).Bytes()
}

func (r *RedisArticleStorage) Save(ctx context.Context, article floridaman.Article) {
	key := r.generateKey(article.Title)
	j, _ := json.Marshal(article)
	r.client.Set(ctx, key, string(j), 0)
}

func (r *RedisArticleStorage) ExistsByTitle(ctx context.Context, title string) (bool, error) {
	key := r.generateKey(title)

	n, err := r.client.Exists(ctx, key).Result()
	if err != nil {
		return false, err
	}

	return n > 0, nil
}

func (r *RedisArticleStorage) generateKey(title string) string {
	return fmt.Sprintf("fm:%s", util.SHA1String(title))
}

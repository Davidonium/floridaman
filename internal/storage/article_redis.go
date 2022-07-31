package storage

import (
	"context"
	"encoding/json"

	"github.com/go-redis/redis/v8"

	"github.com/davidonium/floridaman/internal/floridaman"
)

type RedisArticleReader struct {
	client *redis.Client
}

func NewRedisArticleReader(client *redis.Client) *RedisArticleReader {
	return &RedisArticleReader{client: client}
}

func (r *RedisArticleReader) Random(ctx context.Context) (floridaman.Article, error) {
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

func (r *RedisArticleReader) RawRandom(ctx context.Context) ([]byte, error) {
	key, err := r.client.RandomKey(ctx).Result()
	if err != nil {
		return nil, err
	}

	return r.client.Get(ctx, key).Bytes()
}

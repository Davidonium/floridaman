package main

import (
	"context"
	"log"
	"os"

	"github.com/go-redis/redis/v8"
	"github.com/vartanbeno/go-reddit/v2/reddit"

	"github.com/davidonium/floridaman/internal/storage"

	"github.com/davidonium/floridaman/internal/floridaman"
)

func ReadRedditArticles(logger *log.Logger) {
	credentials := reddit.Credentials{
		ID:       os.Getenv("REDDIT_CLIENT_ID"),
		Secret:   os.Getenv("REDDIT_CLIENT_SECRET"),
		Username: os.Getenv("REDDIT_CLIENT_USERNAME"),
		Password: os.Getenv("REDDIT_CLIENT_PASSWORD"),
	}
	redditClient, err := reddit.NewClient(
		credentials,
		reddit.WithUserAgent("script:github.com/davidonium/floridaman:v0.1 (by tindrem)"),
	)
	if err != nil {
		logger.Fatalln("failed to create reddit client", err)
	}

	client := redis.NewClient(
		&redis.Options{
			Addr:     os.Getenv("REDIS_ADDR"),
			Password: "",
			DB:       0,
		},
	)

	if err != nil {
		logger.Fatalln("failed to create redis client", err)
	}

	articleStorage := storage.NewRedisArticleStorage(client)

	after := ""

	for {
		logger.Printf("requesting /r/FloridaMan/top after=%s\n", after)

		ctx := context.Background()

		opts := &reddit.ListPostOptions{
			ListOptions: reddit.ListOptions{
				Limit: 100,
				After: after,
			},
			Time: "all",
		}
		posts, response, err := redditClient.Subreddit.TopPosts(ctx, "floridaman", opts)
		if err != nil {
			logger.Fatalf("failed to fetch /r/FloridaMan?after=%s err: %v\n", after, err)
		}

		for _, post := range posts {
			article := articleFromReddit(post)

			ok, err := articleStorage.ExistsByTitle(ctx, article.Title)
			if err != nil {
				logger.Printf("failed to check that article \"%s\" exists, skipping: %v\n", article.Title, err)
				continue
			}
			if ok {
				logger.Printf("floridaman article \"%s\" already exists\n", article.Title)
			} else {
				articleStorage.Save(ctx, article)
			}
		}

		after = response.After

		if after == "" {
			break
		}
	}

	logger.Println("finished reading reddit floridaman articles")
}

func articleFromReddit(post *reddit.Post) floridaman.Article {
	return floridaman.Article{
		Title:  post.Title,
		Link:   post.URL,
		Source: "reddit",
	}
}

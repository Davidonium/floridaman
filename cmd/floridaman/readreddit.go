package main

import (
	"context"
	"encoding/json"
	"fmt"
	"log"
	"os"

	"github.com/go-redis/redis/v8"
	"github.com/vartanbeno/go-reddit/v2/reddit"

	"github.com/davidonium/floridaman/internal/floridaman"
	"github.com/davidonium/floridaman/internal/util"
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
		logger.Fatalln("failed to create bot", err)
	}

	first := true
	after := ""

	for first || len(after) > 0 {
		if first {
			first = false
		}

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
			fma := articleFromReddit(post)

			h := util.SHA1String(fma.Title)
			key := fmt.Sprintf("fm:%s", h)

			ex, err := client.Exists(ctx, key).Result()
			if err != nil {
				logger.Printf("failed to check that key \"%s\" exists, skipping\n", key)
				continue
			}
			if ex > 0 {
				logger.Printf("floridaman article with key \"%s\" already exists\n", key)
			} else {
				j, _ := json.Marshal(fma)
				client.Set(ctx, key, string(j), 0)
			}
		}

		after = response.After
	}
}

func articleFromReddit(post *reddit.Post) floridaman.Article {
	return floridaman.Article{
		Title:  post.Title,
		Link:   post.URL,
		Source: "reddit",
	}
}

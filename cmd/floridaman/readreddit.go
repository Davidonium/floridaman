package main

import (
	"context"
	"encoding/json"
	"fmt"
	"log"
	"os"

	"github.com/go-redis/redis/v8"
	"github.com/vartanbeno/go-reddit/v2/reddit"

	"github.com/davidonium/floridaman"
)

func ReadRedditArticles(logger *log.Logger) {
	creds := reddit.Credentials{
		ID:       os.Getenv("REDDIT_CLIENT_ID"),
		Secret:   os.Getenv("REDDIT_CLIENT_SECRET"),
		Username: os.Getenv("REDDIT_CLIENT_USERNAME"),
		Password: os.Getenv("REDDIT_CLIENT_PASSWORD"),
	}
	redditClient, err := reddit.NewClient(
		creds,
		reddit.WithUserAgent("script:github.com/davidonium/floridaman:v0.1 (by tindrem)"),
	)
	if err != nil {
		logger.Fatalln("Failed to create reddit client", err)
	}

	client := redis.NewClient(&redis.Options{
		Addr:     os.Getenv("REDIS_ADDR"),
		Password: "",
		DB:       0,
	})

	if err != nil {
		logger.Fatalln("Failed to create bot", err)
	}

	first := true
	np := ""

	for first || len(np) > 0 {
		if first {
			first = false
		}

		logger.Printf("requesting /r/FloridaMan/top after=%s\n", np)

		opts := &reddit.ListPostOptions{
			ListOptions: reddit.ListOptions{
				Limit: 100,
				After: np,
			},
			Time: "all",
		}
		posts, _, err := redditClient.Subreddit.TopPosts(context.Background(), "floridaman", opts)
		if err != nil {
			log.Fatalf("Failed to fetch /r/rFloridaMan?after=%s err: %v\n", np, err)
		}

		postNum := len(posts)
		if postNum > 0 {
			for _, post := range posts {
				fma := floridaman.NewArticleFromReddit(post)

				h := floridaman.SHA1String(fma.Title)
				key := fmt.Sprintf("fm:%s", h)

				ex, _ := client.Exists(context.Background(), key).Result()
				if ex > 0 {
					logger.Printf("Floridaman article with key \"%s\" already exists\n", key)
				} else {
					j, _ := json.Marshal(fma)
					client.Set(context.Background(), key, string(j), 0)
				}
			}

			np = posts[postNum-1].FullID
		} else {
			np = ""
		}
	}
}

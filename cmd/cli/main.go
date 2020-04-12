package main

import (
	"encoding/json"
	"fmt"
	"log"
	"os"
	"time"

	"github.com/go-redis/redis/v7"
	"github.com/joho/godotenv"
	"github.com/turnage/graw/reddit"

	"github.com/davidonium/floridaman"
)

func main() {
	logger := log.New(os.Stdout, "", log.LstdFlags)

	_, ok := os.LookupEnv("REDDIT_CLIENT_ID")
	if !ok {
		godotenv.Load()
	}

	bcfg := reddit.BotConfig{
		Agent: "script:github.com/davidonium/floridaman:v0.1 (by tindrem)",
		App: reddit.App{
			ID:       os.Getenv("REDDIT_CLIENT_ID"),
			Secret:   os.Getenv("REDDIT_CLIENT_SECRET"),
			Username: os.Getenv("REDDIT_CLIENT_USERNAME"),
			Password: os.Getenv("REDDIT_CLIENT_PASSWORD"),
		},
		Rate: time.Second,
	}

	client := redis.NewClient(&redis.Options{
		Addr:     os.Getenv("REDIS_ADDR"),
		Password: "",
		DB:       0,
	})

	bot, err := reddit.NewBot(bcfg)

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

		params := map[string]string{
			"after": np,
			"limit": "100",
			"t":     "all",
		}
		harvest, err := bot.ListingWithParams("/r/FloridaMan/top", params)
		if err != nil {
			log.Fatalf("Failed to fetch /r/FloridaMan?after=%s %v\n", np, err)
		}

		posts := harvest.Posts

		if len(posts) > 0 {
			for _, post := range posts {
				fma := floridaman.NewArticleFromReddit(post)

				h := floridaman.SHA1String(fma.Title)
				key := fmt.Sprintf("fm:%s", h)

				ex, _ := client.Exists(key).Result()
				if ex > 0 {
					logger.Printf("Floridaman article with key \"%s\" already exists\n", key)
				} else {
					j, _ := json.Marshal(fma)
					client.Set(key, string(j), 0)
				}
			}
			lp := posts[len(posts)-1]
			np = lp.Name
		} else {
			np = ""
		}
	}
}

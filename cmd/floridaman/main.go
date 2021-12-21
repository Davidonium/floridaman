package main

import (
	"errors"
	"flag"
	"fmt"
	"log"
	"os"

	"github.com/joho/godotenv"
)

func main() {
	flag.Parse()

	args := flag.Args()

	if len(args) < 1 {
		fmt.Fprintln(os.Stderr, "missing command, options are \"serve\" or \"readreddit\"")
		os.Exit(1)
	}

	err := godotenv.Load()
	if err != nil {
		if !errors.Is(err, os.ErrNotExist) {
			log.Fatalln("Failed to load dotenv environment variables, err:", err)
		}
	}

	logger := log.New(os.Stdout, "", log.LstdFlags)

	switch args[0] {
	case "serve":
		ApiServerListen(logger)
	case "readreddit":
		ReadRedditArticles(logger)
	default:
		fmt.Fprintf(os.Stderr, "unknown command %s, options are \"serve\" or \"readreddit\"", args[0])
		os.Exit(1)
	}
}

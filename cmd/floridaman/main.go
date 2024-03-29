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
		printErrf("missing command, options are \"serve\" or \"readreddit\"")
		os.Exit(1)
	}

	err := godotenv.Load()
	if err != nil && !errors.Is(err, os.ErrNotExist) {
		printErrf("failed to load dotenv environment variables, err: %v\n", err)
		os.Exit(1)
	}

	logger := log.New(os.Stdout, "", log.LstdFlags)

	switch args[0] {
	case "serve":
		HTTPServerListen(logger)
	case "readreddit":
		ReadRedditArticles(logger)
	default:
		printErrf("unknown command %s, options are \"serve\" or \"readreddit\"", args[0])
		os.Exit(1)
	}
}

func printErrf(msg string, args ...any) {
	_, _ = fmt.Fprintf(os.Stderr, msg+"\n", args...)
}

package main

import (
	"flag"
	"fmt"
	"os"
)

func main() {
	flag.Parse()

	args := flag.Args()

	if len(args) < 1 {
		fmt.Fprintln(os.Stderr, "missing command, options are \"serve\" or \"readreddit\"")
		os.Exit(1)
	}

	switch args[0] {
	case "serve":
		ApiServerListen()
	case "readreddit":
		ReadRedditArticles()
	default:
		fmt.Fprintf(os.Stderr, "unknown command %s, options are \"serve\" or \"readreddit\"", args[0])
		os.Exit(1)
	}
}

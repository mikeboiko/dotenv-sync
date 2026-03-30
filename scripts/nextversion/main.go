package main

import (
	"context"
	"flag"
	"fmt"
	"os"
	"path/filepath"

	"dotenv-sync/internal/release"
)

func main() {
	var bump string
	var dir string
	flag.StringVar(&bump, "bump", "", "semantic version bump to apply: major, minor, or patch")
	flag.StringVar(&dir, "dir", ".", "repository directory to inspect")
	flag.Parse()

	if bump == "" {
		fmt.Fprintln(os.Stderr, "missing required --bump value")
		os.Exit(1)
	}

	repoDir, err := filepath.Abs(dir)
	if err != nil {
		fmt.Fprintln(os.Stderr, err)
		os.Exit(1)
	}
	next, err := release.NextVersionForRepo(context.Background(), repoDir, bump)
	if err != nil {
		fmt.Fprintln(os.Stderr, err)
		os.Exit(1)
	}
	fmt.Fprintln(os.Stdout, next)
}

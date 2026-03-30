package main

import (
	"context"
	"flag"
	"fmt"
	"io"
	"os"
	"path/filepath"

	"dotenv-sync/internal/release"
)

const alreadyReleasedExitCode = 2

func main() {
	os.Exit(run(context.Background(), os.Args[1:], os.Stdout, os.Stderr))
}

func run(ctx context.Context, args []string, stdout, stderr io.Writer) int {
	var dir string
	flags := flag.NewFlagSet("nextversion", flag.ContinueOnError)
	flags.SetOutput(stderr)
	flags.StringVar(&dir, "dir", ".", "repository directory to inspect")
	if err := flags.Parse(args); err != nil {
		return 1
	}
	if flags.NArg() != 0 {
		fmt.Fprintln(stderr, "nextversion does not accept positional arguments")
		return 1
	}

	repoDir, err := filepath.Abs(dir)
	if err != nil {
		fmt.Fprintln(stderr, err)
		return 1
	}
	tag, released, err := release.ReleaseTagForRef(ctx, repoDir, "HEAD")
	if err != nil {
		fmt.Fprintln(stderr, err)
		return 1
	}
	if released {
		fmt.Fprintln(stdout, tag)
		fmt.Fprintf(stderr, "commit already released by tag %s\n", tag)
		return alreadyReleasedExitCode
	}

	next, err := release.NextPatchVersionForRepo(ctx, repoDir)
	if err != nil {
		fmt.Fprintln(stderr, err)
		return 1
	}
	exists, err := release.VersionTagExists(ctx, repoDir, next)
	if err != nil {
		fmt.Fprintln(stderr, err)
		return 1
	}
	if exists {
		fmt.Fprintf(stderr, "next release tag %s already exists; inspect repository tags before retrying\n", next)
		return 1
	}
	fmt.Fprintln(stdout, next)
	return 0
}

package main

import (
	"flag"
	"fmt"
	"io"
	"os"
	"path/filepath"

	"dotenv-sync/internal/aur"
)

func main() {
	os.Exit(run(os.Args[1:], os.Stdout, os.Stderr))
}

func run(args []string, stdout, stderr io.Writer) int {
	var (
		licenseSHA256 string
		readmeSHA256  string
		pkgRel        int
		version       string
		amd64SHA256   string
		aarch64SHA256 string
		outputDir     string
	)

	flags := flag.NewFlagSet("aurpkg", flag.ContinueOnError)
	flags.SetOutput(stderr)
	flags.StringVar(&licenseSHA256, "license-sha256", "", "sha256 checksum for LICENSE at the tagged upstream revision")
	flags.StringVar(&readmeSHA256, "readme-sha256", "", "sha256 checksum for README.md at the tagged upstream revision")
	flags.IntVar(&pkgRel, "pkgrel", 1, "Arch package release number")
	flags.StringVar(&version, "version", "", "release tag to package, for example v1.2.3")
	flags.StringVar(&amd64SHA256, "x86_64-sha256", "", "sha256 checksum for the linux amd64 release asset")
	flags.StringVar(&aarch64SHA256, "aarch64-sha256", "", "sha256 checksum for the linux arm64 release asset")
	flags.StringVar(&outputDir, "output-dir", ".", "directory where PKGBUILD and .SRCINFO will be written")
	if err := flags.Parse(args); err != nil {
		return 1
	}
	if flags.NArg() != 0 {
		fmt.Fprintln(stderr, "aurpkg does not accept positional arguments")
		return 1
	}
	if version == "" || licenseSHA256 == "" || readmeSHA256 == "" || amd64SHA256 == "" || aarch64SHA256 == "" {
		fmt.Fprintln(stderr, "aurpkg requires --version, --license-sha256, --readme-sha256, --x86_64-sha256, and --aarch64-sha256")
		return 1
	}

	pkg, err := aur.NewDotenvSyncBin(version, pkgRel, licenseSHA256, readmeSHA256, amd64SHA256, aarch64SHA256)
	if err != nil {
		fmt.Fprintln(stderr, err)
		return 1
	}

	if err := os.MkdirAll(outputDir, 0o755); err != nil {
		fmt.Fprintln(stderr, err)
		return 1
	}
	if err := os.WriteFile(filepath.Join(outputDir, "PKGBUILD"), []byte(pkg.PKGBUILD()), 0o644); err != nil {
		fmt.Fprintln(stderr, err)
		return 1
	}
	if err := os.WriteFile(filepath.Join(outputDir, ".SRCINFO"), []byte(pkg.SRCINFO()), 0o644); err != nil {
		fmt.Fprintln(stderr, err)
		return 1
	}
	fmt.Fprintf(stdout, "wrote PKGBUILD and .SRCINFO to %s\n", outputDir)
	return 0
}

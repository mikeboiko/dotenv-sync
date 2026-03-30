package dotenvsync

import (
	"fmt"
	"runtime"
	"strings"
)

var (
	Version   = "dev"
	Commit    = "none"
	BuildTime = "unknown"
)

type Metadata struct {
	Version   string
	Commit    string
	BuildTime string
	Platform  string
}

func Current() Metadata {
	version := strings.TrimSpace(Version)
	if version == "" {
		version = "dev"
	}
	commit := strings.TrimSpace(Commit)
	if commit == "" {
		commit = "none"
	}
	buildTime := strings.TrimSpace(BuildTime)
	if buildTime == "" {
		buildTime = "unknown"
	}
	return Metadata{
		Version:   version,
		Commit:    commit,
		BuildTime: buildTime,
		Platform:  runtime.GOOS + "/" + runtime.GOARCH,
	}
}

func Short(name string) string {
	meta := Current()
	return strings.TrimSpace(name) + " " + meta.Version
}

func Detailed() string {
	meta := Current()
	return fmt.Sprintf("Version: %s\nCommit: %s\nBuilt: %s\nPlatform: %s", meta.Version, meta.Commit, meta.BuildTime, meta.Platform)
}

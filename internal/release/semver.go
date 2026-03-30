package release

import (
	"context"
	"fmt"
	"os/exec"
	"regexp"
	"strconv"
	"strings"
)

var semverTagPattern = regexp.MustCompile(`^v(\d+)\.(\d+)\.(\d+)$`)

type SemVer struct {
	Major int
	Minor int
	Patch int
}

func LatestVersionFromTags(tags []string) (string, bool) {
	var latest SemVer
	found := false
	for _, tag := range tags {
		version, ok := Parse(tag)
		if !ok {
			continue
		}
		if !found || version.Compare(latest) > 0 {
			latest = version
			found = true
		}
	}
	if !found {
		return "", false
	}
	return latest.String(), true
}

func NextVersion(currentTag, bump string) (string, error) {
	version, ok := Parse(currentTag)
	if !ok {
		return "", fmt.Errorf("invalid semantic version tag %q", currentTag)
	}
	switch bump {
	case "patch":
		version.Patch++
	case "minor":
		version.Minor++
		version.Patch = 0
	case "major":
		version.Major++
		version.Minor = 0
		version.Patch = 0
	default:
		return "", fmt.Errorf("unsupported bump %q", bump)
	}
	return version.String(), nil
}

func AssetName(version, goos, goarch string) string {
	if goos == "windows" {
		return fmt.Sprintf("ds_%s_windows_%s.zip", version, goarch)
	}
	return fmt.Sprintf("ds_%s_%s_%s.tar.gz", version, goos, goarch)
}

func LatestVersionFromRepo(ctx context.Context, dir string) (string, bool, error) {
	cmd := exec.CommandContext(ctx, "git", "tag", "--merged", "HEAD", "--list")
	cmd.Dir = dir
	out, err := cmd.CombinedOutput()
	if err != nil {
		return "", false, fmt.Errorf("list git tags: %w (%s)", err, strings.TrimSpace(string(out)))
	}
	lines := strings.Fields(string(out))
	version, ok := LatestVersionFromTags(lines)
	return version, ok, nil
}

func NextVersionForRepo(ctx context.Context, dir, bump string) (string, error) {
	current := "v0.0.0"
	if latest, ok, err := LatestVersionFromRepo(ctx, dir); err != nil {
		return "", err
	} else if ok {
		current = latest
	}
	return NextVersion(current, bump)
}

func ValidateReleaseBranch(currentRef, defaultBranch string) error {
	currentRef = strings.TrimSpace(currentRef)
	defaultBranch = strings.TrimSpace(defaultBranch)
	if currentRef == "" || defaultBranch == "" {
		return fmt.Errorf("release branch validation requires both current and default branch names")
	}
	if currentRef != defaultBranch {
		return fmt.Errorf("release workflow must run from %s, got %s", defaultBranch, currentRef)
	}
	return nil
}

func Parse(tag string) (SemVer, bool) {
	matches := semverTagPattern.FindStringSubmatch(tag)
	if matches == nil {
		return SemVer{}, false
	}
	major, err := strconv.Atoi(matches[1])
	if err != nil {
		return SemVer{}, false
	}
	minor, err := strconv.Atoi(matches[2])
	if err != nil {
		return SemVer{}, false
	}
	patch, err := strconv.Atoi(matches[3])
	if err != nil {
		return SemVer{}, false
	}
	return SemVer{Major: major, Minor: minor, Patch: patch}, true
}

func (v SemVer) Compare(other SemVer) int {
	if v.Major != other.Major {
		return compareInt(v.Major, other.Major)
	}
	if v.Minor != other.Minor {
		return compareInt(v.Minor, other.Minor)
	}
	return compareInt(v.Patch, other.Patch)
}

func (v SemVer) String() string {
	return fmt.Sprintf("v%d.%d.%d", v.Major, v.Minor, v.Patch)
}

func compareInt(left, right int) int {
	switch {
	case left < right:
		return -1
	case left > right:
		return 1
	default:
		return 0
	}
}

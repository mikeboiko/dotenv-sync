package bitwarden

import (
	"context"
	"errors"
	"fmt"
	"os/exec"
	"strings"
)

var ErrBinaryMissing = errors.New("rbw binary missing")

type RBWClient struct {
	Bin string
}

func NewRBWClient() *RBWClient {
	return &RBWClient{Bin: "rbw"}
}

func (c *RBWClient) Run(ctx context.Context, args ...string) (string, error) {
	if _, err := exec.LookPath(c.Bin); err != nil {
		return "", ErrBinaryMissing
	}
	cmd := exec.CommandContext(ctx, c.Bin, args...)
	out, err := cmd.CombinedOutput()
	text := strings.TrimSpace(string(out))
	if err != nil {
		if text == "" {
			return text, fmt.Errorf("rbw %s failed: %w", strings.Join(args, " "), err)
		}
		return text, fmt.Errorf("%s", text)
	}
	return text, nil
}

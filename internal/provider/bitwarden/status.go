package bitwarden

import (
	"context"
	"strings"

	"dotenv-sync/internal/provider"
)

func checkReadinessWithClient(ctx context.Context, client *RBWClient) (provider.Status, error) {
	status := provider.Status{Provider: "bitwarden", CLIInstalled: true}
	out, err := client.Run(ctx, "status")
	if err != nil {
		if err == ErrBinaryMissing {
			status.CLIInstalled = false
			status.Code = "E001"
			status.Problem = "rbw CLI not installed"
			status.Impact = "commands cannot reach Bitwarden"
			status.Action = "install rbw or add it to PATH"
			return status, nil
		}
	}
	lower := strings.ToLower(strings.TrimSpace(out))
	switch {
	case strings.Contains(lower, "unlocked"):
		status.Authenticated = true
		status.Unlocked = true
		status.Message = "rbw is ready"
	case strings.Contains(lower, "locked"):
		status.Authenticated = true
		status.Code = "E003"
		status.Problem = "Bitwarden database is locked"
		status.Impact = "sync cannot resolve provider-managed schema keys"
		status.Action = "run 'rbw unlock' and retry"
	case strings.Contains(lower, "unauth"), strings.Contains(lower, "logged out"), strings.Contains(lower, "login"):
		status.Code = "E002"
		status.Problem = "not logged in to Bitwarden through rbw"
		status.Impact = "provider-backed commands cannot continue"
		status.Action = "run 'rbw login' and retry"
	default:
		if err != nil {
			status.Code = "E003"
			status.Problem = "Bitwarden readiness could not be verified"
			status.Impact = "provider-backed commands cannot continue"
			status.Action = "check rbw status and retry"
		}
	}
	return status, nil
}

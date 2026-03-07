package bitwarden

import (
	"context"
	"strings"

	"dotenv-sync/internal/provider"
)

func checkReadinessWithClient(ctx context.Context, client *RBWClient) (provider.Status, error) {
	status := provider.Status{Provider: "bitwarden", CLIInstalled: true}
	out, err := client.Run(ctx, "unlocked")
	if err != nil {
		if err == ErrBinaryMissing {
			status.CLIInstalled = false
			status.Code = "E001"
			status.Problem = "rbw CLI not installed"
			status.Impact = "commands cannot reach Bitwarden"
			status.Action = "install rbw or add it to PATH"
			return status, nil
		}
		if isUnsupportedSubcommand(out) {
			return checkLegacyStatus(ctx, client)
		}
		out, err = client.Run(ctx, "list")
		if err == nil {
			return readyStatus(status), nil
		}
		return classifyReadinessFailure(status, out), nil
	}
	return readyStatus(status), nil
}

func checkLegacyStatus(ctx context.Context, client *RBWClient) (provider.Status, error) {
	status := provider.Status{Provider: "bitwarden", CLIInstalled: true}
	out, err := client.Run(ctx, "status")
	if err != nil && err == ErrBinaryMissing {
		status.CLIInstalled = false
		status.Code = "E001"
		status.Problem = "rbw CLI not installed"
		status.Impact = "commands cannot reach Bitwarden"
		status.Action = "install rbw or add it to PATH"
		return status, nil
	}
	if err == nil {
		return classifyStatusText(status, out), nil
	}
	if isUnsupportedSubcommand(out) {
		status.Code = "E003"
		status.Problem = "Bitwarden readiness could not be verified"
		status.Impact = "provider-backed commands cannot continue"
		status.Action = "run 'rbw unlocked' or 'rbw list' manually and retry"
		return status, nil
	}
	return classifyReadinessFailure(status, out), nil
}

func readyStatus(status provider.Status) provider.Status {
	status.Authenticated = true
	status.Unlocked = true
	status.Message = "rbw is ready"
	return status
}

func classifyStatusText(status provider.Status, raw string) provider.Status {
	lower := strings.ToLower(strings.TrimSpace(raw))
	switch {
	case strings.Contains(lower, "unlocked"):
		return readyStatus(status)
	case strings.Contains(lower, "locked"):
		status.Authenticated = true
		status.Code = "E003"
		status.Problem = "Bitwarden database is locked"
		status.Impact = "sync cannot resolve provider-managed schema keys"
		status.Action = "run 'rbw unlock' and retry"
	case strings.Contains(lower, "unauth"), strings.Contains(lower, "logged out"), strings.Contains(lower, "not logged in"), strings.Contains(lower, "login"), strings.Contains(lower, "register"):
		status.Code = "E002"
		status.Problem = "not logged in to Bitwarden through rbw"
		status.Impact = "provider-backed commands cannot continue"
		status.Action = "run 'rbw login' and retry"
	default:
		status.Code = "E003"
		status.Problem = "Bitwarden readiness could not be verified"
		status.Impact = "provider-backed commands cannot continue"
		status.Action = "run 'rbw unlock' if locked, or 'rbw login' if logged out"
	}
	return status
}

func classifyReadinessFailure(status provider.Status, raw string) provider.Status {
	return classifyStatusText(status, raw)
}

func isUnsupportedSubcommand(raw string) bool {
	lower := strings.ToLower(raw)
	return strings.Contains(lower, "unrecognized subcommand") || strings.Contains(lower, "unknown subcommand")
}

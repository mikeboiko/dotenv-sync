package cli

import (
	"errors"
	"fmt"
	"io"
	"os"

	"dotenv-sync/internal/config"
	"dotenv-sync/internal/provider"
	"dotenv-sync/internal/provider/bitwarden"
	"dotenv-sync/internal/report"
	"github.com/spf13/cobra"
)

type streams struct {
	stdout io.Writer
	stderr io.Writer
}

type rootOptions struct {
	configPath string
	schemaPath string
	envPath    string
}

func Execute(args []string, stdout, stderr io.Writer) int {
	cmd := NewRootCommand(streams{stdout: stdout, stderr: stderr})
	cmd.SetArgs(args)
	cmd.SilenceErrors = true
	cmd.SilenceUsage = true
	err := cmd.Execute()
	if err != nil {
		var appErr *report.AppError
		if errors.As(err, &appErr) {
			_, _ = fmt.Fprintln(stderr, report.ErrorBlock(appErr))
		}
		return report.ExitCode(err)
	}
	return report.ExitSuccess
}

func NewRootCommand(s streams) *cobra.Command {
	opts := &rootOptions{}
	cmd := &cobra.Command{
		Use:   "ds",
		Short: "Sync .env files from .env.example and rbw",
	}
	cmd.PersistentFlags().StringVar(&opts.configPath, "config", "", "config file path")
	cmd.PersistentFlags().StringVar(&opts.schemaPath, "schema", "", "schema file path")
	cmd.PersistentFlags().StringVar(&opts.envPath, "env", "", "env file path")
	cmd.AddCommand(newSyncCommand(s, opts))
	cmd.AddCommand(newDiffCommand(s, opts))
	cmd.AddCommand(newValidateCommand(s, opts))
	cmd.AddCommand(newDoctorCommand(s, opts))
	cmd.AddCommand(newInitCommand(s, opts))
	cmd.AddCommand(newMissingCommand(s, opts))
	cmd.AddCommand(newReverseCommand(s, opts))
	return cmd
}

func loadConfig(opts *rootOptions) (config.Config, error) {
	cwd, err := os.Getwd()
	if err != nil {
		return config.Config{}, err
	}
	cfg, err := config.Load(cwd, config.LoadOptions{ConfigPath: opts.configPath, SchemaPath: opts.schemaPath, EnvPath: opts.envPath})
	if err != nil {
		return config.Config{}, report.NewAppError("E007", report.ExitOperational, "config file invalid", "commands cannot resolve file locations or mappings", "correct .envsync.yaml and retry", err)
	}
	return cfg, nil
}

func providerFor(cfg config.Config) provider.Provider {
	return bitwarden.NewAdapter(cfg)
}

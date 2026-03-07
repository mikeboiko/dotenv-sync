package provider

import "context"

type Status struct {
	Provider      string
	CLIInstalled  bool
	Authenticated bool
	Unlocked      bool
	Code          string
	Message       string
	Problem       string
	Impact        string
	Action        string
}

type Resolution struct {
	Key       string
	Ref       string
	Value     string
	Source    string
	IssueCode string
}

type Provider interface {
	Name() string
	CheckReadiness(ctx context.Context) (Status, error)
	Resolve(ctx context.Context, key, ref string) (Resolution, error)
	ResolveMany(ctx context.Context, refs map[string]string) (map[string]Resolution, error)
}

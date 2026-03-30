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

type EnvPayload struct {
	ItemName    string
	StorageMode string
	Exists      bool
	Format      string
	Notes       string
	Password    string
	Env         map[string]string
}

type WriteResult struct {
	ItemName string
	Created  bool
	Updated  bool
}

type Provider interface {
	Name() string
	CheckReadiness(ctx context.Context) (Status, error)
	Resolve(ctx context.Context, key, ref string) (Resolution, error)
	ResolveMany(ctx context.Context, refs map[string]string) (map[string]Resolution, error)
}

type PushProvider interface {
	Provider
	LoadEnvPayload(ctx context.Context) (EnvPayload, error)
	StoreEnvPayload(ctx context.Context, payload EnvPayload) (WriteResult, error)
}

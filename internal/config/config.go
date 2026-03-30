package config

import (
	"fmt"
	"os"
	"path/filepath"
	"strings"

	"gopkg.in/yaml.v3"
)

type Config struct {
	Provider   string            `yaml:"provider"`
	SchemaFile string            `yaml:"schema_file"`
	EnvFile    string            `yaml:"env_file"`
	ItemName   string            `yaml:"item_name"`
	Vault      string            `yaml:"vault"`
	Mapping    map[string]string `yaml:"mapping"`
	ConfigFile string            `yaml:"-"`
	BaseDir    string            `yaml:"-"`
}

type LoadOptions struct {
	ConfigPath string
	SchemaPath string
	EnvPath    string
}

func Default(baseDir string) Config {
	return Config{
		Provider:   "bitwarden",
		SchemaFile: filepath.Join(baseDir, ".env.example"),
		EnvFile:    filepath.Join(baseDir, ".env"),
		ItemName:   defaultItemName(baseDir),
		ConfigFile: filepath.Join(baseDir, ".envsync.yaml"),
		Mapping:    map[string]string{},
		BaseDir:    baseDir,
	}
}

func Load(baseDir string, opts LoadOptions) (Config, error) {
	baseDir, err := filepath.Abs(baseDir)
	if err != nil {
		return Config{}, err
	}
	cfg := Default(baseDir)
	if opts.ConfigPath != "" {
		cfg.ConfigFile = resolvePath(baseDir, opts.ConfigPath)
	}
	if data, err := os.ReadFile(cfg.ConfigFile); err == nil {
		if err := yaml.Unmarshal(data, &cfg); err != nil {
			return Config{}, fmt.Errorf("parse config: %w", err)
		}
	} else if !os.IsNotExist(err) {
		return Config{}, err
	}
	cfg.BaseDir = baseDir
	cfg.ItemName = strings.TrimSpace(cfg.ItemName)
	if cfg.ItemName == "" {
		cfg.ItemName = defaultItemName(baseDir)
	}
	cfg.ConfigFile = resolvePath(baseDir, cfg.ConfigFile)
	cfg.SchemaFile = resolvePath(baseDir, cfg.SchemaFile)
	cfg.EnvFile = resolvePath(baseDir, cfg.EnvFile)
	if opts.SchemaPath != "" {
		cfg.SchemaFile = resolvePath(baseDir, opts.SchemaPath)
	}
	if opts.EnvPath != "" {
		cfg.EnvFile = resolvePath(baseDir, opts.EnvPath)
	}
	if cfg.Provider == "" {
		cfg.Provider = "bitwarden"
	}
	if cfg.Mapping == nil {
		cfg.Mapping = map[string]string{}
	}
	if filepath.Clean(cfg.SchemaFile) == filepath.Clean(cfg.EnvFile) {
		return Config{}, fmt.Errorf("schema_file and env_file must differ")
	}
	return cfg, nil
}

func resolvePath(baseDir, value string) string {
	if value == "" {
		return value
	}
	if filepath.IsAbs(value) {
		return filepath.Clean(value)
	}
	if strings.HasPrefix(value, baseDir+string(filepath.Separator)) {
		return filepath.Clean(value)
	}
	return filepath.Join(baseDir, value)
}

func defaultItemName(baseDir string) string {
	root := repoRoot(baseDir)
	name := filepath.Base(root)
	if name == "." || name == string(filepath.Separator) || name == "" {
		name = filepath.Base(baseDir)
	}
	if name == "." || name == string(filepath.Separator) || name == "" {
		return "dotenv-sync"
	}
	return name
}

func repoRoot(start string) string {
	dir := filepath.Clean(start)
	for {
		if _, err := os.Stat(filepath.Join(dir, ".git")); err == nil {
			return dir
		}
		parent := filepath.Dir(dir)
		if parent == dir {
			return filepath.Clean(start)
		}
		dir = parent
	}
}

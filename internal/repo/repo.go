package repo

import (
	"errors"
	"fmt"
	"os"
	"path/filepath"
	"regexp"

	"github.com/jmmarotta/openplan/internal/plan"
)

const defaultConfigFilename = "openplan.jsonc"

var (
	errNotInitialized     = errors.New("openplan is not initialized")
	errAlreadyInitialized = errors.New("openplan is already initialized")
	prefixPattern         = regexp.MustCompile(`^[A-Z][A-Z0-9]*$`)
)

var configFilenames = []string{"openplan.jsonc", "openplan.json"}

// Discover walks upward from start until it finds a directory containing
// `.plans/`, then loads that repository context
func Discover(start string) (Context, error) {
	absStart, err := filepath.Abs(start)
	if err != nil {
		return Context{}, err
	}

	current := absStart
	// Walk parents until a repository is found or the filesystem root is reached
	for {
		plansDir := filepath.Join(current, ".plans")
		info, err := os.Stat(plansDir)
		if err == nil && info.IsDir() {
			return Load(current)
		}

		parent := filepath.Dir(current)
		if parent == current {
			return Context{}, fmt.Errorf("%w from %s", errNotInitialized, absStart)
		}
		current = parent
	}
}

// Init creates a new `.plans/` directory and config file rooted at root
func Init(root string, cfg Config) (Context, error) {
	absRoot, err := filepath.Abs(root)
	if err != nil {
		return Context{}, err
	}

	cfg = DefaultConfig(cfg.Prefix)
	if !prefixPattern.MatchString(cfg.Prefix) {
		return Context{}, fmt.Errorf("invalid prefix %q", cfg.Prefix)
	}

	plansDir := filepath.Join(absRoot, ".plans")
	if _, err := os.Stat(plansDir); err == nil {
		return Context{}, fmt.Errorf("%w at %s", errAlreadyInitialized, absRoot)
	} else if !errors.Is(err, os.ErrNotExist) {
		return Context{}, err
	}

	if err := os.Mkdir(plansDir, 0o755); err != nil {
		return Context{}, err
	}

	ctx := Context{
		Root:       absRoot,
		PlansDir:   plansDir,
		ConfigPath: filepath.Join(plansDir, defaultConfigFilename),
		Config:     cfg,
	}

	data, err := encodeConfig(cfg)
	if err != nil {
		return Context{}, err
	}

	tmpPath := ctx.ConfigPath + ".tmp"
	if err := os.WriteFile(tmpPath, data, 0o644); err != nil {
		return Context{}, err
	}
	if err := os.Rename(tmpPath, ctx.ConfigPath); err != nil {
		return Context{}, err
	}

	return ctx, nil
}

// Load reads an existing OpenPlan repository rooted at root
func Load(root string) (Context, error) {
	absRoot, err := filepath.Abs(root)
	if err != nil {
		return Context{}, err
	}

	plansDir := filepath.Join(absRoot, ".plans")
	info, err := os.Stat(plansDir)
	if err != nil {
		if errors.Is(err, os.ErrNotExist) {
			return Context{}, fmt.Errorf("%w at %s", errNotInitialized, absRoot)
		}
		return Context{}, err
	}
	if !info.IsDir() {
		return Context{}, fmt.Errorf("%w at %s", errNotInitialized, absRoot)
	}

	configPath, err := discoverConfigPath(plansDir)
	if err != nil {
		return Context{}, ConfigError{Path: filepath.Join(plansDir, defaultConfigFilename), Err: err}
	}
	data, err := os.ReadFile(configPath)
	if err != nil {
		return Context{}, ConfigError{Path: configPath, Err: err}
	}

	cfg, err := decodeConfig(data)
	if err != nil {
		return Context{}, ConfigError{Path: configPath, Err: err}
	}
	if !prefixPattern.MatchString(cfg.Prefix) {
		return Context{}, ConfigError{Path: configPath, Err: fmt.Errorf("invalid prefix %q", cfg.Prefix)}
	}

	return Context{
		Root:       absRoot,
		PlansDir:   plansDir,
		ConfigPath: configPath,
		Config:     cfg,
	}, nil
}

// PlanPath returns the canonical path for the given plan ID inside the repo
func (ctx Context) PlanPath(id string) string {
	return filepath.Join(ctx.PlansDir, plan.FilenameForID(id))
}

func discoverConfigPath(plansDir string) (string, error) {
	for _, name := range configFilenames {
		path := filepath.Join(plansDir, name)
		info, err := os.Stat(path)
		if err == nil && !info.IsDir() {
			return path, nil
		}
		if err == nil {
			continue
		}
		if !errors.Is(err, os.ErrNotExist) {
			return "", err
		}
	}

	return "", os.ErrNotExist
}

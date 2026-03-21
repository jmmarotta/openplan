package repo

import (
	"errors"
	"os"
	"path/filepath"
	"testing"
)

func TestInitLoadAndDiscover(t *testing.T) {
	root := t.TempDir()
	ctx, err := Init(root, DefaultConfig("abc"))
	if err != nil {
		t.Fatalf("Init returned error: %v", err)
	}

	if ctx.Config.Prefix != "ABC" {
		t.Fatalf("Init prefix = %q, want %q", ctx.Config.Prefix, "ABC")
	}
	if got, want := filepath.Base(ctx.ConfigPath), defaultConfigFilename; got != want {
		t.Fatalf("Init config file = %q, want %q", got, want)
	}

	loaded, err := Load(root)
	if err != nil {
		t.Fatalf("Load returned error: %v", err)
	}
	if loaded.Config.Prefix != "ABC" {
		t.Fatalf("Load prefix = %q, want %q", loaded.Config.Prefix, "ABC")
	}

	nested := filepath.Join(root, "a", "b")
	if err := os.MkdirAll(nested, 0o755); err != nil {
		t.Fatalf("MkdirAll returned error: %v", err)
	}

	discovered, err := Discover(nested)
	if err != nil {
		t.Fatalf("Discover returned error: %v", err)
	}
	if discovered.Root != root {
		t.Fatalf("Discover root = %q, want %q", discovered.Root, root)
	}
}

func TestInitAlreadyInitialized(t *testing.T) {
	root := t.TempDir()
	if _, err := Init(root, DefaultConfig("OPN")); err != nil {
		t.Fatalf("first Init returned error: %v", err)
	}

	_, err := Init(root, DefaultConfig("OPN"))
	if !errors.Is(err, errAlreadyInitialized) {
		t.Fatalf("Init error = %v, want wrapped already-initialized error", err)
	}
}

func TestLoadInvalidConfig(t *testing.T) {
	root := t.TempDir()
	plansDir := filepath.Join(root, ".plans")
	if err := os.MkdirAll(plansDir, 0o755); err != nil {
		t.Fatalf("MkdirAll returned error: %v", err)
	}
	if err := os.WriteFile(filepath.Join(plansDir, defaultConfigFilename), []byte("{"), 0o644); err != nil {
		t.Fatalf("WriteFile returned error: %v", err)
	}

	_, err := Load(root)
	var cfgErr ConfigError
	if !errors.As(err, &cfgErr) {
		t.Fatalf("Load error = %v, want ConfigError", err)
	}
}

func TestLoadJSONConfig(t *testing.T) {
	root := t.TempDir()
	plansDir := filepath.Join(root, ".plans")
	if err := os.MkdirAll(plansDir, 0o755); err != nil {
		t.Fatalf("MkdirAll returned error: %v", err)
	}
	if err := os.WriteFile(filepath.Join(plansDir, "openplan.json"), []byte("{\n  \"prefix\": \"OPN\"\n}\n"), 0o644); err != nil {
		t.Fatalf("WriteFile returned error: %v", err)
	}

	ctx, err := Load(root)
	if err != nil {
		t.Fatalf("Load returned error: %v", err)
	}
	if got, want := filepath.Base(ctx.ConfigPath), "openplan.json"; got != want {
		t.Fatalf("Load config file = %q, want %q", got, want)
	}
}

func TestLoadJSONCConfig(t *testing.T) {
	root := t.TempDir()
	plansDir := filepath.Join(root, ".plans")
	if err := os.MkdirAll(plansDir, 0o755); err != nil {
		t.Fatalf("MkdirAll returned error: %v", err)
	}
	data := []byte("{\n  // repo-local prefix\n  \"prefix\": \"OPN\",\n}\n")
	if err := os.WriteFile(filepath.Join(plansDir, defaultConfigFilename), data, 0o644); err != nil {
		t.Fatalf("WriteFile returned error: %v", err)
	}

	ctx, err := Load(root)
	if err != nil {
		t.Fatalf("Load returned error: %v", err)
	}
	if got, want := filepath.Base(ctx.ConfigPath), defaultConfigFilename; got != want {
		t.Fatalf("Load config file = %q, want %q", got, want)
	}
}

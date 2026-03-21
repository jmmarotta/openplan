package repo

import (
	"encoding/json"
	"fmt"
	"strings"

	"github.com/tidwall/jsonc"
)

type Config struct {
	Prefix string `json:"prefix"`
}

type Context struct {
	Root       string `json:"root"`
	PlansDir   string `json:"plansDir"`
	ConfigPath string `json:"configPath"`
	Config     Config `json:"config"`
}

type ConfigError struct {
	Path string
	Err  error
}

func (e ConfigError) Error() string {
	return fmt.Sprintf("read config %s: %v", e.Path, e.Err)
}

func (e ConfigError) Unwrap() error {
	return e.Err
}

func DefaultConfig(prefix string) Config {
	prefix = strings.ToUpper(strings.TrimSpace(prefix))
	if prefix == "" {
		prefix = "OPN"
	}
	return Config{Prefix: prefix}
}

func encodeConfig(cfg Config) ([]byte, error) {
	data, err := json.MarshalIndent(cfg, "", "  ")
	if err != nil {
		return nil, err
	}
	return append(data, '\n'), nil
}

func decodeConfig(data []byte) (Config, error) {
	var cfg Config
	if err := json.Unmarshal(jsonc.ToJSON(data), &cfg); err != nil {
		return Config{}, err
	}
	return cfg, nil
}

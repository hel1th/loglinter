package config

import (
	"encoding/json"
	"os"
	"path/filepath"
)

type Config struct {
	Enabled bool `json:"enabled"`

	Rules RulesConfig `json:"rules"`

	CustomSensitivePatterns []string `json:"custom-sensitive-patterns"`
}

type RulesConfig struct {
	LowercaseStart   RuleConfig `json:"lowercase-start"`
	EnglishOnly      RuleConfig `json:"english-only"`
	NoSpecialSymbols RuleConfig `json:"no-special-symbols"`
	NoSensitiveData  RuleConfig `json:"no-sensitive-data"`
}

type RuleConfig struct {
	Enabled bool `json:"enabled"`
}

func DefaultConfig() *Config {
	return &Config{
		Enabled: true,
		Rules: RulesConfig{
			LowercaseStart:   RuleConfig{Enabled: true},
			EnglishOnly:      RuleConfig{Enabled: true},
			NoSpecialSymbols: RuleConfig{Enabled: true},
			NoSensitiveData:  RuleConfig{Enabled: true},
		},
		CustomSensitivePatterns: []string{},
	}
}

func LoadConfig(path string) (*Config, error) {
	if path == "" {
		path = ".loglinter.json"
	}

	cfg := DefaultConfig()

	if _, err := os.Stat(path); os.IsNotExist(err) {
		return cfg, err
	}

	data, err := os.ReadFile(path)
	if err != nil {
		return nil, err
	}

	if err := json.Unmarshal(data, cfg); err != nil {
		return nil, err
	}

	return cfg, nil
}

func SaveConfig(cfg *Config, path string) error {
	dir := filepath.Dir(path)
	if err := os.MkdirAll(dir, 0o755); err != nil {
		return err
	}

	data, err := json.MarshalIndent(cfg, "", "  ")
	if err != nil {
		return err
	}

	return os.WriteFile(path, data, 0o644)
}

func (c *Config) IsEnabled() bool {
	if c == nil {
		return false
	}

	return c.Enabled
}

func (c *Config) GetEnabledRules() []string {
	enabled := make([]string, 0)

	if c == nil {
		return enabled
	}

	if c.Rules.LowercaseStart.Enabled {
		enabled = append(enabled, "lowercase-start")
	}
	if c.Rules.EnglishOnly.Enabled {
		enabled = append(enabled, "english-only")
	}
	if c.Rules.NoSpecialSymbols.Enabled {
		enabled = append(enabled, "no-special-symbols")
	}
	if c.Rules.NoSensitiveData.Enabled {
		enabled = append(enabled, "no-sensitive-data")
	}

	return enabled
}

func (c *Config) GetDisabledRules() []string {
	disabled := make([]string, 0)

	if c == nil {
		return disabled
	}

	if !c.Rules.LowercaseStart.Enabled {
		disabled = append(disabled, "lowercase-start")
	}
	if !c.Rules.EnglishOnly.Enabled {
		disabled = append(disabled, "english-only")
	}
	if !c.Rules.NoSpecialSymbols.Enabled {
		disabled = append(disabled, "no-special-symbols")
	}
	if !c.Rules.NoSensitiveData.Enabled {
		disabled = append(disabled, "no-sensitive-data")
	}

	return disabled
}

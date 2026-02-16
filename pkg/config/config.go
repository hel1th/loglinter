package config

import (
	"encoding/json"
	"os"
	"path/filepath"
)

type Config struct {
	Enabled bool `json:"enabled"`

	Rules RulesConfig `json:"rules"`

	CustomSensitivePatterns []string `json:"custom_sensitive_patterns"`
}

type RulesConfig struct {
	LowercaseStart   RuleConfig `json:"lowercase_start"`
	EnglishOnly      RuleConfig `json:"english_only"`
	NoSpecialSymbols RuleConfig `json:"no_special_symbols"`
	NoSensitiveData  RuleConfig `json:"no_sensitive_data"`
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

	if _, err := os.Stat(path); os.IsNotExist(err) {
		return DefaultConfig(), nil
	}

	data, err := os.ReadFile(path)
	if err != nil {
		return nil, err
	}

	var cfg Config
	if err := json.Unmarshal(data, &cfg); err != nil {
		return nil, err
	}

	return &cfg, nil
}

func SaveConfig(cfg *Config, path string) error {
	dir := filepath.Dir(path)
	if err := os.MkdirAll(dir, 0755); err != nil {
		return err
	}

	data, err := json.MarshalIndent(cfg, "", "  ")
	if err != nil {
		return err
	}

	return os.WriteFile(path, data, 0644)
}

func (c *Config) GetEnabledRules() []string {
	var enabled []string

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
	var disabled []string

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

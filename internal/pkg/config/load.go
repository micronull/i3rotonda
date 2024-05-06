package config

import (
	"fmt"
	"log/slog"
	"os"
	"path/filepath"

	"gopkg.in/yaml.v3"
)

const (
	defaultConfigFilename = "config.yml"
	defaultDir            = "i3rotonda"
)

func Load() (Config, error) {
	usrCfg, err := userConfigFilePath()
	if err != nil {
		return Config{}, fmt.Errorf("could not get user config file path: %w", err)
	}

	cfg := Config{}

	if _, err := os.Stat(usrCfg); err != nil {
		slog.Warn("config file not found, using defaults", "path", usrCfg)
	} else {
		b, err := os.ReadFile(usrCfg)
		if err != nil {
			return Config{}, fmt.Errorf("could not read config file: %s: %w", usrCfg, err)
		}

		if err := yaml.Unmarshal(b, &cfg); err != nil {
			return Config{}, fmt.Errorf("could not parse config file: %s: %w", usrCfg, err)
		}

		slog.Info("loaded config", "path", usrCfg)
	}

	return cfg, nil
}

func userConfigFilePath() (string, error) {
	dir, err := os.UserConfigDir()
	if err != nil {
		return "", fmt.Errorf("could not determine user config dir: %w", err)
	}

	dir = filepath.Join(dir, defaultDir, defaultConfigFilename)

	return dir, nil
}

package config

import (
	"encoding/json"
	"fmt"
	"os"
)

func ParseLLPkgConfig(configPath string) (LLPkgConfig, error) {
	var config LLPkgConfig
	file, err := os.Open(configPath)
	if err != nil {
		return config, fmt.Errorf("failed to open config file: %w", err)
	}
	defer file.Close()

	decoder := json.NewDecoder(file)
	err = decoder.Decode(&config)
	if err != nil {
		return config, fmt.Errorf("failed to decode config file: %w", err)
	}

	// set default values
	config = fillDefaults(config)

	return config, nil
}

func fillDefaults(config LLPkgConfig) LLPkgConfig {
	if config.Upstream.Installer.Name == "" {
		config.Upstream.Installer.Name = "conan"
	}
	return config
}

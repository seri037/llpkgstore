package config

import (
	"fmt"
	"slices"
)

// Check the validity of an LLPkgConfig
func ValidateLLPkgConfig(config LLPkgConfig) error {
	return validateUpstreamConfig(config.Upstream)
}

func validateUpstreamConfig(config UpstreamConfig) error {
	// 1. check if upstream installer is valid
	if config.Installer.Name == "" {
		return fmt.Errorf("invalid configuration: upstream.installer.name is required")
	}
	if !slices.Contains(ValidInstallers, config.Installer.Name) {
		return fmt.Errorf("invalid configuration: upstream.installer.name %s is not supported", config.Installer.Name)
	}

	// 2. check if package is valid
	if config.Package.Name == "" {
		return fmt.Errorf("invalid configuration: upstream.package.name is required")
	}
	if config.Package.Version == "" {
		return fmt.Errorf("invalid configuration: upstream.package.version is required")
	}

	return nil
}

package config

import (
	"errors"

	"github.com/goplus/llpkgstore/upstream"
	"github.com/goplus/llpkgstore/upstream/installer/conan"
)

var ValidInstallers = []string{"conan"}

// Represents an specific llpkg.cfg, can be parse from a file
type LLPkgConfig struct {
	Upstream UpstreamConfig `json:"upstream"`
}

// Represents an "upstream" field in llpkg.cfg
type UpstreamConfig struct {
	Installer InstallerConfig `json:"installer"`
	Package   PackageConfig   `json:"package"`
}

// Represents an "installer" field in llpkg.cfg
type InstallerConfig struct {
	Name   string            `json:"name"`
	Config map[string]string `json:"config,omitempty"`
}

// Represents a "package" field in llpkg.cfg
type PackageConfig struct {
	Name    string `json:"name"`
	Version string `json:"version"`
}

// Creates an Upstream from a UpstreamConfig
func NewUpstreamFromConfig(upstreamConfig UpstreamConfig) (*upstream.Upstream, error) {
	switch upstreamConfig.Installer.Name {
	case "conan":
		return &upstream.Upstream{
			Installer: conan.NewConanInstaller(upstreamConfig.Installer.Config),
			Pkg: upstream.Package{
				Name:    upstreamConfig.Package.Name,
				Version: upstreamConfig.Package.Version,
			},
		}, nil
	default:
		return nil, errors.New("unknown upstream installer: " + upstreamConfig.Installer.Name)
	}
}

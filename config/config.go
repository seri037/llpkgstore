package config

import (
	"errors"

	"github.com/goplus/llpkgstore/upstream"
	"github.com/goplus/llpkgstore/upstream/installer/conan"
)

var ValidInstallers = []string{"conan"}

type LLpkgConfig struct {
	Upstream UpstreamConfig `json:"upstream"`
}

type UpstreamConfig struct {
	Installer InstallerConfig `json:"installer"`
	Package   PackageConfig   `json:"package"`
}

type InstallerConfig struct {
	Name   string            `json:"name"`
	Config map[string]string `json:"config,omitempty"`
}

type PackageConfig struct {
	Name    string `json:"name"`
	Version string `json:"version"`
}

func NewUpstreamFromConfig(upstreamConfig UpstreamConfig) (*upstream.Upstream, error) {
	switch upstreamConfig.Installer.Name {
	case "conan":
		return upstream.NewUpstream(conan.NewConanInstaller(upstreamConfig.Installer.Config), upstream.Package{
			Name:    upstreamConfig.Package.Name,
			Version: upstreamConfig.Package.Version,
		}), nil
	default:
		return nil, errors.New("unknown upstream installer: " + upstreamConfig.Installer.Name)
	}
}

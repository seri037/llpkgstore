package upstream

// An Installer can install and search binaries from a specific remote.
type Installer interface {
	Name() string
	Config() map[string]string
	// Installs a binary according to pkg, and stores the result (usually .pc files) in outputDir.
	Install(pkg Package, dir string) error
	Search(pkg Package) (string, error)
}

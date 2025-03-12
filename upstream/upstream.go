package upstream

// An Upstream represents a binary and its installation method.
type Upstream struct {
	Installer Installer
	Pkg       Package
}

// A Package defines a binary by its name and version.
//
// It's meaningful only when using an installer.
type Package struct {
	Name    string
	Version string
}

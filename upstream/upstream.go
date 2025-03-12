package upstream

type Upstream struct {
	Installer Installer
	Pkg       Package
}

type Package struct {
	Name    string
	Version string
}

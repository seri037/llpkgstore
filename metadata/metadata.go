package metadata

type CVersion = string
type GoVersion = string

// MetadataMap represents llpkgstore.json
type MetadataMap map[string]*Metadata

type Metadata struct {
	Versions map[CVersion][]GoVersion `json:"versions"`
}

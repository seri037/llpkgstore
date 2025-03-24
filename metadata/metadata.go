package metadata

// MetadataMap represents llpkgstore.json
type MetadataMap map[string]*Metadata

type Metadata struct {
	VersionMappings map[string][]string `json:"versions"` // c_ver -> module_vers
}

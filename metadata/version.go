package metadata

import (
	"fmt"

	"golang.org/x/mod/semver"
)

// LatestMappedModuleVersion returns the latest mapped module version according to the C version
func (m *metadataMgr) LatestMappedModuleVersion(name, cversion string) (string, error) {
	cToGoVersions, ok := m.cToGoVersionsMaps[name]
	if !ok {
		return "", fmt.Errorf("no version mappings for %s", name)
	}

	goVersions, ok := cToGoVersions[cversion]
	if !ok {
		return "", fmt.Errorf("no version mappings for %s %s", name, cversion)
	}

	if len(goVersions) > 0 {
		lastestGoVersion := goVersions[0]
		for _, goVersion := range goVersions {
			if semver.Compare(goVersion, lastestGoVersion) > 0 {
				lastestGoVersion = goVersion
			}
		}
		return lastestGoVersion, nil
	}

	return "", fmt.Errorf("no version mappings for %s %s", name, cversion)
}

// MappedModuleVersions returns all mapped module versions according to the C version
func (m *metadataMgr) MappedModuleVersions(name, cversion string) ([]string, error) {
	cToGoVersions, ok := m.cToGoVersionsMaps[name]
	if !ok {
		return nil, fmt.Errorf("no version mappings for %s", name)
	}

	versions, ok := cToGoVersions[cversion]
	if !ok {
		return nil, fmt.Errorf("no version mappings for %s %s", name, cversion)
	}

	return versions, nil
}

// AllCToGoVersions returns all version mappings in the format of Name -> CVersion -> GoVersions
func (m *metadataMgr) AllCToGoVersions() (map[string]map[string][]string, error) {
	err := m.update()
	if err != nil {
		return nil, err
	}
	return m.cToGoVersionsMaps, nil
}

// AllGoToCVersion returns all version mappings in the format of Name -> GoVersion -> CVersion
func (m *metadataMgr) AllGoToCVersion() (map[string]map[string]string, error) {
	err := m.update()
	if err != nil {
		return nil, err
	}
	return m.goToCVersionMaps, nil
}

// CToGoVersionsByName returns a module's C-to-Go mappings
func (m *metadataMgr) CToGoVersionsByName(name string) (map[string][]string, error) {
	// First try to find the version mappings in the cache
	cToGoVersions, ok := m.cToGoVersionsMaps[name]
	if !ok {
		// If the version mappings are not in the cache, update the cache
		err := m.update()
		if err != nil {
			return nil, err
		}

		// Find the version mappings again
		cToGoVersions, ok = m.cToGoVersionsMaps[name]
		if !ok {
			return nil, ErrMetadataNotInCache
		}
	}

	return cToGoVersions, nil
}

// GoToCVersionByName returns a module's Go-to-C mappings
func (m *metadataMgr) GoToCVersionByName(name string) (map[string]string, error) {
	// First try to find the version mappings in the cache
	goToCVersion, ok := m.goToCVersionMaps[name]
	if !ok {
		// If the version mappings are not in the cache, update the cache
		err := m.update()
		if err != nil {
			return nil, err
		}

		// Find the version mappings again
		goToCVersion, ok = m.goToCVersionMaps[name]
		if !ok {
			return nil, ErrMetadataNotInCache
		}
	}

	return goToCVersion, nil
}

// AllVersionMappings returns all up-to-date version mappings in a primitive format
func (m *metadataMgr) AllVersionMappings() (map[string][]VersionMapping, error) {
	err := m.update()
	if err != nil {
		return nil, err
	}

	return m.allCachedVersionMappings(), nil
}

// VersionMappingsByName returns a module's version mappings in a primitive format
func (m *metadataMgr) VersionMappingsByName(name string) ([]VersionMapping, error) {
	// First try to find the version mappings in the cache
	versionMappings, ok := m.allCachedVersionMappings()[name]
	if !ok {
		// If the version mappings are not in the cache, update the cache
		err := m.update()
		if err != nil {
			return nil, err
		}

		// Find the version mappings again
		versionMappings, ok = m.allCachedVersionMappings()[name]
		if !ok {
			return nil, ErrMetadataNotInCache
		}
	}

	return versionMappings, nil
}

func (m *metadataMgr) allCachedVersionMappings() map[string][]VersionMapping {
	allCachedMetadata := m.allCachedMetadata()

	allVersionMappings := map[string][]VersionMapping{}
	for name, info := range allCachedMetadata {
		allVersionMappings[name] = info.VersionMappings
	}

	return allVersionMappings
}

func (m *metadataMgr) buildVersionsHash() error {
	m.cToGoVersionsMaps = make(map[string]map[string][]string)
	m.goToCVersionMaps = make(map[string]map[string]string)

	allCachedVersionMappings := m.allCachedVersionMappings()

	for name, versionMappings := range allCachedVersionMappings {
		if m.cToGoVersionsMaps[name] == nil {
			m.cToGoVersionsMaps[name] = make(map[string][]string)
		}
		if m.goToCVersionMaps[name] == nil {
			m.goToCVersionMaps[name] = make(map[string]string)
		}
		for _, versionMapping := range versionMappings {
			m.cToGoVersionsMaps[name][versionMapping.CVersion] = versionMapping.GoVersions
			for _, goVersion := range versionMapping.GoVersions {
				m.goToCVersionMaps[name][goVersion] = versionMapping.CVersion
			}
		}
	}

	return nil
}

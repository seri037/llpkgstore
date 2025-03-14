package metadata

import (
	"fmt"

	"golang.org/x/mod/semver"
)

// LatestMappedModuleVersion returns the latest mapped module version according to the C version
func (m *metadataMgr) LatestMappedModuleVersion(name, cversion string) (string, error) {
	// 使用扁平化结构进行查询 - 构建复合键
	cKey := name + "/" + cversion

	// 直接从扁平哈希表中获取Go版本列表
	goVersions, ok := m.flatCToGo[cKey]
	if !ok {
		// 如果没找到，尝试更新缓存
		err := m.update()
		if err != nil {
			return "", err
		}

		// 再次查询
		goVersions, ok = m.flatCToGo[cKey]
		if !ok {
			return "", fmt.Errorf("no version mappings for %s %s", name, cversion)
		}
	}

	if len(goVersions) > 0 {
		// 找出语义版本最高的Go版本
		latestGoVersion := goVersions[0]
		for _, goVersion := range goVersions {
			if semver.Compare(goVersion, latestGoVersion) > 0 {
				latestGoVersion = goVersion
			}
		}
		return latestGoVersion, nil
	}

	return "", fmt.Errorf("no version mappings for %s %s", name, cversion)
}

// MappedModuleVersions returns all mapped module versions according to the C version
func (m *metadataMgr) MappedModuleVersions(name, cversion string) ([]string, error) {
	// 使用扁平化结构进行查询 - 构建复合键
	cKey := name + "/" + cversion

	// 直接从扁平哈希表中获取Go版本列表
	versions, ok := m.flatCToGo[cKey]
	if !ok {
		// 如果没找到，尝试更新缓存
		err := m.update()
		if err != nil {
			return nil, err
		}

		// 再次查询
		versions, ok = m.flatCToGo[cKey]
		if !ok {
			return nil, fmt.Errorf("no version mappings for %s %s", name, cversion)
		}
	}

	// 返回副本避免修改内部数据
	result := make([]string, len(versions))
	copy(result, versions)

	return result, nil
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

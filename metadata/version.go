package metadata

import (
	"fmt"

	"golang.org/x/mod/semver"
)

// LatestCVer 通过 clibName 获取最新的 C 版本
func (m *metadataMgr) LatestCVer(name string) (string, error) {
	// 获取该模块的所有版本映射
	versionMappings, err := m.VersionMappingsByName(name)
	if err != nil {
		return "", err
	}

	if len(versionMappings) == 0 {
		return "", fmt.Errorf("no version mappings for %s", name)
	}

	// 查找最新的 C 版本（假设所有 C 版本都遵循语义版本格式）
	latestCVersion := versionMappings[0].CVersion
	for _, mapping := range versionMappings {
		if semver.Compare(ensureSemverPrefix(mapping.CVersion), ensureSemverPrefix(latestCVersion)) > 0 {
			latestCVersion = mapping.CVersion
		}
	}

	return latestCVersion, nil
}

// LatestGoVer 通过 clibName 获取最新的 Go 版本
func (m *metadataMgr) LatestGoVer(name string) (string, error) {
	// 获取该模块的所有 Go 版本
	allGoVersions, err := m.AllGoVersFromName(name)
	if err != nil {
		return "", err
	}

	if len(allGoVersions) == 0 {
		return "", fmt.Errorf("no Go versions found for %s", name)
	}

	// 查找最新的 Go 版本
	latestGoVersion := allGoVersions[0]
	for _, goVersion := range allGoVersions {
		if semver.Compare(goVersion, latestGoVersion) > 0 {
			latestGoVersion = goVersion
		}
	}

	return latestGoVersion, nil
}

// LatestGoVerFromCVer 通过 clibName 和 C 版本获取最新的 Go 版本
func (m *metadataMgr) LatestGoVerFromCVer(name, cVer string) (string, error) {
	// 使用扁平化结构进行查询 - 构建复合键
	cKey := name + "/" + cVer

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
			return "", fmt.Errorf("no version mappings for %s %s", name, cVer)
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

	return "", fmt.Errorf("no version mappings for %s %s", name, cVer)
}

// GoVersFromCVer 通过 clibName 和 C 版本获取 Go 版本列表
func (m *metadataMgr) GoVersFromCVer(name, cVer string) ([]string, error) {
	// 使用扁平化结构进行查询 - 构建复合键
	cKey := name + "/" + cVer

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
			return nil, fmt.Errorf("no version mappings for %s %s", name, cVer)
		}
	}

	// 返回副本避免修改内部数据
	result := make([]string, len(versions))
	copy(result, versions)

	return result, nil
}

// CVerFromGoVer 通过 clibName 和 Go 版本获取 C 版本
func (m *metadataMgr) CVerFromGoVer(name, goVer string) (string, error) {
	// 使用扁平化结构进行查询 - 构建复合键
	goKey := name + "/" + goVer

	// 直接从扁平哈希表中获取 C 版本
	cVersion, ok := m.flatGoToC[goKey]
	if !ok {
		// 如果没找到，尝试更新缓存
		err := m.update()
		if err != nil {
			return "", err
		}

		// 再次查询
		cVersion, ok = m.flatGoToC[goKey]
		if !ok {
			return "", fmt.Errorf("no C version found for %s %s", name, goVer)
		}
	}

	return cVersion, nil
}

// AllGoVersFromName 通过 clibName 获取 Go 所有的版本列表
func (m *metadataMgr) AllGoVersFromName(name string) ([]string, error) {
	// 获取该模块的所有版本映射
	versionMappings, err := m.VersionMappingsByName(name)
	if err != nil {
		return nil, err
	}

	// 使用 map 来去重
	goVersionsMap := make(map[string]struct{})
	for _, mapping := range versionMappings {
		for _, goVer := range mapping.GoVersions {
			goVersionsMap[goVer] = struct{}{}
		}
	}

	// 从 map 转换回切片
	goVersions := make([]string, 0, len(goVersionsMap))
	for goVer := range goVersionsMap {
		goVersions = append(goVersions, goVer)
	}

	return goVersions, nil
}

// AllCVersFromName 通过 clibName 获取 C 所有的版本列表
func (m *metadataMgr) AllCVersFromName(name string) ([]string, error) {
	// 获取该模块的所有版本映射
	versionMappings, err := m.VersionMappingsByName(name)
	if err != nil {
		return nil, err
	}

	// 提取所有唯一的 C 版本
	cVersions := make([]string, 0, len(versionMappings))
	for _, mapping := range versionMappings {
		cVersions = append(cVersions, mapping.CVersion)
	}

	return cVersions, nil
}

// AllVersionMappingsFromName 返回原始的版本映射
func (m *metadataMgr) AllVersionMappingsFromName(name string) ([]VersionMapping, error) {
	// 重用现有的 VersionMappingsByName 方法
	return m.VersionMappingsByName(name)
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

// 辅助函数，确保语义版本号有前缀 "v"
func ensureSemverPrefix(version string) string {
	if len(version) > 0 && version[0] != 'v' {
		return "v" + version
	}
	return version
}

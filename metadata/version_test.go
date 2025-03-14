package metadata

import (
	"encoding/json"
	"net/http"
	"net/http/httptest"
	"testing"

	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/require"
)

var testVersionData = MetadataMap{
	"test-module": Metadata{
		VersionMappings: []VersionMapping{
			{CVersion: "1.7.18", GoVersions: []string{"v1.2.0", "v1.2.1"}},
			{CVersion: "1.7.19", GoVersions: []string{"v1.3.0"}},
		},
	},
	"empty-module": Metadata{VersionMappings: []VersionMapping{}},
}

// TestBuildVersionsHash verifies that the buildVersionsHash method correctly constructs the cToGoVersionsMaps and goToCVersionMaps based on the testVersionData.
func TestBuildVersionsHash(t *testing.T) {
	server := newTestServer(testVersionData)
	defer server.Close()

	originalURL := remoteMetadataURL
	defer func() { remoteMetadataURL = originalURL }()
	remoteMetadataURL = server.URL

	tmpDir := t.TempDir()
	mgr, _ := NewMetadataMgr(tmpDir)
	require.NoError(t, mgr.buildVersionsHash())

	assert.Equal(t, []string{"v1.2.0", "v1.2.1"}, mgr.cToGoVersionsMaps["test-module"]["1.7.18"])
	assert.Equal(t, []string{"v1.3.0"}, mgr.cToGoVersionsMaps["test-module"]["1.7.19"])
	assert.Equal(t, "1.7.18", mgr.goToCVersionMaps["test-module"]["v1.2.1"])
	assert.Equal(t, "1.7.19", mgr.goToCVersionMaps["test-module"]["v1.3.0"])
}

// TestLatestMappedModuleVersion checks the LatestMappedModuleVersion method for valid/invalid inputs and error handling.
func TestLatestMappedModuleVersion(t *testing.T) {
	server := newTestServer(testVersionData)
	defer server.Close()

	originalURL := remoteMetadataURL
	defer func() { remoteMetadataURL = originalURL }()
	remoteMetadataURL = server.URL

	tmpDir := t.TempDir()
	mgr, _ := NewMetadataMgr(tmpDir)

	ver, err := mgr.LatestMappedModuleVersion("test-module", "1.7.18")
	require.NoError(t, err)
	assert.Equal(t, "v1.2.1", ver)

	_, err = mgr.LatestMappedModuleVersion("test-module", "v3.0")
	assert.Error(t, err)

	_, err = mgr.LatestMappedModuleVersion("empty-module", "v1.0")
	assert.Error(t, err)
}

// TestMappedModuleVersions validates the MappedModuleVersions method returns correct Go versions for valid C versions and errors for invalid cases.
func TestMappedModuleVersions(t *testing.T) {
	server := newTestServer(testVersionData)
	defer server.Close()

	originalURL := remoteMetadataURL
	defer func() { remoteMetadataURL = originalURL }()
	remoteMetadataURL = server.URL

	tmpDir := t.TempDir()
	mgr, _ := NewMetadataMgr(tmpDir)

	versions, err := mgr.MappedModuleVersions("test-module", "1.7.18")
	require.NoError(t, err)
	assert.ElementsMatch(t, []string{"v1.2.0", "v1.2.1"}, versions)

	_, err = mgr.MappedModuleVersions("test-module", "v3.0")
	assert.Error(t, err)
}

// TestAllCToGoVersions verifies the AllCToGoVersions method returns complete C-to-Go version mappings.
func TestAllCToGoVersions(t *testing.T) {
	server := newTestServer(testVersionData)
	defer server.Close()

	originalURL := remoteMetadataURL
	defer func() { remoteMetadataURL = originalURL }()
	remoteMetadataURL = server.URL

	tmpDir := t.TempDir()
	mgr, _ := NewMetadataMgr(tmpDir)

	ctogo, err := mgr.AllCToGoVersions()
	require.NoError(t, err)
	assert.Equal(t, []string{"v1.2.0", "v1.2.1"}, ctogo["test-module"]["1.7.18"])
}

// TestAllGoToCVersion checks the AllGoToCVersion method returns correct Go-to-C version mappings.
func TestAllGoToCVersion(t *testing.T) {
	server := newTestServer(testVersionData)
	defer server.Close()

	originalURL := remoteMetadataURL
	defer func() { remoteMetadataURL = originalURL }()
	remoteMetadataURL = server.URL

	tmpDir := t.TempDir()
	mgr, _ := NewMetadataMgr(tmpDir)

	gotoc, err := mgr.AllGoToCVersion()
	require.NoError(t, err)
	assert.Equal(t, "1.7.19", gotoc["test-module"]["v1.3.0"])
}

// TestCToGoVersionsByName validates the CToGoVersionsByName method retrieves C-to-Go mappings for valid modules and returns errors for invalid ones.
func TestCToGoVersionsByName(t *testing.T) {
	server := newTestServer(testVersionData)
	defer server.Close()

	originalURL := remoteMetadataURL
	defer func() { remoteMetadataURL = originalURL }()
	remoteMetadataURL = server.URL

	tmpDir := t.TempDir()
	mgr, _ := NewMetadataMgr(tmpDir)

	mapping, err := mgr.CToGoVersionsByName("test-module")
	require.NoError(t, err)
	assert.Equal(t, []string{"v1.2.0", "v1.2.1"}, mapping["1.7.18"])

	_, err = mgr.CToGoVersionsByName("invalid-module")
	assert.Error(t, err)
}

// TestGoToCVersionByName checks the GoToCVersionByName method retrieves Go-to-C mappings for valid modules and handles invalid cases.
func TestGoToCVersionByName(t *testing.T) {
	server := newTestServer(testVersionData)
	defer server.Close()

	originalURL := remoteMetadataURL
	defer func() { remoteMetadataURL = originalURL }()
	remoteMetadataURL = server.URL

	tmpDir := t.TempDir()
	mgr, _ := NewMetadataMgr(tmpDir)

	mapping, err := mgr.GoToCVersionByName("test-module")
	require.NoError(t, err)
	assert.Equal(t, "1.7.18", mapping["v1.2.0"])

	_, err = mgr.GoToCVersionByName("invalid-module")
	assert.Error(t, err)
}

// TestVersionMappingsByName verifies the VersionMappingsByName method returns correct version mappings including empty module cases.
func TestVersionMappingsByName(t *testing.T) {
	server := newTestServer(testVersionData)
	defer server.Close()

	originalURL := remoteMetadataURL
	defer func() { remoteMetadataURL = originalURL }()
	remoteMetadataURL = server.URL

	tmpDir := t.TempDir()
	mgr, _ := NewMetadataMgr(tmpDir)

	mappings, err := mgr.VersionMappingsByName("test-module")
	require.NoError(t, err)
	assert.Len(t, mappings, 2)

	_, err = mgr.VersionMappingsByName("empty-module")
	assert.NoError(t, err)
}

// newTestServer creates a mock HTTP server for testing that returns the provided MetadataMap when requested.
func newTestServer(data MetadataMap) *httptest.Server {
	return httptest.NewServer(http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		json.NewEncoder(w).Encode(data)
	}))
}

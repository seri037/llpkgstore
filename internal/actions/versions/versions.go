package versions

import (
	"encoding/json"
	"io"
	"log"
	"os"
	"slices"

	"github.com/goplus/llpkgstore/metadata"
	"golang.org/x/mod/semver"
)

// Versions is a mapping table implement for Github Action only.
// It's recommend to use another implement in llgo for common usage.
type Versions struct {
	metadata.MetadataMap

	fileName string
}

// appendVersion appends a version to an array, panic if the specified version has already existed.
func appendVersion(arr []string, elem string) []string {
	if slices.Contains(arr, elem) {
		log.Fatalf("version %s has already existed", elem)
	}
	return append(arr, elem)
}

// Read reads version mappings from a file and initializes the Versions struct
func Read(fileName string) *Versions {
	// read or create a file
	f, err := os.OpenFile(fileName, os.O_CREATE|os.O_RDONLY, 0644)
	if err != nil {
		panic(err)
	}
	defer f.Close()

	b, err := io.ReadAll(f)
	if err != nil {
		panic(err)
	}

	m := metadata.MetadataMap{}

	if len(b) > 0 {
		json.Unmarshal(b, &m)
	}

	return &Versions{
		MetadataMap: m,
		fileName:    f.Name(),
	}
}

func (v *Versions) cVersions(clib string) (ret []string) {
	versions := v.MetadataMap[clib]
	if versions == nil {
		return
	}
	for version := range versions.Versions {
		ret = append(ret, version)
	}
	return
}

func (v *Versions) CVersions(clib string) (ret []string) {
	versions := v.MetadataMap[clib]
	if versions == nil {
		return
	}
	for version := range versions.Versions {
		ret = append(ret, ToSemVer(version))
	}
	return
}

func (v *Versions) GoVersions(clib string) (ret []string) {
	versions := v.MetadataMap[clib]
	if versions == nil {
		return
	}
	for _, goversion := range versions.Versions {
		ret = append(ret, goversion...)
	}
	return
}

func (v *Versions) LatestGoVersionForCVersion(clib, cver string) string {
	version := v.MetadataMap[clib]
	if version == nil {
		return ""
	}
	goVersions := version.Versions[cver]
	if len(goVersions) == 0 {
		return ""
	}
	semver.Sort(goVersions)
	return goVersions[len(goVersions)-1]
}

func (v *Versions) SearchBySemVer(clib, semver string) string {
	for _, version := range v.cVersions(clib) {
		if ToSemVer(version) == semver {
			return version
		}
	}
	return ""
}

// LatestGoVersion returns the latest Go version associated with the given C library
func (v *Versions) LatestGoVersion(clib string) string {
	versions := v.GoVersions(clib)
	if len(versions) == 0 {
		return ""
	}
	semver.Sort(versions)
	return versions[len(versions)-1]
}

// Write records a new Go version mapping for a C library version and persists to file
func (v *Versions) Write(clib, clibVersion, mappedVersion string) {
	clibVersions := v.MetadataMap[clib]
	if clibVersions == nil {
		clibVersions = &metadata.Metadata{
			Versions: map[metadata.CVersion][]metadata.GoVersion{},
		}
		v.MetadataMap[clib] = clibVersions
	}
	versions := clibVersions.Versions[clibVersion]

	versions = appendVersion(versions, mappedVersion)

	clibVersions.Versions[clibVersion] = versions
	// sync to disk
	b, _ := json.Marshal(&v.MetadataMap)

	os.WriteFile(v.fileName, []byte(b), 0644)
}

func (v *Versions) String() string {
	b, _ := json.Marshal(&v.MetadataMap)
	return string(b)
}

// Structs and utilities for retrieving package metadata.
package gonepm

import (
	"encoding/json"
	"fmt"
	"log"
	"net/http"
	"sync"
	"time"
)

// ShortMetadata contains the abbreviated metadata retrieved from npm with
// "Accept: application/vnd.npm.install-v1+json" header.
type ShortMetadata struct {
	DistTags struct {
		Latest string `json:"latest"`
	} `json:"dist-tags"`
	Modified time.Time               `json:"modified"`
	Name     string                  `json:"name"`
	Versions map[string]ShortVersion `json:"versions"`
}

// ShortVersion contains the version information in the abbreviated metadata.
type ShortVersion struct {
	// HasShrinkwrap bool `json:"_hasShrinkwrap"`
	// Directories struct {
	// } `json:"directories"`
	Dist struct {
		Shasum  string `json:"shasum"`
		Tarball string `json:"tarball"`
	} `json:"dist"`
	Name    string `json:"name"`
	Version string `json:"version"`
}

// FullMetadata contains the full metadata retrieved from npm.
type FullMetadata struct {
	ID       string                 `json:"_id"`
	Rev      string                 `json:"_rev"`
	Name     string                 `json:"name"`
	Time     map[string]string      `json:"time"`
	Versions map[string]FullVersion `json:"versions"`
}

// FullVersion
type FullVersion struct {
	ID     string `json:"_id"`
	Shasum string `json:"_shasum"`
	Dist   struct {
		// Shasum  string `json:"shasum"`
		Tarball string `json:"tarball"`
	} `json:"dist"`
	Name    string `json:"name"`
	Version string `json:"version"`
}

// PackageMetadata returns the full metadata of a package.
func (r *Registry) PackageMetadata(packageName string) (FullMetadata, error) {

	var mt FullMetadata
	if packageName == "" {
		return mt, fmt.Errorf("gonepm.PackageMetadata: empty package name")
	}

	resp, err := http.Get(fmt.Sprintf("%s/%s", r.BaseURL, packageName))
	if err != nil {
		return mt, fmt.Errorf("gonepm.PackageMetadata: %s", err.Error())
	}
	defer resp.Body.Close()

	if err := json.NewDecoder(resp.Body).Decode(&mt); err != nil {
		return mt, fmt.Errorf("gonepm.PackageMetadata: %s", err.Error())
	}
	return mt, nil
}

// ShortPackageMetadata returns the abbreviated metadata of a package.
func (r *Registry) ShortPackageMetadata(packageName string) (ShortMetadata, error) {

	var sm ShortMetadata
	if packageName == "" {
		return sm, fmt.Errorf("gonepm.ShortPackageMetadata: empty package name")
	}

	// Create request.
	req, err := http.NewRequest("GET", fmt.Sprintf("%s/%s", r.BaseURL, packageName), nil)
	if err != nil {
		return sm, fmt.Errorf("gonepm.ShortPackageMetadata: %s", err.Error())
	}
	// Set "Accept: application/vnd.npm.install-v1+json" header.
	req.Header.Set("Accept", "application/vnd.npm.install-v1+json")

	client := &http.Client{}

	resp, err := client.Do(req)
	if err != nil {
		return sm, fmt.Errorf("gonepm.ShortPackageMetadata: %s", err.Error())
	}
	defer resp.Body.Close()

	if err := json.NewDecoder(resp.Body).Decode(&sm); err != nil {
		return sm, fmt.Errorf("gonepm.ShortPackageMetadata: %s", err.Error())
	}
	return sm, nil
}

// PackageSizes downloads full metadata for a specific package from a
// registry. Calls HEAD on the tarball addresses and returns a map of
// [version]tarballSize.
func (r *Registry) PackageSizes(packageName string) (map[string]int64, error) {

	sizes := make(map[string]int64)
	if packageName == "" {
		return sizes, fmt.Errorf("gonepm.GetSizes: empty package name")
	}

	fullMT, err := r.PackageMetadata(packageName)
	if err != nil {
		return sizes, fmt.Errorf("gonepm.GetSizes: %s", err.Error())
	}

	var wg sync.WaitGroup
	var mapLock sync.Mutex

	for i, v := range fullMT.Versions {
		wg.Add(1)
		go func(index string, ver FullVersion) {
			defer wg.Done()
			// URL is v.Dist.Tarball
			resp, err := http.Head(ver.Dist.Tarball)
			if err != nil {
				log.Printf("error getting size for package %s at %s: %s", packageName, ver.Dist.Tarball, err.Error())
				return
			}
			// Map is not thread safe so we lock and unlock.
			mapLock.Lock()
			sizes[index] = resp.ContentLength
			mapLock.Unlock()
		}(i, v)
	}
	wg.Wait()
	return sizes, nil
}

// Functions to compare sizes of all version of the same package in different
// registries.
package gonepm

import (
	"fmt"
	"log"
	"sync"
)

// Result holds the package sizes from different registries.
type result struct {
	Reg *Registry
	// URLs  map[string]string
	Sizes map[string]int64
}

// Compare downloads the full metadata for a specific package from all
// registries, performs a HEAD on the tarballs and records the sizes if they
// do not match. The output along with any errors is written to an io.Writer and
// it returns the number of differences.
func Compare(registries []*Registry, packageName string) (int, error) {
	var results []result
	if packageName == "" {
		return 0, fmt.Errorf("gonepm.Compare: package name is empty")
	}
	log.Println("Comparing package", packageName)

	var wg sync.WaitGroup

	for _, reg := range registries {
		wg.Add(1)
		go func(r *Registry) {
			defer wg.Done()
			log.Printf("Retrieving package %s sizes from %s\n", packageName, r.BaseURL)
			sizes, err := r.PackageSizes(packageName)
			if err != nil {
				log.Printf("gonepm.Compare: %s", err.Error())
			}
			log.Printf("Retrieved package %s sizes from %s\n", packageName, r.BaseURL)
			results = append(results, result{Reg: r, Sizes: sizes})
		}(reg)
	}

	wg.Wait()

	counter := 0
	// Start comparing.
	maxResult := max(results)
	log.Printf("Registry with most versions: %s", maxResult.Reg)
	log.Printf("Comparing %d versions.\n", len(maxResult.Sizes))
	for i := 1; i < len(results); i++ {
		for version, size := range results[0].Sizes {
			// If file was not found, size is -1. This reduces false positives
			// when a package/version is missing from one registry.
			// Ditto for size 0.
			if size != results[i].Sizes[version] && size != -1 && results[i].Sizes[version] != -1 && size*results[i].Sizes[version] != 0 {
				log.Printf("Sizes do not match for package %s version %s.\n", packageName, version)
				log.Printf("Size: %d - URL: %s.\n", size,
					fmt.Sprintf("%s/%s/-/%s@%s.tgz", results[0].Reg.BaseURL, packageName, packageName, version))
				log.Printf("Size: %d - URL: %s.\n", results[i].Sizes[version],
					fmt.Sprintf("%s/%s/-/%s@%s.tgz", results[i].Reg.BaseURL, packageName, packageName, version))
				counter++
			}
		}
	}
	log.Println("Finished comparing sizes for", packageName)
	return counter, nil
}

// max returns the registry with the largest dataset, npm is selected if multiple
// are equal (and npm is one of them).
func max(results []result) result {
	var max result
	maxSize := 0
	for _, res := range results {
		// If sizes are the same and we are comparing npm, use npm.
		if len(res.Sizes) == maxSize && res.Reg.BaseURL == NPMAddress {
			max = res
			continue
		}
		// If not npm, check for max size.
		if len(res.Sizes) > maxSize {
			maxSize = len(res.Sizes)
			max = res
		}
	}
	return max
}

// ComparePackages runs Compare on a slice of packages. Returns the total
// number of differences.
func ComparePackages(registries []*Registry, packageList []string) (int, error) {
	log.Println("Starting ComparePackages.")
	if len(packageList) == 0 {
		return 0, fmt.Errorf("gonepm.ComparePackages: packageList is empty")
	}

	counter := 0
	for _, pk := range packageList {
		c, err := Compare(registries, pk)
		if err != nil {
			log.Printf("Error comparing package %s: %v.\n", pk, err)
			continue
		}
		counter += c
	}
	log.Printf("Done with ComparePackages. Compared %d packages.\n", len(packageList))
	return counter, nil
}

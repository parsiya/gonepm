package gonepm

import (
	"encoding/json"
	"fmt"
	"net/http"
	"net/url"
	"strconv"
	"strings"
	"time"
)

const (
	npmSearchEndpoint = "/-/v1/search"
	NPMAddress        = "https://registry.npmjs.com"
)

// SearchResults contains the result from an npm search with
// https://registry.npmjs.com/-/v1/search?text=...
// https://github.com/npm/registry/blob/master/docs/REGISTRY-API.md#get-v1search
type SearchResults struct {
	Objects []struct {
		Package struct {
			Name        string    `json:"name"`
			Version     string    `json:"version"`
			Description string    `json:"description"`
			Keywords    []string  `json:"keywords"`
			Date        time.Time `json:"date"`
			Links       struct {
				Npm        string `json:"npm"`
				Homepage   string `json:"homepage"`
				Repository string `json:"repository"`
				Bugs       string `json:"bugs"`
			} `json:"links"`
			Publisher struct {
				Username string `json:"username"`
				Email    string `json:"email"`
			} `json:"publisher"`
			Maintainers []struct {
				Username string `json:"username"`
				Email    string `json:"email"`
			} `json:"maintainers"`
		} `json:"package"`
		Score struct {
			Final  float64 `json:"final"`
			Detail struct {
				Quality     float64 `json:"quality"`
				Popularity  float64 `json:"popularity"`
				Maintenance float64 `json:"maintenance"`
			} `json:"detail"`
		} `json:"score"`
		SearchScore float64 `json:"searchScore"`
	} `json:"objects"`
	Total int    `json:"total"`
	Time  string `json:"time"`
}

// Search allows fine grained searching.
// See https://github.com/npm/registry/blob/master/docs/REGISTRY-API.md#get-v1search.
// text: string to search for in package name.
// size: number of results. Default/min: 20 and max: 250.
// from: offset to return the results from.
// quality, popularity, and maintenance cannot be zero at the same time. If so,
// popularity = 1.0 is used by default.
func (r *Registry) Search(text string, size int, from int, quality, popularity, maintenance float64) (SearchResults, error) {
	var sr SearchResults
	if text == "" {
		return sr, fmt.Errorf("gonepm.Search: search string is empty")
	}

	// Create request.
	req, err := http.NewRequest("GET", fmt.Sprintf("%s%s", r.BaseURL, npmSearchEndpoint), nil)
	if err != nil {
		return sr, fmt.Errorf("gonepm.Search: %s", err.Error())
	}
	// Size 0 returns nothing in objects.
	if size <= 0 {
		size = 20
	}

	// One of quality, popularity, and maintenance must be non-zero.
	if quality == 0 && popularity == 0 && maintenance == 0 {
		// If all three are zero, go with popularity.
		popularity = 1.0
	}

	// Add params.
	q := url.Values{}
	q.Set("text", text)
	q.Set("size", strconv.Itoa(size))
	q.Set("from", strconv.Itoa(from))
	q.Set("quality", strconv.FormatFloat(quality, 'f', -1, 64))
	q.Set("popularity", strconv.FormatFloat(popularity, 'f', -1, 64))
	q.Set("maintenance", strconv.FormatFloat(maintenance, 'f', -1, 64))

	req.URL.RawQuery = q.Encode()

	client := &http.Client{}

	resp, err := client.Do(req)
	if err != nil {
		return sr, fmt.Errorf("gonepm.Search: %s", err.Error())
	}
	defer resp.Body.Close()

	if err := json.NewDecoder(resp.Body).Decode(&sr); err != nil {
		return sr, fmt.Errorf("gonepm.QuickSearch: %s", err.Error())
	}
	return sr, nil
}

// QuickSearch queries npm for a string.
func (r *Registry) QuickSearch(text string) (SearchResults, error) {
	return r.Search(text, 0, 0, 0, 0, 0)
}

// AuthorSearch returns all packages by a specific author.
func (r *Registry) AuthorSearch(author string) (SearchResults, error) {
	return r.QuickSearch("author:" + author)
}

// MaintainerSearch returns all packages by a specific Maintainer.
func (r *Registry) MaintainerSearch(maintainer string) (SearchResults, error) {
	return r.QuickSearch("maintainer:" + maintainer)
}

// KeywordSearch returns all packages with a specific keyword.
// Query can have multiple keywords:
// OR: ",". E.g. "word1,word2"
// And: "+". E.g. "word1+word2"
// Exclude: ",-". E.g. "word1,-word2"
func (r *Registry) KeywordSearch(keywords string) (SearchResults, error) {
	// Remove all space.
	keywords = strings.Replace(keywords, " ", "", -1)
	return r.QuickSearch("keywords:" + keywords)
}

// QuickSearchOld queries npm for a string. Initial version.
func (r *Registry) QuickSearchOld(text string) (SearchResults, error) {
	var sr SearchResults
	if text == "" {
		return sr, fmt.Errorf("gonepm.QuickSearch: search string is empty")
	}

	q := fmt.Sprintf("%s%s?text=%s", r.BaseURL, npmSearchEndpoint, text)

	resp, err := http.Get(q)
	if err != nil {
		return sr, fmt.Errorf("gonepm.QuickSearch: %s", err.Error())
	}
	defer resp.Body.Close()

	if err := json.NewDecoder(resp.Body).Decode(&sr); err != nil {
		return sr, fmt.Errorf("gonepm.QuickSearch: %s", err.Error())
	}
	return sr, nil
}

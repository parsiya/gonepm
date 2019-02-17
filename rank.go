// Functions and methods for getting npm packages by rank for several metrics.
package gonepm

import (
	"fmt"
	"io/ioutil"
	"log"
	"net/http"
	"regexp"
	"strings"
)

const (
	AvankaMostDependent = "https://gist.githubusercontent.com/anvaka/8e8fa57c7ee1350e3491/raw/8bafd425f48f713c1629bad6ef199fddd9fb2216/01.most-dependent-upon.md"
	AvankaTopPageRank   = "https://gist.githubusercontent.com/anvaka/8e8fa57c7ee1350e3491/raw/dbab7af56bd09a458a99618ce7a9cc0c62d852f1/03.pagerank.md"
)

// ParseAvankaList parses one of Avanka's npm rank lists at
// https://gist.github.com/anvaka/8e8fa57c7ee1350e3491 and returns a []string of
// package names.
// # Top 1000 most depended-upon packages
// 0. [lodash](https://www.npmjs.org/package/lodash) - 50647
// 1. [request](https://www.npmjs.org/package/request) - 29350
// 2. [chalk](https://www.npmjs.org/package/chalk) - 26737
// 3. [commander](https://www.npmjs.org/package/commander) - 23111
// 4. [express](https://www.npmjs.org/package/express) - 21166
// 5. [async](https://www.npmjs.org/package/async) - 20462
// 6. [react](https://www.npmjs.org/package/react) - 18994
// 7. [debug](https://www.npmjs.org/package/debug) - 17441
// We grab everything between brackets.
func ParseAvankaList(URL string) ([]string, error) {
	log.Println("Getting packages from", URL)
	var pk []string
	resp, err := http.Get(URL)
	if err != nil {
		return pk, fmt.Errorf("gonepm.Top1000Dependents: %s", err.Error())
	}
	defer resp.Body.Close()

	allResp, err := ioutil.ReadAll(resp.Body)
	if err != nil {
		return pk, fmt.Errorf("gonepm.Top1000Dependents: %s", err.Error())
	}
	re := regexp.MustCompile(`\[(.*?)\]`)
	pk = re.FindAllString(string(allResp), -1)

	// This matches the brackets too, so we need to remove them.
	for i := range pk {
		pk[i] = strings.Trim(pk[i], "[]")
	}
	log.Printf("Processed %d packages.\n", len(pk))
	return pk, nil
}

// Top1000Dependents returns a []string of top most dependent packages from
// Avanka's list.
func Top1000Dependents() ([]string, error) {
	return ParseAvankaList(AvankaMostDependent)
}

// Top1000PageRank returns a []string of top npm packages by page rank
// from Avanka's list.
func TopPageRankAvanka() ([]string, error) {
	return ParseAvankaList(AvankaTopPageRank)
}

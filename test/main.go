package main

import (
	"fmt"
	"io"
	"log"
	"os"

	"github.com/parsiya/gonepm"
)

func main() {

	// reg, err := gonepm.NewRegistry("https://registry.npmjs.com/")
	// if err != nil {
	// 	panic(err)
	// }

	// packageName := "lodash"

	// fmt.Println(reg)

	// fmt.Println(reg.RegInfo())

	// res1, err := reg.QuickSearch("meow")
	// if err != nil {
	// 	panic(err)
	// }
	// fmt.Println(res1.Objects[0].Package.Publisher.Email)

	// res, err := reg.Search("meow", 0, 0, 0, 0, 0)
	// if err != nil {
	// 	panic(err)
	// }
	// fmt.Println(res.Objects[0].Package.Publisher.Email)

	// Get a package.
	// fullMT, err := reg.PackageMetadata(packageName)
	// if err != nil {
	// 	panic(err)
	// }

	// sizes := make(map[string]int64)

	// for i, v := range fullMT.Versions {
	// 	// URL is v.Dist.Tarball
	// 	res, err := http.Head(v.Dist.Tarball)
	// 	if err != nil {
	// 		panic(err)
	// 	}
	// 	sizes[i] = res.ContentLength
	// }

	// f, _ := os.Create("result")
	// defer f.Close()

	// tao, err := gonepm.NewRegistry("http://registry.npm.taobao.org")
	// if err != nil {
	// 	panic(err)
	// }
	// fullTao, err := tao.PackageMetadata(packageName)
	// if err != nil {
	// 	panic(err)
	// }

	// sizesTao := make(map[string]int64)

	// for i, v := range fullTao.Versions {
	// 	// URL is v.Dist.Tarball
	// 	res, err := http.Head(v.Dist.Tarball)
	// 	if err != nil {
	// 		panic(err)
	// 	}
	// 	sizesTao[i] = res.ContentLength
	// }

	// if len(sizesTao) != len(sizes) {
	// 	fmt.Println("sizes do not match")
	// 	fmt.Println("len(sizes)", len(sizes))
	// 	fmt.Println("len(sizesTao)", len(sizesTao))
	// }

	// for i, v := range sizesTao {

	// 	if v != sizes[i] {
	// 		f.WriteString(fmt.Sprintf("%s: %d\n", i, v))
	// 	}
	// }

	// **************

	// Create log file.
	f, err := os.Create("top1000-run-1.txt")
	if err != nil {
		panic(err)
	}
	defer f.Close()

	mw := io.MultiWriter(f, os.Stdout)

	log.SetOutput(mw)

	// Get top 1000 dependent packages from github gist.
	packages, err := gonepm.Top1000Dependents()
	if err != nil {
		panic(err)
	}

	// Define registries.
	registries := []string{
		"https://registry.npmjs.com/",
		"http://registry.npm.taobao.org",
		"https://registry.cnpmjs.org/",
	}

	// Make registries.
	regz, err := gonepm.MakeRegistries(registries)
	if err != nil {
		panic(err)
	}

	// packages := []string{
	// 	"lodash",
	// 	// "request",
	// 	// "chalk",
	// 	// "commander",
	// 	// "express",
	// 	// "react",
	// 	// "async",
	// 	// "debug",
	// 	// "bluebird",
	// 	// "yargs",
	// 	// "q",
	// 	// "vue",
	// 	// "gulp",
	// 	// "@angular/core",
	// 	// "optimist",
	// 	// "co",
	// }

	// Do the thingie.
	counter, err := gonepm.ComparePackages(regz, packages)
	if err != nil {
		panic(err)
	}
	fmt.Println("Number of differences:", counter)
	fmt.Println("Done")

	// Check why q is not working.
	// reg, err := gonepm.NewRegistry(gonepm.NPMAddress)
	// if err != nil {
	// 	panic(err)
	// }
	// mt, err := reg.PackageMetadata("q")
	// if err != nil {
	// 	panic(err)
	// }

	// resp, err := http.Get("https://registry.npmjs.com/q")
	// if err != nil {
	// 	panic(err)
	// }
	// defer resp.Body.Close()
	// f, err := os.Create("q-resp.json")
	// if err != nil {
	// 	panic(err)
	// }
	// defer f.Close()
	// io.Copy(f, resp.Body)

}

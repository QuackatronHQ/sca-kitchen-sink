package main

import (
	"encoding/json"
	"fmt"
	"os"
)

func main() {
	if len(os.Args) != 2 {
		fmt.Fprintf(os.Stderr, "Usage: %s <path-to-deps.json>\n", os.Args[0])
		os.Exit(1)
	}
	path := os.Args[1]
	direct, err := RootDependencies(path)
	// print the direct dependencies
	for dep := range direct {
		fmt.Println(dep)
	}
	if err != nil {
		fmt.Fprintf(os.Stderr, "Error: %v\n", err)
		os.Exit(1)
	}
	for dep := range direct {
		fmt.Println(dep)
	}
}

func RootDependencies(path string) (map[string]struct{}, error) {
	type pkg struct {
		Dependencies map[string]string `json:"dependencies"`
	}
	type libraryItem struct {
		Type string `json:"type"`
	}
	type depsJSON struct {
		Targets   map[string]map[string]pkg `json:"targets"`
		Libraries map[string]libraryItem    `json:"libraries"`
	}

	depsFile, err := os.ReadFile(path)
	if err != nil {
		return nil, fmt.Errorf("failed to read deps.json: %w", err)
	}
	var depsData depsJSON
	err = json.Unmarshal(depsFile, &depsData)
	if err != nil {
		return nil, fmt.Errorf("failed to unmarshal deps.json: %w", err)
	}

	projectNames := make(map[string]struct{})
	directDepMap := make(map[string]struct{})

	// get name of project from libraries
	for libName, libItem := range depsData.Libraries {
		if libItem.Type == "project" {
			projectNames[libName] = struct{}{}
		}
	}
	// now iterate over targets and get direct dependencies.
	// dependencies for package names that are not in projectNames are considered direct dependencies.
	for _, target := range depsData.Targets {
		for pkgName, pkgItem := range target {
			if _, exists := projectNames[pkgName]; exists {
				// this is a package entry. All its dependencies are direct dependencies.
				for depName := range pkgItem.Dependencies {
					directDepMap[depName] = struct{}{}
				}
			}
		}
	}
	fmt.Println("Direct dependencies: ", len(directDepMap))
	for dep := range directDepMap {
		fmt.Println(dep)
	}
	return directDepMap, nil
}

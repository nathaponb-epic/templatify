package utils

import (
	"fmt"
	"io/fs"
	"os"
	"path/filepath"
	"testing"
)

type Confix struct {
	Name      string
	Domain    string
	Path      string
	AppPrefix string
	RootPath  struct {
		Cmd      string
		CSS      string
		Script   string
		Font     string
		Constant string
	}
	IgnoreDir  []string
	IgnoreFile []string
}

func TestWalker(t *testing.T) {

	config := Confix{
		Name:      "MyApp",
		Domain:    "https://example.com",
		Path:      "/myID",
		AppPrefix: "myprefix",
		RootPath: struct {
			Cmd      string
			CSS      string
			Script   string
			Font     string
			Constant string
		}{
			Cmd:      "/cmd",
			CSS:      "/css",
			Script:   "/scripts",
			Font:     "/fonts",
			Constant: "/constants",
		},
		IgnoreDir:  []string{"cmd", "pkg", "fonts", "images", "mail"},
		IgnoreFile: []string{"file1", "file2"},
	}

	currentDir, err := os.Getwd()
	if err != nil {
		t.Error(err)
	}

	parrentDir := filepath.Dir(currentDir)
	rootDir := filepath.Dir(parrentDir)

	fmt.Println(rootDir)

	// if there is potentially error on traverse file

	err = filepath.Walk(rootDir, visitFileX(&config))
	if err != nil {
		t.Error(err)
	}

	// [opt] check visited files and ignore if correct
}

func visitFileX(config *Confix) filepath.WalkFunc {
	return func(path string, info fs.FileInfo, err error) error {
		if err != nil {
			fmt.Printf("error accessing path %s: %v\n", path, err)
			return nil
		}

		if info.IsDir() {

			var skip bool
			for _, v := range config.IgnoreDir {
				if info.Name() == v {
					skip = true
				}
			}

			if skip {
				fmt.Printf("Skipped dir: %s\n", path)
				return filepath.SkipDir
			}

		}

		var shouldSkipThisFile bool
		for _, v := range config.IgnoreFile {
			if info.Name() == v {
				shouldSkipThisFile = true
			}
		}

		if shouldSkipThisFile {
			fmt.Printf("Skipped file: %s\n", path)
			return nil
		}

		fmt.Printf("Processed file: %s\n", path)
		return nil
	}
}

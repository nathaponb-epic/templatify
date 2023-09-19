package utils

import (
	"fmt"
	"os"
	"path/filepath"
)

type CMD struct {
	Configulation []Configuration `mapstructure:"cmd"`
}
type Configuration struct {
	Name      string `mapstructure:"name"`
	Domain    string `mapstructure:"domain"`
	Path      string `mapstructure:"path"`
	AppFolder string `mapstructure:"app_folder"`
	RootPath  struct {
		Image    string `mapstructure:"image"`
		CSS      string `mapstructure:"css"`
		Script   string `mapstructure:"script"`
		Font     string `mapstructure:"font"`
		Constant string `mapstructure:"constant"`
	} `mapstructure:"root_path"`
	IgnoreDir  []string `mapstructure:"ignore_dir"`
	IgnoreFile []string `mapstructure:"ignore_file"`
}

func Walker(config Configuration) error {

	rootDir, _ := os.Getwd() // get current directory as root walk dir

	err := filepath.Walk(rootDir, visitFile(config))
	if err != nil {
		return err
	}

	return nil
}

func visitFile(config Configuration) filepath.WalkFunc {

	return func(path string, info os.FileInfo, err error) error {

		if err != nil {
			fmt.Printf("error accessing path %s: %v\n", path, err)
			return nil
		}

		// skip directory that define in ignore_dir attribute
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

		// skip file that define in ignore_file attribute
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

		// process file based on file extension
		if filepath.Ext(path) != "" {

			fileType := filepath.Ext(path)[1:]

			switch fileType {
			case "html":
				processHTML(path, &config)
			case "css":
				processCSS(path, &config)
			case "js":
				processJS(path, &config)
			}
		}
		return nil
	}
}

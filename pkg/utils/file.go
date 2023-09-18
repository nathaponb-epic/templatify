package utils

import (
	"fmt"
	"path/filepath"
	"strings"
)

var supportFileType map[string]string

func verifyFileType(refPath string, config Configuration) string {

	fileSurname := filepath.Ext(refPath)
	if fileSurname == "" {
		return refPath
	}

	if fileSurname[1:] == "html" {
		return refPath
	}

	supportFileType = map[string]string{
		"png":    config.Image,
		"svg":    config.Image,
		"ico":    config.Image,
		"gif":    config.Image,
		"jpeg":   config.Image,
		"js":     config.Script,
		"css":    config.Script,
		"json":   config.Script,
		"ttf":    config.Font,
		"woff":   config.Font,
		"woff2":  config.Font,
		"icloud": config.Font,
	}

	testFileType := fileSurname[1:]

	if testFileType[len(testFileType)-1] == '\'' || testFileType[len(testFileType)-1] == '"' {
		testFileType = testFileType[:len(testFileType)-1]
	}

	if supportFileType[testFileType] == "" {
		return refPath
	}

	// preserve the sub-dir path
	finalPath := preserveSubDir(refPath)
	if finalPath == "" {
		return refPath
	}

	return fmt.Sprintf("%s%s%s", config.Domain, config.Path, finalPath)

}

func preserveSubDir(refPath string) string {

	slashedPaths := strings.Split(refPath, "/")

	// find the file type
	lastSlash := slashedPaths[len(slashedPaths)-1]
	dotType := filepath.Ext(lastSlash)

	if dotType != "" {

		fileType := dotType[1:]

		// if filetype is empty or is not valid url return out as it
		if fileType == "" {
			return refPath
		}

		// exeption for last index symbol '
		if fileType[len(fileType)-1] == '\'' || fileType[len(fileType)-1] == '"' {
			fileType = fileType[:len(fileType)-1]
		}

		// get default root path of fileType
		var rootPath string
		for k, v := range supportFileType {
			if k == fileType {
				rootPath = v[1:] // slice out forward slash
			}
		}

		// use rootPath to find the index of exact name in full path
		var startIndex int
		var found bool
		for i, v := range slashedPaths {
			if v == rootPath {
				startIndex = i
				found = true
			}
		}

		if !found {
			// return with default file type root path
			return fmt.Sprintf("%s/%s", supportFileType[fileType], lastSlash)
		} else {

			target := slashedPaths[startIndex:]

			join := strings.Join(target, "/")

			if join[len(join)-1] == '\'' || join[len(join)-1] == '"' {
				join = join[:len(join)-1]
			}

			return "/" + join
		}

	}

	return ""
	// if cannot find file type, possible something like Lisense file

}

// func isValidURL(s string) bool {

// 	_, err := url.ParseRequestURI(s)
// 	if err != nil {
// 		return false
// 	}

// 	return true
// 	// return u.Scheme != "" && strings.Contains(u.Scheme, "http")
// }

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
		"png":    config.RootPath.Image,
		"svg":    config.RootPath.Image,
		"ico":    config.RootPath.Image,
		"gif":    config.RootPath.Image,
		"jpeg":   config.RootPath.Image,
		"js":     config.RootPath.Script,
		"css":    config.RootPath.Script,
		"json":   config.RootPath.Script,
		"ttf":    config.RootPath.Font,
		"woff":   config.RootPath.Font,
		"woff2":  config.RootPath.Font,
		"icloud": config.RootPath.Font,
	}

	testFileType := fileSurname[1:]

	// if testFileType[len(testFileType)-1] == '\'' || testFileType[len(testFileType)-1] == '"' || testFileType[len(testFileType)-1] == '`' {
	// 	testFileType = testFileType[:len(testFileType)-1]
	// }
	testFileType = unQuoteSuffix(testFileType)

	if supportFileType[testFileType] == "" {
		return refPath
	}

	// preserve the sub-dir path
	finalPath := preserveSubDir(refPath, config)
	if finalPath == "" {
		return refPath
	}

	return fmt.Sprintf("%s%s%s", config.Domain, config.Path, finalPath)

}

func preserveSubDir(refPath string, config Configuration) string {

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

		fileType = unQuoteSuffix(fileType)

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
			tmpExactRootPath := config.TemplateExacRoot
			paths := strings.Split(refPath, tmpExactRootPath)
			if len(paths) > 1 {
				return paths[1]
			}

		} else {

			target := slashedPaths[startIndex:]

			join := strings.Join(target, "/")

			// if join[len(join)-1] == '\'' || join[len(join)-1] == '"' {
			// 	join = join[:len(join)-1]
			// }
			join = unQuoteSuffix(join)

			return "/" + join
		}

	}

	return ""
	// if cannot find file type, possible something like Lisense file

}

func unQuoteSuffix(s string) string {
	if s == "" {
		return s
	}

	if s[len(s)-1] == '\'' || s[len(s)-1] == '"' || s[len(s)-1] == '`' {
		s = s[:len(s)-1]
	}

	return s
}

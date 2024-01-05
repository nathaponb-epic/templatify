package utils

import (
	"fmt"
	"io"
	"log"
	"os"
	"path/filepath"
	"regexp"
	"strings"

	"golang.org/x/net/html"
)

func processHTML(fileAbsPath string, config *Configuration) {

	// handle slash path for every OS
	fileAbsPath = filepath.ToSlash(fileAbsPath)

	// read file content
	htmlFile, err := os.Open(fileAbsPath)
	if err != nil {
		log.Fatal(err)
	}
	defer htmlFile.Close()

	htmlByte, _ := io.ReadAll(htmlFile)

	doc, err := html.Parse(strings.NewReader(string(htmlByte)))
	if err != nil {
		log.Fatal(err)
	}

	// recursively find specific node element on list linked node
	var findAndReplace func(*html.Node)

	findAndReplace = func(n *html.Node) {

		if n.Type == html.ElementNode {
			// Check if the element is an <img>, <a>, or <script> tag
			if n.Data == "img" || n.Data == "link" || n.Data == "script" || n.Data == "style" {
				// Iterate through the element's attributes
				for i, attr := range n.Attr {

					if (n.Data == "img" && attr.Key == "src") ||
						(n.Data == "link" && attr.Key == "href") ||
						(n.Data == "script" && attr.Key == "src") {

						refPath := n.Attr[i].Val

						// Do update ref path
						newRefPath := verifyFileType(refPath, *config)
						n.Attr[i].Val = newRefPath
					}

				}

			}
			// process javascript code embeded in html file <script/>
			if n.Data == "script" {

				//*: refactor this logic onto a function use with processJS
				if n.FirstChild != nil {

					jsContent := n.FirstChild.Data

					// pattern := `\$\.getJSON\(([^,]+),`
					// regex := regexp.MustCompile(pattern)

					// matches := regex.FindAllStringSubmatch(jsContent, -1)

					// kv := make(map[string]string)

					// for _, match := range matches {
					// 	if len(match) > 1 {
					// 		kv[match[1]] = ""
					// 	}
					// }

					kv := regexJs(jsContent)

					for k, _ := range kv {

						slashes := strings.Split(k, "/")
						fileName := slashes[len(slashes)-1]

						replaceVal := verifyFileType(fileName[:len(fileName)-1], *config)

						kv[k] = replaceVal
					}

					for find, replace := range kv {
						jsContent = strings.Replace(jsContent, find, fmt.Sprintf("'%s'", replace), -1)
					}

					n.FirstChild.Data = jsContent
				}

			}
			if n.Data == "style" {

				styleContent := n.FirstChild.Data

				// fmt.Println(styleContent)
				kv := make(map[string]string)

				pattern := `src: url\(([^)]+)\)`

				re := regexp.MustCompile(pattern)

				matches := re.FindAllStringSubmatch(styleContent, -1)

				for _, match := range matches {
					if len(match) > 1 {
						kv[match[1]] = ""
					}
				}

				for k, _ := range kv {
					replaceVal := verifyFileType(k[:len(k)-1], *config)
					kv[k] = replaceVal
				}

				for find, replace := range kv {
					styleContent = strings.Replace(styleContent, find, fmt.Sprintf("'%s'", replace), -1)
				}

				n.FirstChild.Data = styleContent
			}
		}
		// Recursively traverse child nodes, until there is no nextSibling
		for c := n.FirstChild; c != nil; c = c.NextSibling {
			findAndReplace(c)
		}
	}

	findAndReplace(doc)

	var modifiedHTML strings.Builder

	html.Render(&modifiedHTML, doc)
	modifiedContent := modifiedHTML.String()
	h := html.UnescapeString(modifiedContent)

	err = os.WriteFile(fileAbsPath, []byte(h), 0644)
	if err != nil {
		fmt.Println("Error writing HTML file:", err)
		return
	}

	fmt.Printf("✔ Updated File: %s\n", fileAbsPath)
}

func processCSS(fileAbsPath string, config *Configuration) {

	// handle slash path for every OS
	fileAbsPath = filepath.ToSlash(fileAbsPath)

	// open file and read file content
	cssFile, err := os.Open(fileAbsPath)
	if err != nil {
		log.Fatal(err)
	}
	defer cssFile.Close()

	cssByte, err := io.ReadAll(cssFile)
	if err != nil {
		log.Fatal(err)
	}
	cssContent := string(cssByte)

	// find all the targets of replacing by regex
	patternSrc := `src: url\((.*?)\)`
	regexSrc := regexp.MustCompile(patternSrc)
	patternSrcMatches := regexSrc.FindAllStringSubmatch(cssContent, -1)

	kv := make(map[string]string)

	for _, match := range patternSrcMatches {
		if len(match) > 1 {
			kv[match[1]] = ""
		}
	}

	// capture background-color path
	patternBg := `background-image\s*:\s*url\((.*?)\)`
	regexBg := regexp.MustCompile(patternBg)
	patternBgMatches := regexBg.FindAllStringSubmatch(cssContent, -1)

	for _, match := range patternBgMatches {
		if len(match) > 1 && len(match[1]) < 100 {
			kv[match[1]] = ""
		}
	}

	for k, _ := range kv {
		// slashes := strings.Split(k, "/")
		// fileName := slashes[len(slashes)-1]

		replaceVal := verifyFileType(k, *config)

		kv[k] = replaceVal
	}

	// find and replace based on matches and its coresponding value from verifyFileType
	for find, replace := range kv {
		cssContent = strings.Replace(cssContent, find, fmt.Sprintf("'%s'", replace), -1)
	}

	err = os.WriteFile(fileAbsPath, []byte(cssContent), 0644)
	if err != nil {
		log.Fatal(err)
	}

	fmt.Printf("✔ Updated File: %s\n", fileAbsPath)
}

func processJS(fileAbsPath string, config *Configuration) {

	// handle slash path for every OS
	fileAbsPath = filepath.ToSlash(fileAbsPath)

	// open file and read file content
	jsFile, err := os.Open(fileAbsPath)
	if err != nil {
		log.Fatal(err)
	}
	defer jsFile.Close()

	jsByte, err := io.ReadAll(jsFile)
	if err != nil {
		log.Fatal(err)
	}

	jsContent := string(jsByte)

	// if the JS file named prefix.js do some extra work
	filePaths := strings.Split(fileAbsPath, "/")
	fileNameWithExt := strings.Split(filePaths[len(filePaths)-1], ".")

	if fileNameWithExt[0] == "prefix" {

		appPrefixPattern := `const PREFIX = (.*?);`
		appFolderpattern := `const FOLDER = (.*?);`

		appPrefixRegex := regexp.MustCompile(appPrefixPattern)
		appFolderRegex := regexp.MustCompile(appFolderpattern)

		appPrefixmatches := appPrefixRegex.FindAllStringSubmatch(jsContent, -1)
		appFoldermatches := appFolderRegex.FindAllStringSubmatch(jsContent, -1)

		kv := make(map[string]string)

		// prepare value from user define yaml config file
		for _, match := range appPrefixmatches {
			if len(match) > 1 {
				kv[match[1]] = config.AppPrefix
			}
		}

		for _, match := range appFoldermatches {
			if len(match) > 1 {
				kv[match[1]] = config.AppFolder
			}
		}

		for find, replace := range kv {
			jsContent = strings.Replace(jsContent, find, fmt.Sprintf("'%s'", replace), -1)
		}

	}

	kv := regexJs(jsContent)

	// verify ref file type from all the matches before replacing
	for k, _ := range kv {
		replaceVal := verifyFileType(k, *config)
		kv[k] = replaceVal
	}

	for find, replace := range kv {

		// replace = `"` + replace + `"`
		// parse literal string to raw string
		replace = fmt.Sprintf(`"%s"`, replace)

		// if cmd is localify append PREFIX + FOLDER +
		if config.Name == "localify" {

			//* slice out the prefix since in js file/tag the prefix is set as variable
			// - find root_path by file extension
			// - use it to slice out the
			replace = sliceoutPrefixString(replace)

			pathPrefix := "PREFIX + FOLDER +"

			replace = fmt.Sprintf("%s%s", pathPrefix, replace)
		}

		jsContent = strings.Replace(jsContent, find, replace, -1)
	}

	err = os.WriteFile(fileAbsPath, []byte(jsContent), 0644)
	if err != nil {
		log.Fatal(err)
	}

	fmt.Printf("✔ Updated File: %s\n", fileAbsPath)

}

func sliceoutPrefixString(s string) string {
	// get file extension
	fileExt := filepath.Ext(s)

	fileExt = unQuoteSuffix(fileExt)

	// get root_path
	rootPath := supportFileType[fileExt[1:]]
	if rootPath == "" {
		return s
	}

	var targetIndex int

	// slice the string by slash /
	paths := strings.Split(s, "/")

	for i, v := range paths {
		// remove /
		if v == rootPath[1:] {
			targetIndex = i
		}
	}

	newPath := paths[targetIndex:]
	join := strings.Join(newPath, "/")

	join = unQuoteSuffix(join)

	return fmt.Sprintf(`"/%s"`, join)

}

// find the matches of target javascript content
func regexJs(content string) map[string]string {

	kv := make(map[string]string)

	// pattern getJSON
	patternGetJSON := `\$\.getJSON\(([^,]+),`
	regexGetJSON := regexp.MustCompile(patternGetJSON)

	patternGetJSONMatches := regexGetJSON.FindAllStringSubmatch(content, -1)

	for _, match := range patternGetJSONMatches {
		if len(match) > 1 {
			// fmt.Printf("getJSON value: %s\n", match[1])
			kv[match[1]] = ""
		}
	}

	// pattern imageUrl
	patternImgUrl := `imageUrl:(.*?),`
	regexImgUrl := regexp.MustCompile(patternImgUrl)

	patternImgUrlMatches := regexImgUrl.FindAllStringSubmatch(content, -1)

	// combine every matches
	for _, match := range patternImgUrlMatches {
		if len(match) > 1 {
			kv[match[1]] = ""
		}
	}

	// pattern .attr src
	patternAttr := `\$\([^)]+\)\.attr\('src', (.*?);`
	regexAttr := regexp.MustCompile(patternAttr)
	patternAttrMatches := regexAttr.FindAllStringSubmatch(content, -1)

	for _, match := range patternAttrMatches {
		if len(match) > 1 {
			kv[match[1][:len(match[1])-1]] = ""
		}
	}

	return kv

}

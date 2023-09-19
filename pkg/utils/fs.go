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
			if n.Data == "img" || n.Data == "link" || n.Data == "script" {
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
				if n.FirstChild != nil {

					jsContent := n.FirstChild.Data

					pattern := `\$\.getJSON\(([^,]+),`
					regex := regexp.MustCompile(pattern)

					matches := regex.FindAllStringSubmatch(jsContent, -1)

					kv := make(map[string]string)

					for _, match := range matches {
						if len(match) > 1 {
							kv[match[1]] = ""
						}
					}

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

	fmt.Printf("Updated file: %s\n", fileAbsPath)
}

func processCSS(fileAbsPath string, config *Configuration) {

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
		if len(match) > 1 {
			kv[match[1]] = ""
		}
	}

	for k, _ := range kv {
		slashes := strings.Split(k, "/")
		fileName := slashes[len(slashes)-1]

		replaceVal := verifyFileType(fileName[:len(fileName)-1], *config)

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

	fmt.Printf("Updated file: %s\n", fileAbsPath)
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

		pattern := `const FOLDER = (.*?);`
		regex := regexp.MustCompile(pattern)
		matches := regex.FindAllStringSubmatch(jsContent, -1)

		kv := make(map[string]string)

		// prepare value from user define yaml config file
		for _, match := range matches {
			if len(match) > 1 {
				kv[match[1]] = config.AppPrefix
			}
		}

		for find, replace := range kv {
			jsContent = strings.Replace(jsContent, find, fmt.Sprintf("'%s'", replace), -1)
		}

	}

	// pattern getJSON
	patternGetJSON := `\$\.getJSON\(([^,]+),`
	regexGetJSON := regexp.MustCompile(patternGetJSON)

	patternGetJSONMatches := regexGetJSON.FindAllStringSubmatch(jsContent, -1)

	kv := make(map[string]string)

	for _, match := range patternGetJSONMatches {
		if len(match) > 1 {
			kv[match[1]] = ""
		}
	}

	// pattern imageUrl
	patternImgUrl := `imageUrl:(.*?),`
	regexImgUrl := regexp.MustCompile(patternImgUrl)

	patternImgUrlMatches := regexImgUrl.FindAllStringSubmatch(jsContent, -1)

	// combine every matches
	for _, match := range patternImgUrlMatches {
		if len(match) > 1 {
			kv[match[1]] = ""
		}
	}

	// pattern .attr src
	patternAttr := `\$\([^)]+\)\.attr\('src', (.*?);`
	regexAttr := regexp.MustCompile(patternAttr)
	patternAttrMatches := regexAttr.FindAllStringSubmatch(jsContent, -1)

	for _, match := range patternAttrMatches {
		if len(match) > 1 {
			kv[match[1][:len(match[1])-1]] = ""
		}
	}

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

			pathPrefix := "PREFIX + FOLDER +"

			replace = fmt.Sprintf("%s%s", pathPrefix, replace)
		}

		jsContent = strings.Replace(jsContent, find, replace, -1)
	}

	err = os.WriteFile(fileAbsPath, []byte(jsContent), 0644)
	if err != nil {
		log.Fatal(err)
	}

	fmt.Printf("Updated file: %s\n", fileAbsPath)

}

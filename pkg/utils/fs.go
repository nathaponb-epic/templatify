package utils

import (
	"fmt"
	"io"
	"log"
	"os"
	"regexp"
	"strings"

	"golang.org/x/net/html"
)

func processHTML(filePathAbs string, config *Configuration) {

	// read file content
	htmlFile, err := os.Open(filePathAbs)
	if err != nil {
		log.Fatal(err)
	}

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

	err = os.WriteFile(filePathAbs, []byte(h), 0644)
	if err != nil {
		fmt.Println("Error writing HTML file:", err)
		return
	}

	fmt.Printf("Updated file: %s\n", filePathAbs)
}

func processCSS(filePathAbs string, config *Configuration) {

}

func processJS(filePathAbs string, config *Configuration) {

}

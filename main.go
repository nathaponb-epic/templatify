package main

import "github.com/nathaponb-epic/templatify/cmd"

func main() {

	cmd.Execute()

}

// func testQuote() {

// 	str := `"images/logo.png"`
// 	prefix := "PREFIX + FOLDER +"

// 	// fmt.Sprintf(``, a)

// 	output := fmt.Sprintf("%s %s", prefix, str)

// 	err := os.WriteFile("testQuote.js", []byte(output), 0644)
// 	if err != nil {
// 		fmt.Println("Error writing HTML file:", err)
// 		return
// 	}

// 	fmt.Println("successfully write js file")
// }

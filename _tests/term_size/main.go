package main

import (
	"bufio"
	"os"
)

func main() {

	scanner := bufio.NewScanner(os.Stdin)

	var text string
	for text != "t" { // break the loop if text == "q"
		scanner.Scan()
		print("\033[14t")
		text = scanner.Text()
	}

	println(text)
}

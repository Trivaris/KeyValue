package main

import (
	"fmt"
	"os"
)

func main() {
	if len(os.Args) < 2 {
		printHelp()
		return
	}
	for _, arg := range os.Args {
		fmt.Print(arg)
	}
	fmt.Println("")
}

func printHelp() {
	fmt.Println("`Help Message`")
}
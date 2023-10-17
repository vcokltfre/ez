package main

import (
	"fmt"
	"os"
)

func main() {
	if len(os.Args) < 2 {
		fmt.Printf("Usage: %s <filename> [opts]", os.Args[1])
		return
	}

	data, err := os.ReadFile(os.Args[1])
	if err != nil {
		fmt.Println(err)
		return
	}

	fmt.Println(string(data)) // TODO: Run code
}

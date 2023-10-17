package main

import (
	"fmt"
	"os"

	"github.com/vcokltfre/ez/ez"
)

func main() {
	if len(os.Args) < 2 {
		fmt.Printf("Usage: %s <filename> [opts]", os.Args[0])
		return
	}

	data, err := os.ReadFile(os.Args[1])
	if err != nil {
		fmt.Println(err)
		return
	}

	err = ez.Run(string(data), os.Args[1])
	if err != nil {
		fmt.Println(err)
		return
	}
}

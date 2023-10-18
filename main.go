package main

import (
	"fmt"
	"os"
	"strconv"
	"strings"

	"github.com/vcokltfre/ez/ez"
)

func main() {
	if len(os.Args) < 2 {
		fmt.Printf("Usage: %s <filename> [opts]", os.Args[0])
		return
	}

	opts := make(map[string]string)
	for i := 2; i < len(os.Args); i++ {
		parts := strings.Split(os.Args[i], "=")
		if len(parts) != 2 {
			fmt.Println("Invalid option:", os.Args[i])
			return
		}

		opts[parts[0]] = parts[1]
	}

	data, err := os.ReadFile(os.Args[1])
	if err != nil {
		fmt.Println(err)
		return
	}

	memory := 1 << 16
	if val, ok := opts["memory"]; ok {
		memory, err = strconv.Atoi(val)
		if err != nil {
			fmt.Println(err)
			return
		}
	}

	err = ez.Run(string(data), os.Args[1], memory)
	if err != nil {
		fmt.Println(err)
		return
	}
}

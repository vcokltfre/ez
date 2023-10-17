package ez

import (
	"fmt"

	"github.com/vcokltfre/ez/ez/lexer"
)

func Run(code, filename string) error {
	tokens, err := lexer.Lex(code, filename)
	if err != nil {
		return err
	}

	fmt.Println(tokens)

	return err
}

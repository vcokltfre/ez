package ez

import (
	"fmt"

	"github.com/vcokltfre/ez/ez/lexer"
	"github.com/vcokltfre/ez/ez/parser"
)

func Run(code, filename string) error {
	tokens, err := lexer.Lex(code, filename)
	if err != nil {
		return err
	}

	program, err := parser.Parse(tokens)
	if err != nil {
		return err
	}

	fmt.Println(*program)

	return nil
}

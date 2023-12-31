package ez

import (
	"github.com/vcokltfre/ez/ez/lexer"
	"github.com/vcokltfre/ez/ez/parser"
	"github.com/vcokltfre/ez/ez/vm"
)

func Run(code, filename string, memory int) error {
	tokens, err := lexer.Lex(code, filename)
	if err != nil {
		return err
	}

	program, err := parser.Parse(tokens)
	if err != nil {
		return err
	}

	executor := vm.New(memory)

	return executor.Run(program)
}

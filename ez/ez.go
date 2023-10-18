package ez

import (
	"fmt"
	"strconv"

	"github.com/vcokltfre/ez/ez/lexer"
	"github.com/vcokltfre/ez/ez/parser"
	"github.com/vcokltfre/ez/ez/vm"
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

	executor := vm.New(65536)

	executor.RegisterFunc("debug", -1, func(ctx lexer.TokenContext, args ...parser.Value) error {
		for _, arg := range args {
			var val int64
			if arg.Type == parser.ValueTypeInt {
				val, _ = strconv.ParseInt(arg.Value, 10, 64)
			} else {
				val = executor.Variables[arg.Value]
			}
			fmt.Printf("Debug: %s (%s): %d\n", arg.Value, arg.Type, val)
		}

		return nil
	})

	return executor.Run(program)
}

package vm

import (
	"fmt"
	"os"
	"os/exec"
	"strconv"

	"github.com/vcokltfre/ez/ez/lexer"
	"github.com/vcokltfre/ez/ez/parser"
)

type ExternalFunc struct {
	ArgCount    int
	ArgValidate bool
	Fn          func(lexer.TokenContext, ...parser.Value) error
}

type VM struct {
	Memory    []int64
	Variables map[string]int64
	Funcs     map[string]ExternalFunc

	program *parser.Program
	jumps   map[string]int
	index   int
}

func (vm *VM) setValue(stmt parser.VarDeclValue) error {
	if stmt.Value.Type == parser.ValueTypeInt {
		vm.Variables[stmt.Name], _ = strconv.ParseInt(stmt.Value.Value, 10, 64)
		return nil
	}

	val, ok := vm.Variables[stmt.Value.Value]
	if !ok {
		return stmt.Value.Token.Context.Error("runtime", "variable does not exist")
	}

	vm.Variables[stmt.Name] = val

	return nil
}

func (vm *VM) setValueFromOp(stmt parser.VarDeclExpr) error {
	var lhs int64
	if stmt.Expr.Lhs.Type == parser.ValueTypeInt {
		lhs, _ = strconv.ParseInt(stmt.Expr.Lhs.Value, 10, 64)
	} else {
		val, ok := vm.Variables[stmt.Expr.Lhs.Value]
		if !ok {
			return stmt.Expr.Lhs.Token.Context.Error("runtime", "variable does not exist")
		}

		lhs = val
	}

	var rhs int64
	if stmt.Expr.Rhs.Type == parser.ValueTypeInt {
		rhs, _ = strconv.ParseInt(stmt.Expr.Rhs.Value, 10, 64)
	} else {
		val, ok := vm.Variables[stmt.Expr.Rhs.Value]
		if !ok {
			return stmt.Expr.Rhs.Token.Context.Error("runtime", "variable does not exist")
		}

		rhs = val
	}

	switch stmt.Expr.Op {
	case "+":
		vm.Variables[stmt.Name] = lhs + rhs
	case "-":
		vm.Variables[stmt.Name] = lhs - rhs
	case "*":
		vm.Variables[stmt.Name] = lhs * rhs
	case "/":
		vm.Variables[stmt.Name] = lhs / rhs
	case "%":
		vm.Variables[stmt.Name] = lhs % rhs
	default:
		panic("invalid operator: " + stmt.Expr.Op)
	}

	return nil
}

func (vm *VM) ifStmt(stmt parser.If) error {
	var lhs int64
	if stmt.Cond.Lhs.Type == parser.ValueTypeInt {
		lhs, _ = strconv.ParseInt(stmt.Cond.Lhs.Value, 10, 64)
	} else {
		val, ok := vm.Variables[stmt.Cond.Lhs.Value]
		if !ok {
			return stmt.Cond.Lhs.Token.Context.Error("runtime", "variable does not exist")
		}

		lhs = val
	}

	var rhs int64
	if stmt.Cond.Rhs.Type == parser.ValueTypeInt {
		rhs, _ = strconv.ParseInt(stmt.Cond.Rhs.Value, 10, 64)
	} else {
		val, ok := vm.Variables[stmt.Cond.Rhs.Value]
		if !ok {
			return stmt.Cond.Rhs.Token.Context.Error("runtime", "variable does not exist")
		}

		rhs = val
	}

	var result bool
	switch stmt.Cond.Op {
	case "==":
		result = lhs == rhs
	case "!=":
		result = lhs != rhs
	case ">":
		result = lhs > rhs
	case "<":
		result = lhs < rhs
	case ">=":
		result = lhs >= rhs
	case "<=":
		result = lhs <= rhs
	default:
		panic("invalid operator: " + stmt.Cond.Op)
	}

	if result {
		if _, ok := vm.jumps[stmt.Goto.Name]; !ok {
			return stmt.Goto.Token.Context.Error("runtime", "label does not exist")
		}

		vm.index = vm.jumps[stmt.Goto.Name]
	}

	return nil
}

func (vm *VM) goTo(stmt parser.Goto) error {
	if _, ok := vm.jumps[stmt.Name]; !ok {
		return stmt.Token.Context.Error("runtime", "label does not exist")
	}

	vm.index = vm.jumps[stmt.Name]

	return nil
}

func (vm *VM) call(stmt parser.Call) error {
	callFn, ok := vm.Funcs[stmt.Name]
	if !ok {
		return stmt.Token.Context.Error("runtime", "function does not exist")
	}

	if len(stmt.Values) != callFn.ArgCount && callFn.ArgCount != -1 {
		return stmt.Token.Context.Error("runtime", "incorrect number of arguments")
	}

	if callFn.ArgValidate {
		for _, val := range stmt.Values {
			if val.Type == parser.ValueTypeVar {
				if _, ok := vm.Variables[val.Value]; !ok {
					return val.Token.Context.Error("runtime", "variable does not exist")
				}
			}
		}
	}

	return callFn.Fn(stmt.Token.Context, stmt.Values...)
}

func (vm *VM) RegisterFunc(name string, argCount int, argValidate bool, fn func(lexer.TokenContext, ...parser.Value) error) {
	vm.Funcs[name] = ExternalFunc{
		ArgCount:    argCount,
		ArgValidate: argValidate,
		Fn:          fn,
	}
}

func (vm *VM) Run(program *parser.Program) error {
	vm.program = program

	for i, stmt := range vm.program.Stmts {
		if stmt.Type() == parser.StmtTypeLabel {
			vm.jumps[stmt.(parser.Label).Name] = i
		}
	}

	for vm.index < len(vm.program.Stmts) {
		stmt := vm.program.Stmts[vm.index]

		switch stmt.Type() {
		case parser.StmtTypeVarDeclValue:
			err := vm.setValue(stmt.(parser.VarDeclValue))
			if err != nil {
				return err
			}
		case parser.StmtTypeVarDeclExpr:
			err := vm.setValueFromOp(stmt.(parser.VarDeclExpr))
			if err != nil {
				return err
			}
		case parser.StmtTypeIf:
			err := vm.ifStmt(stmt.(parser.If))
			if err != nil {
				return err
			}
		case parser.StmtTypeLabel:
		case parser.StmtTypeGoto:
			err := vm.goTo(stmt.(parser.Goto))
			if err != nil {
				return err
			}
		case parser.StmtTypeCall:
			err := vm.call(stmt.(parser.Call))
			if err != nil {
				return err
			}
		}

		vm.index++
	}

	return nil
}

func denyStringType(vals ...parser.Value) error {
	for _, val := range vals {
		if val.Type == parser.ValueTypeStr {
			return val.Token.Context.Error("runtime", "expected identifier or literal int not literal str")
		}
	}

	return nil
}

func New(memsize int) *VM {
	vm := &VM{
		Memory:    make([]int64, memsize),
		Variables: make(map[string]int64),
		Funcs:     make(map[string]ExternalFunc),

		jumps: make(map[string]int),
	}

	// call showc <var>
	vm.RegisterFunc("showc", 1, true, func(ctx lexer.TokenContext, args ...parser.Value) error {
		if err := denyStringType(args...); err != nil {
			return err
		}

		var val int64
		if args[0].Type == parser.ValueTypeInt {
			val, _ = strconv.ParseInt(args[0].Value, 10, 64)
		} else {
			var ok bool
			val, ok = vm.Variables[args[0].Value]
			if !ok {
				return args[0].Token.Context.Error("runtime", "variable does not exist")
			}
		}

		fmt.Printf("%c", byte(val))

		return nil
	})

	// call shown <var>
	vm.RegisterFunc("shown", 1, true, func(ctx lexer.TokenContext, args ...parser.Value) error {
		if err := denyStringType(args...); err != nil {
			return err
		}

		var val int64
		if args[0].Type == parser.ValueTypeInt {
			val, _ = strconv.ParseInt(args[0].Value, 10, 64)
		} else {
			var ok bool
			val, ok = vm.Variables[args[0].Value]
			if !ok {
				return args[0].Token.Context.Error("runtime", "variable does not exist")
			}
		}

		fmt.Printf("%d", val)

		return nil
	})

	// call input <var>
	vm.RegisterFunc("input", 1, false, func(ctx lexer.TokenContext, args ...parser.Value) error {
		if err := denyStringType(args...); err != nil {
			return err
		}

		if args[0].Type != parser.ValueTypeVar {
			return args[0].Token.Context.Error("runtime", "expected identifier not literal")
		}

		char := make([]byte, 1)
		_, err := os.Stdin.Read(char)
		if err != nil {
			panic(err)
		}

		vm.Variables[args[0].Value] = int64(char[0])

		return nil
	})

	// call memset <addr> <value>
	vm.RegisterFunc("memset", 2, true, func(ctx lexer.TokenContext, args ...parser.Value) error {
		if err := denyStringType(args...); err != nil {
			return err
		}

		var addr int64
		if args[0].Type == parser.ValueTypeInt {
			addr, _ = strconv.ParseInt(args[0].Value, 10, 64)
		} else {
			addr = vm.Variables[args[0].Value]
		}

		if addr < 0 || addr >= int64(len(vm.Memory)) {
			return ctx.Error("runtime", "invalid memory address")
		}

		var val int64
		if args[1].Type == parser.ValueTypeInt {
			val, _ = strconv.ParseInt(args[1].Value, 10, 64)
		} else {
			val = vm.Variables[args[1].Value]
		}

		vm.Memory[addr] = val

		return nil
	})

	// call memget <addr> <var>
	vm.RegisterFunc("memget", 2, false, func(ctx lexer.TokenContext, args ...parser.Value) error {
		if err := denyStringType(args...); err != nil {
			return err
		}

		var addr int64
		if args[0].Type == parser.ValueTypeInt {
			addr, _ = strconv.ParseInt(args[0].Value, 10, 64)
		} else {
			var ok bool
			addr, ok = vm.Variables[args[0].Value]
			if !ok {
				return args[0].Token.Context.Error("runtime", "variable does not exist")
			}
		}

		if addr < 0 || addr >= int64(len(vm.Memory)) {
			return ctx.Error("runtime", "invalid memory address")
		}

		if args[1].Type != parser.ValueTypeVar {
			return args[1].Token.Context.Error("runtime", "expected identifier not literal")
		}

		vm.Variables[args[1].Value] = vm.Memory[addr]

		return nil
	})

	// call debug ...vars
	vm.RegisterFunc("debug", -1, true, func(ctx lexer.TokenContext, args ...parser.Value) error {
		for _, arg := range args {
			if arg.Type == parser.ValueTypeStr {
				fmt.Printf("Debug: %s (str): %s\n", arg.Value, arg.Value)
				continue
			}

			var val int64
			if arg.Type == parser.ValueTypeInt {
				val, _ = strconv.ParseInt(arg.Value, 10, 64)
			} else {
				val = vm.Variables[arg.Value]
			}
			fmt.Printf("Debug: %s (%s): %d\n", arg.Value, arg.Type, val)
		}

		return nil
	})

	// call vm_no_input_buffering
	vm.RegisterFunc("vm_no_input_buffering", 0, false, func(ctx lexer.TokenContext, args ...parser.Value) error {
		return exec.Command("stty", "-F", "/dev/tty", "cbreak", "min", "1").Run()
	})

	// call read_file <filename> <addr> <length_var>
	vm.RegisterFunc("read_file", 3, false, func(ctx lexer.TokenContext, args ...parser.Value) error {
		if err := denyStringType(args[1:]...); err != nil {
			return err
		}

		file, addr, length := args[0], args[1], args[2]

		if file.Type != parser.ValueTypeStr {
			return file.Token.Context.Error("runtime", "expected string literal")
		}

		if length.Type != parser.ValueTypeVar {
			return length.Token.Context.Error("runtime", "expected identifier")
		}

		var address int64
		if addr.Type == parser.ValueTypeInt {
			address, _ = strconv.ParseInt(addr.Value, 10, 64)
		} else {
			var ok bool
			address, ok = vm.Variables[addr.Value]
			if !ok {
				return addr.Token.Context.Error("runtime", "variable does not exist")
			}
		}

		if address < 0 || address >= int64(len(vm.Memory)) {
			return ctx.Error("runtime", "invalid memory address")
		}

		data, err := os.ReadFile(file.Value)
		if err != nil {
			return file.Token.Context.Error("runtime", err.Error())
		}

		if address+int64(len(data)) >= int64(len(vm.Memory)) {
			return ctx.Error("runtime", "file too large")
		}

		for i, b := range data {
			vm.Memory[address+int64(i)] = int64(b)
		}

		vm.Variables[length.Value] = int64(len(data))

		return nil
	})

	// call write_file <filename> <addr> <length>
	vm.RegisterFunc("write_file", 3, false, func(ctx lexer.TokenContext, args ...parser.Value) error {
		if err := denyStringType(args[1:]...); err != nil {
			return err
		}

		file, addr, length := args[0], args[1], args[2]

		if file.Type != parser.ValueTypeStr {
			return file.Token.Context.Error("runtime", "expected string literal")
		}

		var address int64
		if addr.Type == parser.ValueTypeInt {
			address, _ = strconv.ParseInt(addr.Value, 10, 64)
		} else {
			var ok bool
			address, ok = vm.Variables[addr.Value]
			if !ok {
				return addr.Token.Context.Error("runtime", "variable does not exist")
			}
		}

		if address < 0 || address >= int64(len(vm.Memory)) {
			return ctx.Error("runtime", "invalid memory address")
		}

		var flen int64
		if length.Type == parser.ValueTypeInt {
			flen, _ = strconv.ParseInt(length.Value, 10, 64)
		} else {
			var ok bool
			flen, ok = vm.Variables[length.Value]
			if !ok {
				return length.Token.Context.Error("runtime", "variable does not exist")
			}
		}

		if address+flen >= int64(len(vm.Memory)) {
			return ctx.Error("runtime", fmt.Sprintf("memory out of bounds (%d)", address+flen))
		}

		data := make([]byte, flen)
		for i := int64(0); i < flen; i++ {
			data[i] = byte(vm.Memory[address+i])
		}

		err := os.WriteFile(file.Value, data, 0644)
		if err != nil {
			return file.Token.Context.Error("runtime", err.Error())
		}

		return nil
	})

	vm.Variables["__memsize"] = int64(memsize)

	return vm
}

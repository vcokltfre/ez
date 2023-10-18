package vm

import (
	"fmt"
	"os"
	"strconv"

	"github.com/vcokltfre/ez/ez/parser"
)

type VM struct {
	Memory    []int64
	Variables map[string]int64

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
		vm.index = vm.jumps[stmt.Goto]
	}

	return nil
}

func (vm *VM) showNum(stmt parser.ShowNum) error {
	var val int64
	if stmt.Value.Type == parser.ValueTypeInt {
		val, _ = strconv.ParseInt(stmt.Value.Value, 10, 64)
	} else {
		var ok bool
		val, ok = vm.Variables[stmt.Value.Value]
		if !ok {
			return stmt.Value.Token.Context.Error("runtime", "variable does not exist")
		}
	}

	fmt.Println(val)

	return nil
}

func (vm *VM) showChar(stmt parser.ShowChar) error {
	var val int64
	if stmt.Value.Type == parser.ValueTypeInt {
		val, _ = strconv.ParseInt(stmt.Value.Value, 10, 64)
	} else {
		var ok bool
		val, ok = vm.Variables[stmt.Value.Value]
		if !ok {
			return stmt.Value.Token.Context.Error("runtime", "variable does not exist")
		}
	}

	fmt.Printf("%c", val)

	return nil
}

func (vm *VM) input(stmt parser.Input) error {
	char := make([]byte, 1)
	_, err := os.Stdin.Read(char)
	if err != nil {
		panic(err)
	}

	vm.Variables[stmt.Name] = int64(char[0])

	return nil
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
		case parser.StmtTypeShownum:
			err := vm.showNum(stmt.(parser.ShowNum))
			if err != nil {
				return err
			}
		case parser.StmtTypeShowchar:
			err := vm.showChar(stmt.(parser.ShowChar))
			if err != nil {
				return err
			}
		case parser.StmtTypeInput:
			err := vm.input(stmt.(parser.Input))
			if err != nil {
				return err
			}
		}

		vm.index++
	}

	return nil
}

func New(memsize int) *VM {
	return &VM{
		Memory:    make([]int64, memsize),
		Variables: make(map[string]int64),

		jumps: make(map[string]int),
	}
}

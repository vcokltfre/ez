package lexer

import (
	"bytes"
	"errors"
	"fmt"
	"os"
	"strings"
)

type TokenType string

const (
	TTKeywordIf       TokenType = "if"
	TTKeywordEnd      TokenType = "end"
	TTKeywordGoto     TokenType = "goto"
	TTKeywordShownum  TokenType = "shownum"
	TTKeywordShowchar TokenType = "showchar"
	TTKeywordInput    TokenType = "input"

	TTIdentifier TokenType = "identifier"
	TTLiteralInt TokenType = "literal_int"
	TTLiteralStr TokenType = "literal_str"

	TTOpAssign TokenType = "op_assign"

	TTOpLt  TokenType = "op_lt"
	TTOpGt  TokenType = "op_gt"
	TTOpLte TokenType = "op_lte"
	TTOpGte TokenType = "op_gte"
	TTOpEq  TokenType = "op_eq"
	TTOpNeq TokenType = "op_neq"

	TTOpAdd TokenType = "op_add"
	TTOpSub TokenType = "op_sub"
	TTOpMul TokenType = "op_mul"
	TTOpDiv TokenType = "op_div"
	TTOpMod TokenType = "op_mod"
	TTOpPow TokenType = "op_pow"

	TTLabel TokenType = "label"

	TTEndStmt TokenType = "end_stmt"
)

var Keywords = map[string]TokenType{
	"if":    TTKeywordIf,
	"end":   TTKeywordEnd,
	"goto":  TTKeywordGoto,
	"shown": TTKeywordShownum,
	"showc": TTKeywordShowchar,
	"input": TTKeywordInput,
}

func IsKeyword(word string) bool {
	_, ok := Keywords[word]
	return ok
}

var Operators = map[string]TokenType{
	"=":  TTOpAssign,
	"<":  TTOpLt,
	">":  TTOpGt,
	"<=": TTOpLte,
	">=": TTOpGte,
	"==": TTOpEq,
	"!=": TTOpNeq,
	"+":  TTOpAdd,
	"-":  TTOpSub,
	"*":  TTOpMul,
	"/":  TTOpDiv,
	"%":  TTOpMod,
	"^":  TTOpPow,
}

func IsOperator(op string) bool {
	_, ok := Operators[op]
	return ok
}

type TokenContext struct {
	Line   int
	Column int
	Index  int
	File   string
}

func (ctx *TokenContext) Error(step, message string, tip ...string) error {
	fd, err := os.ReadFile(ctx.File)
	if err != nil {
		panic(err)
	}

	lines := bytes.Split(fd, []byte("\n"))
	line := lines[ctx.Line-1]

	padding := strings.Repeat(" ", ctx.Column-1)

	tipData := ""
	if len(tip) > 0 {
		tipData = fmt.Sprintf("\n|  \033[92m? %s", tip[0])
	}

	data := fmt.Sprintf(
		"| An error occurred during %s:\n|  File %s, line %d col %d\n|\n|  \033[96m%s\033[39m\n|  %s\033[93m^ %s\033[39m%s",
		step,
		ctx.File,
		ctx.Line,
		ctx.Column,
		string(line),
		padding,
		message,
		tipData,
	)

	return errors.New(data)
}

type Token struct {
	Type    TokenType
	Length  int
	Data    string
	Context TokenContext
}

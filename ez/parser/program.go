package parser

import "github.com/vcokltfre/ez/ez/lexer"

type StmtType string

const (
	StmtTypeVarDeclValue StmtType = "var_decl_value"
	StmtTypeVarDeclExpr  StmtType = "var_decl_expr"
	StmtTypeIf           StmtType = "if"
	StmtTypeLabel        StmtType = "label"
	StmtTypeGoto         StmtType = "goto"
	StmtTypeCall         StmtType = "call"
)

type ValueType string

const (
	ValueTypeInt ValueType = "int"
	ValueTypeVar ValueType = "var"
)

func valueTypeFromToken(t lexer.TokenType) ValueType {
	switch t {
	case lexer.TTLiteralInt:
		return ValueTypeInt
	case lexer.TTIdentifier:
		return ValueTypeVar
	default:
		return ""
	}
}

type Value struct {
	Type  ValueType
	Value string
	Token lexer.Token
}

type VarDeclValue struct {
	Name  string
	Value Value
}

func (v VarDeclValue) Type() StmtType {
	return StmtTypeVarDeclValue
}

type OpExpr struct {
	Op  string
	Lhs Value
	Rhs Value
}

type VarDeclExpr struct {
	Name string
	Expr OpExpr
}

func (v VarDeclExpr) Type() StmtType {
	return StmtTypeVarDeclExpr
}

type If struct {
	Cond OpExpr
	Goto Goto
}

func (i If) Type() StmtType {
	return StmtTypeIf
}

type Label struct {
	Name string
}

func (l Label) Type() StmtType {
	return StmtTypeLabel
}

type Goto struct {
	Name  string
	Token lexer.Token
}

func (g Goto) Type() StmtType {
	return StmtTypeGoto
}

type Call struct {
	Name   string
	Values []Value
	Token  lexer.Token
}

func (c Call) Type() StmtType {
	return StmtTypeCall
}

type Stmt interface {
	Type() StmtType
}

type Program struct {
	Stmts []Stmt
}

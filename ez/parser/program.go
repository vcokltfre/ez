package parser

import "github.com/vcokltfre/ez/ez/lexer"

type StmtType string

const (
	StmtTypeVarDeclValue StmtType = "var_decl_value"
	StmtTypeVarDeclExpr  StmtType = "var_decl_expr"
	StmtTypeIf           StmtType = "if"
	StmtTypeLabel        StmtType = "label"
	StmtTypeShownum      StmtType = "shownum"
	StmtTypeShowchar     StmtType = "showchar"
	StmtTypeInput        StmtType = "input"
	StmtTypeGoto         StmtType = "goto"
)

type ValueType string

const (
	ValueTypeInt ValueType = "int"
	ValueTypeStr ValueType = "str"
	ValueTypeVar ValueType = "var"
)

func valueTypeFromToken(t lexer.TokenType) ValueType {
	switch t {
	case lexer.TTLiteralInt:
		return ValueTypeInt
	case lexer.TTLiteralStr:
		return ValueTypeStr
	case lexer.TTIdentifier:
		return ValueTypeVar
	default:
		return ""
	}
}

type Value struct {
	Type  ValueType
	Value string
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
	Goto string
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

type ShowNum struct {
	Value Value
}

func (s ShowNum) Type() StmtType {
	return StmtTypeShownum
}

type ShowChar struct {
	Value Value
}

func (s ShowChar) Type() StmtType {
	return StmtTypeShowchar
}

type Input struct {
	Name string
}

func (i Input) Type() StmtType {
	return StmtTypeInput
}

type Goto struct {
	Name string
}

func (g Goto) Type() StmtType {
	return StmtTypeGoto
}

type Stmt interface {
	Type() StmtType
}

type Program struct {
	Stmts []Stmt
}

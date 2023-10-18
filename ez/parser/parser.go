package parser

import (
	"github.com/vcokltfre/ez/ez/lexer"
)

const STEP = "parsing"

func matchTokenPattern(pattern ...[]lexer.TokenType) func([]lexer.Token) bool {
	return func(tokens []lexer.Token) bool {
		if len(tokens) < len(pattern) {
			return false
		}

		for i, p := range pattern {
			anyValid := false

			for _, t := range p {
				if tokens[i].Type == t {
					anyValid = true
					break
				}
			}

			if !anyValid {
				return false
			}
		}

		return true
	}
}

var (
	matchVarDeclValue = matchTokenPattern(
		[]lexer.TokenType{lexer.TTIdentifier},
		[]lexer.TokenType{lexer.TTOpAssign},
		[]lexer.TokenType{lexer.TTLiteralInt, lexer.TTLiteralStr, lexer.TTIdentifier},
		[]lexer.TokenType{lexer.TTEndStmt},
	)
	matchVarDeclExpr = matchTokenPattern(
		[]lexer.TokenType{lexer.TTIdentifier},
		[]lexer.TokenType{lexer.TTOpAssign},
		[]lexer.TokenType{lexer.TTLiteralInt, lexer.TTLiteralStr, lexer.TTIdentifier},
		[]lexer.TokenType{lexer.TTOpAdd, lexer.TTOpSub, lexer.TTOpMul, lexer.TTOpDiv, lexer.TTOpMod, lexer.TTOpPow},
		[]lexer.TokenType{lexer.TTLiteralInt, lexer.TTLiteralStr, lexer.TTIdentifier},
		[]lexer.TokenType{lexer.TTEndStmt},
	)
	matchIf = matchTokenPattern(
		[]lexer.TokenType{lexer.TTKeywordIf},
		[]lexer.TokenType{lexer.TTLiteralInt, lexer.TTLiteralStr, lexer.TTIdentifier},
		[]lexer.TokenType{lexer.TTOpLt, lexer.TTOpGt, lexer.TTOpLte, lexer.TTOpGte, lexer.TTOpEq, lexer.TTOpNeq},
		[]lexer.TokenType{lexer.TTLiteralInt, lexer.TTLiteralStr, lexer.TTIdentifier},
		[]lexer.TokenType{lexer.TTKeywordGoto},
		[]lexer.TokenType{lexer.TTIdentifier},
		[]lexer.TokenType{lexer.TTEndStmt},
	)
	matchLabel = matchTokenPattern(
		[]lexer.TokenType{lexer.TTLabel},
		[]lexer.TokenType{lexer.TTEndStmt},
	)
	matchShownum = matchTokenPattern(
		[]lexer.TokenType{lexer.TTKeywordShownum},
		[]lexer.TokenType{lexer.TTLiteralInt, lexer.TTIdentifier},
		[]lexer.TokenType{lexer.TTEndStmt},
	)
	matchShowchar = matchTokenPattern(
		[]lexer.TokenType{lexer.TTKeywordShowchar},
		[]lexer.TokenType{lexer.TTLiteralInt, lexer.TTIdentifier},
		[]lexer.TokenType{lexer.TTEndStmt},
	)
	matchInput = matchTokenPattern(
		[]lexer.TokenType{lexer.TTKeywordInput},
		[]lexer.TokenType{lexer.TTIdentifier},
		[]lexer.TokenType{lexer.TTEndStmt},
	)
	matchGoto = matchTokenPattern(
		[]lexer.TokenType{lexer.TTKeywordGoto},
		[]lexer.TokenType{lexer.TTIdentifier},
		[]lexer.TokenType{lexer.TTEndStmt},
	)
)

func parseVarDeclValue(tokens []lexer.Token) VarDeclValue {
	assignTo := tokens[0]
	assignFrom := tokens[2]

	return VarDeclValue{
		Name: assignTo.Data,
		Value: Value{
			Type:  valueTypeFromToken(assignFrom.Type),
			Value: assignFrom.Data,
		},
	}
}

func parseVarDeclExpr(tokens []lexer.Token) VarDeclExpr {
	assignTo := tokens[0]
	lhsValue := tokens[2]
	op := tokens[3]
	rhsValue := tokens[4]

	return VarDeclExpr{
		Name: assignTo.Data,
		Expr: OpExpr{
			Op: op.Data,
			Lhs: Value{
				Type:  valueTypeFromToken(lhsValue.Type),
				Value: lhsValue.Data,
			},
			Rhs: Value{
				Type:  valueTypeFromToken(rhsValue.Type),
				Value: rhsValue.Data,
			},
		},
	}
}

func parseIf(tokens []lexer.Token) If {
	condLhs := tokens[1]
	condOp := tokens[2]
	condRhs := tokens[3]
	gotoLabel := tokens[5]

	return If{
		Cond: OpExpr{
			Op: condOp.Data,
			Lhs: Value{
				Type:  valueTypeFromToken(condLhs.Type),
				Value: condLhs.Data,
			},
			Rhs: Value{
				Type:  valueTypeFromToken(condRhs.Type),
				Value: condRhs.Data,
			},
		},
		Goto: gotoLabel.Data,
	}
}

func parseLabel(tokens []lexer.Token) Label {
	label := tokens[0]

	return Label{
		Name: label.Data,
	}
}

func parseShownum(tokens []lexer.Token) ShowNum {
	value := tokens[1]

	return ShowNum{
		Value: Value{
			Type:  valueTypeFromToken(value.Type),
			Value: value.Data,
		},
	}
}

func parseShowchar(tokens []lexer.Token) ShowChar {
	value := tokens[1]

	return ShowChar{
		Value: Value{
			Type:  valueTypeFromToken(value.Type),
			Value: value.Data,
		},
	}
}

func parseInput(tokens []lexer.Token) Input {
	variable := tokens[1]

	return Input{
		Name: variable.Data,
	}
}

func parseGoto(tokens []lexer.Token) Goto {
	label := tokens[1]

	return Goto{
		Name: label.Data,
	}
}

func Parse(tokens []lexer.Token) (*Program, error) {
	program := &Program{}
	index := 0

	for index < len(tokens) {
		if matchVarDeclValue(tokens[index:]) {
			program.Stmts = append(program.Stmts, parseVarDeclValue(tokens))
			index += 4
			continue
		} else if matchVarDeclExpr(tokens[index:]) {
			program.Stmts = append(program.Stmts, parseVarDeclExpr(tokens))
			index += 6
			continue
		} else if matchIf(tokens[index:]) {
			program.Stmts = append(program.Stmts, parseIf(tokens))
			index += 7
			continue
		} else if matchLabel(tokens[index:]) {
			program.Stmts = append(program.Stmts, parseLabel(tokens))
			index += 2
			continue
		} else if matchShownum(tokens[index:]) {
			program.Stmts = append(program.Stmts, parseShownum(tokens))
			index += 3
			continue
		} else if matchShowchar(tokens[index:]) {
			program.Stmts = append(program.Stmts, parseShowchar(tokens))
			index += 3
			continue
		} else if matchInput(tokens[index:]) {
			program.Stmts = append(program.Stmts, parseInput(tokens))
			index += 3
			continue
		} else if matchGoto(tokens[index:]) {
			program.Stmts = append(program.Stmts, parseGoto(tokens))
			index += 3
			continue
		} else if tokens[index].Type == lexer.TTEndStmt {
			index++
			continue
		}

		return nil, tokens[index].Context.Error(STEP, "Invalid statement")
	}

	return program, nil
}

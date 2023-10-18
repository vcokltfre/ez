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
		[]lexer.TokenType{lexer.TTLiteralInt, lexer.TTIdentifier},
		[]lexer.TokenType{lexer.TTEndStmt},
	)
	matchVarDeclExpr = matchTokenPattern(
		[]lexer.TokenType{lexer.TTIdentifier},
		[]lexer.TokenType{lexer.TTOpAssign},
		[]lexer.TokenType{lexer.TTLiteralInt, lexer.TTIdentifier},
		[]lexer.TokenType{lexer.TTOpAdd, lexer.TTOpSub, lexer.TTOpMul, lexer.TTOpDiv, lexer.TTOpMod, lexer.TTOpPow},
		[]lexer.TokenType{lexer.TTLiteralInt, lexer.TTIdentifier},
		[]lexer.TokenType{lexer.TTEndStmt},
	)
	matchIf = matchTokenPattern(
		[]lexer.TokenType{lexer.TTKeywordIf},
		[]lexer.TokenType{lexer.TTLiteralInt, lexer.TTIdentifier},
		[]lexer.TokenType{lexer.TTOpLt, lexer.TTOpGt, lexer.TTOpLte, lexer.TTOpGte, lexer.TTOpEq, lexer.TTOpNeq},
		[]lexer.TokenType{lexer.TTLiteralInt, lexer.TTIdentifier},
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

func matchCall(tokens []lexer.Token) bool {
	if tokens[0].Type != lexer.TTKeywordCall {
		return false
	}

	if tokens[1].Type != lexer.TTIdentifier {
		return false
	}

	index := 2
	for index < len(tokens) && (tokens[index].Type == lexer.TTLiteralInt || tokens[index].Type == lexer.TTIdentifier) {
		index++
	}

	if index >= len(tokens) || tokens[index].Type != lexer.TTEndStmt {
		return false
	}

	return true
}

func parseVarDeclValue(tokens []lexer.Token) VarDeclValue {
	assignTo := tokens[0]
	assignFrom := tokens[2]

	return VarDeclValue{
		Name: assignTo.Data,
		Value: Value{
			Type:  valueTypeFromToken(assignFrom.Type),
			Value: assignFrom.Data,
			Token: assignFrom,
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
				Token: lhsValue,
			},
			Rhs: Value{
				Type:  valueTypeFromToken(rhsValue.Type),
				Value: rhsValue.Data,
				Token: rhsValue,
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
				Token: condLhs,
			},
			Rhs: Value{
				Type:  valueTypeFromToken(condRhs.Type),
				Value: condRhs.Data,
				Token: condRhs,
			},
		},
		Goto: Goto{
			Name:  gotoLabel.Data,
			Token: gotoLabel,
		},
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
			Token: value,
		},
	}
}

func parseShowchar(tokens []lexer.Token) ShowChar {
	value := tokens[1]

	return ShowChar{
		Value: Value{
			Type:  valueTypeFromToken(value.Type),
			Value: value.Data,
			Token: value,
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
		Name:  label.Data,
		Token: label,
	}
}

func parseCall(tokens []lexer.Token) Call {
	name := tokens[1]

	values := []Value{}

	index := 2
	for tokens[index].Type != lexer.TTEndStmt {
		values = append(values, Value{
			Type:  valueTypeFromToken(tokens[index].Type),
			Value: tokens[index].Data,
			Token: tokens[index],
		})
		index++
	}

	return Call{
		Name:   name.Data,
		Values: values,
		Token:  name,
	}
}

func Parse(tokens []lexer.Token) (*Program, error) {
	program := &Program{}
	index := 0

	for index < len(tokens) {
		if matchVarDeclValue(tokens[index:]) {
			program.Stmts = append(program.Stmts, parseVarDeclValue(tokens[index:]))
			index += 4
			continue
		} else if matchVarDeclExpr(tokens[index:]) {
			program.Stmts = append(program.Stmts, parseVarDeclExpr(tokens[index:]))
			index += 6
			continue
		} else if matchIf(tokens[index:]) {
			program.Stmts = append(program.Stmts, parseIf(tokens[index:]))
			index += 7
			continue
		} else if matchLabel(tokens[index:]) {
			program.Stmts = append(program.Stmts, parseLabel(tokens[index:]))
			index += 2
			continue
		} else if matchShownum(tokens[index:]) {
			program.Stmts = append(program.Stmts, parseShownum(tokens[index:]))
			index += 3
			continue
		} else if matchShowchar(tokens[index:]) {
			program.Stmts = append(program.Stmts, parseShowchar(tokens[index:]))
			index += 3
			continue
		} else if matchInput(tokens[index:]) {
			program.Stmts = append(program.Stmts, parseInput(tokens[index:]))
			index += 3
			continue
		} else if matchGoto(tokens[index:]) {
			program.Stmts = append(program.Stmts, parseGoto(tokens[index:]))
			index += 3
			continue
		} else if matchCall(tokens[index:]) {
			stmt := parseCall(tokens[index:])
			program.Stmts = append(program.Stmts, stmt)
			index += 2 + len(stmt.Values)
			continue
		} else if tokens[index].Type == lexer.TTEndStmt {
			index++
			continue
		}

		return nil, tokens[index].Context.Error(STEP, "Invalid statement")
	}

	return program, nil
}

package lexer

import (
	"fmt"
	"strconv"
	"strings"
)

const STEP = "lexing"

var (
	matchDecimalInt = match(`^\d+\b`)
	matchHexInt     = match(`^0x[0-9a-fA-F]+\b`)
	matchIdentifier = match(`^[a-zA-Z_][a-zA-Z0-9_]*\b`)
	matchLabel      = match(`^:[a-zA-Z_][a-zA-Z0-9_]*\b`)
)

func skipWhitespace(code string) int {
	index := 0

	for index < len(code) {
		if code[index] == ' ' || code[index] == '\t' {
			index++
			continue
		}
		break
	}

	return index
}

func getIntLiteral(code string, ctx TokenContext) (*Token, error) {
	if match := matchHexInt(code); match != nil {
		val, _ := strconv.ParseInt(*match, 16, 64)
		decimal := fmt.Sprintf("%d", val)

		return &Token{
			Type:    TTLiteralInt,
			Length:  len(*match),
			Data:    decimal,
			Context: ctx,
		}, nil
	}

	if match := matchDecimalInt(code); match != nil {
		return &Token{
			Type:    TTLiteralInt,
			Length:  len(*match),
			Data:    *match,
			Context: ctx,
		}, nil
	}

	return nil, ctx.Error(STEP, "Invalid integer literal", "Integer literals must be in decimal or hexadecimal format")
}

func getIdentifier(code string, ctx TokenContext) (*Token, error) {
	if match := matchIdentifier(code); match != nil {
		if IsKeyword(*match) {
			return &Token{
				Type:    Keywords[*match],
				Length:  len(*match),
				Data:    *match,
				Context: ctx,
			}, nil
		}

		return &Token{
			Type:    TTIdentifier,
			Length:  len(*match),
			Data:    *match,
			Context: ctx,
		}, nil
	}

	return nil, ctx.Error(STEP, "Invalid identifier")
}

func getOperator(code string, ctx TokenContext) (*Token, error) {
	for _, op := range orderedOperators {
		if strings.HasPrefix(code, op) {
			return &Token{
				Type:    Operators[op],
				Length:  len(op),
				Data:    op,
				Context: ctx,
			}, nil
		}
	}

	return nil, ctx.Error(STEP, "Invalid operator")
}

func getLabel(code string, ctx TokenContext) (*Token, error) {
	if match := matchLabel(code); match != nil {
		return &Token{
			Type:    TTLabel,
			Length:  len(*match),
			Data:    (*match)[1:],
			Context: ctx,
		}, nil
	}

	return nil, ctx.Error(STEP, "Invalid label")
}

func Lex(code, filename string) ([]Token, error) {
	index := 0
	line := 1
	column := 1

	tokens := []Token{}

	for index < len(code) {
		current := code[index]

		context := TokenContext{
			Line:   line,
			Column: column,
			Index:  index,
			File:   filename,
		}

		if current == '\n' {
			tokens = append(tokens, Token{
				Type:    TTEndStmt,
				Context: context,
			})
			index++
			line++
			column = 1
			continue
		}

		whitespace := skipWhitespace(code[index:])
		if whitespace > 0 {
			index += whitespace
			column += whitespace
			continue
		}

		if '0' <= current && current <= '9' {
			token, err := getIntLiteral(code[index:], context)
			if err != nil {
				return nil, err
			}

			tokens = append(tokens, *token)
			index += token.Length
			column += token.Length
			continue
		}

		if 'a' <= current && current <= 'z' || 'A' <= current && current <= 'Z' || current == '_' {
			token, err := getIdentifier(code[index:], context)
			if err != nil {
				return nil, err
			}

			tokens = append(tokens, *token)
			index += token.Length
			column += token.Length
			continue
		}

		if isOperatorCharacter(string(current)) {
			token, err := getOperator(code[index:], context)
			if err != nil {
				return nil, err
			}

			tokens = append(tokens, *token)
			index += token.Length
			column += token.Length
			continue
		}

		if current == ':' {
			token, err := getLabel(code[index:], context)
			if err != nil {
				return nil, err
			}

			tokens = append(tokens, *token)
			index += token.Length
			column += token.Length
			continue
		}

		return nil, context.Error(STEP, fmt.Sprintf("Unexpected character: %s", string(current)))
	}

	return tokens, nil
}

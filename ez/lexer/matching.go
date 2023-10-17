package lexer

import "regexp"

type matcher func(string) *string

func match(pattern string) matcher {
	rePattern := regexp.MustCompile(pattern)

	return func(code string) *string {
		match := rePattern.FindString(code)

		if match == "" {
			return nil
		}

		return &match
	}
}

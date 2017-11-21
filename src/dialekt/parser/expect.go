package parser

import "fmt"

func expect(tok *Token, types ...TokenType) {
	if tok == nil {
		panic(fmt.Errorf("unexpected end of input, expected %s", formatTokenTypes(types)))
	}

	for _, t := range types {
		if tok.Type == t {
			return
		}
	}

	panic(fmt.Errorf(
		"unexpected %s, expected %s",
		tok.Type,
		formatTokenTypes(types),
	))
}

func formatTokenTypes(types []TokenType) (s string) {
	l := len(types)

	for i, t := range types {
		if i > 0 {
			if i+1 == l {
				s += " or "
			} else {
				s += ", "
			}
		}

		s += t.String()
	}

	return
}

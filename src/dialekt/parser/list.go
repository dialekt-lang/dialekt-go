package parser

import (
	"io"

	"github.com/dialekt-lang/dialekt-go/src/dialekt/ast"
)

// ParseList parses a whitespace separated sequence of tags.
func ParseList(r io.RuneReader) (tags []ast.Tag, err error) {
	defer func() {
		if r := recover(); r != nil {
			if e, ok := r.(error); ok {
				err = e
			} else {
				panic(r)
			}
		}
	}()

	l := newLexer(r)

	for {
		tok := l.next()
		if tok == nil {
			return
		}

		expect(tok, TagToken)
		tags = append(tags, ast.Tag{
			Name: tok.Values[0],
		})
	}
}

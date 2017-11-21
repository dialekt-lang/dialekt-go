package parser

import (
	"bufio"
	"bytes"
	"errors"
	"io"
	"strings"
	"unicode"
)

// Lex tokenises input from r and writes each token to c.
func Lex(r io.Reader, c chan<- Token) error {
	l := &lexer{
		r:    bufio.NewReader(r),
		c:    c,
		line: 1,
		off:  -1,
	}

	return l.run()
}

type lexer struct {
	r         *bufio.Reader
	c         chan<- Token
	buf       bytes.Buffer
	cur, prev rune
	off       int
	line, col uint
	tok       Token
}

const wildcardRune = '*'

type state func() state

func (l *lexer) run() (err error) {
	defer func() {
		close(l.c)

		if r := recover(); r != nil {
			if e, ok := r.(error); ok {
				err = e
			} else {
				panic(r)
			}
		}
	}()

	s := l.begin
	for s != nil {
		s = s()
	}

	return
}

func (l *lexer) begin() state {
	for {
		if !l.advance() {
			return nil
		}

		if unicode.IsSpace(l.cur) {
			continue
		}

		switch l.cur {
		case '(':
			l.emitToken(OpenParenToken)

		case ')':
			l.emitToken(CloseParenToken)

		case '"':
			return l.quotedString

		default:
			return l.unquotedString
		}
	}
}

func (l *lexer) quotedString() state {
	l.startToken()

	for {
		l.mustAdvance()

		switch l.cur {
		case '"':
			l.captureValue()

			if len(l.tok.Values) == 1 {
				l.endToken(TagToken, 1)
			} else {
				l.endToken(PatternToken, 1)
			}

			return l.begin

		case wildcardRune:
			l.captureValue()

		case '\\':
			l.mustAdvance()
			l.buf.WriteRune(l.cur)

		default:
			l.buf.WriteRune(l.cur)
		}
	}
}

func (l *lexer) unquotedString() state {
	l.startToken()

	for {
		switch l.cur {
		case wildcardRune:
			l.captureValue()

		case '"':
			l.endUnquotedString()
			return l.quotedString

		case '(':
			l.endUnquotedString()
			l.emitToken(OpenParenToken)
			return l.begin

		case ')':
			l.endUnquotedString()
			l.emitToken(CloseParenToken)
			return l.begin

		default:
			l.buf.WriteRune(l.cur)
		}

		if !l.advance() {
			l.endUnquotedString()
			return nil
		}

		if unicode.IsSpace(l.cur) {
			l.endUnquotedString()
			return l.begin
		}
	}
}

func (l *lexer) endUnquotedString() {
	if len(l.tok.Values) == 0 {
		switch strings.ToLower(l.buf.String()) {
		case "and":
			l.endToken(LogicalAndToken, 0)
		case "or":
			l.endToken(LogicalOrToken, 0)
		case "not":
			l.endToken(LogicalNotToken, 0)
		default:
			l.captureValue()
			l.endToken(TagToken, 0)
		}
	} else {
		l.captureValue()
		l.endToken(PatternToken, 0)
	}
}

func (l *lexer) advance() bool {
	l.prev = l.cur
	r, size, err := l.r.ReadRune()

	if err == io.EOF {
		l.off++
		return false
	} else if err != nil {
		panic(err)
	} else if size == 1 && r == unicode.ReplacementChar {
		panic(errors.New("invalid UTF-8 rune"))
	}

	l.cur = r
	l.col++
	l.off += size

	if l.prev == '\n' || (l.prev == '\r' && l.cur != '\n') {
		l.line++
		l.col = 1
	}

	return true
}

func (l *lexer) mustAdvance() {
	if !l.advance() {
		panic(io.EOF)
	}
}

func (l *lexer) startToken() {
	l.tok = Token{
		StartOffset: uint(l.off),
		Line:        l.line,
		Column:      l.col,
	}
}

func (l *lexer) endToken(t TokenType, off int) {
	l.tok.Type = t
	l.tok.EndOffset = uint(l.off + off)

	l.c <- l.tok
}

func (l *lexer) emitToken(t TokenType) {
	l.startToken()
	l.endToken(t, 1)
}

func (l *lexer) captureValue() {
	l.tok.Values = append(l.tok.Values, l.buf.String())
	l.buf.Reset()
}

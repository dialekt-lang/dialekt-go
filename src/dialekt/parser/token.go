package parser

// TokenType is an enumeration of the various types of tokens.
type TokenType uint

const (
	// LogicalAndToken is a token for the "AND" keyword.
	LogicalAndToken = iota

	// LogicalOrToken is a token for the "OR" keyword.
	LogicalOrToken

	// LogicalNotToken is a token for the "NOT" keyword.
	LogicalNotToken

	// TagToken is a token for a literal tag string.
	TagToken

	// PatternToken is a token for a tag pattern.
	PatternToken

	// OpenParenToken is a token for the opening parenthesis.
	OpenParenToken

	// CloseParenToken is a token for the closing parenthesis.
	CloseParenToken
)

// Token is a structure for tokens produced by the lexer.
type Token struct {
	Type                   TokenType
	Values                 []string
	StartOffset, EndOffset uint
	Line, Column           uint
}

func (t Token) String() string {
	switch t.Type {
	case LogicalAndToken:
		return "AND operator"
	case LogicalOrToken:
		return "OR operator"
	case LogicalNotToken:
		return "NOT operator"
	case TagToken:
		return "tag"
	case PatternToken:
		return "pattern"
	case OpenParenToken:
		return "opening parenthesis"
	case CloseParenToken:
		return "closing parenthesis"
	}

	panic("unknown token type")
}

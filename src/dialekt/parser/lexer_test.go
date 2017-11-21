package parser

import (
	"bytes"

	. "github.com/onsi/ginkgo"
	. "github.com/onsi/ginkgo/extensions/table"
	. "github.com/onsi/gomega"
)

var _ = Describe("lexer", func() {
	DescribeTable(
		"it produces the expected token stream",
		func(input string, expected []Token) {
			buf := bytes.NewBufferString(input)

			l := newLexer(buf)
			var output []Token

			for {
				tok := l.next()
				if tok == nil {
					Expect(output).To(Equal(expected))
					return
				}

				output = append(output, *tok)
			}
		},

		Entry(
			"empty expression",
			"",
			nil,
		),
		Entry(
			"whitespace only",
			" \n \t ",
			nil,
		),
		Entry(
			"line counter tracks LF",
			"\"foo\nbar\" baz",
			[]Token{
				{TagToken, []string{"foo\nbar"}, 0, 9, 1, 1},
				{TagToken, []string{"baz"}, 10, 13, 2, 6},
			},
		),
		Entry(
			"line counter tracks CRLF",
			"\"foo\r\nbar\" baz",
			[]Token{
				{TagToken, []string{"foo\r\nbar"}, 0, 10, 1, 1},
				{TagToken, []string{"baz"}, 11, 14, 2, 6},
			},
		),
		Entry(
			"line counter tracks standalone CR",
			"\"foo\rbar\" baz",
			[]Token{
				{TagToken, []string{"foo\rbar"}, 0, 9, 1, 1},
				{TagToken, []string{"baz"}, 10, 13, 2, 6},
			},
		),

		Entry(
			"unquoted tag",
			`foo-bar`,
			[]Token{
				{TagToken, []string{"foo-bar"}, 0, 7, 1, 1},
			},
		),
		Entry(
			"unquoted tag with leading hyphen",
			`-foo`,
			[]Token{
				{TagToken, []string{"-foo"}, 0, 4, 1, 1},
			},
		),
		Entry(
			"adjacent unquoted tags",
			`foo bar`,
			[]Token{
				{TagToken, []string{`foo`}, 0, 3, 1, 1},
				{TagToken, []string{`bar`}, 4, 7, 1, 5},
			},
		),
		Entry(
			"whitespace around unquoted tags",
			" \t\nfoo\tbar\nbaz\t ",
			[]Token{
				{TagToken, []string{`foo`}, 3, 6, 2, 1},
				{TagToken, []string{`bar`}, 7, 10, 2, 5},
				{TagToken, []string{`baz`}, 11, 14, 3, 1},
			},
		),

		Entry(
			"unquoted pattern with leading wildcard",
			`*foo`,
			[]Token{
				{PatternToken, []string{"", "foo"}, 0, 4, 1, 1},
			},
		),
		Entry(
			"unquoted pattern with trailing wildcard",
			`foo*`,
			[]Token{
				{PatternToken, []string{"foo", ""}, 0, 4, 1, 1},
			},
		),
		Entry(
			"unquoted pattern with leading / trailing wildcard",
			`*foo*`,
			[]Token{
				{PatternToken, []string{"", "foo", ""}, 0, 5, 1, 1},
			},
		),
		Entry(
			"unquoted pattern with enclosed wildcard",
			`foo*bar`,
			[]Token{
				{PatternToken, []string{"foo", "bar"}, 0, 7, 1, 1},
			},
		),

		Entry(
			"quoted tag",
			`"foo bar"`,
			[]Token{
				{TagToken, []string{`foo bar`}, 0, 9, 1, 1},
			},
		),
		Entry(
			"quoted tag with escaped quote",
			`"foo \"the\" bar"`,
			[]Token{
				{TagToken, []string{`foo "the" bar`}, 0, 17, 1, 1},
			},
		),
		Entry(
			"quoted tag with escaped backslash",
			`"foo\\bar"`,
			[]Token{
				{TagToken, []string{`foo\bar`}, 0, 10, 1, 1},
			},
		),
		Entry(
			"quoted tag with escaped asterisk",
			`"foo\*bar"`,
			[]Token{
				{TagToken, []string{`foo*bar`}, 0, 10, 1, 1},
			},
		),
		Entry(
			"quoted tag with parents",
			`"foo(bar)baz"`,
			[]Token{
				{TagToken, []string{`foo(bar)baz`}, 0, 13, 1, 1},
			},
		),

		Entry(
			"adjacent quoted tags",
			`"foo""bar"`,
			[]Token{
				{TagToken, []string{`foo`}, 0, 5, 1, 1},
				{TagToken, []string{`bar`}, 5, 10, 1, 6},
			},
		),
		Entry(
			"unquoted tag followed by quoted tag",
			`foo"bar"`,
			[]Token{
				{TagToken, []string{`foo`}, 0, 3, 1, 1},
				{TagToken, []string{`bar`}, 3, 8, 1, 4},
			},
		),
		Entry(
			"quoted tag followed by unquoted tag",
			`"foo"bar`,
			[]Token{
				{TagToken, []string{`foo`}, 0, 5, 1, 1},
				{TagToken, []string{`bar`}, 5, 8, 1, 6},
			},
		),

		Entry(
			"quoted pattern with leading wildcard",
			`"*foo"`,
			[]Token{
				{PatternToken, []string{"", "foo"}, 0, 6, 1, 1},
			},
		),
		Entry(
			"quoted pattern with trailing wildcard",
			`"foo*"`,
			[]Token{
				{PatternToken, []string{"foo", ""}, 0, 6, 1, 1},
			},
		),
		Entry(
			"quoted pattern with leading / trailing wildcard",
			`"*foo*"`,
			[]Token{
				{PatternToken, []string{"", "foo", ""}, 0, 7, 1, 1},
			},
		),
		Entry(
			"quoted pattern with enclosed wildcard",
			`"foo*bar"`,
			[]Token{
				{PatternToken, []string{"foo", "bar"}, 0, 9, 1, 1},
			},
		),

		Entry(
			"logical and",
			`and`,
			[]Token{
				{LogicalAndToken, nil, 0, 3, 1, 1},
			},
		),
		Entry(
			"logical and with mixed case",
			`aNd`,
			[]Token{
				{LogicalAndToken, nil, 0, 3, 1, 1},
			},
		),
		Entry(
			"logical or",
			`or`,
			[]Token{
				{LogicalOrToken, nil, 0, 2, 1, 1},
			},
		),
		Entry(
			"logical or with mixed case",
			`oR`,
			[]Token{
				{LogicalOrToken, nil, 0, 2, 1, 1},
			},
		),
		Entry(
			"logical not",
			`not`,
			[]Token{
				{LogicalNotToken, nil, 0, 3, 1, 1},
			},
		),
		Entry(
			"logical not with mixed case",
			`nOt`,
			[]Token{
				{LogicalNotToken, nil, 0, 3, 1, 1},
			},
		),

		Entry(
			"open paren",
			`(`,
			[]Token{
				{OpenParenToken, nil, 0, 1, 1, 1},
			},
		),
		Entry(
			"close paren",
			`)`,
			[]Token{
				{CloseParenToken, nil, 0, 1, 1, 1},
			},
		),
		Entry(
			"parens interrupt tags",
			`foo(bar)baz`,
			[]Token{
				{TagToken, []string{"foo"}, 0, 3, 1, 1},
				{OpenParenToken, nil, 3, 4, 1, 4},
				{TagToken, []string{"bar"}, 4, 7, 1, 5},
				{CloseParenToken, nil, 7, 8, 1, 8},
				{TagToken, []string{"baz"}, 8, 11, 1, 9},
			},
		),
	)
})

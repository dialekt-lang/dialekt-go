package parser_test

import (
	"bytes"

	"github.com/dialekt-lang/dialekt-go/src/dialekt/ast"
	. "github.com/dialekt-lang/dialekt-go/src/dialekt/parser"
	. "github.com/onsi/ginkgo"
	. "github.com/onsi/ginkgo/extensions/table"
	. "github.com/onsi/gomega"
)

var _ = Describe("ParseList", func() {
	DescribeTable(
		"it produces the expected slice of tags",
		func(input string, expected []ast.Tag) {
			buf := bytes.NewBufferString(input)
			tags, err := ParseList(buf)

			Expect(err).ShouldNot(HaveOccurred())
			Expect(tags).To(Equal(expected))
		},

		Entry(
			"empty list",
			"",
			nil,
		),
		Entry(
			"single tag",
			"foo",
			[]ast.Tag{
				{Name: "foo"},
			},
		),
		Entry(
			"multiple tags",
			`foo "bar baz"`,
			[]ast.Tag{
				{Name: "foo"},
				{Name: "bar baz"},
			},
		),
	)

	DescribeTable(
		"returns an error when invalid input is provided",
		func(input string, expected string) {
			buf := bytes.NewBufferString(input)
			_, err := ParseList(buf)
			Expect(err).Should(HaveOccurred())
			Expect(err).To(MatchError(expected))
		},

		Entry(
			"pattern",
			"foo*",
			"unexpected pattern, expected tag",
		),
		Entry(
			"logical and",
			"and",
			"unexpected AND operator, expected tag",
		),
		Entry(
			"logical or",
			"or",
			"unexpected OR operator, expected tag",
		),
		Entry(
			"logical not",
			"not",
			"unexpected NOT operator, expected tag",
		),
		Entry(
			"open paren",
			"(",
			"unexpected opening parenthesis, expected tag",
		),
		Entry(
			"close paren",
			")",
			"unexpected closing parenthesis, expected tag",
		),
	)
})

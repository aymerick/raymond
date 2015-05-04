package parser

import (
	"log"
	"testing"

	"github.com/aymerick/raymond/ast"
)

const (
	VERBOSE = false
)

type parserTest struct {
	name  string
	input string
	node  ast.Node
}

var parserTests = []parserTest{
	{"Content", "Hello", ast.NewContentNode(0, "Hello")},
	{"Comment", "{{! This is a comment }}", ast.NewCommentNode(0, " This is a comment ")},
	{"Comment dash", "{{!-- This is a comment --}}", ast.NewCommentNode(0, " This is a comment ")},
}

func equal(a, b ast.Node) bool {
	return (a.String() == b.String())
}

func TestParser(t *testing.T) {
	for _, test := range parserTests {
		if VERBOSE {
			log.Printf("\n\n**********************************")
			log.Printf("Testing: %s", test.name)
		}

		node, err := Parse(test.input)
		if (err != nil) || (node == nil) || !equal(node, test.node) {
			t.Errorf("Test '%s' failed\ninput:\n\t'%s'\nexpected\n\t%q\ngot\n\t%q\nerror:\n\t%s", test.name, test.input, test.node, node, err)
		}
	}
}

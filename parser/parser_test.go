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
	name   string
	input  string
	output string
}

var parserTests = []parserTest{
	{"Content", "Hello", "CONTENT[Hello]\n"},
	{"Comment", "{{! This is a comment }}", "{{! 'This is a comment' }}\n"},
	{"Comment dash", "{{!-- This is a comment --}}", "{{! 'This is a comment' }}\n"},
}

func TestParser(t *testing.T) {
	for _, test := range parserTests {
		if VERBOSE {
			log.Printf("\n\n**********************************")
			log.Printf("Testing: %s", test.name)
		}

		node, err := Parse(test.input)
		output := ast.PrintNode(node)

		if (err != nil) || (test.output != output) {
			t.Errorf("Test '%s' failed\ninput:\n\t'%s'\nexpected\n\t%q\ngot\n\t%q\nerror:\n\t%s", test.name, test.input, test.output, output, err)
		}
	}
}

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
	//
	// Next tests come from:
	//   https://github.com/wycats/handlebars.js/blob/master/spec/parser.js
	//
	{"parses simple mustaches (1)", `{{123}}`, "{{ NUMBER{123} [] }}\n"},
	{"parses simple mustaches (2)", `{{"foo"}}`, "{{ \"foo\" [] }}\n"},
	{"parses simple mustaches (3)", `{{false}}`, "{{ BOOLEAN{false} [] }}\n"},
	{"parses simple mustaches (4)", `{{true}}`, "{{ BOOLEAN{true} [] }}\n"},
	{"parses simple mustaches (5)", `{{foo}}`, "{{ PATH:foo [] }}\n"},
	{"parses simple mustaches (6)", `{{foo?}}`, "{{ PATH:foo? [] }}\n"},
	{"parses simple mustaches (7)", `{{foo_}}`, "{{ PATH:foo_ [] }}\n"},
	{"parses simple mustaches (8)", `{{foo-}}`, "{{ PATH:foo- [] }}\n"},
	{"parses simple mustaches (9)", `{{foo:}}`, "{{ PATH:foo: [] }}\n"},

	{"parses simple mustaches with data", `{{@foo}}`, "{{ @PATH:foo [] }}\n"},
	{"parses simple mustaches with data paths", `{{@../foo}}`, "{{ @PATH:foo [] }}\n"},
	{"parses mustaches with paths", `{{foo/bar}}`, "{{ PATH:foo/bar [] }}\n"},
	{"parses mustaches with this/foo", `{{this/foo}}`, "{{ PATH:foo [] }}\n"},
	{"parses mustaches with - in a path", `{{foo-bar}}`, "{{ PATH:foo-bar [] }}\n"},
	{"parses mustaches with parameters", `{{foo bar}}`, "{{ PATH:foo [PATH:bar] }}\n"},
	{"parses mustaches with string parameters", `{{foo bar "baz" }}`, "{{ PATH:foo [PATH:bar, \"baz\"] }}\n"},
	{"parses mustaches with NUMBER parameters", `{{foo 1}}`, "{{ PATH:foo [NUMBER{1}] }}\n"},
	{"parses mustaches with BOOLEAN parameters (1)", `{{foo true}}`, "{{ PATH:foo [BOOLEAN{true}] }}\n"},
	{"parses mustaches with BOOLEAN parameters (2)", `{{foo false}}`, "{{ PATH:foo [BOOLEAN{false}] }}\n"},
	{"parses mustaches with DATA parameters", `{{foo @bar}}`, "{{ PATH:foo [@PATH:bar] }}\n"},

	{"parses mustaches with hash arguments (01)", `{{foo bar=baz}}`, "{{ PATH:foo [] HASH{bar=PATH:baz} }}\n"},
	{"parses mustaches with hash arguments (02)", `{{foo bar=1}}`, "{{ PATH:foo [] HASH{bar=NUMBER{1}} }}\n"},
	{"parses mustaches with hash arguments (03)", `{{foo bar=true}}`, "{{ PATH:foo [] HASH{bar=BOOLEAN{true}} }}\n"},
	{"parses mustaches with hash arguments (04)", `{{foo bar=false}}`, "{{ PATH:foo [] HASH{bar=BOOLEAN{false}} }}\n"},
	{"parses mustaches with hash arguments (05)", `{{foo bar=@baz}}`, "{{ PATH:foo [] HASH{bar=@PATH:baz} }}\n"},
	{"parses mustaches with hash arguments (06)", `{{foo bar=baz bat=bam}}`, "{{ PATH:foo [] HASH{bar=PATH:baz, bat=PATH:bam} }}\n"},
	{"parses mustaches with hash arguments (07)", `{{foo bar=baz bat="bam"}}`, "{{ PATH:foo [] HASH{bar=PATH:baz, bat=\"bam\"} }}\n"},
	{"parses mustaches with hash arguments (08)", `{{foo bat='bam'}}`, "{{ PATH:foo [] HASH{bat=\"bam\"} }}\n"},
	{"parses mustaches with hash arguments (09)", `{{foo omg bar=baz bat="bam"}}`, "{{ PATH:foo [PATH:omg] HASH{bar=PATH:baz, bat=\"bam\"} }}\n"},
	{"parses mustaches with hash arguments (10)", `{{foo omg bar=baz bat="bam" baz=1}}`, "{{ PATH:foo [PATH:omg] HASH{bar=PATH:baz, bat=\"bam\", baz=NUMBER{1}} }}\n"},
	{"parses mustaches with hash arguments (11)", `{{foo omg bar=baz bat="bam" baz=true}}`, "{{ PATH:foo [PATH:omg] HASH{bar=PATH:baz, bat=\"bam\", baz=BOOLEAN{true}} }}\n"},
	{"parses mustaches with hash arguments (12)", `{{foo omg bar=baz bat="bam" baz=false}}`, "{{ PATH:foo [PATH:omg] HASH{bar=PATH:baz, bat=\"bam\", baz=BOOLEAN{false}} }}\n"},

	{"parses contents followed by a mustache", `foo bar {{baz}}`, "CONTENT[ 'foo bar ' ]\n{{ PATH:baz [] }}\n"},

	{"parses a partial (1)", `{{> foo }}`, "{{> PARTIAL:foo }}\n"},
	{"parses a partial (2)", `{{> "foo" }}`, "{{> PARTIAL:foo }}\n"},
	{"parses a partial (3)", `{{> 1 }}`, "{{> PARTIAL:1 }}\n"},
	{"parses a partial with context", `{{> foo bar}}`, "{{> PARTIAL:foo PATH:bar }}\n"},
	{"parses a partial with hash", `{{> foo bar=bat}}`, "{{> PARTIAL:foo HASH{bar=PATH:bat} }}\n"},
	{"parses a partial with context and hash", `{{> foo bar bat=baz}}`, "{{> PARTIAL:foo PATH:bar HASH{bat=PATH:baz} }}\n"},
	{"parses a partial with a complex name", `{{> shared/partial?.bar}}`, "{{> PARTIAL:shared/partial?.bar }}\n"},

	{"parses a comment", `{{! this is a comment }}`, "{{! ' this is a comment ' }}\n"},
	{"parses a multi-line comment", "{{!\nthis is a multi-line comment\n}}", "{{! '\nthis is a multi-line comment\n' }}\n"},
}

func TestParser(t *testing.T) {
	for _, test := range parserTests {
		if VERBOSE {
			log.Printf("\n\n**********************************")
			log.Printf("Testing: %s", test.name)
		}

		output := ""

		node, err := Parse(test.input)
		if err == nil {
			output = ast.PrintNode(node)
		}

		if (err != nil) || (test.output != output) {
			t.Errorf("Test '%s' failed\ninput:\n\t'%s'\nexpected\n\t%q\ngot\n\t%q\nerror:\n\t%s", test.name, test.input, test.output, output, err)
		}
	}
}

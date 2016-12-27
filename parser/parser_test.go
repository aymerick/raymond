package parser

import (
	"fmt"
	"regexp"
	"testing"

	"github.com/gobuffalo/ray/ast"
	"github.com/gobuffalo/ray/lexer"
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

	{"parses an inverse section", `{{#foo}} bar {{^}} baz {{/foo}}`, "BLOCK:\n  PATH:foo []\n  PROGRAM:\n    CONTENT[ ' bar ' ]\n  {{^}}\n    CONTENT[ ' baz ' ]\n"},
	{"parses an inverse (else-style) section", `{{#foo}} bar {{else}} baz {{/foo}}`, "BLOCK:\n  PATH:foo []\n  PROGRAM:\n    CONTENT[ ' bar ' ]\n  {{^}}\n    CONTENT[ ' baz ' ]\n"},
	{"parses multiple inverse sections", `{{#foo}} bar {{else if bar}}{{else}} baz {{/foo}}`, "BLOCK:\n  PATH:foo []\n  PROGRAM:\n    CONTENT[ ' bar ' ]\n  {{^}}\n    BLOCK:\n      PATH:if [PATH:bar]\n      PROGRAM:\n      {{^}}\n        CONTENT[ ' baz ' ]\n"},
	{"parses empty blocks", `{{#foo}}{{/foo}}`, "BLOCK:\n  PATH:foo []\n  PROGRAM:\n"},
	{"parses empty blocks with empty inverse section", `{{#foo}}{{^}}{{/foo}}`, "BLOCK:\n  PATH:foo []\n  PROGRAM:\n  {{^}}\n"},
	{"parses empty blocks with empty inverse (else-style) section", `{{#foo}}{{else}}{{/foo}}`, "BLOCK:\n  PATH:foo []\n  PROGRAM:\n  {{^}}\n"},
	{"parses non-empty blocks with empty inverse section", `{{#foo}} bar {{^}}{{/foo}}`, "BLOCK:\n  PATH:foo []\n  PROGRAM:\n    CONTENT[ ' bar ' ]\n  {{^}}\n"},
	{"parses non-empty blocks with empty inverse (else-style) section", `{{#foo}} bar {{else}}{{/foo}}`, "BLOCK:\n  PATH:foo []\n  PROGRAM:\n    CONTENT[ ' bar ' ]\n  {{^}}\n"},
	{"parses empty blocks with non-empty inverse section", `{{#foo}}{{^}} bar {{/foo}}`, "BLOCK:\n  PATH:foo []\n  PROGRAM:\n  {{^}}\n    CONTENT[ ' bar ' ]\n"},
	{"parses empty blocks with non-empty inverse (else-style) section", `{{#foo}}{{else}} bar {{/foo}}`, "BLOCK:\n  PATH:foo []\n  PROGRAM:\n  {{^}}\n    CONTENT[ ' bar ' ]\n"},
	{"parses a standalone inverse section", `{{^foo}}bar{{/foo}}`, "BLOCK:\n  PATH:foo []\n  {{^}}\n    CONTENT[ 'bar' ]\n"},
	{"parses block with block params", `{{#foo as |bar baz|}}content{{/foo}}`, "BLOCK:\n  PATH:foo []\n  PROGRAM:\n    BLOCK PARAMS: [ bar baz ]\n    CONTENT[ 'content' ]\n"},
	{"parses inverse block with block params", `{{^foo as |bar baz|}}content{{/foo}}`, "BLOCK:\n  PATH:foo []\n  {{^}}\n    BLOCK PARAMS: [ bar baz ]\n    CONTENT[ 'content' ]\n"},
	{"parses chained inverse block with block params", `{{#foo}}{{else foo as |bar baz|}}content{{/foo}}`, "BLOCK:\n  PATH:foo []\n  PROGRAM:\n  {{^}}\n    BLOCK:\n      PATH:foo []\n      PROGRAM:\n        BLOCK PARAMS: [ bar baz ]\n        CONTENT[ 'content' ]\n"},
}

func TestParser(t *testing.T) {
	t.Parallel()

	for _, test := range parserTests {
		output := ""

		node, err := Parse(test.input)
		if err == nil {
			output = ast.Print(node)
		}

		if (err != nil) || (test.output != output) {
			t.Errorf("Test '%s' failed\ninput:\n\t'%s'\nexpected\n\t%q\ngot\n\t%q\nerror:\n\t%s", test.name, test.input, test.output, output, err)
		}
	}
}

var parserErrorTests = []parserTest{
	{"lexer error", `{{! unclosed comment`, "Lexer error"},
	{"syntax error", `foo{{^}}`, "Syntax error"},

	{"open raw block must be closed", `{{{{raw foo}} bar {{{{/raw}}}}`, "Expecting CloseRawBlock"},
	{"end raw block must be closed", `{{{{raw foo}}}} bar {{{{/raw}}`, "Expecting CloseRawBlock"},

	{"raw block names must match (1)", `{{{{1}}}}{{foo}}{{{{/raw}}}}`, "1 doesn't match raw"},
	{"raw block names must match (2)", `{{{{raw}}}}{{foo}}{{{{/1}}}}`, "raw doesn't match 1"},
	{"raw block names must match (3)", `{{{{goodbyes}}}}test{{{{/hellos}}}}`, "goodbyes doesn't match hellos"},

	{"open block must be closed", `{{#foo bar}}}{{/foo}}`, "Expecting Close"},
	{"end block must be closed", `{{#foo bar}}{{/foo}}}`, "Expecting Close"},
	{"an open block must have a end block", `{{#foo}}test`, "Expecting OpenEndBlock"},

	{"block names must match (1)", `{{#1 bar}}{{/foo}}`, "1 doesn't match foo"},
	{"block names must match (2)", `{{#foo bar}}{{/1}}`, "foo doesn't match 1"},
	{"block names must match (3)", `{{#foo}}test{{/bar}}`, "foo doesn't match bar"},

	{"an mustache must terminate with a close mustache", `{{foo}}}`, "Expecting Close"},
	{"an unescaped mustache must terminate with a close unescaped mustache", `{{{foo}}`, "Expecting CloseUnescaped"},

	{"an partial must terminate with a close mustache", `{{> foo}}}`, "Expecting Close"},
	{"a subexpression must terminate with a close subexpression", `{{foo (false}}`, "Expecting CloseSexpr"},

	{"raises on missing hash value (1)", `{{foo bar=}}`, "Parse error on line 1"},
	{"raises on missing hash value (2)", `{{foo bar=baz bim=}}`, "Parse error on line 1"},

	{"block param must have at least one param", `{{#foo as ||}}content{{/foo}}`, "Expecting ID"},
	{"open block params must be closed", `{{#foo as |}}content{{/foo}}`, "Expecting ID"},

	{"a path must start with an ID", `{{#/}}content{{/foo}}`, "Expecting ID"},
	{"a path must end with an ID", `{{foo/bar/}}`, "Expecting ID"},

	//
	// Next tests come from:
	//   https://github.com/wycats/handlebars.js/blob/master/spec/parser.js
	//
	{"throws on old inverse section", `{{else foo}}bar{{/foo}}`, ""},

	{"raises if there's a parser error (1)", `foo{{^}}bar`, "Parse error on line 1"},
	{"raises if there's a parser error (2)", `{{foo}`, "Parse error on line 1"},
	{"raises if there's a parser error (3)", `{{foo &}}`, "Parse error on line 1"},
	{"raises if there's a parser error (4)", `{{#goodbyes}}{{/hellos}}`, "Parse error on line 1"},
	{"raises if there's a parser error (5)", `{{#goodbyes}}{{/hellos}}`, "goodbyes doesn't match hellos"},

	{"should handle invalid paths (1)", `{{foo/../bar}}`, `Invalid path: foo/..`},
	{"should handle invalid paths (2)", `{{foo/./bar}}`, `Invalid path: foo/.`},
	{"should handle invalid paths (3)", `{{foo/this/bar}}`, `Invalid path: foo/this`},

	{"knows how to report the correct line number in errors (1)", "hello\nmy\n{{foo}", "Parse error on line 3"},
	{"knows how to report the correct line number in errors (2)", "hello\n\nmy\n\n{{foo}", "Parse error on line 5"},

	{"knows how to report the correct line number in errors when the first character is a newline", "\n\nhello\n\nmy\n\n{{foo}", "Parse error on line 7"},
}

func TestParserErrors(t *testing.T) {
	t.Parallel()

	for _, test := range parserErrorTests {
		node, err := Parse(test.input)
		if err == nil {
			output := ast.Print(node)
			tokens := lexer.Collect(test.input)

			t.Errorf("Test '%s' failed - Error expected\ninput:\n\t'%s'\ngot\n\t%q\ntokens:\n\t%q", test.name, test.input, output, tokens)
		} else if test.output != "" {
			matched, errMatch := regexp.MatchString(regexp.QuoteMeta(test.output), fmt.Sprint(err))
			if errMatch != nil {
				panic("Failed to match regexp")
			}

			if !matched {
				t.Errorf("Test '%s' failed - Incorrect error returned\ninput:\n\t'%s'\nexpected\n\t%q\ngot\n\t%q", test.name, test.input, test.output, err)
			}
		}
	}
}

// package example
func Example() {
	source := "You know {{nothing}} John Snow"

	// parse template
	program, err := Parse(source)
	if err != nil {
		panic(err)
	}

	// print AST
	output := ast.Print(program)

	fmt.Print(output)
	// CONTENT[ 'You know ' ]
	// {{ PATH:nothing [] }}
	// CONTENT[ ' John Snow' ]
}

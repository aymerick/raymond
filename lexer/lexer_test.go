package lexer

import (
	"fmt"
	"testing"
)

type lexTest struct {
	name   string
	input  string
	tokens []Token
}

// helpers
func tokContent(val string) Token { return Token{TokenContent, val, 0, 1} }
func tokID(val string) Token      { return Token{TokenID, val, 0, 1} }
func tokSep(val string) Token     { return Token{TokenSep, val, 0, 1} }
func tokString(val string) Token  { return Token{TokenString, val, 0, 1} }
func tokNumber(val string) Token  { return Token{TokenNumber, val, 0, 1} }
func tokInverse(val string) Token { return Token{TokenInverse, val, 0, 1} }
func tokBool(val string) Token    { return Token{TokenBoolean, val, 0, 1} }
func tokError(val string) Token   { return Token{TokenError, val, 0, 1} }
func tokComment(val string) Token { return Token{TokenComment, val, 0, 1} }

var tokEOF = Token{TokenEOF, "", 0, 1}
var tokEquals = Token{TokenEquals, "=", 0, 1}
var tokData = Token{TokenData, "@", 0, 1}
var tokOpen = Token{TokenOpen, "{{", 0, 1}
var tokOpenAmp = Token{TokenOpen, "{{&", 0, 1}
var tokOpenPartial = Token{TokenOpenPartial, "{{>", 0, 1}
var tokClose = Token{TokenClose, "}}", 0, 1}
var tokOpenStrip = Token{TokenOpen, "{{~", 0, 1}
var tokCloseStrip = Token{TokenClose, "~}}", 0, 1}
var tokOpenUnescaped = Token{TokenOpenUnescaped, "{{{", 0, 1}
var tokCloseUnescaped = Token{TokenCloseUnescaped, "}}}", 0, 1}
var tokOpenUnescapedStrip = Token{TokenOpenUnescaped, "{{~{", 0, 1}
var tokCloseUnescapedStrip = Token{TokenCloseUnescaped, "}~}}", 0, 1}
var tokOpenBlock = Token{TokenOpenBlock, "{{#", 0, 1}
var tokOpenEndBlock = Token{TokenOpenEndBlock, "{{/", 0, 1}
var tokOpenInverse = Token{TokenOpenInverse, "{{^", 0, 1}
var tokOpenInverseChain = Token{TokenOpenInverseChain, "{{else", 0, 1}
var tokOpenSexpr = Token{TokenOpenSexpr, "(", 0, 1}
var tokCloseSexpr = Token{TokenCloseSexpr, ")", 0, 1}
var tokOpenBlockParams = Token{TokenOpenBlockParams, "as |", 0, 1}
var tokCloseBlockParams = Token{TokenCloseBlockParams, "|", 0, 1}
var tokOpenRawBlock = Token{TokenOpenRawBlock, "{{{{", 0, 1}
var tokCloseRawBlock = Token{TokenCloseRawBlock, "}}}}", 0, 1}
var tokOpenEndRawBlock = Token{TokenOpenEndRawBlock, "{{{{/", 0, 1}

var lexTests = []lexTest{
	{"empty", "", []Token{tokEOF}},
	{"spaces", " \t\n", []Token{tokContent(" \t\n"), tokEOF}},
	{"content", `now is the time`, []Token{tokContent(`now is the time`), tokEOF}},

	{
		`does not tokenizes identifier starting with true as boolean`,
		`{{ foo truebar }}`,
		[]Token{tokOpen, tokID("foo"), tokID("truebar"), tokClose, tokEOF},
	},
	{
		`does not tokenizes identifier starting with false as boolean`,
		`{{ foo falsebar }}`,
		[]Token{tokOpen, tokID("foo"), tokID("falsebar"), tokClose, tokEOF},
	},
	{
		`tokenizes raw block`,
		`{{{{foo}}}} {{{{/foo}}}}`,
		[]Token{tokOpenRawBlock, tokID("foo"), tokCloseRawBlock, tokContent(" "), tokOpenEndRawBlock, tokID("foo"), tokCloseRawBlock, tokEOF},
	},
	{
		`tokenizes raw block with mustaches in content`,
		`{{{{foo}}}}{{bar}}{{{{/foo}}}}`,
		[]Token{tokOpenRawBlock, tokID("foo"), tokCloseRawBlock, tokContent("{{bar}}"), tokOpenEndRawBlock, tokID("foo"), tokCloseRawBlock, tokEOF},
	},
	{
		`tokenizes @../foo`,
		`{{@../foo}}`,
		[]Token{tokOpen, tokData, tokID(".."), tokSep("/"), tokID("foo"), tokClose, tokEOF},
	},
	{
		`tokenizes escaped mustaches`,
		"\\{{bar}}",
		[]Token{tokContent("{{bar}}"), tokEOF},
	},
	{
		`tokenizes strip mustaches`,
		`{{~ foo ~}}`,
		[]Token{tokOpenStrip, tokID("foo"), tokCloseStrip, tokEOF},
	},
	{
		`tokenizes unescaped strip mustaches`,
		`{{~{ foo }~}}`,
		[]Token{tokOpenUnescapedStrip, tokID("foo"), tokCloseUnescapedStrip, tokEOF},
	},

	//
	// Next tests come from:
	//   https://github.com/wycats/handlebars.js/blob/master/spec/tokenizer.js
	//
	{
		`tokenizes a simple mustache as "OPEN ID CLOSE"`,
		`{{foo}}`,
		[]Token{tokOpen, tokID("foo"), tokClose, tokEOF},
	},
	{
		`supports unescaping with &`,
		`{{&bar}}`,
		[]Token{tokOpenAmp, tokID("bar"), tokClose, tokEOF},
	},
	{
		`supports unescaping with {{{`,
		`{{{bar}}}`,
		[]Token{tokOpenUnescaped, tokID("bar"), tokCloseUnescaped, tokEOF},
	},
	{
		`supports escaping delimiters`,
		"{{foo}} \\{{bar}} {{baz}}",
		[]Token{tokOpen, tokID("foo"), tokClose, tokContent(" "), tokContent("{{bar}} "), tokOpen, tokID("baz"), tokClose, tokEOF},
	},
	{
		`supports escaping multiple delimiters`,
		"{{foo}} \\{{bar}} \\{{baz}}",
		[]Token{tokOpen, tokID("foo"), tokClose, tokContent(" "), tokContent("{{bar}} "), tokContent("{{baz}}"), tokEOF},
	},
	{
		`supports escaping a triple stash`,
		"{{foo}} \\{{{bar}}} {{baz}}",
		[]Token{tokOpen, tokID("foo"), tokClose, tokContent(" "), tokContent("{{{bar}}} "), tokOpen, tokID("baz"), tokClose, tokEOF},
	},
	{
		`supports escaping escape character`,
		"{{foo}} \\\\{{bar}} {{baz}}",
		[]Token{tokOpen, tokID("foo"), tokClose, tokContent(" \\"), tokOpen, tokID("bar"), tokClose, tokContent(" "), tokOpen, tokID("baz"), tokClose, tokEOF},
	},
	{
		`supports escaping multiple escape characters`,
		"{{foo}} \\\\{{bar}} \\\\{{baz}}",
		[]Token{tokOpen, tokID("foo"), tokClose, tokContent(" \\"), tokOpen, tokID("bar"), tokClose, tokContent(" \\"), tokOpen, tokID("baz"), tokClose, tokEOF},
	},
	{
		`supports escaped mustaches after escaped escape characters`,
		"{{foo}} \\\\{{bar}} \\{{baz}}",
		// NOTE: JS implementation returns:
		//   ['OPEN', 'ID', 'CLOSE', 'CONTENT', 'OPEN', 'ID', 'CLOSE', 'CONTENT', 'CONTENT', 'CONTENT'],
		// WTF is the last CONTENT ?
		[]Token{tokOpen, tokID("foo"), tokClose, tokContent(" \\"), tokOpen, tokID("bar"), tokClose, tokContent(" "), tokContent("{{baz}}"), tokEOF},
	},
	{
		`supports escaped escape characters after escaped mustaches`,
		"{{foo}} \\{{bar}} \\\\{{baz}}",
		// NOTE: JS implementation returns:
		//   []Token{tokOpen, tokID("foo"), tokClose, tokContent(" "), tokContent("{{bar}} "), tokContent("\\"), tokOpen, tokID("baz"), tokClose, tokEOF},
		[]Token{tokOpen, tokID("foo"), tokClose, tokContent(" "), tokContent("{{bar}} \\"), tokOpen, tokID("baz"), tokClose, tokEOF},
	},
	{
		`supports escaped escape character on a triple stash`,
		"{{foo}} \\\\{{{bar}}} {{baz}}",
		[]Token{tokOpen, tokID("foo"), tokClose, tokContent(" \\"), tokOpenUnescaped, tokID("bar"), tokCloseUnescaped, tokContent(" "), tokOpen, tokID("baz"), tokClose, tokEOF},
	},
	{
		`tokenizes a simple path`,
		`{{foo/bar}}`,
		[]Token{tokOpen, tokID("foo"), tokSep("/"), tokID("bar"), tokClose, tokEOF},
	},
	{
		`allows dot notation (1)`,
		`{{foo.bar}}`,
		[]Token{tokOpen, tokID("foo"), tokSep("."), tokID("bar"), tokClose, tokEOF},
	},
	{
		`allows dot notation (2)`,
		`{{foo.bar.baz}}`,
		[]Token{tokOpen, tokID("foo"), tokSep("."), tokID("bar"), tokSep("."), tokID("baz"), tokClose, tokEOF},
	},
	{
		`allows path literals with []`,
		`{{foo.[bar]}}`,
		[]Token{tokOpen, tokID("foo"), tokSep("."), tokID("[bar]"), tokClose, tokEOF},
	},
	{
		`allows multiple path literals on a line with []`,
		`{{foo.[bar]}}{{foo.[baz]}}`,
		[]Token{tokOpen, tokID("foo"), tokSep("."), tokID("[bar]"), tokClose, tokOpen, tokID("foo"), tokSep("."), tokID("[baz]"), tokClose, tokEOF},
	},
	{
		`tokenizes {{.}} as OPEN ID CLOSE`,
		`{{.}}`,
		[]Token{tokOpen, tokID("."), tokClose, tokEOF},
	},
	{
		`tokenizes a path as "OPEN (ID SEP)* ID CLOSE"`,
		`{{../foo/bar}}`,
		[]Token{tokOpen, tokID(".."), tokSep("/"), tokID("foo"), tokSep("/"), tokID("bar"), tokClose, tokEOF},
	},
	{
		`tokenizes a path with .. as a parent path`,
		`{{../foo.bar}}`,
		[]Token{tokOpen, tokID(".."), tokSep("/"), tokID("foo"), tokSep("."), tokID("bar"), tokClose, tokEOF},
	},
	{
		`tokenizes a path with this/foo as OPEN ID SEP ID CLOSE`,
		`{{this/foo}}`,
		[]Token{tokOpen, tokID("this"), tokSep("/"), tokID("foo"), tokClose, tokEOF},
	},
	{
		`tokenizes a simple mustache with spaces as "OPEN ID CLOSE"`,
		`{{  foo  }}`,
		[]Token{tokOpen, tokID("foo"), tokClose, tokEOF},
	},
	{
		`tokenizes a simple mustache with line breaks as "OPEN ID ID CLOSE"`,
		"{{  foo  \n   bar }}",
		[]Token{tokOpen, tokID("foo"), tokID("bar"), tokClose, tokEOF},
	},
	{
		`tokenizes raw content as "CONTENT"`,
		`foo {{ bar }} baz`,
		[]Token{tokContent("foo "), tokOpen, tokID("bar"), tokClose, tokContent(" baz"), tokEOF},
	},
	{
		`tokenizes a partial as "OPEN_PARTIAL ID CLOSE"`,
		`{{> foo}}`,
		[]Token{tokOpenPartial, tokID("foo"), tokClose, tokEOF},
	},
	{
		`tokenizes a partial with context as "OPEN_PARTIAL ID ID CLOSE"`,
		`{{> foo bar }}`,
		[]Token{tokOpenPartial, tokID("foo"), tokID("bar"), tokClose, tokEOF},
	},
	{
		`tokenizes a partial without spaces as "OPEN_PARTIAL ID CLOSE"`,
		`{{>foo}}`,
		[]Token{tokOpenPartial, tokID("foo"), tokClose, tokEOF},
	},
	{
		`tokenizes a partial space at the }); as "OPEN_PARTIAL ID CLOSE"`,
		`{{>foo  }}`,
		[]Token{tokOpenPartial, tokID("foo"), tokClose, tokEOF},
	},
	{
		`tokenizes a partial space at the }); as "OPEN_PARTIAL ID CLOSE"`,
		`{{>foo/bar.baz  }}`,
		[]Token{tokOpenPartial, tokID("foo"), tokSep("/"), tokID("bar"), tokSep("."), tokID("baz"), tokClose, tokEOF},
	},
	{
		`tokenizes a comment as "COMMENT"`,
		`foo {{! this is a comment }} bar {{ baz }}`,
		[]Token{tokContent("foo "), tokComment("{{! this is a comment }}"), tokContent(" bar "), tokOpen, tokID("baz"), tokClose, tokEOF},
	},
	{
		`tokenizes a block comment as "COMMENT"`,
		`foo {{!-- this is a {{comment}} --}} bar {{ baz }}`,
		[]Token{tokContent("foo "), tokComment("{{!-- this is a {{comment}} --}}"), tokContent(" bar "), tokOpen, tokID("baz"), tokClose, tokEOF},
	},
	{
		`tokenizes a block comment with whitespace as "COMMENT"`,
		"foo {{!-- this is a\n{{comment}}\n--}} bar {{ baz }}",
		[]Token{tokContent("foo "), tokComment("{{!-- this is a\n{{comment}}\n--}}"), tokContent(" bar "), tokOpen, tokID("baz"), tokClose, tokEOF},
	},
	{
		`tokenizes open and closing blocks as OPEN_BLOCK, ID, CLOSE ..., OPEN_ENDBLOCK ID CLOSE`,
		`{{#foo}}content{{/foo}}`,
		[]Token{tokOpenBlock, tokID("foo"), tokClose, tokContent("content"), tokOpenEndBlock, tokID("foo"), tokClose, tokEOF},
	},
	{
		`tokenizes inverse sections as "INVERSE"`,
		`{{^}}`,
		[]Token{tokInverse("{{^}}"), tokEOF},
	},
	{
		`tokenizes inverse sections as "INVERSE" with alternate format`,
		`{{else}}`,
		[]Token{tokInverse("{{else}}"), tokEOF},
	},
	{
		`tokenizes inverse sections as "INVERSE" with spaces`,
		`{{ else }}`,
		[]Token{tokInverse("{{ else }}"), tokEOF},
	},
	{
		`tokenizes inverse sections with ID as "OPEN_INVERSE ID CLOSE"`,
		`{{^foo}}`,
		[]Token{tokOpenInverse, tokID("foo"), tokClose, tokEOF},
	},
	{
		`tokenizes inverse sections with ID and spaces as "OPEN_INVERSE ID CLOSE"`,
		`{{^ foo  }}`,
		[]Token{tokOpenInverse, tokID("foo"), tokClose, tokEOF},
	},
	{
		`tokenizes mustaches with params as "OPEN ID ID ID CLOSE"`,
		`{{ foo bar baz }}`,
		[]Token{tokOpen, tokID("foo"), tokID("bar"), tokID("baz"), tokClose, tokEOF},
	},
	{
		`tokenizes mustaches with String params as "OPEN ID ID STRING CLOSE"`,
		`{{ foo bar "baz" }}`,
		[]Token{tokOpen, tokID("foo"), tokID("bar"), tokString("baz"), tokClose, tokEOF},
	},
	{
		`tokenizes mustaches with String params using single quotes as "OPEN ID ID STRING CLOSE"`,
		`{{ foo bar 'baz' }}`,
		[]Token{tokOpen, tokID("foo"), tokID("bar"), tokString("baz"), tokClose, tokEOF},
	},
	{
		`tokenizes String params with spaces inside as "STRING"`,
		`{{ foo bar "baz bat" }}`,
		[]Token{tokOpen, tokID("foo"), tokID("bar"), tokString("baz bat"), tokClose, tokEOF},
	},
	{
		`tokenizes String params with escapes quotes as STRING`,
		`{{ foo "bar\"baz" }}`,
		[]Token{tokOpen, tokID("foo"), tokString(`bar"baz`), tokClose, tokEOF},
	},
	{
		`tokenizes String params using single quotes with escapes quotes as STRING`,
		`{{ foo 'bar\'baz' }}`,
		[]Token{tokOpen, tokID("foo"), tokString(`bar'baz`), tokClose, tokEOF},
	},
	{
		`tokenizes numbers`,
		`{{ foo 1 }}`,
		[]Token{tokOpen, tokID("foo"), tokNumber("1"), tokClose, tokEOF},
	},
	{
		`tokenizes floats`,
		`{{ foo 1.1 }}`,
		[]Token{tokOpen, tokID("foo"), tokNumber("1.1"), tokClose, tokEOF},
	},
	{
		`tokenizes negative numbers`,
		`{{ foo -1 }}`,
		[]Token{tokOpen, tokID("foo"), tokNumber("-1"), tokClose, tokEOF},
	},
	{
		`tokenizes negative floats`,
		`{{ foo -1.1 }}`,
		[]Token{tokOpen, tokID("foo"), tokNumber("-1.1"), tokClose, tokEOF},
	},
	{
		`tokenizes boolean true`,
		`{{ foo true }}`,
		[]Token{tokOpen, tokID("foo"), tokBool("true"), tokClose, tokEOF},
	},
	{
		`tokenizes boolean false`,
		`{{ foo false }}`,
		[]Token{tokOpen, tokID("foo"), tokBool("false"), tokClose, tokEOF},
	},
	// SKIP: 'tokenizes undefined and null'
	{
		`tokenizes hash arguments (1)`,
		`{{ foo bar=baz }}`,
		[]Token{tokOpen, tokID("foo"), tokID("bar"), tokEquals, tokID("baz"), tokClose, tokEOF},
	},
	{
		`tokenizes hash arguments (2)`,
		`{{ foo bar baz=bat }}`,
		[]Token{tokOpen, tokID("foo"), tokID("bar"), tokID("baz"), tokEquals, tokID("bat"), tokClose, tokEOF},
	},
	{
		`tokenizes hash arguments (3)`,
		`{{ foo bar baz=1 }}`,
		[]Token{tokOpen, tokID("foo"), tokID("bar"), tokID("baz"), tokEquals, tokNumber("1"), tokClose, tokEOF},
	},
	{
		`tokenizes hash arguments (4)`,
		`{{ foo bar baz=true }}`,
		[]Token{tokOpen, tokID("foo"), tokID("bar"), tokID("baz"), tokEquals, tokBool("true"), tokClose, tokEOF},
	},
	{
		`tokenizes hash arguments (5)`,
		`{{ foo bar baz=false }}`,
		[]Token{tokOpen, tokID("foo"), tokID("bar"), tokID("baz"), tokEquals, tokBool("false"), tokClose, tokEOF},
	},
	{
		`tokenizes hash arguments (6)`,
		"{{ foo bar\n  baz=bat }}",
		[]Token{tokOpen, tokID("foo"), tokID("bar"), tokID("baz"), tokEquals, tokID("bat"), tokClose, tokEOF},
	},
	{
		`tokenizes hash arguments (7)`,
		`{{ foo bar baz="bat" }}`,
		[]Token{tokOpen, tokID("foo"), tokID("bar"), tokID("baz"), tokEquals, tokString("bat"), tokClose, tokEOF},
	},
	{
		`tokenizes hash arguments (8)`,
		`{{ foo bar baz="bat" bam=wot }}`,
		[]Token{tokOpen, tokID("foo"), tokID("bar"), tokID("baz"), tokEquals, tokString("bat"), tokID("bam"), tokEquals, tokID("wot"), tokClose, tokEOF},
	},
	{
		`tokenizes hash arguments (9)`,
		`{{foo omg bar=baz bat="bam"}}`,
		[]Token{tokOpen, tokID("foo"), tokID("omg"), tokID("bar"), tokEquals, tokID("baz"), tokID("bat"), tokEquals, tokString("bam"), tokClose, tokEOF},
	},
	{
		`tokenizes special @ identifiers (1)`,
		`{{ @foo }}`,
		[]Token{tokOpen, tokData, tokID("foo"), tokClose, tokEOF},
	},
	{
		`tokenizes special @ identifiers (2)`,
		`{{ foo @bar }}`,
		[]Token{tokOpen, tokID("foo"), tokData, tokID("bar"), tokClose, tokEOF},
	},
	{
		`tokenizes special @ identifiers (3)`,
		`{{ foo bar=@baz }}`,
		[]Token{tokOpen, tokID("foo"), tokID("bar"), tokEquals, tokData, tokID("baz"), tokClose, tokEOF},
	},
	{
		`does not time out in a mustache with a single } followed by EOF`,
		`{{foo}`,
		[]Token{tokOpen, tokID("foo"), tokError("Unexpected character in expression: '}'")},
	},
	{
		`does not time out in a mustache when invalid ID characters are used`,
		`{{foo & }}`,
		[]Token{tokOpen, tokID("foo"), tokError("Unexpected character in expression: '&'")},
	},
	{
		`tokenizes subexpressions (1)`,
		`{{foo (bar)}}`,
		[]Token{tokOpen, tokID("foo"), tokOpenSexpr, tokID("bar"), tokCloseSexpr, tokClose, tokEOF},
	},
	{
		`tokenizes subexpressions (2)`,
		`{{foo (a-x b-y)}}`,
		[]Token{tokOpen, tokID("foo"), tokOpenSexpr, tokID("a-x"), tokID("b-y"), tokCloseSexpr, tokClose, tokEOF},
	},
	{
		`tokenizes nested subexpressions`,
		`{{foo (bar (lol rofl)) (baz)}}`,
		[]Token{tokOpen, tokID("foo"), tokOpenSexpr, tokID("bar"), tokOpenSexpr, tokID("lol"), tokID("rofl"), tokCloseSexpr, tokCloseSexpr, tokOpenSexpr, tokID("baz"), tokCloseSexpr, tokClose, tokEOF},
	},
	{
		`tokenizes nested subexpressions: literals`,
		`{{foo (bar (lol true) false) (baz 1) (blah 'b') (blorg "c")}}`,
		[]Token{tokOpen, tokID("foo"), tokOpenSexpr, tokID("bar"), tokOpenSexpr, tokID("lol"), tokBool("true"), tokCloseSexpr, tokBool("false"), tokCloseSexpr, tokOpenSexpr, tokID("baz"), tokNumber("1"), tokCloseSexpr, tokOpenSexpr, tokID("blah"), tokString("b"), tokCloseSexpr, tokOpenSexpr, tokID("blorg"), tokString("c"), tokCloseSexpr, tokClose, tokEOF},
	},
	{
		`tokenizes block params (1)`,
		`{{#foo as |bar|}}`,
		[]Token{tokOpenBlock, tokID("foo"), tokOpenBlockParams, tokID("bar"), tokCloseBlockParams, tokClose, tokEOF},
	},
	{
		`tokenizes block params (2)`,
		`{{#foo as |bar baz|}}`,
		[]Token{tokOpenBlock, tokID("foo"), tokOpenBlockParams, tokID("bar"), tokID("baz"), tokCloseBlockParams, tokClose, tokEOF},
	},
	{
		`tokenizes block params (3)`,
		`{{#foo as | bar baz |}}`,
		[]Token{tokOpenBlock, tokID("foo"), tokOpenBlockParams, tokID("bar"), tokID("baz"), tokCloseBlockParams, tokClose, tokEOF},
	},
	{
		`tokenizes block params (4)`,
		`{{#foo as as | bar baz |}}`,
		[]Token{tokOpenBlock, tokID("foo"), tokID("as"), tokOpenBlockParams, tokID("bar"), tokID("baz"), tokCloseBlockParams, tokClose, tokEOF},
	},
	{
		`tokenizes block params (5)`,
		`{{else foo as |bar baz|}}`,
		[]Token{tokOpenInverseChain, tokID("foo"), tokOpenBlockParams, tokID("bar"), tokID("baz"), tokCloseBlockParams, tokClose, tokEOF},
	},
}

func collect(t *lexTest) []Token {
	var result []Token

	l := scanWithName(t.input, t.name)
	for {
		token := l.NextToken()
		result = append(result, token)

		if token.Kind == TokenEOF || token.Kind == TokenError {
			break
		}
	}

	return result
}

func equal(i1, i2 []Token, checkPos bool) bool {
	if len(i1) != len(i2) {
		return false
	}

	for k := range i1 {
		if i1[k].Kind != i2[k].Kind {
			return false
		}

		if checkPos && i1[k].Pos != i2[k].Pos {
			return false
		}

		if i1[k].Val != i2[k].Val {
			return false
		}
	}

	return true
}

func TestLexer(t *testing.T) {
	t.Parallel()

	for _, test := range lexTests {
		tokens := collect(&test)
		if !equal(tokens, test.tokens, false) {
			t.Errorf("Test '%s' failed\ninput:\n\t'%s'\nexpected\n\t%v\ngot\n\t%+v\n", test.name, test.input, test.tokens, tokens)
		}
	}
}

// @todo Test errors:
//   `{{{{raw foo`

// package example
func Example() {
	source := "You know {{nothing}} John Snow"

	output := ""

	lex := Scan(source)
	for {
		// consume next token
		token := lex.NextToken()

		output += fmt.Sprintf(" %s", token)

		// stops when all tokens have been consumed, or on error
		if token.Kind == TokenEOF || token.Kind == TokenError {
			break
		}
	}

	fmt.Print(output)
	// Output: Content{"You know "} Open{"{{"} ID{"nothing"} Close{"}}"} Content{" John Snow"} EOF
}

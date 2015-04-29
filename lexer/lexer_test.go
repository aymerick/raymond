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

var tokenName = map[TokenKind]string{
	TokenError:             "Error",
	TokenEOF:               "EOF",
	TokenContent:           "Content",
	TokenComment:           "Comment",
	TokenOpen:              "Open",
	TokenClose:             "Close",
	TokenOpenUnescaped:     "OpenUnescaped",
	TokenCloseUnescaped:    "CloseUnescaped",
	TokenOpenBlock:         "OpenBlock",
	TokenOpenEndBlock:      "OpenEndBlock",
	TokenOpenRawBlock:      "OpenRawBlock",
	TokenCloseRawBlock:     "CloseRawBlock",
	TokenEndRawBlock:       "EndRawBlock",
	TokenOpenBlockParams:   "OpenBlockParams",
	TokenCloseBlockPaarams: "CloseBlockPaarams",
	TokenInverse:           "Inverse",
	TokenOpenInverse:       "OpenInverse",
	TokenOpenInverseChain:  "OpenInverseChain",
	TokenOpenPartial:       "OpenPartial",
	TokenOpenSexpr:         "OpenSexpr",
	TokenCloseSexpr:        "CloseSexpr",
	TokenID:                "ID",
	TokenEquals:            "Equals",
	TokenString:            "String",
	TokenNumber:            "Number",
	TokenBoolean:           "Boolean",
	TokenUndefined:         "Undefined",
	TokenNull:              "Null",
	TokenData:              "Data",
	TokenSep:               "Sep",
}

func (k TokenKind) String() string {
	s := tokenName[k]
	if s == "" {
		return fmt.Sprintf("Token-%d", int(k))
	}
	return s
}

func (t Token) String() string {
	result := fmt.Sprintf("%d:%s", t.pos, t.kind)

	if (t.kind >= TokenContent) && len(t.val) > 0 {
		if len(t.val) > 20 {
			result += fmt.Sprintf("{%.20q...}", t.val)
		} else {
			result += fmt.Sprintf("{%q}", t.val)
		}
	}

	return result
}

func tokEOF(pos int) Token   { return Token{TokenEOF, pos, ""} }
func tokOpen(pos int) Token  { return Token{TokenOpen, pos, "{{"} }
func tokClose(pos int) Token { return Token{TokenClose, pos, "{{"} }

var lexTests = []lexTest{
	// cf. https://github.com/golang/go/blob/master/src/text/template/parse/lex_test.go
	{"empty", "", []Token{tokEOF(0)}},
	{"spaces", " \t\n", []Token{{TokenContent, 0, " \t\n"}, tokEOF(3)}},
	{"content", `now is the time`, []Token{{TokenContent, 0, `now is the time`}, tokEOF(15)}},
	// cf. https://github.com/wycats/handlebars.js/blob/master/spec/tokenizer.js
	{"OPEN ID CLOSE", `{{foo}}`, []Token{tokOpen(0), Token{TokenID, 2, "foo"}, tokClose(5), tokEOF(7)}},
}

func collect(t *lexTest) []Token {
	var result []Token

	l := Scan(t.input, t.name)
	for {
		token := l.NextToken()
		result = append(result, token)

		if token.kind == TokenEOF || token.kind == TokenError {
			break
		}
	}

	return result
}

func equal(i1, i2 []Token) bool {
	if len(i1) != len(i2) {
		return false
	}

	for k := range i1 {
		if i1[k].kind != i2[k].kind {
			return false
		}

		if i1[k].pos != i2[k].pos {
			return false
		}

		if i1[k].val != i2[k].val {
			return false
		}
	}

	return true
}

func TestLexer(t *testing.T) {
	for _, test := range lexTests {
		tokens := collect(&test)
		if !equal(tokens, test.tokens) {
			t.Errorf("[%s] Failed to scan '%s'\nexpected\n\t%v\ngot\n\t%+v\n", test.name, test.input, test.tokens, tokens)
		}
	}
}

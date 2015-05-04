package lexer

import (
	"fmt"
	"regexp"
	"strings"
	"unicode"
	"unicode/utf8"
)

// Reference: https://github.com/wycats/handlebars.js/blob/master/src/handlebars.l

const (
	TokenError TokenKind = iota
	TokenEOF

	// mustache delimiters
	TokenOpen             // OPEN: <mu>"{{"{LEFT_STRIP}?"&" - 22. OPEN: <mu>"{{"{LEFT_STRIP}?
	TokenClose            // CLOSE: <mu>{RIGHT_STRIP}?"}}"
	TokenOpenRawBlock     // OPEN_RAW_BLOCK: <mu>"{{{{"
	TokenCloseRawBlock    // CLOSE_RAW_BLOCK: <mu>"}}}}"
	TokenOpenUnescaped    // OPEN_UNESCAPED: <mu>"{{"{LEFT_STRIP}?"{"
	TokenCloseUnescaped   // CLOSE_UNESCAPED: <mu>"}"{RIGHT_STRIP}?"}}"
	TokenOpenBlock        // OPEN_BLOCK: <mu>"{{"{LEFT_STRIP}?"#"
	TokenOpenEndBlock     // OPEN_ENDBLOCK: <mu>"{{"{LEFT_STRIP}?"/"
	TokenInverse          // INVERSE: <mu>"{{"{LEFT_STRIP}?"^"\s*{RIGHT_STRIP}?"}}" - 15. INVERSE: <mu>"{{"{LEFT_STRIP}?\s*"else"\s*{RIGHT_STRIP}?"}}"
	TokenOpenInverse      // OPEN_INVERSE: <mu>"{{"{LEFT_STRIP}?"^"
	TokenOpenInverseChain // OPEN_INVERSE_CHAIN: <mu>"{{"{LEFT_STRIP}?\s*"else"
	TokenOpenPartial      // OPEN_PARTIAL: <mu>"{{"{LEFT_STRIP}?">"
	TokenEndRawBlock      // END_RAW_BLOCK: <raw>"{{{{/"[^\s!"#%-,\.\/;->@\[-\^`\{-~]+/[=}\s\/.]"}}}}"
	TokenComment          // COMMENT: <com>[\s\S]*?"--"{RIGHT_STRIP}?"}}" - 20. begin 'com': <mu>"{{"{LEFT_STRIP}?"!--" - 21. COMMENT: <mu>"{{"{LEFT_STRIP}?"!"[\s\S]*?"}}"

	// inside mustaches
	TokenOpenSexpr        // OPEN_SEXPR: <mu>"("
	TokenCloseSexpr       // CLOSE_SEXPR: <mu>")"
	TokenEquals           // EQUALS: <mu>"="
	TokenData             // DATA: <mu>"@"
	TokenSep              // SEP: <mu>[\/.]
	TokenOpenBlockParams  // OPEN_BLOCK_PARAMS: <mu>"as"\s+"|"
	TokenCloseBlockParams // CLOSE_BLOCK_PARAMS <mu>"|"

	// tokens with content
	TokenContent // begin 'mu', begin 'emu', CONTENT: [^\x00]*?/("{{") - 02. CONTENT: [^\x00]+ - 03. CONTENT: <emu>[^\x00]{2,}?/("{{"|"\\{{"|"\\\\{{"|<<EOF>>) - 05: CONTENT: <raw>[^\x00]*?/("{{{{/")
	TokenID      // ID: <mu>".." - 25. ID: <mu>"."/{LOOKAHEAD} - 39. ID: <mu>{ID} - 40. ID: <mu>'['[^\]]*']'
	TokenString  // STRING: <mu>'"'("\\"["]|[^"])*'"' - 30. STRING: <mu>"'"("\\"[']|[^'])*"'"
	TokenNumber  // NUMBER: <mu>\-?[0-9]+(?:\.[0-9]+)?/{LITERAL_LOOKAHEAD}
	TokenBoolean // BOOLEAN: <mu>"true"/{LITERAL_LOOKAHEAD} - 33. BOOLEAN: <mu>"false"/{LITERAL_LOOKAHEAD}
)

const (
	// mustache detection
	ESCAPED_OPEN_MUSTACHE = "\\{{"
	OPEN_MUSTACHE         = "{{"
	CLOSE_MUSTACHE        = "}}"
	CLOSE_STRIP_MUSTACHE  = "~}}"
)

const eof = -1

type TokenKind int

type Token struct {
	kind TokenKind // Token kind
	pos  int       // Position in input string
	val  string    // Token value
}

// function that returns the next lexer function
type lexFunc func(*Lexer) lexFunc

// Lexical analyzer
type Lexer struct {
	input    string     // input to scan
	name     string     // lexer name, used for testing purpose
	tokens   chan Token // channel of scanned tokens
	nextFunc lexFunc    // the next function to execute

	pos   int // current scan position in input string
	width int // size of last rune scanned from input string
	start int // start position of the token we are scanning

	closeComment *regexp.Regexp // regexp to scan close of current comment
}

var (
	// characters not allowed in an identifier
	unallowedIDChars = " \n\t!\"#%&'()*+,./;<=>@[\\]^`{|}~"

	// regular expressions
	rOpenRaw        = regexp.MustCompile(`^{{{{`)
	rCloseRaw       = regexp.MustCompile(`^}}}}`)
	rOpenUnescaped  = regexp.MustCompile(`^{{~?{`)
	rCloseUnescaped = regexp.MustCompile(`^}~?}}`)
	rOpenBlock      = regexp.MustCompile(`^{{~?#`)
	rOpenEndBlock   = regexp.MustCompile(`^{{~?/`)
	rOpenPartial    = regexp.MustCompile(`^{{~?>`)
	// {{^}} or {{else}}
	rInverse          = regexp.MustCompile(`^({{~?\^\s*~?}}|{{~?\s*else\s*~?}})`)
	rOpenInverse      = regexp.MustCompile(`^{{~?\^`)
	rOpenInverseChain = regexp.MustCompile(`^{{~?\s*else`)
	// {{ or {{&
	rOpen            = regexp.MustCompile(`^{{~?&?`)
	rClose           = regexp.MustCompile(`^~?}}`)
	rOpenBlockParams = regexp.MustCompile(`^as\s+\|`)
	// {{!--  ... --}}
	rOpenCommentDash  = regexp.MustCompile(`^{{~?!--\s*`)
	rCloseCommentDash = regexp.MustCompile(`^\s*--~?}}`)
	// {{! ... }}
	rOpenComment  = regexp.MustCompile(`^{{~?!\s*`)
	rCloseComment = regexp.MustCompile(`^\s*~?}}`)

	rID    = regexp.MustCompile(`^[^` + regexp.QuoteMeta(unallowedIDChars) + `]+`)
	rDotID = regexp.MustCompile(`^\.[\s` + regexp.QuoteMeta("=~}/)|") + `]`)
)

// scans given input
func Scan(input string, name string) *Lexer {
	result := &Lexer{
		input:  input,
		name:   name,
		tokens: make(chan Token),
	}

	go result.run()

	return result
}

// returns the next scanned token
func (l *Lexer) NextToken() Token {
	result := <-l.tokens

	return result
}

// starts lexical analysis
func (l *Lexer) run() {
	for l.nextFunc = lexContent; l.nextFunc != nil; {
		l.nextFunc = l.nextFunc(l)
	}
}

// returns next character from input, or eof of there is nothing left to scan
func (l *Lexer) next() rune {
	if l.pos >= len(l.input) {
		l.width = 0
		return eof
	}

	r, w := utf8.DecodeRuneInString(l.input[l.pos:])
	l.width = w
	l.pos += l.width

	return r
}

// emits a new scanned token
func (l *Lexer) emit(kind TokenKind) {
	l.tokens <- Token{kind, l.start, l.input[l.start:l.pos]}

	// scanning a new token
	l.start = l.pos
}

// emits a scanned string
func (l *Lexer) emitString(delimiter rune) {
	str := l.input[l.start:l.pos]

	// replace escaped delimiters
	str = strings.Replace(str, "\\"+string(delimiter), string(delimiter), -1)

	l.tokens <- Token{TokenString, l.start, str}

	// scanning a new token
	l.start = l.pos
}

// returns but does not consume the next character in the input
func (l *Lexer) peek() rune {
	r := l.next()
	l.backup()
	return r
}

// steps back one character
// @warning Can only be called once per call of next
func (l *Lexer) backup() {
	l.pos -= l.width
}

// skips all characters that have been scanned up to current position
func (l *Lexer) ignore() {
	l.start = l.pos
}

// scans the next character if it is included in given string
func (l *Lexer) accept(valid string) bool {
	if strings.IndexRune(valid, l.next()) >= 0 {
		return true
	}

	l.backup()

	return false
}

// scans all following characters that are part of given string
func (l *Lexer) acceptRun(valid string) {
	for strings.IndexRune(valid, l.next()) >= 0 {
	}

	l.backup()
}

func (l *Lexer) errorf(format string, args ...interface{}) lexFunc {
	l.tokens <- Token{TokenError, l.start, fmt.Sprintf(format, args...)}
	return nil
}

// returns true if content at current scanning position starts with given string
func (l *Lexer) isString(str string) bool {
	return strings.HasPrefix(l.input[l.pos:], str)
}

// returns the first string from current scanning position that matches given regular expression
func (l *Lexer) findRegexp(r *regexp.Regexp) string {
	return r.FindString(l.input[l.pos:])
}

// scanning content (ie: not between mustaches)
func lexContent(l *Lexer) lexFunc {
	var next lexFunc

	// find opening mustaches
	if l.isString(ESCAPED_OPEN_MUSTACHE) {
		// check \\{{
		l.backup()
		if r := l.next(); r != '\\' {
			// \{{
			next = lexEscapedOpenMustache
		}
	} else if str := l.findRegexp(rOpenCommentDash); str != "" {
		// {{!--
		l.closeComment = rCloseCommentDash

		next = lexComment
	} else if str := l.findRegexp(rOpenComment); str != "" {
		// {{!
		l.closeComment = rCloseComment

		next = lexComment
	} else if l.isString(OPEN_MUSTACHE) {
		// {{
		next = lexOpenMustache
	}

	if next != nil {
		// emit scanned content
		if l.pos > l.start {
			l.emit(TokenContent)
		}

		// scan next token
		return next
	}

	// scan next rune
	if l.next() == eof {
		// emit scanned content
		if l.pos > l.start {
			l.emit(TokenContent)
		}

		// this is over
		l.emit(TokenEOF)
		return nil
	}

	// continue content scanning
	return lexContent
}

// scanning \{{
func lexEscapedOpenMustache(l *Lexer) lexFunc {
	// ignore escape character
	l.next()
	l.ignore()

	// scan mustaches
	for l.peek() == '{' {
		l.next()
	}

	return lexContent
}

// scanning {{
func lexOpenMustache(l *Lexer) lexFunc {
	var str string
	var tok TokenKind

	nextFunc := lexExpression

	if str = l.findRegexp(rOpenRaw); str != "" {
		tok = TokenOpenRawBlock
	} else if str = l.findRegexp(rOpenUnescaped); str != "" {
		tok = TokenOpenUnescaped
	} else if str = l.findRegexp(rOpenBlock); str != "" {
		tok = TokenOpenBlock
	} else if str = l.findRegexp(rOpenEndBlock); str != "" {
		tok = TokenOpenEndBlock
	} else if str = l.findRegexp(rOpenPartial); str != "" {
		tok = TokenOpenPartial
	} else if str = l.findRegexp(rInverse); str != "" {
		tok = TokenInverse
		nextFunc = lexContent
	} else if str = l.findRegexp(rOpenInverse); str != "" {
		tok = TokenOpenInverse
	} else if str = l.findRegexp(rOpenInverseChain); str != "" {
		tok = TokenOpenInverseChain
	} else if str = l.findRegexp(rOpen); str != "" {
		tok = TokenOpen
	} else {
		// this is rotten
		panic("Current pos MUST be an opening mustache")
	}

	l.pos += len(str)
	l.emit(tok)

	return nextFunc
}

// scanning }} or ~}}
func lexCloseMustache(l *Lexer) lexFunc {
	var str string
	var tok TokenKind

	if str = l.findRegexp(rCloseRaw); str != "" {
		tok = TokenCloseRawBlock
	} else if str = l.findRegexp(rCloseUnescaped); str != "" {
		tok = TokenCloseUnescaped
	} else if str = l.findRegexp(rClose); str != "" {
		tok = TokenClose
	} else {
		// this is rotten
		panic("Current pos MUST be a closing mustache")
	}

	l.pos += len(str)
	l.emit(tok)

	return lexContent
}

// scanning inside mustaches
func lexExpression(l *Lexer) lexFunc {
	// search close mustache delimiter
	if l.isString(CLOSE_MUSTACHE) || l.isString(CLOSE_STRIP_MUSTACHE) {
		if l.pos > l.start {
			// emit scanned content
			l.emit(TokenContent)
		}

		return lexCloseMustache
	}

	// search some patterns before advancing scanning position

	// "as |"
	if str := l.findRegexp(rOpenBlockParams); str != "" {
		l.pos += len(str)
		l.emit(TokenOpenBlockParams)
		return lexExpression
	}

	// ..
	if l.isString("..") {
		l.pos += len("..")
		l.emit(TokenID)
		return lexExpression
	}

	// .
	if str := l.findRegexp(rDotID); str != "" {
		l.pos += len(".")
		l.emit(TokenID)
		return lexExpression
	}

	// true
	if l.isString("true") {
		l.pos += len("true")
		l.emit(TokenBoolean)
		return lexExpression
	}

	// false
	if l.isString("false") {
		l.pos += len("false")
		l.emit(TokenBoolean)
		return lexExpression
	}

	// let's scan next character
	switch r := l.next(); {
	case r == eof:
		return l.errorf("Unclosed expression")
	case isIgnorable(r):
		return lexIgnorable
	case r == '(':
		l.emit(TokenOpenSexpr)
	case r == ')':
		l.emit(TokenCloseSexpr)
	case r == '=':
		l.emit(TokenEquals)
	case r == '@':
		l.emit(TokenData)
	case r == '"' || r == '\'':
		l.backup()
		return lexString
	case r == '/' || r == '.':
		l.emit(TokenSep)
	case r == '|':
		l.emit(TokenCloseBlockParams)
	case r == '+' || r == '-' || (r >= '0' && r <= '9'):
		l.backup()
		return lexNumber
	case r == '[':
		return lexPathLiteral
	case strings.IndexRune(unallowedIDChars, r) < 0:
		l.backup()
		return lexIdentifier
	default:
		return l.errorf("Unexpected character in expression: %#U", r)
	}

	return lexExpression
}

// scanning {{!-- or {{!
func lexComment(l *Lexer) lexFunc {
	if str := l.findRegexp(l.closeComment); str != "" {
		l.pos += len(str)
		l.emit(TokenComment)

		return lexContent
	}

	if r := l.next(); r == eof {
		return l.errorf("Unclosed comment")
	}

	return lexComment
}

// scans all following ignorable characters
func lexIgnorable(l *Lexer) lexFunc {
	for isIgnorable(l.peek()) {
		l.next()
	}
	l.ignore()

	return lexExpression
}

func lexString(l *Lexer) lexFunc {
	// get string delimiter
	delim := l.next()
	var prev rune = 0

	// ignore delimiter
	l.ignore()

	for {
		r := l.next()
		if r == eof || r == '\n' {
			return l.errorf("Unterminated string")
		}

		if (r == delim) && (prev != '\\') {
			break
		}

		prev = r
	}

	// remove end delimiter
	l.backup()

	// emit string
	l.emitString(delim)

	// skip end delimiter
	l.next()
	l.ignore()

	return lexExpression
}

// lexNumber scans a number: decimal, octal, hex, float, or imaginary. This
// isn't a perfect number scanner - for instance it accepts "." and "0x0.2"
// and "089" - but when it's wrong the input is invalid and the parser (via
// strconv) will notice.
//
// @note borrowed from https://github.com/golang/go/tree/master/src/text/template/parse/lex.go
func lexNumber(l *Lexer) lexFunc {
	if !l.scanNumber() {
		return l.errorf("bad number syntax: %q", l.input[l.start:l.pos])
	}
	if sign := l.peek(); sign == '+' || sign == '-' {
		// Complex: 1+2i. No spaces, must end in 'i'.
		if !l.scanNumber() || l.input[l.pos-1] != 'i' {
			return l.errorf("bad number syntax: %q", l.input[l.start:l.pos])
		}
		l.emit(TokenNumber)
	} else {
		l.emit(TokenNumber)
	}
	return lexExpression
}

// @note borrowed from https://github.com/golang/go/tree/master/src/text/template/parse/lex.go
func (l *Lexer) scanNumber() bool {
	// Optional leading sign.
	l.accept("+-")

	// Is it hex?
	digits := "0123456789"

	if l.accept("0") && l.accept("xX") {
		digits = "0123456789abcdefABCDEF"
	}

	l.acceptRun(digits)

	if l.accept(".") {
		l.acceptRun(digits)
	}

	if l.accept("eE") {
		l.accept("+-")
		l.acceptRun("0123456789")
	}

	// Is it imaginary?
	l.accept("i")

	// Next thing mustn't be alphanumeric.
	if isAlphaNumeric(l.peek()) {
		l.next()
		return false
	}

	return true
}

// scans an ID
func lexIdentifier(l *Lexer) lexFunc {
	str := l.findRegexp(rID)
	if len(str) == 0 {
		// this is rotten
		panic("Identifier expected")
	}

	l.pos += len(str)
	l.emit(TokenID)

	return lexExpression
}

// scans an [ID]
func lexPathLiteral(l *Lexer) lexFunc {
	for {
		r := l.next()
		if r == eof || r == '\n' {
			return l.errorf("Unterminated path literal")
		}

		if r == ']' {
			break
		}
	}

	l.emit(TokenID)

	return lexExpression
}

// returns true if given character is ignorable (ie. whitespace of line feed)
func isIgnorable(r rune) bool {
	return r == ' ' || r == '\t' || r == '\n'
}

// isAlphaNumeric reports whether r is an alphabetic, digit, or underscore.
// @note borrowed from https://github.com/golang/go/tree/master/src/text/template/parse/lex.go
func isAlphaNumeric(r rune) bool {
	return r == '_' || unicode.IsLetter(r) || unicode.IsDigit(r)
}

package lexer

import (
	"fmt"
	"regexp"
	"strings"
	"unicode/utf8"
)

// Reference: https://github.com/wycats/handlebars.js/blob/master/src/handlebars.l

const (
	TokenError TokenKind = iota
	TokenEOF

	// mustache delimiters
	TokenOpen             // [x] 19. OPEN: <mu>"{{"{LEFT_STRIP}?"&" - 22. OPEN: <mu>"{{"{LEFT_STRIP}?
	TokenClose            // [x] 28. CLOSE: <mu>{RIGHT_STRIP}?"}}"
	TokenOpenRawBlock     // [x] 09. OPEN_RAW_BLOCK: <mu>"{{{{"
	TokenCloseRawBlock    // [x] 10. CLOSE_RAW_BLOCK: <mu>"}}}}"
	TokenOpenUnescaped    // [x] 18. OPEN_UNESCAPED: <mu>"{{"{LEFT_STRIP}?"{"
	TokenCloseUnescaped   // [x] 27. CLOSE_UNESCAPED: <mu>"}"{RIGHT_STRIP}?"}}"
	TokenOpenBlock        // [x] 12. OPEN_BLOCK: <mu>"{{"{LEFT_STRIP}?"#"
	TokenOpenEndBlock     // [x] 13. OPEN_ENDBLOCK: <mu>"{{"{LEFT_STRIP}?"/"
	TokenInverse          // [x] 14. INVERSE: <mu>"{{"{LEFT_STRIP}?"^"\s*{RIGHT_STRIP}?"}}" - 15. INVERSE: <mu>"{{"{LEFT_STRIP}?\s*"else"\s*{RIGHT_STRIP}?"}}"
	TokenOpenInverse      // [x] 16. OPEN_INVERSE: <mu>"{{"{LEFT_STRIP}?"^"
	TokenOpenInverseChain // [x] 17. OPEN_INVERSE_CHAIN: <mu>"{{"{LEFT_STRIP}?\s*"else"
	TokenOpenPartial      // [x] 11. OPEN_PARTIAL: <mu>"{{"{LEFT_STRIP}?">"
	TokenEndRawBlock      // [ ] 04. END_RAW_BLOCK: <raw>"{{{{/"[^\s!"#%-,\.\/;->@\[-\^`\{-~]+/[=}\s\/.]"}}}}"
	TokenComment          // [ ] 06. COMMENT: <com>[\s\S]*?"--"{RIGHT_STRIP}?"}}" - 20. begin 'com': <mu>"{{"{LEFT_STRIP}?"!--" - 21. COMMENT: <mu>"{{"{LEFT_STRIP}?"!"[\s\S]*?"}}"

	// inside mustaches
	TokenOpenSexpr        // [x] 07. OPEN_SEXPR: <mu>"("
	TokenCloseSexpr       // [x] 08. CLOSE_SEXPR: <mu>")"
	TokenEquals           // [x] 23. EQUALS: <mu>"="
	TokenData             // [x] 31. DATA: <mu>"@"
	TokenSep              // [x] 26. SEP: <mu>[\/.]
	TokenOpenBlockParams  // [x] 37. OPEN_BLOCK_PARAMS: <mu>"as"\s+"|"
	TokenCloseBlockParams // [x] 38. CLOSE_BLOCK_PARAMS <mu>"|"
	// TokenUndefined         // [ ] 34. UNDEFINED: <mu>"undefined"/{LITERAL_LOOKAHEAD}
	// TokenNull              // [ ] 35. NULL: <mu>"null"/{LITERAL_LOOKAHEAD}

	// tokens with content
	TokenContent // [ ] 01. begin 'mu', begin 'emu', CONTENT: [^\x00]*?/("{{") - 02. CONTENT: [^\x00]+ - 03. CONTENT: <emu>[^\x00]{2,}?/("{{"|"\\{{"|"\\\\{{"|<<EOF>>) - 05: CONTENT: <raw>[^\x00]*?/("{{{{/")
	TokenID      // [x] 24. ID: <mu>".." - 25. ID: <mu>"."/{LOOKAHEAD} - 39. ID: <mu>{ID} - 40. ID: <mu>'['[^\]]*']'
	TokenString  // [x] 29. STRING: <mu>'"'("\\"["]|[^"])*'"' - 30. STRING: <mu>"'"("\\"[']|[^'])*"'"
	TokenNumber  // [ ] 36. NUMBER: <mu>\-?[0-9]+(?:\.[0-9]+)?/{LITERAL_LOOKAHEAD}
	TokenBoolean // [x] 32. BOOLEAN: <mu>"true"/{LITERAL_LOOKAHEAD} - 33. BOOLEAN: <mu>"false"/{LITERAL_LOOKAHEAD}
)

const (
	// mustache detection
	OPEN_MUSTACHE        = "{{"
	CLOSE_MUSTACHE       = "}}"
	CLOSE_STRIP_MUSTACHE = "~}}"
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
}

var (
	// characters not allowed in an identifier
	unallowedIDChars = " \t!\"#%&'()*+,./;<=>@[\\]^`{|}~"

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
	rOpen  = regexp.MustCompile(`^{{~?&?`)
	rClose = regexp.MustCompile(`^~?}}`)

	rOpenBlockParams = regexp.MustCompile(`^as\s+\|`)

	rID = regexp.MustCompile(`^[^` + regexp.QuoteMeta(unallowedIDChars) + `]+`)
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

func (l *Lexer) errorf(format string, args ...interface{}) lexFunc {
	l.tokens <- Token{TokenError, l.start, fmt.Sprintf(format, args...)}
	return nil
}

// returns true if content at current scanning position starts with given string
func (l *Lexer) isString(str string) bool {
	return strings.HasPrefix(l.input[l.pos:], str)
}

// returns true if current scanning position matches given regular expression
func (l *Lexer) matchRegexp(r *regexp.Regexp) bool {
	return r.MatchString(l.input[l.pos:])
}

// returns the first string from current scanning position that matches given regular expression
func (l *Lexer) findRegexp(r *regexp.Regexp) string {
	return r.FindString(l.input[l.pos:])
}

// scanning content (ie: not between mustaches)
func lexContent(l *Lexer) lexFunc {
	for {
		// search open mustache delimiter
		if l.isString(OPEN_MUSTACHE) {
			if l.pos > l.start {
				// emit scanned content
				l.emit(TokenContent)
			}

			return lexOpenMustache
		}

		// scan next rune
		if l.next() == eof {
			break
		}
	}

	// emit scanned content
	if l.pos > l.start {
		l.emit(TokenContent)
	}

	l.emit(TokenEOF)

	return nil
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
	if l.matchRegexp(rOpenBlockParams) {
		l.pos += len(l.findRegexp(rOpenBlockParams))
		l.emit(TokenOpenBlockParams)
		return lexExpression
	}

	if l.isString("true") {
		l.pos += len("true")
		l.emit(TokenBoolean)
		return lexExpression
	}

	if l.isString("false") {
		l.pos += len("false")
		l.emit(TokenBoolean)
		return lexExpression
	}

	// ok, let's scan next character
	switch r := l.next(); {
	case r == eof:
		return l.errorf("Unclosed expression")
	case isSpace(r):
		return lexSpaces
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
	case strings.IndexRune(unallowedIDChars, r) < 0:
		l.backup()
		return lexIdentifier
	default:
		return l.errorf("Unexpected character in expression: %#U", r)
	}

	return lexExpression
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

// scans all following space characters and ignore them
func lexSpaces(l *Lexer) lexFunc {
	for isSpace(l.peek()) {
		l.next()
	}
	l.ignore()

	return lexExpression
}

// scans a quoted string
func lexString(l *Lexer) lexFunc {
	// get string delimiter
	delim := l.next()

	// ignore delimiter
	l.ignore()

Loop:
	for {
		switch l.next() {
		case '\\':
			if r := l.next(); r != eof && r != '\n' {
				break
			}
			fallthrough
		case eof, '\n':
			return l.errorf("Unterminated string")
		case delim:
			break Loop
		}
	}

	// remove end delimiter
	l.backup()

	// emit string
	l.emit(TokenString)

	// skip end delimiter
	l.next()
	l.ignore()

	return lexExpression
}

// returns true if given character is a whitespace
func isSpace(r rune) bool {
	return r == ' ' || r == '\t'
}

package lexer

import (
	"fmt"
	"strings"
	"unicode/utf8"
)

// Reference: https://github.com/wycats/handlebars.js/blob/master/src/handlebars.l

type TokenKind int

const (
	TokenError TokenKind = iota
	TokenEOF

	TokenOpen              // 19. OPEN: <mu>"{{"{LEFT_STRIP}?"&" - 22. OPEN: <mu>"{{"{LEFT_STRIP}?
	TokenClose             // 28. CLOSE: <mu>{RIGHT_STRIP}?"}}"
	TokenOpenRawBlock      // 09. OPEN_RAW_BLOCK: <mu>"{{{{"
	TokenCloseRawBlock     // 10. CLOSE_RAW_BLOCK: <mu>"}}}}"
	TokenOpenUnescaped     // 18. OPEN_UNESCAPED: <mu>"{{"{LEFT_STRIP}?"{"
	TokenCloseUnescaped    // 27. CLOSE_UNESCAPED: <mu>"}"{RIGHT_STRIP}?"}}"
	TokenOpenBlock         // 12. OPEN_BLOCK: <mu>"{{"{LEFT_STRIP}?"#"
	TokenOpenEndBlock      // 13. OPEN_ENDBLOCK: <mu>"{{"{LEFT_STRIP}?"/"
	TokenOpenSexpr         // 07. OPEN_SEXPR: <mu>"("
	TokenCloseSexpr        // 08. CLOSE_SEXPR: <mu>")"
	TokenInverse           // 14. INVERSE: <mu>"{{"{LEFT_STRIP}?"^"\s*{RIGHT_STRIP}?"}}" - 15. INVERSE: <mu>"{{"{LEFT_STRIP}?\s*"else"\s*{RIGHT_STRIP}?"}}"
	TokenOpenInverse       // 16. OPEN_INVERSE: <mu>"{{"{LEFT_STRIP}?"^"
	TokenOpenInverseChain  // 17. OPEN_INVERSE_CHAIN: <mu>"{{"{LEFT_STRIP}?\s*"else"
	TokenOpenPartial       // 11. OPEN_PARTIAL: <mu>"{{"{LEFT_STRIP}?">"
	TokenEndRawBlock       // 04. END_RAW_BLOCK: <raw>"{{{{/"[^\s!"#%-,\.\/;->@\[-\^`\{-~]+/[=}\s\/.]"}}}}"
	TokenOpenBlockParams   // 37. OPEN_BLOCK_PARAMS: <mu>"as"\s+"|"
	TokenCloseBlockPaarams // 38. CLOSE_BLOCK_PARAMS <mu>"|"
	TokenEquals            // 23. EQUALS: <mu>"="
	TokenData              // 31. DATA: <mu>"@"
	TokenSep               // 26. SEP: <mu>[\/.]
	TokenUndefined         // 34. UNDEFINED: <mu>"undefined"/{LITERAL_LOOKAHEAD}
	TokenNull              // 35. NULL: <mu>"null"/{LITERAL_LOOKAHEAD}

	// tokens with content
	TokenContent // 01. begin 'mu', begin 'emu', CONTENT: [^\x00]*?/("{{") - 02. CONTENT: [^\x00]+ - 03. CONTENT: <emu>[^\x00]{2,}?/("{{"|"\\{{"|"\\\\{{"|<<EOF>>) - 05: CONTENT: <raw>[^\x00]*?/("{{{{/")
	TokenComment // 06. COMMENT: <com>[\s\S]*?"--"{RIGHT_STRIP}?"}}" - 20. begin 'com': <mu>"{{"{LEFT_STRIP}?"!--" - 21. COMMENT: <mu>"{{"{LEFT_STRIP}?"!"[\s\S]*?"}}"
	TokenID      // 24. ID: <mu>".." - 25. ID: <mu>"."/{LOOKAHEAD} - 39. ID: <mu>{ID} - 40. ID: <mu>'['[^\]]*']'
	TokenString  // 29. STRING: <mu>'"'("\\"["]|[^"])*'"' - 30. STRING: <mu>"'"("\\"[']|[^'])*"'"
	TokenNumber  // 36. NUMBER: <mu>\-?[0-9]+(?:\.[0-9]+)?/{LITERAL_LOOKAHEAD}
	TokenBoolean // 32. BOOLEAN: <mu>"true"/{LITERAL_LOOKAHEAD} - 33. BOOLEAN: <mu>"false"/{LITERAL_LOOKAHEAD}
)

const (
	OPEN_MUSTACHE  = "{{"
	CLOSE_MUSTACHE = "}}"
)

type Token struct {
	kind TokenKind // Token kind
	pos  int       // Position in input string
	val  string    // Token value
}

const eof = -1

// function that returns the next lexer state
type stateFn func(*Lexer) stateFn

// Lexical analyzer
type Lexer struct {
	input  string     // input to scan
	name   string     // lexer name, used for testing purpose
	tokens chan Token // channel of scanned tokens
	state  stateFn    // the next function to execute

	pos     int // current scan position in input string
	width   int // size of last rune scanned from input string
	start   int // start position of the token we are scanning
	lastPos int // position of last token retrieved by consummer
}

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
	l.lastPos = result.pos

	return result
}

// runs lexical analysis
func (l *Lexer) run() {
	for l.state = lexContent; l.state != nil; {
		l.state = l.state(l)
	}
}

// returns next rune from input, or eof of there is nothing left to scan
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

func (l *Lexer) emit(kind TokenKind) {
	l.tokens <- Token{kind, l.start, l.input[l.start:l.pos]}

	// starting next token
	l.start = l.pos
}

func (l *Lexer) errorf(format string, args ...interface{}) stateFn {
	l.tokens <- Token{TokenError, l.start, fmt.Sprintf(format, args...)}
	return nil
}

func lexContent(l *Lexer) stateFn {
	for {
		// search open mustache delimiter
		if strings.HasPrefix(l.input[l.pos:], OPEN_MUSTACHE) {
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

// Scanning {{
func lexOpenMustache(l *Lexer) stateFn {
	l.pos += len(OPEN_MUSTACHE)
	l.emit(TokenOpen)

	return lexInsideMustache
}

// Scanning inside {{ ... }}
func lexInsideMustache(l *Lexer) stateFn {
	// @todo !!!
	return l.errorf("NOT IMPLEMENTED")
}

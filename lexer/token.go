package lexer

import "fmt"

const (
	TokenError TokenKind = iota
	TokenEOF

	// mustache delimiters
	TokenOpen             // OPEN
	TokenClose            // CLOSE
	TokenOpenRawBlock     // OPEN_RAW_BLOCK
	TokenCloseRawBlock    // CLOSE_RAW_BLOCK
	TokenOpenEndRawBlock  // END_RAW_BLOCK
	TokenOpenUnescaped    // OPEN_UNESCAPED
	TokenCloseUnescaped   // CLOSE_UNESCAPED
	TokenOpenBlock        // OPEN_BLOCK
	TokenOpenEndBlock     // OPEN_ENDBLOCK
	TokenInverse          // INVERSE
	TokenOpenInverse      // OPEN_INVERSE
	TokenOpenInverseChain // OPEN_INVERSE_CHAIN
	TokenOpenPartial      // OPEN_PARTIAL
	TokenComment          // COMMENT

	// inside mustaches
	TokenOpenSexpr        // OPEN_SEXPR
	TokenCloseSexpr       // CLOSE_SEXPR
	TokenEquals           // EQUALS
	TokenData             // DATA
	TokenSep              // SEP
	TokenOpenBlockParams  // OPEN_BLOCK_PARAMS
	TokenCloseBlockParams // CLOSE_BLOCK_PARAMS

	// tokens with content
	TokenContent // CONTENT
	TokenID      // ID
	TokenString  // STRING
	TokenNumber  // NUMBER
	TokenBoolean // BOOLEAN
)

const (
	// Option to generate token position in its string representation
	DUMP_TOKEN_POS = false

	// Option to generate values for all token kinds for their string representations
	DUMP_ALL_TOKENS_VAL = true
)

// TokenKind represents a Token type.
type TokenKind int

// Token represents a scanned token.
type Token struct {
	Kind TokenKind // Token kind
	Val  string    // Token value

	Pos  int // Byte position in input string
	Line int // Line number in input string
}

// tokenName permits to display token name given token type
var tokenName = map[TokenKind]string{
	TokenError:            "Error",
	TokenEOF:              "EOF",
	TokenContent:          "Content",
	TokenComment:          "Comment",
	TokenOpen:             "Open",
	TokenClose:            "Close",
	TokenOpenUnescaped:    "OpenUnescaped",
	TokenCloseUnescaped:   "CloseUnescaped",
	TokenOpenBlock:        "OpenBlock",
	TokenOpenEndBlock:     "OpenEndBlock",
	TokenOpenRawBlock:     "OpenRawBlock",
	TokenCloseRawBlock:    "CloseRawBlock",
	TokenOpenEndRawBlock:  "OpenEndRawBlock",
	TokenOpenBlockParams:  "OpenBlockParams",
	TokenCloseBlockParams: "CloseBlockParams",
	TokenInverse:          "Inverse",
	TokenOpenInverse:      "OpenInverse",
	TokenOpenInverseChain: "OpenInverseChain",
	TokenOpenPartial:      "OpenPartial",
	TokenOpenSexpr:        "OpenSexpr",
	TokenCloseSexpr:       "CloseSexpr",
	TokenID:               "ID",
	TokenEquals:           "Equals",
	TokenString:           "String",
	TokenNumber:           "Number",
	TokenBoolean:          "Boolean",
	TokenData:             "Data",
	TokenSep:              "Sep",
}

// String returns the token kind string representation for debugging.
func (k TokenKind) String() string {
	s := tokenName[k]
	if s == "" {
		return fmt.Sprintf("Token-%d", int(k))
	}
	return s
}

// String returns the token string representation for debugging.
func (t Token) String() string {
	result := ""

	if DUMP_TOKEN_POS {
		result += fmt.Sprintf("%d:", t.Pos)
	}

	result += fmt.Sprintf("%s", t.Kind)

	if (DUMP_ALL_TOKENS_VAL || (t.Kind >= TokenContent)) && len(t.Val) > 0 {
		if len(t.Val) > 100 {
			result += fmt.Sprintf("{%.20q...}", t.Val)
		} else {
			result += fmt.Sprintf("{%q}", t.Val)
		}
	}

	return result
}

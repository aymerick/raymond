package lexer

import "fmt"

const (
	TokenError TokenKind = iota
	TokenEOF

	// mustache delimiters
	TokenOpen             // OPEN: <mu>"{{"{LEFT_STRIP}?"&" - 22. OPEN: <mu>"{{"{LEFT_STRIP}?
	TokenClose            // CLOSE: <mu>{RIGHT_STRIP}?"}}"
	TokenOpenRawBlock     // OPEN_RAW_BLOCK: <mu>"{{{{"
	TokenCloseRawBlock    // CLOSE_RAW_BLOCK: <mu>"}}}}"
	TokenOpenEndRawBlock  // END_RAW_BLOCK: <raw>"{{{{/"[^\s!"#%-,\.\/;->@\[-\^`\{-~]+/[=}\s\/.]"}}}}"
	TokenOpenUnescaped    // OPEN_UNESCAPED: <mu>"{{"{LEFT_STRIP}?"{"
	TokenCloseUnescaped   // CLOSE_UNESCAPED: <mu>"}"{RIGHT_STRIP}?"}}"
	TokenOpenBlock        // OPEN_BLOCK: <mu>"{{"{LEFT_STRIP}?"#"
	TokenOpenEndBlock     // OPEN_ENDBLOCK: <mu>"{{"{LEFT_STRIP}?"/"
	TokenInverse          // INVERSE: <mu>"{{"{LEFT_STRIP}?"^"\s*{RIGHT_STRIP}?"}}" - 15. INVERSE: <mu>"{{"{LEFT_STRIP}?\s*"else"\s*{RIGHT_STRIP}?"}}"
	TokenOpenInverse      // OPEN_INVERSE: <mu>"{{"{LEFT_STRIP}?"^"
	TokenOpenInverseChain // OPEN_INVERSE_CHAIN: <mu>"{{"{LEFT_STRIP}?\s*"else"
	TokenOpenPartial      // OPEN_PARTIAL: <mu>"{{"{LEFT_STRIP}?">"
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
	DUMP_TOKEN_POS      = false
	DUMP_ALL_TOKENS_VAL = true
)

type TokenKind int

type Token struct {
	Kind TokenKind // Token kind
	Pos  int       // Position in input string
	Val  string    // Token value
}

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

func (k TokenKind) String() string {
	s := tokenName[k]
	if s == "" {
		return fmt.Sprintf("Token-%d", int(k))
	}
	return s
}

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

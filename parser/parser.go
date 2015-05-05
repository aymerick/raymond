package parser

import (
	"errors"
	"fmt"
	"regexp"
	"strconv"
	"strings"

	"github.com/aymerick/raymond/ast"
	"github.com/aymerick/raymond/lexer"
)

// References:
//   - https://github.com/wycats/handlebars.js/blob/master/src/handlebars.yy
//   - https://github.com/golang/go/blob/master/src/text/template/parse/parse.go

// Grammar parser
type Parser struct {
	// Lexer
	lex *lexer.Lexer

	// Root node
	root ast.Node

	// Tokens parsed but not consumed yet
	tokens []*lexer.Token
}

var (
	rOpenComment  = regexp.MustCompile(`^\{\{~?!-?-?`)
	rCloseComment = regexp.MustCompile(`-?-?~?\}\}$`)
)

// instanciate a new parser
func New(input string) *Parser {
	return &Parser{
		lex: lexer.Scan(input),
	}
}

// parse given input and returns the ast root node
func Parse(input string) (ast.Node, error) {
	return New(input).ParseProgram()
}

// program : statement*
func (p *Parser) ParseProgram() (ast.Node, error) {
	result := ast.NewProgramNode(p.lex.Pos())

	for !p.over() {
		node, err := p.parseStatement()
		if err != nil {
			return nil, err
		}

		result.Statements = append(result.Statements, node)
	}

	return result, p.err()
}

// statement : mustache | block | rawBlock | partial | content | COMMENT
func (p *Parser) parseStatement() (ast.Node, error) {
	var result ast.Node
	var err error

	tok := p.next()

	switch tok.Kind {
	case lexer.TokenContent:
		result = p.parseContent()
	case lexer.TokenComment:
		result = p.parseComment()
	case lexer.TokenOpenRawBlock:
		result, err = p.parseRawBlock()
		if err != nil {
			return nil, err
		}
	default:
		return nil, errors.New(fmt.Sprintf("Failed to parse statement: %s", tok))
	}

	return result, p.err()
}

// content : CONTENT
func (p *Parser) parseContent() ast.Node {
	tok := p.shift()

	return ast.NewContentNode(tok.Pos, tok.Val)
}

// COMMENT
func (p *Parser) parseComment() ast.Node {
	tok := p.shift()

	value := rOpenComment.ReplaceAllString(tok.Val, "")
	value = rCloseComment.ReplaceAllString(value, "")

	return ast.NewCommentNode(tok.Pos, strings.TrimSpace(value))
}

// rawBlock : openRawBlock content endRawBlock
// openRawBlock : OPEN_RAW_BLOCK helperName param* hash? CLOSE_RAW_BLOCK
// endRawBlock : OPEN_EN_RAW_BLOCK helperName CLOSE_RAW_BLOCK
func (p *Parser) parseRawBlock() (ast.Node, error) {
	// OPEN_RAW_BLOCK
	p.shift()

	// // helperName
	// helper, err := p.parseHelperName()
	// if err != nil {
	// 	return nil, err
	// }

	// param*

	// hash?

	// CLOSE_RAW_BLOCK

	// content

	// OPEN_EN_RAW_BLOCK

	// helperName

	// CLOSE_RAW_BLOCK

	// @todo !!!
	return nil, errors.New("NOT IMPLEMENTED")
}

// block : openBlock program inverseChain? closeBlock
//       | openInverse program inverseAndProgram? closeBlock
func (p *Parser) parseBlock() (ast.Node, error) {
	// @todo !!!
	return nil, errors.New("NOT IMPLEMENTED")
}

// openBlock : OPEN_BLOCK helperName param* hash? blockParams? CLOSE
func (p *Parser) parseOpenBlock() (ast.Node, error) {
	// @todo !!!
	return nil, errors.New("NOT IMPLEMENTED")
}

// openInverse : OPEN_INVERSE helperName param* hash? blockParams? CLOSE
func (p *Parser) parseOpenInverse() (ast.Node, error) {
	// @todo !!!
	return nil, errors.New("NOT IMPLEMENTED")
}

// openInverseChain : OPEN_INVERSE_CHAIN helperName param* hash? blockParams? CLOSE
func (p *Parser) parseOpenInverseChain() (ast.Node, error) {
	// @todo !!!
	return nil, errors.New("NOT IMPLEMENTED")
}

// inverseAndProgram : INVERSE program
func (p *Parser) parseInverseAndProgram() (ast.Node, error) {
	// @todo !!!
	return nil, errors.New("NOT IMPLEMENTED")
}

// inverseChain : openInverseChain program inverseChain?
//              | inverseAndProgram
func (p *Parser) parseInverseChain() (ast.Node, error) {
	// @todo !!!
	return nil, errors.New("NOT IMPLEMENTED")
}

// closeBlock : OPEN_ENDBLOCK helperName CLOSE
func (p *Parser) parseCloseBlock() (ast.Node, error) {
	// @todo !!!
	return nil, errors.New("NOT IMPLEMENTED")
}

// mustache : OPEN helperName param* hash? CLOSE
//          | OPEN_UNESCAPED helperName param* hash? CLOSE_UNESCAPED
func (p *Parser) parseMustache() (ast.Node, error) {
	// @todo !!!
	return nil, errors.New("NOT IMPLEMENTED")
}

// partial : OPEN_PARTIAL partialName param* hash? CLOSE
func (p *Parser) parsePartial() (ast.Node, error) {
	// @todo !!!
	return nil, errors.New("NOT IMPLEMENTED")
}

// param : helperName
//       | sexpr
func (p *Parser) parseParam() (ast.Node, error) {
	// @todo !!!
	return nil, errors.New("NOT IMPLEMENTED")
}

// sexpr : OPEN_SEXPR helperName param* hash? CLOSE_SEXPR
func (p *Parser) parseSexpr() (ast.Node, error) {
	// @todo !!!
	return nil, errors.New("NOT IMPLEMENTED")
}

// hash : hashSegment+
func (p *Parser) parseHash() (ast.Node, error) {
	// @todo !!!
	return nil, errors.New("NOT IMPLEMENTED")
}

// hashSegment : ID EQUALS param
func (p *Parser) parseHashSegment() (ast.Node, error) {
	// @todo !!!
	return nil, errors.New("NOT IMPLEMENTED")
}

// blockParams : OPEN_BLOCK_PARAMS ID+ CLOSE_BLOCK_PARAMS
func (p *Parser) parseBlockParams() (ast.Node, error) {
	// @todo !!!
	return nil, errors.New("NOT IMPLEMENTED")
}

// helperName : path | dataName | STRING | NUMBER | BOOLEAN | UNDEFINED | NULL
func (p *Parser) parseHelperName() (ast.Node, error) {
	var result ast.Node
	var err error

	tok := p.next()

	switch tok.Kind {
	case lexer.TokenBoolean:
		p.shift()
		result = ast.NewBooleanNode(tok.Pos, (tok.Val == "true"))
	case lexer.TokenNumber:
		p.shift()
		val, err := strconv.Atoi(tok.Val)
		if err != nil {
			return nil, err
		}
		result = ast.NewNumberNode(tok.Pos, val)
	case lexer.TokenString:
		p.shift()
		result = ast.NewStringNode(tok.Pos, tok.Val)
	case lexer.TokenData:
		result, err = p.parseDataName()
		if err != nil {
			return nil, err
		}
	default:
		result, err = p.parsePath()
		if err != nil {
			return nil, err
		}
	}

	return result, p.err()
}

// partialName : helperName | sexpr
func (p *Parser) parsePartialName() (ast.Node, error) {
	// @todo !!!
	return nil, errors.New("NOT IMPLEMENTED")
}

// dataName : DATA pathSegments
func (p *Parser) parseDataName() (ast.Node, error) {
	// @todo !!!
	return nil, errors.New("NOT IMPLEMENTED")
}

// path : pathSegments
func (p *Parser) parsePath() (ast.Node, error) {
	// @todo !!!
	return nil, errors.New("NOT IMPLEMENTED")
}

// pathSegments : pathSegments SEP ID
//              | ID
func (p *Parser) parsePathSegments() (ast.Node, error) {
	// @todo !!!
	return nil, errors.New("NOT IMPLEMENTED")
}

// Ensure there is at least a token to parse
func (p *Parser) ensure() {
	if len(p.tokens) == 0 {
		// fetch next token
		tok := p.lex.NextToken()

		// queue it
		p.tokens = append(p.tokens, &tok)
	}
}

// Returns next token without removing it from tokens buffer
func (p *Parser) next() *lexer.Token {
	p.ensure()

	return p.tokens[0]
}

// Returns next token and remove it from the tokens buffer
func (p *Parser) shift() *lexer.Token {
	var result *lexer.Token

	p.ensure()

	result, p.tokens = p.tokens[0], p.tokens[1:]

	return result
}

// Returns true if parsing is over
func (p *Parser) over() bool {
	tok := p.next()
	return (tok.Kind == lexer.TokenEOF) || (tok.Kind == lexer.TokenError)
}

// Returns lexer error, or nil if no error
func (p *Parser) err() error {
	if token := p.next(); token.Kind == lexer.TokenError {
		return errors.New(fmt.Sprintf("Lexer error: %s", token.String()))
	} else {
		return nil
	}
}

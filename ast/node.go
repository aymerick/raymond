package ast

import (
	"bytes"
	"fmt"
)

// References:
//   - https://github.com/wycats/handlebars.js/blob/master/lib/handlebars/compiler/ast.js
//   - https://github.com/wycats/handlebars.js/blob/master/docs/compiler-api.md
//   - https://github.com/golang/go/blob/master/src/text/template/parse/node.go

// AST node interface
type Node interface {
	// node type
	Type() NodeType

	// location of node in original input string
	Location() Loc

	// accepts visitor
	Accept(Visitor)
}

// AST visitor interface
type Visitor interface {
	visitProgram(node *Program)

	// statements
	visitMustache(node *MustacheStatement)
	visitBlock(node *BlockStatement)
	visitPartial(node *PartialStatement)
	visitContent(node *ContentStatement)
	visitComment(node *CommentStatement)

	// expressions
	visitSubExpression(node *SubExpression)
	visitPath(node *PathExpression)

	// literals
	visitString(node *StringLiteral)
	visitBoolean(node *BooleanLiteral)
	visitNumber(node *NumberLiteral)

	// miscellaneous
	visitHash(node *Hash)
	visitHashPair(node *HashPair)
}

// NodeType
type NodeType int

func (t NodeType) Type() NodeType {
	return t
}

const (
	NodeProgram NodeType = iota

	// statements
	NodeMustache
	NodeBlock
	NodePartial
	NodeContent
	NodeComment

	// expressions
	NodeSubExpression
	NodePath

	// literals
	NodeBoolean
	NodeNumber
	NodeString

	// miscellaneous
	NodeHash
	NodeHashPair
)

// Location
type Loc struct {
	Pos  int // Byte position
	Line int // Line number
}

func (l Loc) Location() Loc {
	return l
}

//
// Program
//

type Program struct {
	NodeType
	Loc

	Body        []Node // [ Statement ... ]
	BlockParams []string
}

func NewProgram(pos int, line int) *Program {
	return &Program{
		NodeType: NodeProgram,
		Loc:      Loc{pos, line},
	}
}

func (node *Program) String() string {
	b := new(bytes.Buffer)

	for _, n := range node.Body {
		fmt.Fprint(b, n)
	}

	return b.String()
}

func (node *Program) Accept(visitor Visitor) {
	visitor.visitProgram(node)
}

func (node *Program) AddStatement(statement Node) {
	node.Body = append(node.Body, statement)
}

//
// Mustache Statement
//

type MustacheStatement struct {
	NodeType
	Loc

	Path   Node   // PathExpression
	Params []Node // [ Expression ... ]
	Hash   Node   // Hash
}

func NewMustacheStatement(pos int, line int) *MustacheStatement {
	return &MustacheStatement{
		NodeType: NodeMustache,
		Loc:      Loc{pos, line},
	}
}

func (node *MustacheStatement) Accept(visitor Visitor) {
	visitor.visitMustache(node)
}

//
// Block Statement
//

type BlockStatement struct {
	NodeType
	Loc

	Path    Node   // PathExpression
	Params  []Node // [ Expression ... ]
	Hash    Node   // Hash
	Program Node   // Program
	Inverse Node   // Program
}

func NewBlockStatement(pos int, line int) *BlockStatement {
	return &BlockStatement{
		NodeType: NodeBlock,
		Loc:      Loc{pos, line},
	}
}

func (node *BlockStatement) Accept(visitor Visitor) {
	visitor.visitBlock(node)
}

//
// Partial Statement
//

type PartialStatement struct {
	NodeType
	Loc

	Name   Node   // PathExpression | SubExpression
	Params []Node // [ Expression ... ]
	Hash   Node   // Hash
}

func NewPartialStatement(pos int, line int) *PartialStatement {
	return &PartialStatement{
		NodeType: NodePartial,
		Loc:      Loc{pos, line},
	}
}

func (node *PartialStatement) Accept(visitor Visitor) {
	visitor.visitPartial(node)
}

//
// Content Statement
//

type ContentStatement struct {
	NodeType
	Loc

	Value string
}

func NewContentStatement(pos int, line int, val string) *ContentStatement {
	return &ContentStatement{
		NodeType: NodeContent,
		Loc:      Loc{pos, line},

		Value: val,
	}
}

func (node *ContentStatement) Accept(visitor Visitor) {
	visitor.visitContent(node)
}

//
// Comment Statement
//

type CommentStatement struct {
	NodeType
	Loc

	Value string
}

func NewCommentStatement(pos int, line int, val string) *CommentStatement {
	return &CommentStatement{
		NodeType: NodeComment,
		Loc:      Loc{pos, line},

		Value: val,
	}
}

func (node *CommentStatement) Accept(visitor Visitor) {
	visitor.visitComment(node)
}

//
// SubExpression
//

type SubExpression struct {
	NodeType
	Loc

	Path   Node   // PathExpression
	Params []Node // [ Expression ... ]
	Hash   Node   // Hash
}

func NewSubExpression(pos int, line int) *SubExpression {
	return &SubExpression{
		NodeType: NodeSubExpression,
		Loc:      Loc{pos, line},
	}
}

func (node *SubExpression) Accept(visitor Visitor) {
	visitor.visitSubExpression(node)
}

//
// Path Expression
//

type PathExpression struct {
	NodeType
	Loc

	Original string
	Depth    int
	Parts    []string
	Data     bool
}

func NewPathExpression(pos int, line int, data bool) *PathExpression {
	result := &PathExpression{
		NodeType: NodePath,
		Loc:      Loc{pos, line},

		Data: data,
	}

	if data {
		result.Original = "@"
	}

	return result
}

func (node *PathExpression) Accept(visitor Visitor) {
	visitor.visitPath(node)
}

// Adds path part
func (node *PathExpression) Part(part string) {
	node.Original += part

	switch part {
	case "..":
		node.Depth += 1
	case ".", "this":
		// NOOP
	default:
		node.Parts = append(node.Parts, part)
	}
}

// Adds path separator
func (node *PathExpression) Sep(separator string) {
	node.Original += separator
}

//
// String Literal
//

type StringLiteral struct {
	NodeType
	Loc

	Value string
}

func NewStringLiteral(pos int, line int, val string) *StringLiteral {
	return &StringLiteral{
		NodeType: NodeString,
		Loc:      Loc{pos, line},

		Value: val,
	}
}

func (node *StringLiteral) Accept(visitor Visitor) {
	visitor.visitString(node)
}

//
// Boolean Literal
//

type BooleanLiteral struct {
	NodeType
	Loc

	Value    bool
	Original string
}

func NewBooleanLiteral(pos int, line int, val bool, original string) *BooleanLiteral {
	return &BooleanLiteral{
		NodeType: NodeBoolean,
		Loc:      Loc{pos, line},

		Value:    val,
		Original: original,
	}
}

func (node *BooleanLiteral) Accept(visitor Visitor) {
	visitor.visitBoolean(node)
}

func (node *BooleanLiteral) String() string {
	if node.Value {
		return "true"
	} else {
		return "false"
	}
}

//
// Number Literal
//

type NumberLiteral struct {
	NodeType
	Loc

	Value    int
	Original string
}

func NewNumberLiteral(pos int, line int, val int, original string) *NumberLiteral {
	return &NumberLiteral{
		NodeType: NodeNumber,
		Loc:      Loc{pos, line},

		Value:    val,
		Original: original,
	}
}

func (node *NumberLiteral) Accept(visitor Visitor) {
	visitor.visitNumber(node)
}

//
// Hash
//

type Hash struct {
	NodeType
	Loc

	Pairs []Node // [ HashPair ... ]
}

func NewHash(pos int, line int) *Hash {
	return &Hash{
		NodeType: NodeHash,
		Loc:      Loc{pos, line},
	}
}

func (node *Hash) Accept(visitor Visitor) {
	visitor.visitHash(node)
}

//
// HashPair
//

type HashPair struct {
	NodeType
	Loc

	Key string
	Val Node // Expression
}

func NewHashPair(pos int, line int) *HashPair {
	return &HashPair{
		NodeType: NodeHashPair,
		Loc:      Loc{pos, line},
	}
}

func (node *HashPair) Accept(visitor Visitor) {
	visitor.visitHashPair(node)
}

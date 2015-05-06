package ast

import (
	"bytes"
	"fmt"
)

// References:
//   - https://github.com/wycats/handlebars.js/blob/master/lib/handlebars/compiler/ast.js
//   - https://github.com/wycats/handlebars.js/blob/master/docs/compiler-api.md
//   - https://github.com/golang/go/blob/master/src/text/template/parse/node.go

type NodeType int

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

// AST node interface
type Node interface {
	Type() NodeType

	// byte position of start of node in full original input string
	Position() Pos

	// accepts visitor
	Accept(Visitor)
}

type Pos int

func (p Pos) Position() Pos {
	return p
}

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

//
// Program
//

type Program struct {
	NodeType
	Pos

	Statements []Node // [ Statement ... ]
}

func NewProgram(pos int) *Program {
	return &Program{
		NodeType: NodeProgram,
		Pos:      Pos(pos),
	}
}

func (node *Program) String() string {
	b := new(bytes.Buffer)

	for _, n := range node.Statements {
		fmt.Fprint(b, n)
	}

	return b.String()
}

func (node *Program) Accept(visitor Visitor) {
	visitor.visitProgram(node)
}

func (node *Program) AddStatement(statement Node) {
	node.Statements = append(node.Statements, statement)
}

//
// Mustache Statement
//

type MustacheStatement struct {
	NodeType
	Pos

	Path   Node   // PathExpression
	Params []Node // [ Expression ... ]
	Hash   Node   // Hash
}

func NewMustacheStatement(pos int) *MustacheStatement {
	return &MustacheStatement{
		NodeType: NodeMustache,
		Pos:      Pos(pos),
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
	Pos

	Path    Node   // PathExpression
	Params  []Node // [ Expression ... ]
	Hash    Node   // Hash
	Program Node   // Program
}

func NewBlockStatement(pos int) *BlockStatement {
	return &BlockStatement{
		NodeType: NodeBlock,
		Pos:      Pos(pos),
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
	Pos

	Name   Node   // PathExpression | SubExpression
	Params []Node // [ Expression ... ]
	Hash   Node   // Hash
}

func NewPartialStatement(pos int) *PartialStatement {
	return &PartialStatement{
		NodeType: NodePartial,
		Pos:      Pos(pos),
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
	Pos

	Value string
}

func NewContentStatement(pos int, val string) *ContentStatement {
	return &ContentStatement{
		NodeType: NodeContent,
		Pos:      Pos(pos),
		Value:    val,
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
	Pos

	Value string
}

func NewCommentStatement(pos int, val string) *CommentStatement {
	return &CommentStatement{
		NodeType: NodeComment,
		Pos:      Pos(pos),
		Value:    val,
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
	Pos

	Path   Node   // PathExpression
	Params []Node // [ Expression ... ]
	Hash   Node   // Hash
}

func NewSubExpression(pos int) *SubExpression {
	return &SubExpression{
		NodeType: NodeSubExpression,
		Pos:      Pos(pos),
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
	Pos

	Original string
	Parts    []string
	Data     bool
}

func NewPathExpression(pos int, data bool) *PathExpression {
	result := &PathExpression{
		NodeType: NodePath,
		Pos:      Pos(pos),
		Data:     data,
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

	node.Parts = append(node.Parts, part)
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
	Pos

	Value string
}

func NewStringLiteral(pos int, val string) *StringLiteral {
	return &StringLiteral{
		NodeType: NodeString,
		Pos:      Pos(pos),
		Value:    val,
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
	Pos

	Value bool
}

func NewBooleanLiteral(pos int, val bool) *BooleanLiteral {
	return &BooleanLiteral{
		NodeType: NodeBoolean,
		Pos:      Pos(pos),
		Value:    val,
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
	Pos

	Value int
}

func NewNumberLiteral(pos int, val int) *NumberLiteral {
	return &NumberLiteral{
		NodeType: NodeNumber,
		Pos:      Pos(pos),
		Value:    val,
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
	Pos

	Pairs []Node // [ HashPair ... ]
}

func NewHash(pos int) *Hash {
	return &Hash{
		NodeType: NodeHash,
		Pos:      Pos(pos),
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
	Pos

	Key string
	Val Node // Expression
}

func NewHashPair(pos int) *HashPair {
	return &HashPair{
		NodeType: NodeHashPair,
		Pos:      Pos(pos),
	}
}

func (node *HashPair) Accept(visitor Visitor) {
	visitor.visitHashPair(node)
}

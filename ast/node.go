package ast

import (
	"bytes"
	"fmt"
)

// References:
//   - https://github.com/wycats/handlebars.js/blob/master/lib/handlebars/compiler/ast.js
//   - https://github.com/golang/go/blob/master/src/text/template/parse/node.go

type NodeType int

// AST visitor interface
type Visitor interface {
	visitProgram(node *ProgramNode)
	visitContent(node *ContentNode)
	visitComment(node *CommentNode)
	visitBoolean(node *BooleanNode)
	visitNumber(node *NumberNode)
	visitString(node *StringNode)
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
	NodeContent
	NodeComment
	NodeBoolean
	NodeNumber
	NodeString
)

//
// Program
//

type ProgramNode struct {
	NodeType
	Pos
	Statements []Node
}

func NewProgramNode(pos int) *ProgramNode {
	return &ProgramNode{
		NodeType: NodeProgram,
		Pos:      Pos(pos),
	}
}

func (node *ProgramNode) String() string {
	b := new(bytes.Buffer)

	for _, n := range node.Statements {
		fmt.Fprint(b, n)
	}

	return b.String()
}

func (node *ProgramNode) Accept(visitor Visitor) {
	visitor.visitProgram(node)

	for _, n := range node.Statements {
		n.Accept(visitor)
	}
}

//
// Content
//

type ContentNode struct {
	NodeType
	Pos
	Value string
}

func NewContentNode(pos int, val string) *ContentNode {
	return &ContentNode{
		NodeType: NodeContent,
		Pos:      Pos(pos),
		Value:    val,
	}
}

func (node *ContentNode) Accept(visitor Visitor) {
	visitor.visitContent(node)
}

//
// Comment
//

type CommentNode struct {
	NodeType
	Pos
	Value string
}

func NewCommentNode(pos int, val string) *CommentNode {
	return &CommentNode{
		NodeType: NodeComment,
		Pos:      Pos(pos),
		Value:    val,
	}
}

func (node *CommentNode) Accept(visitor Visitor) {
	visitor.visitComment(node)
}

//
// Boolean
//

type BooleanNode struct {
	NodeType
	Pos
	Value bool
}

func NewBooleanNode(pos int, val bool) *BooleanNode {
	return &BooleanNode{
		NodeType: NodeBoolean,
		Pos:      Pos(pos),
		Value:    val,
	}
}

func (node *BooleanNode) Accept(visitor Visitor) {
	visitor.visitBoolean(node)
}

//
// Number
//

type NumberNode struct {
	NodeType
	Pos
	Value int
}

func NewNumberNode(pos int, val int) *NumberNode {
	return &NumberNode{
		NodeType: NodeNumber,
		Pos:      Pos(pos),
		Value:    val,
	}
}

func (node *NumberNode) Accept(visitor Visitor) {
	visitor.visitNumber(node)
}

//
// String
//

type StringNode struct {
	NodeType
	Pos
	Value string
}

func NewStringNode(pos int, val string) *StringNode {
	return &StringNode{
		NodeType: NodeString,
		Pos:      Pos(pos),
		Value:    val,
	}
}

func (node *StringNode) Accept(visitor Visitor) {
	visitor.visitString(node)
}

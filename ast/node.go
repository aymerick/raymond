package ast

import "fmt"

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

	// string representation, used for debugging
	String() string

	// accepts visitor
	Accept(Visitor)
}

// AST visitor interface
type Visitor interface {
	VisitProgram(*Program) interface{}

	// statements
	VisitMustache(*MustacheStatement) interface{}
	VisitBlock(*BlockStatement) interface{}
	VisitPartial(*PartialStatement) interface{}
	VisitContent(*ContentStatement) interface{}
	VisitComment(*CommentStatement) interface{}

	// expressions
	VisitSubExpression(*SubExpression) interface{}
	VisitPath(*PathExpression) interface{}

	// literals
	VisitString(*StringLiteral) interface{}
	VisitBoolean(*BooleanLiteral) interface{}
	VisitNumber(*NumberLiteral) interface{}

	// miscellaneous
	VisitHash(*Hash) interface{}
	VisitHashPair(*HashPair) interface{}
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
	return fmt.Sprintf("Program{Pos: %d}", node.Loc.Pos)
}

func (node *Program) Accept(visitor Visitor) {
	visitor.VisitProgram(node)
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

func (node *MustacheStatement) String() string {
	return fmt.Sprintf("Mustache{Pos: %d}", node.Loc.Pos)
}

func (node *MustacheStatement) Accept(visitor Visitor) {
	visitor.VisitMustache(node)
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

func (node *BlockStatement) String() string {
	return fmt.Sprintf("Block{Pos: %d}", node.Loc.Pos)
}

func (node *BlockStatement) Accept(visitor Visitor) {
	visitor.VisitBlock(node)
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

func (node *PartialStatement) String() string {
	return fmt.Sprintf("Partial{Name:%s, Pos:%d}", node.Name, node.Loc.Pos)
}

func (node *PartialStatement) Accept(visitor Visitor) {
	visitor.VisitPartial(node)
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

func (node *ContentStatement) String() string {
	return fmt.Sprintf("Content{Value:'%s', Pos:%d}", node.Value, node.Loc.Pos)
}

func (node *ContentStatement) Accept(visitor Visitor) {
	visitor.VisitContent(node)
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

func (node *CommentStatement) String() string {
	return fmt.Sprintf("Comment{Value:'%s', Pos:%d}", node.Value, node.Loc.Pos)
}

func (node *CommentStatement) Accept(visitor Visitor) {
	visitor.VisitComment(node)
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

func (node *SubExpression) String() string {
	return fmt.Sprintf("Sexp{Path:%s, Pos:%d}", node.Path, node.Loc.Pos)
}

func (node *SubExpression) Accept(visitor Visitor) {
	visitor.VisitSubExpression(node)
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

func (node *PathExpression) String() string {
	return fmt.Sprintf("Path{Original:'%s', Pos:%d}", node.Original, node.Loc.Pos)
}

func (node *PathExpression) Accept(visitor Visitor) {
	visitor.VisitPath(node)
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

func (node *StringLiteral) String() string {
	return fmt.Sprintf("String{Value:'%s', Pos:%d}", node.Value, node.Loc.Pos)
}

func (node *StringLiteral) Accept(visitor Visitor) {
	visitor.VisitString(node)
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

func (node *BooleanLiteral) Canonical() string {
	if node.Value {
		return "true"
	} else {
		return "false"
	}
}

func (node *BooleanLiteral) String() string {
	return fmt.Sprintf("Boolean{Value:%s, Pos:%d}", node.Canonical(), node.Loc.Pos)
}

func (node *BooleanLiteral) Accept(visitor Visitor) {
	visitor.VisitBoolean(node)
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

func (node *NumberLiteral) String() string {
	return fmt.Sprintf("Number{Value:%d, Pos:%d}", node.Value, node.Loc.Pos)
}

func (node *NumberLiteral) Accept(visitor Visitor) {
	visitor.VisitNumber(node)
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

func (node *Hash) String() string {
	result := fmt.Sprintf("Hash{[", node.Loc.Pos)

	for i, p := range node.Pairs {
		if i > 0 {
			result += ", "
		}
		result += p.String()
	}

	return result + fmt.Sprintf("], Pos:%d}", node.Loc.Pos)
}

func (node *Hash) Accept(visitor Visitor) {
	visitor.VisitHash(node)
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

func (node *HashPair) String() string {
	return node.Key + "=" + node.Val.String()
}

func (node *HashPair) Accept(visitor Visitor) {
	visitor.VisitHashPair(node)
}

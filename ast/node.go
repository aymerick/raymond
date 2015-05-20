package ast

import (
	"fmt"
	"strconv"
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

	// string representation, used for debugging
	String() string

	// accepts visitor
	Accept(Visitor) interface{}
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
	VisitExpression(*Expression) interface{}
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
	NodeExpression
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

func (node *Program) Accept(visitor Visitor) interface{} {
	return visitor.VisitProgram(node)
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

	Unescaped  bool
	Expression *Expression
}

func NewMustacheStatement(pos int, line int, unescaped bool) *MustacheStatement {
	return &MustacheStatement{
		NodeType:  NodeMustache,
		Loc:       Loc{pos, line},
		Unescaped: unescaped,
	}
}

func (node *MustacheStatement) String() string {
	return fmt.Sprintf("Mustache{Pos: %d}", node.Loc.Pos)
}

func (node *MustacheStatement) Accept(visitor Visitor) interface{} {
	return visitor.VisitMustache(node)
}

//
// Block Statement
//

type BlockStatement struct {
	NodeType
	Loc

	Expression *Expression

	Program Node // Program
	Inverse Node // Program
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

func (node *BlockStatement) Accept(visitor Visitor) interface{} {
	return visitor.VisitBlock(node)
}

//
// Partial Statement
//

type PartialStatement struct {
	NodeType
	Loc

	Name   Node   // PathExpression | SubExpression
	Params []Node // [ Expression ... ]
	Hash   *Hash  // Hash
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

func (node *PartialStatement) Accept(visitor Visitor) interface{} {
	return visitor.VisitPartial(node)
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

func (node *ContentStatement) Accept(visitor Visitor) interface{} {
	return visitor.VisitContent(node)
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

func (node *CommentStatement) Accept(visitor Visitor) interface{} {
	return visitor.VisitComment(node)
}

//
// Expression
//

type Expression struct {
	NodeType
	Loc

	Path   Node   // PathExpression | StringLiteral | BooleanLiteral | NumberLiteral
	Params []Node // [ Expression ... ]
	Hash   *Hash  // Hash
}

func NewExpression(pos int, line int) *Expression {
	return &Expression{
		NodeType: NodeExpression,
		Loc:      Loc{pos, line},
	}
}

func (node *Expression) String() string {
	return fmt.Sprintf("Expr{Path:%s, Pos:%d}", node.Path, node.Loc.Pos)
}

func (node *Expression) Accept(visitor Visitor) interface{} {
	return visitor.VisitExpression(node)
}

func (node *Expression) haveParams() bool {
	return (len(node.Params) > 0) || ((node.Hash != nil) && (len(node.Hash.Pairs) > 0))
}

// return helper name, or an empty string if this expression can't be an helper
func (node *Expression) HelperName() string {
	path, ok := node.Path.(*PathExpression)
	if !ok {
		return ""
	}

	if path.Data || (len(path.Parts) != 1) || (path.Depth > 0) {
		return ""
	}

	return path.Parts[0]
}

// returns path expression representing a field path, or nil if this is not a field path
func (node *Expression) FieldPath() *PathExpression {
	path, ok := node.Path.(*PathExpression)
	if !ok {
		return nil
	}

	return path
}

// returns string representation of literal value, with a boolean set to false if this is not a literal
func (node *Expression) LiteralStr() (string, bool) {
	if node.haveParams() {
		return "", false
	}

	return LiteralStr(node.Path)
}

// returns string representation of expression
func (node *Expression) Str() string {
	if str, ok := HelperNameStr(node.Path); ok {
		return str
	}

	return ""
}

// returns string representation of an helper name, with a boolean set to false if this is not a valid helper name
//
// helperName : path | dataName | STRING | NUMBER | BOOLEAN | UNDEFINED | NULL
func HelperNameStr(node Node) (string, bool) {
	// PathExpression
	if str, ok := PathExpressionStr(node); ok {
		return str, ok
	}

	// Literal
	if str, ok := LiteralStr(node); ok {
		return str, ok
	}

	return "", false
}

// returns string representation of path expression value, with a boolean set to false if this is not a path expression
func PathExpressionStr(node Node) (string, bool) {
	if path, ok := node.(*PathExpression); ok {
		return path.Original, true
	}

	return "", false
}

// returns string representation of literal value, with a boolean set to false if this is not a literal
func LiteralStr(node Node) (string, bool) {
	if lit, ok := node.(*StringLiteral); ok {
		return lit.Value, true
	}

	if lit, ok := node.(*BooleanLiteral); ok {
		return lit.Canonical(), true
	}

	if lit, ok := node.(*NumberLiteral); ok {
		return lit.Canonical(), true
	}

	return "", false
}

//
// SubExpression
//

type SubExpression struct {
	NodeType
	Loc

	Expression *Expression
}

func NewSubExpression(pos int, line int) *SubExpression {
	return &SubExpression{
		NodeType: NodeSubExpression,
		Loc:      Loc{pos, line},
	}
}

func (node *SubExpression) String() string {
	return fmt.Sprintf("Sexp{Path:%s, Pos:%d}", node.Expression.Path, node.Loc.Pos)
}

func (node *SubExpression) Accept(visitor Visitor) interface{} {
	return visitor.VisitSubExpression(node)
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

func (node *PathExpression) Accept(visitor Visitor) interface{} {
	return visitor.VisitPath(node)
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

func (node *StringLiteral) Accept(visitor Visitor) interface{} {
	return visitor.VisitString(node)
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

func (node *BooleanLiteral) String() string {
	return fmt.Sprintf("Boolean{Value:%s, Pos:%d}", node.Canonical(), node.Loc.Pos)
}

func (node *BooleanLiteral) Accept(visitor Visitor) interface{} {
	return visitor.VisitBoolean(node)
}

func (node *BooleanLiteral) Canonical() string {
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

	Value    float64
	IsInt    bool
	Original string
}

func NewNumberLiteral(pos int, line int, val float64, isInt bool, original string) *NumberLiteral {
	return &NumberLiteral{
		NodeType: NodeNumber,
		Loc:      Loc{pos, line},

		Value:    val,
		IsInt:    isInt,
		Original: original,
	}
}

func (node *NumberLiteral) String() string {
	return fmt.Sprintf("Number{Value:%s, Pos:%d}", node.Canonical(), node.Loc.Pos)
}

func (node *NumberLiteral) Accept(visitor Visitor) interface{} {
	return visitor.VisitNumber(node)
}

func (node *NumberLiteral) Canonical() string {
	prec := -1
	if node.IsInt {
		prec = 0
	}
	return strconv.FormatFloat(node.Value, 'f', prec, 64)
}

// Returns an integer or a float
func (node *NumberLiteral) Number() interface{} {
	if node.IsInt {
		return int(node.Value)
	} else {
		return node.Value
	}
}

//
// Hash
//

type Hash struct {
	NodeType
	Loc

	Pairs []*HashPair // [ HashPair ... ]
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

func (node *Hash) Accept(visitor Visitor) interface{} {
	return visitor.VisitHash(node)
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

func (node *HashPair) Accept(visitor Visitor) interface{} {
	return visitor.VisitHashPair(node)
}

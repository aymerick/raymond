// Package ast provides structures to represent a handlebars Abstract Syntax Tree, and a Visitor interface to visit that tree.
package ast

import (
	"fmt"
	"strconv"
	"strings"
)

// References:
//   - https://github.com/wycats/handlebars.js/blob/master/lib/handlebars/compiler/ast.js
//   - https://github.com/wycats/handlebars.js/blob/master/docs/compiler-api.md
//   - https://github.com/golang/go/blob/master/src/text/template/parse/node.go

// Node is an element in the AST.
type Node interface {
	// node type
	Type() NodeType

	// location of node in original input string
	Location() Loc

	// string representation, used for debugging
	String() string

	// accepts visitor
	Accept(Visitor) interface{}

	// writes back to textual form
	Serialize() string
}

// Visitor is the interface to visit an AST.
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

// NodeType represents an AST Node type.
type NodeType int

// Type returns itself, and permits struct includers to satisfy that part of Node interface.
func (t NodeType) Type() NodeType {
	return t
}

const (
	// NodeProgram is the program node
	NodeProgram NodeType = iota // 0

	// NodeMustache is the mustache statement node
	NodeMustache  // 1

	// NodeBlock is the block statement node
	NodeBlock  // 2

	// NodePartial is the partial statement node
	NodePartial  // 3

	// NodeContent is the content statement node
	NodeContent  // 4

	// NodeComment is the comment statement node
	NodeComment  // 5

	// NodeExpression is the expression node
	NodeExpression  // 6

	// NodeSubExpression is the subexpression node
	NodeSubExpression  // 7

	// NodePath is the expression path node
	NodePath  // 8

	// NodeBoolean is the literal boolean node
	NodeBoolean  // 9

	// NodeNumber is the literal number node
	NodeNumber  // 10

	// NodeString is the literal string node
	NodeString  // 11

	// NodeHash is the hash node
	NodeHash  // 12

	// NodeHashPair is the hash pair node
	NodeHashPair  // 13
)

// Loc represents the position of a parsed node in source file.
type Loc struct {
	Pos  int // Byte position
	Line int // Line number
}

// Location returns itself, and permits struct includers to satisfy that part of Node interface.
func (l Loc) Location() Loc {
	return l
}

// Strip describes node whitespace management.
type Strip struct {
	Open  bool
	Close bool

	OpenStandalone   bool
	CloseStandalone  bool
	InlineStandalone bool
}

// NewStrip instantiates a Strip for given open and close mustaches.
func NewStrip(openStr, closeStr string) *Strip {
	return &Strip{
		Open:  (len(openStr) > 2) && openStr[2] == '~',
		Close: (len(closeStr) > 2) && closeStr[len(closeStr)-3] == '~',
	}
}

// NewStripForStr instantiates a Strip for given tag.
func NewStripForStr(str string) *Strip {
	return &Strip{
		Open:  (len(str) > 2) && str[2] == '~',
		Close: (len(str) > 2) && str[len(str)-3] == '~',
	}
}

// String returns a string representation of receiver that can be used for debugging.
func (s *Strip) String() string {
	return fmt.Sprintf("Open: %t, Close: %t, OpenStandalone: %t, CloseStandalone: %t, InlineStandalone: %t", s.Open, s.Close, s.OpenStandalone, s.CloseStandalone, s.InlineStandalone)
}

//
// Program
//

// Program represents a program node.
type Program struct {
	NodeType
	Loc

	Body        []Node // [ Statement ... ]
	BlockParams []string
	Chained     bool

	// whitespace management
	Strip *Strip
}

// NewProgram instantiates a new program node.
func NewProgram(pos int, line int) *Program {
	return &Program{
		NodeType: NodeProgram,
		Loc:      Loc{pos, line},
	}
}

// String returns a string representation of receiver that can be used for debugging.
func (node *Program) String() string {
	return fmt.Sprintf("Program{Pos: %d}", node.Loc.Pos)
}

// Accept is the receiver entry point for visitors.
func (node *Program) Accept(visitor Visitor) interface{} {
	return visitor.VisitProgram(node)
}

// AddStatement adds given statement to program.
func (node *Program) AddStatement(statement Node) {
	node.Body = append(node.Body, statement)
}

// Writes node back to textual form
func (node *Program) Serialize() string {
	content := ""

	// TODO Chained

	for _, node := range node.Body {
		content += node.Serialize()
	}

	return content
}

//
// Mustache Statement
//

// MustacheStatement represents a mustache node.
type MustacheStatement struct {
	NodeType
	Loc

	Expression *Expression

	// whitespace management
	Strip *Strip
}

// NewMustacheStatement instantiates a new mustache node.
func NewMustacheStatement(pos int, line int) *MustacheStatement {
	return &MustacheStatement{
		NodeType: NodeMustache,
		Loc:      Loc{pos, line},
	}
}

// String returns a string representation of receiver that can be used for debugging.
func (node *MustacheStatement) String() string {
	return fmt.Sprintf("Mustache{Pos: %d}", node.Loc.Pos)
}

// Accept is the receiver entry point for visitors.
func (node *MustacheStatement) Accept(visitor Visitor) interface{} {
	return visitor.VisitMustache(node)
}

// Writes node back to textual form
func (node *MustacheStatement) Serialize() string {
	content := "{{"

	if node.Strip != nil && node.Strip.Open {
		content += "~"
	}

	content += node.Expression.Serialize()

	if node.Strip != nil && node.Strip.Close {
		content += "~"
	}

	content += "}}"

	return content
}

//
// Block Statement
//

// BlockStatement represents a block node.
type BlockStatement struct {
	NodeType
	Loc

	Expression *Expression

	Program *Program
	Inverse *Program

	// whitespace management
	OpenStrip    *Strip
	InverseStrip *Strip
	CloseStrip   *Strip
}

// NewBlockStatement instantiates a new block node.
func NewBlockStatement(pos int, line int) *BlockStatement {
	return &BlockStatement{
		NodeType: NodeBlock,
		Loc:      Loc{pos, line},
	}
}

// String returns a string representation of receiver that can be used for debugging.
func (node *BlockStatement) String() string {
	return fmt.Sprintf("Block{Pos: %d}", node.Loc.Pos)
}

// Accept is the receiver entry point for visitors.
func (node *BlockStatement) Accept(visitor Visitor) interface{} {
	return visitor.VisitBlock(node)
}

// Writes node back to textual form
func (node *BlockStatement) Serialize() string {
	program := node.Program
	inverse := node.Inverse

	content := ""

	if program != nil {
		// Open block
		content += "{{"

		if node.OpenStrip != nil && node.OpenStrip.Open {
			content += "~"
		}

		content += "#"

		content += node.Expression.Serialize()

		// Open
		if len(program.BlockParams) > 0 {
			content += " as |"
			content += strings.Join(program.BlockParams, " ")
			content += "|"
		}

		if node.OpenStrip != nil && node.OpenStrip.Close {
			content += "~"
		}

		content += "}}"
		// Content
		content += program.Serialize()
	}

	// Inverse
	if inverse != nil {
		lastInverse := inverse

		if inverse.Chained {
			var inverseBlocks []*BlockStatement

			for lastInverse != nil && lastInverse.Chained {
				lastInverseBlock := lastInverse.Body[len(lastInverse.Body)-1].(*BlockStatement)
				lastInverse = lastInverseBlock.Inverse
				inverseBlocks = append(inverseBlocks, lastInverseBlock)
			}

			for _, block := range inverseBlocks {
				content += "{{"

				if block.OpenStrip != nil && block.OpenStrip.Open {
					content += "~"
				}

				content += "else "

				content += block.Expression.Serialize()

				if len(block.Program.BlockParams) > 0 {
					content += " as |"
					content += strings.Join(block.Program.BlockParams, " ")
					content += "|"
				}

				if block.OpenStrip != nil && block.OpenStrip.Close {
					content += "~"
				}

				content += "}}"

				content += block.Program.Serialize()
			}
		}

		if lastInverse != nil {
			content += "{{"

			if lastInverse.Strip != nil && lastInverse.Strip.Open {
				content += "~"
			}

			content += "^"

			// Pure inverse section
			if program == nil {
				content += node.Expression.Serialize()
				if len(lastInverse.BlockParams) > 0 {
					content += " as |"
					content += strings.Join(lastInverse.BlockParams, " ")
					content += "|"
				}
			}

			if lastInverse.Strip != nil && lastInverse.Strip.Close {
				content += "~"
			}

			content += "}}"

			content += lastInverse.Serialize()
		}

	}

	// Close
	content += "{{"

	if node.CloseStrip != nil && node.CloseStrip.Open {
		content += "~"
	}

	content += "/"

	content += node.Expression.Serialize()

	if node.CloseStrip != nil && node.CloseStrip.Close {
		content += "~"
	}

	content += "}}"

	return content
}

//
// Partial Statement
//

// PartialStatement represents a partial node.
type PartialStatement struct {
	NodeType
	Loc

	Name   Node   // PathExpression | SubExpression
	Params []Node // [ Expression ... ]
	Hash   *Hash

	// whitespace management
	Strip  *Strip
	Indent string
}

// NewPartialStatement instantiates a new partial node.
func NewPartialStatement(pos int, line int) *PartialStatement {
	return &PartialStatement{
		NodeType: NodePartial,
		Loc:      Loc{pos, line},
	}
}

// String returns a string representation of receiver that can be used for debugging.
func (node *PartialStatement) String() string {
	return fmt.Sprintf("Partial{Name:%s, Pos:%d}", node.Name, node.Loc.Pos)
}

// Accept is the receiver entry point for visitors.
func (node *PartialStatement) Accept(visitor Visitor) interface{} {
	return visitor.VisitPartial(node)
}

// Writes node back to textual form
func (node *PartialStatement) Serialize() string {
	content := node.Indent + "{{"

	if node.Strip != nil && node.Strip.Open {
		content += "~"
	}

	content += "> "

	content += node.Name.Serialize()

	for _, p := range node.Params {
		content += " " + p.Serialize()
	}

	if node.Hash != nil {
		content += " " + node.Hash.Serialize()
	}

	if node.Strip != nil && node.Strip.Close {
		content += "~"
	}

	content += "}}"

	return content
}

//
// Content Statement
//

// ContentStatement represents a content node.
type ContentStatement struct {
	NodeType
	Loc

	Value    string
	Original string

	// whitespace management
	RightStripped bool
	LeftStripped  bool
}

// NewContentStatement instantiates a new content node.
func NewContentStatement(pos int, line int, val string) *ContentStatement {
	return &ContentStatement{
		NodeType: NodeContent,
		Loc:      Loc{pos, line},

		Value:    val,
		Original: val,
	}
}

// String returns a string representation of receiver that can be used for debugging.
func (node *ContentStatement) String() string {
	return fmt.Sprintf("Content{Value:'%s', Pos:%d}", node.Value, node.Loc.Pos)
}

// Accept is the receiver entry point for visitors.
func (node *ContentStatement) Accept(visitor Visitor) interface{} {
	return visitor.VisitContent(node)
}

// Writes node back to textual form
func (node *ContentStatement) Serialize() string {
	return node.Original
}

//
// Comment Statement
//

// CommentStatement represents a comment node.
type CommentStatement struct {
	NodeType
	Loc

	Value string

	// whitespace management
	Strip *Strip
}

// NewCommentStatement instantiates a new comment node.
func NewCommentStatement(pos int, line int, val string) *CommentStatement {
	return &CommentStatement{
		NodeType: NodeComment,
		Loc:      Loc{pos, line},

		Value: val,
	}
}

// String returns a string representation of receiver that can be used for debugging.
func (node *CommentStatement) String() string {
	return fmt.Sprintf("Comment{Value:'%s', Pos:%d}", node.Value, node.Loc.Pos)
}

// Accept is the receiver entry point for visitors.
func (node *CommentStatement) Accept(visitor Visitor) interface{} {
	return visitor.VisitComment(node)
}

// Writes node back to textual form
func (node *CommentStatement) Serialize() string {
	content := "{{"

	if node.Strip != nil && node.Strip.Open {
		content += "~"
	}

	content += "!" + node.Value

	if node.Strip != nil && node.Strip.Close {
		content += "~"
	}

	content += "}}"

	return content
}

//
// Expression
//

// Expression represents an expression node.
type Expression struct {
	NodeType
	Loc

	Path   Node   // PathExpression | StringLiteral | BooleanLiteral | NumberLiteral
	Params []Node // [ Expression ... ]
	Hash   *Hash
}

// NewExpression instantiates a new expression node.
func NewExpression(pos int, line int) *Expression {
	return &Expression{
		NodeType: NodeExpression,
		Loc:      Loc{pos, line},
	}
}

// String returns a string representation of receiver that can be used for debugging.
func (node *Expression) String() string {
	return fmt.Sprintf("Expr{Path:%s, Pos:%d}", node.Path, node.Loc.Pos)
}

// Accept is the receiver entry point for visitors.
func (node *Expression) Accept(visitor Visitor) interface{} {
	return visitor.VisitExpression(node)
}

// HelperName returns helper name, or an empty string if this expression can't be a helper.
func (node *Expression) HelperName() string {
	path, ok := node.Path.(*PathExpression)
	if !ok {
		return ""
	}

	if path.Data || (len(path.Parts) != 1) || (path.Depth > 0) || path.Scoped {
		return ""
	}

	return path.Parts[0]
}

// FieldPath returns path expression representing a field path, or nil if this is not a field path.
func (node *Expression) FieldPath() *PathExpression {
	path, ok := node.Path.(*PathExpression)
	if !ok {
		return nil
	}

	return path
}

// LiteralStr returns the string representation of literal value, with a boolean set to false if this is not a literal.
func (node *Expression) LiteralStr() (string, bool) {
	return LiteralStr(node.Path)
}

// Canonical returns the canonical form of expression node as a string.
func (node *Expression) Canonical() string {
	if str, ok := HelperNameStr(node.Path); ok {
		return str
	}

	return ""
}

// Writes node back to textual form
func (node *Expression) Serialize() string {
	content := node.Path.Serialize()

	for _, p := range node.Params {
		content += " " + p.Serialize()
	}

	if node.Hash != nil {
		content += " " + node.Hash.Serialize()
	}

	return content
}

// HelperNameStr returns the string representation of a helper name, with a boolean set to false if this is not a valid helper name.
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

// PathExpressionStr returns the string representation of path expression value, with a boolean set to false if this is not a path expression.
func PathExpressionStr(node Node) (string, bool) {
	if path, ok := node.(*PathExpression); ok {
		result := path.Original

		// "[foo bar]"" => "foo bar"
		if (len(result) >= 2) && (result[0] == '[') && (result[len(result)-1] == ']') {
			result = result[1 : len(result)-1]
		}

		return result, true
	}

	return "", false
}

// LiteralStr returns the string representation of literal value, with a boolean set to false if this is not a literal.
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

// SubExpression represents a subexpression node.
type SubExpression struct {
	NodeType
	Loc

	Expression *Expression
}

// NewSubExpression instantiates a new subexpression node.
func NewSubExpression(pos int, line int) *SubExpression {
	return &SubExpression{
		NodeType: NodeSubExpression,
		Loc:      Loc{pos, line},
	}
}

// String returns a string representation of receiver that can be used for debugging.
func (node *SubExpression) String() string {
	return fmt.Sprintf("Sexp{Path:%s, Pos:%d}", node.Expression.Path, node.Loc.Pos)
}

// Accept is the receiver entry point for visitors.
func (node *SubExpression) Accept(visitor Visitor) interface{} {
	return visitor.VisitSubExpression(node)
}

// Writes node back to textual form
func (node *SubExpression) Serialize() string {
	return node.Expression.Serialize()
}

//
// Path Expression
//

// PathExpression represents a path expression node.
type PathExpression struct {
	NodeType
	Loc

	Original string
	Depth    int
	Parts    []string
	Data     bool
	Scoped   bool
}

// NewPathExpression instantiates a new path expression node.
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

// String returns a string representation of receiver that can be used for debugging.
func (node *PathExpression) String() string {
	return fmt.Sprintf("Path{Original:'%s', Pos:%d}", node.Original, node.Loc.Pos)
}

// Accept is the receiver entry point for visitors.
func (node *PathExpression) Accept(visitor Visitor) interface{} {
	return visitor.VisitPath(node)
}

// Part adds path part.
func (node *PathExpression) Part(part string) {
	node.Original += part

	switch part {
	case "..":
		node.Depth++
		node.Scoped = true
	case ".", "this":
		node.Scoped = true
	default:
		node.Parts = append(node.Parts, part)
	}
}

// Sep adds path separator.
func (node *PathExpression) Sep(separator string) {
	node.Original += separator
}

// IsDataRoot returns true if path expression is @root.
func (node *PathExpression) IsDataRoot() bool {
	return node.Data && (node.Parts[0] == "root")
}

// Writes node back to textual form
func (node *PathExpression) Serialize() string {
	return node.Original
}

//
// String Literal
//

// StringLiteral represents a string node.
type StringLiteral struct {
	NodeType
	Loc

	Value string
}

// NewStringLiteral instantiates a new string node.
func NewStringLiteral(pos int, line int, val string) *StringLiteral {
	return &StringLiteral{
		NodeType: NodeString,
		Loc:      Loc{pos, line},

		Value: val,
	}
}

// String returns a string representation of receiver that can be used for debugging.
func (node *StringLiteral) String() string {
	return fmt.Sprintf("String{Value:'%s', Pos:%d}", node.Value, node.Loc.Pos)
}

// Accept is the receiver entry point for visitors.
func (node *StringLiteral) Accept(visitor Visitor) interface{} {
	return visitor.VisitString(node)
}

// Writes node back to textual form
func (node *StringLiteral) Serialize() string {
	return "\"" + node.Value + "\""
}

//
// Boolean Literal
//

// BooleanLiteral represents a boolean node.
type BooleanLiteral struct {
	NodeType
	Loc

	Value    bool
	Original string
}

// NewBooleanLiteral instantiates a new boolean node.
func NewBooleanLiteral(pos int, line int, val bool, original string) *BooleanLiteral {
	return &BooleanLiteral{
		NodeType: NodeBoolean,
		Loc:      Loc{pos, line},

		Value:    val,
		Original: original,
	}
}

// String returns a string representation of receiver that can be used for debugging.
func (node *BooleanLiteral) String() string {
	return fmt.Sprintf("Boolean{Value:%s, Pos:%d}", node.Canonical(), node.Loc.Pos)
}

// Accept is the receiver entry point for visitors.
func (node *BooleanLiteral) Accept(visitor Visitor) interface{} {
	return visitor.VisitBoolean(node)
}

// Canonical returns the canonical form of boolean node as a string (ie. "true" | "false").
func (node *BooleanLiteral) Canonical() string {
	if node.Value {
		return "true"
	}

	return "false"
}

// Writes node back to textual form
func (node *BooleanLiteral) Serialize() string {
	return node.Original
}

//
// Number Literal
//

// NumberLiteral represents a number node.
type NumberLiteral struct {
	NodeType
	Loc

	Value    float64
	IsInt    bool
	Original string
}

// NewNumberLiteral instantiates a new number node.
func NewNumberLiteral(pos int, line int, val float64, isInt bool, original string) *NumberLiteral {
	return &NumberLiteral{
		NodeType: NodeNumber,
		Loc:      Loc{pos, line},

		Value:    val,
		IsInt:    isInt,
		Original: original,
	}
}

// String returns a string representation of receiver that can be used for debugging.
func (node *NumberLiteral) String() string {
	return fmt.Sprintf("Number{Value:%s, Pos:%d}", node.Canonical(), node.Loc.Pos)
}

// Accept is the receiver entry point for visitors.
func (node *NumberLiteral) Accept(visitor Visitor) interface{} {
	return visitor.VisitNumber(node)
}

// Canonical returns the canonical form of number node as a string (eg: "12", "-1.51").
func (node *NumberLiteral) Canonical() string {
	prec := -1
	if node.IsInt {
		prec = 0
	}
	return strconv.FormatFloat(node.Value, 'f', prec, 64)
}

// Number returns an integer or a float.
func (node *NumberLiteral) Number() interface{} {
	if node.IsInt {
		return int(node.Value)
	}

	return node.Value
}

// Writes node back to textual form
func (node *NumberLiteral) Serialize() string {
	return node.Original
}

//
// Hash
//

// Hash represents a hash node.
type Hash struct {
	NodeType
	Loc

	Pairs []*HashPair
}

// NewHash instantiates a new hash node.
func NewHash(pos int, line int) *Hash {
	return &Hash{
		NodeType: NodeHash,
		Loc:      Loc{pos, line},
	}
}

// String returns a string representation of receiver that can be used for debugging.
func (node *Hash) String() string {
	result := fmt.Sprintf("Hash{[%d", node.Loc.Pos)

	for i, p := range node.Pairs {
		if i > 0 {
			result += ", "
		}
		result += p.String()
	}

	return result + fmt.Sprintf("], Pos:%d}", node.Loc.Pos)
}

// Accept is the receiver entry point for visitors.
func (node *Hash) Accept(visitor Visitor) interface{} {
	return visitor.VisitHash(node)
}

// Writes node back to textual form
func (node *Hash) Serialize() string {
	content := ""

	lastIdx := len(node.Pairs) - 1
	for i, p := range node.Pairs {
		content += p.Serialize()
		if i < lastIdx {
			content += " "
		}
	}

	return content
}

//
// HashPair
//

// HashPair represents a hash pair node.
type HashPair struct {
	NodeType
	Loc

	Key string
	Val Node // Expression
}

// NewHashPair instantiates a new hash pair node.
func NewHashPair(pos int, line int) *HashPair {
	return &HashPair{
		NodeType: NodeHashPair,
		Loc:      Loc{pos, line},
	}
}

// String returns a string representation of receiver that can be used for debugging.
func (node *HashPair) String() string {
	return node.Key + "=" + node.Val.String()
}

// Accept is the receiver entry point for visitors.
func (node *HashPair) Accept(visitor Visitor) interface{} {
	return visitor.VisitHashPair(node)
}

// Writes node back to textual form
func (node *HashPair) Serialize() string {
	content := node.Key + "=" + node.Val.Serialize()

	return content
}

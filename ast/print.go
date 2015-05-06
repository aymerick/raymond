package ast

import (
	"fmt"
	"strings"
)

// Print AST
type PrintVisitor struct {
	buf string

	depth int
}

func NewPrintVisitor() *PrintVisitor {
	return &PrintVisitor{}
}

func PrintNode(node Node) string {
	visitor := NewPrintVisitor()
	node.Accept(visitor)
	return visitor.Output()
}

func (v *PrintVisitor) Output() string {
	return v.buf
}

func (v *PrintVisitor) indent() {
	for i := 0; i < v.depth; {
		v.buf += " "
		i++
	}
}

func (v *PrintVisitor) str(val string) {
	v.buf += val
}

func (v *PrintVisitor) nl() {
	v.str("\n")
}

func (v *PrintVisitor) line(val string) {
	v.indent()
	v.str(val)
	v.nl()
}

func (v *PrintVisitor) printExpression(path Node, params []Node, hash Node) {
	// path
	path.Accept(v)

	// params
	v.str(" [")
	for i, n := range params {
		if i > 0 {
			v.str(", ")
		}
		n.Accept(v)
	}
	v.str("]")

	// hash
	if hash != nil {
		hash.Accept(v)
	}
}

//
// Visitor interface
//

// Statements

func (v *PrintVisitor) visitProgram(node *Program) {
	for _, n := range node.Statements {
		n.Accept(v)
	}
}

func (v *PrintVisitor) visitMustache(node *MustacheStatement) {
	v.indent()
	v.str("{{ ")

	v.printExpression(node.Path, node.Params, node.Hash)

	v.str(" }}")
	v.nl()
}

func (v *PrintVisitor) visitBlock(node *BlockStatement) {
	// @todo !!!
}

func (v *PrintVisitor) visitPartial(node *PartialStatement) {
	// @todo !!!
}

func (v *PrintVisitor) visitContent(node *ContentStatement) {
	v.line("CONTENT[" + node.Value + "]")
}

func (v *PrintVisitor) visitComment(node *CommentStatement) {
	v.line("{{! '" + node.Value + "' }}")
}

// Expressions

func (v *PrintVisitor) visitSubExpression(node *SubExpression) {
	v.printExpression(node.Path, node.Params, node.Hash)
}

func (v *PrintVisitor) visitPath(node *PathExpression) {
	path := strings.Join(node.Parts, "/")

	result := ""
	if node.Data {
		result += "@"
	}

	v.str(result + "PATH:" + path)
}

// Literals

func (v *PrintVisitor) visitString(node *StringLiteral) {
	v.str("\"" + node.Value + "\"")
}

func (v *PrintVisitor) visitBoolean(node *BooleanLiteral) {
	v.str(fmt.Sprintf("BOOLEAN{%s}", node))
}

func (v *PrintVisitor) visitNumber(node *NumberLiteral) {
	v.str(fmt.Sprintf("NUMBER{%d}", node.Value))
}

// Miscellaneous

func (v *PrintVisitor) visitHash(node *Hash) {
	v.str("HASH{")

	for _, p := range node.Pairs {
		p.Accept(v)
	}

	v.str("}")
}

func (v *PrintVisitor) visitHashPair(node *HashPair) {
	v.str(node.Key + "=")
	node.Val.Accept(v)
}

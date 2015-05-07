package ast

import (
	"fmt"
	"strings"
)

// Print AST
type PrintVisitor struct {
	buf   string
	depth int

	original bool
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
		v.buf += "  "
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

func (v *PrintVisitor) printExpression(path Node, params []Node, hash Node, line bool) {
	if line {
		v.indent()
	}

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
		v.str(" ")
		hash.Accept(v)
	}

	if line {
		v.nl()
	}
}

//
// Visitor interface
//

// Statements

func (v *PrintVisitor) visitProgram(node *Program) {
	if len(node.BlockParams) > 0 {
		v.line("BLOCK PARAMS: [ " + strings.Join(node.BlockParams, " ") + " ]")
	}

	for _, n := range node.Body {
		n.Accept(v)
	}
}

func (v *PrintVisitor) visitMustache(node *MustacheStatement) {
	v.indent()
	v.str("{{ ")

	v.printExpression(node.Path, node.Params, node.Hash, false)

	v.str(" }}")
	v.nl()
}

func (v *PrintVisitor) visitBlock(node *BlockStatement) {
	v.line("BLOCK:")
	v.depth++

	v.printExpression(node.Path, node.Params, node.Hash, true)

	if node.Program != nil {
		v.line("PROGRAM:")
		v.depth++
		node.Program.Accept(v)
		v.depth--
	}

	if node.Inverse != nil {
		// if node.Program != nil {
		// 	v.depth++
		// }

		v.line("{{^}}")
		v.depth++
		node.Inverse.Accept(v)
		v.depth--

		// if node.Program != nil {
		// 	v.depth--
		// }
	}
}

func (v *PrintVisitor) visitPartial(node *PartialStatement) {
	v.indent()
	v.str("{{> PARTIAL:")

	v.original = true
	node.Name.Accept(v)
	v.original = false

	if len(node.Params) > 0 {
		v.str(" ")
		node.Params[0].Accept(v)
	}

	// hash
	if node.Hash != nil {
		v.str(" ")
		node.Hash.Accept(v)
	}

	v.str(" }}")
	v.nl()
}

func (v *PrintVisitor) visitContent(node *ContentStatement) {
	v.line("CONTENT[ '" + node.Value + "' ]")
}

func (v *PrintVisitor) visitComment(node *CommentStatement) {
	v.line("{{! '" + node.Value + "' }}")
}

// Expressions

func (v *PrintVisitor) visitSubExpression(node *SubExpression) {
	v.printExpression(node.Path, node.Params, node.Hash, false)
}

func (v *PrintVisitor) visitPath(node *PathExpression) {
	if v.original {
		v.str(node.Original)
	} else {
		path := strings.Join(node.Parts, "/")

		result := ""
		if node.Data {
			result += "@"
		}

		v.str(result + "PATH:" + path)
	}
}

// Literals

func (v *PrintVisitor) visitString(node *StringLiteral) {
	if v.original {
		v.str(node.Value)
	} else {
		v.str("\"" + node.Value + "\"")
	}
}

func (v *PrintVisitor) visitBoolean(node *BooleanLiteral) {
	if v.original {
		v.str(node.Original)
	} else {
		v.str(fmt.Sprintf("BOOLEAN{%s}", node.Canonical()))
	}
}

func (v *PrintVisitor) visitNumber(node *NumberLiteral) {
	if v.original {
		v.str(node.Original)
	} else {
		v.str(fmt.Sprintf("NUMBER{%d}", node.Value))
	}
}

// Miscellaneous

func (v *PrintVisitor) visitHash(node *Hash) {
	v.str("HASH{")

	for i, p := range node.Pairs {
		if i > 0 {
			v.str(", ")
		}
		p.Accept(v)
	}

	v.str("}")
}

func (v *PrintVisitor) visitHashPair(node *HashPair) {
	v.str(node.Key + "=")
	node.Val.Accept(v)
}

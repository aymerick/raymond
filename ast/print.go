package ast

import (
	"fmt"
	"strings"
)

// PrintVisitor implements the `Visitor` interface to print the AST. It is used for unit testing.
type PrintVisitor struct {
	buf   string
	depth int

	original bool
	inBlock  bool
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

//
// Visitor interface
//

// Statements

func (v *PrintVisitor) VisitProgram(node *Program) interface{} {
	if len(node.BlockParams) > 0 {
		v.line("BLOCK PARAMS: [ " + strings.Join(node.BlockParams, " ") + " ]")
	}

	for _, n := range node.Body {
		n.Accept(v)
	}

	return nil
}

func (v *PrintVisitor) VisitMustache(node *MustacheStatement) interface{} {
	v.indent()
	v.str("{{ ")

	node.Expression.Accept(v)

	v.str(" }}")
	v.nl()

	return nil
}

func (v *PrintVisitor) VisitBlock(node *BlockStatement) interface{} {
	v.inBlock = true

	v.line("BLOCK:")
	v.depth++

	node.Expression.Accept(v)

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

	v.inBlock = false

	return nil
}

func (v *PrintVisitor) VisitPartial(node *PartialStatement) interface{} {
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

	return nil
}

func (v *PrintVisitor) VisitContent(node *ContentStatement) interface{} {
	v.line("CONTENT[ '" + node.Value + "' ]")

	return nil
}

func (v *PrintVisitor) VisitComment(node *CommentStatement) interface{} {
	v.line("{{! '" + node.Value + "' }}")

	return nil
}

// Expressions

func (v *PrintVisitor) VisitExpression(node *Expression) interface{} {
	if v.inBlock {
		v.indent()
	}

	// path
	node.Path.Accept(v)

	// params
	v.str(" [")
	for i, n := range node.Params {
		if i > 0 {
			v.str(", ")
		}
		n.Accept(v)
	}
	v.str("]")

	// hash
	if node.Hash != nil {
		v.str(" ")
		node.Hash.Accept(v)
	}

	if v.inBlock {
		v.nl()
	}

	return nil
}

func (v *PrintVisitor) VisitSubExpression(node *SubExpression) interface{} {
	node.Expression.Accept(v)

	return nil
}

func (v *PrintVisitor) VisitPath(node *PathExpression) interface{} {
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

	return nil
}

// Literals

func (v *PrintVisitor) VisitString(node *StringLiteral) interface{} {
	if v.original {
		v.str(node.Value)
	} else {
		v.str("\"" + node.Value + "\"")
	}

	return nil
}

func (v *PrintVisitor) VisitBoolean(node *BooleanLiteral) interface{} {
	if v.original {
		v.str(node.Original)
	} else {
		v.str(fmt.Sprintf("BOOLEAN{%s}", node.Canonical()))
	}

	return nil
}

func (v *PrintVisitor) VisitNumber(node *NumberLiteral) interface{} {
	if v.original {
		v.str(node.Original)
	} else {
		v.str(fmt.Sprintf("NUMBER{%s}", node.Canonical()))
	}

	return nil
}

// Miscellaneous

func (v *PrintVisitor) VisitHash(node *Hash) interface{} {
	v.str("HASH{")

	for i, p := range node.Pairs {
		if i > 0 {
			v.str(", ")
		}
		p.Accept(v)
	}

	v.str("}")

	return nil
}

func (v *PrintVisitor) VisitHashPair(node *HashPair) interface{} {
	v.str(node.Key + "=")
	node.Val.Accept(v)

	return nil
}

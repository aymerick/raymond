package parser

import (
	"regexp"

	"github.com/aymerick/raymond/ast"
)

type WhitespaceVisitor struct {
}

var (
	// @todo multiple: ^\s+/ | else: ^[ \t]*\r?\n?
	rTrimLeft = regexp.MustCompile(`^\s+`)

	// @todo multiple: \s+$ | else: [ \t]+$
	rTrimRight = regexp.MustCompile(`\s+$`)
)

func NewWhitespaceVisitor() *WhitespaceVisitor {
	return &WhitespaceVisitor{}
}

func ProcessWhitespaces(node ast.Node) {
	node.Accept(NewWhitespaceVisitor())
}

func (v *WhitespaceVisitor) trimLeft(node ast.Node) {
	if node.Type() != ast.NodeContent {
		return
	}

	n, _ := node.(*ast.ContentStatement)
	n.Value = rTrimLeft.ReplaceAllString(n.Value, "")
}

func (v *WhitespaceVisitor) trimRight(node ast.Node) {
	if node.Type() != ast.NodeContent {
		return
	}

	n, _ := node.(*ast.ContentStatement)
	n.Value = rTrimRight.ReplaceAllString(n.Value, "")
}

//
// Visitor interface
//

func (v *WhitespaceVisitor) VisitProgram(node *ast.Program) interface{} {
	for i, n := range node.Body {
		strip, _ := n.Accept(v).(*ast.Strip)
		if strip == nil {
			continue
		}

		if strip.Open && (i > 0) {
			v.trimRight(node.Body[i-1])
		}

		if strip.Close && (len(node.Body) > i+1) {
			v.trimLeft(node.Body[i+1])
		}
	}

	return nil
}

func (v *WhitespaceVisitor) VisitMustache(node *ast.MustacheStatement) interface{} {
	return node.Strip
}

func (v *WhitespaceVisitor) VisitBlock(node *ast.BlockStatement) interface{} {
	if node.Program != nil {
		node.Program.Accept(v)
	}

	if node.Inverse != nil {
		node.Inverse.Accept(v)
	}

	program := node.Program
	inverse := node.Inverse

	// @todo Handles chained inverse

	if program == nil {
		program = node.Inverse
		inverse = nil
	}

	p, _ := program.(*ast.Program)
	i, _ := inverse.(*ast.Program)

	if (program != nil) && (node.OpenStrip != nil) && node.OpenStrip.Close {
		v.trimLeft(p.Body[0])
	}

	if inverse != nil {
		if node.InverseStrip != nil {
			if node.InverseStrip.Open {
				v.trimRight(p.Body[len(p.Body)-1])
			}

			if node.InverseStrip.Close {
				v.trimLeft(i.Body[0])
			}
		}

		if (node.CloseStrip != nil) && node.CloseStrip.Open {
			v.trimRight(i.Body[len(i.Body)-1])
		}
	} else if (node.CloseStrip != nil) && node.CloseStrip.Open {
		v.trimRight(p.Body[len(p.Body)-1])
	}

	return ast.NewStripBool((node.OpenStrip != nil) && node.OpenStrip.Open, (node.CloseStrip != nil) && node.CloseStrip.Close)
}

func (v *WhitespaceVisitor) VisitPartial(node *ast.PartialStatement) interface{} {
	return node.Strip
}

func (v *WhitespaceVisitor) VisitComment(node *ast.CommentStatement) interface{} {
	return node.Strip
}

// NOOP
func (v *WhitespaceVisitor) VisitContent(node *ast.ContentStatement) interface{}    { return nil }
func (v *WhitespaceVisitor) VisitExpression(node *ast.Expression) interface{}       { return nil }
func (v *WhitespaceVisitor) VisitSubExpression(node *ast.SubExpression) interface{} { return nil }
func (v *WhitespaceVisitor) VisitPath(node *ast.PathExpression) interface{}         { return nil }
func (v *WhitespaceVisitor) VisitString(node *ast.StringLiteral) interface{}        { return nil }
func (v *WhitespaceVisitor) VisitBoolean(node *ast.BooleanLiteral) interface{}      { return nil }
func (v *WhitespaceVisitor) VisitNumber(node *ast.NumberLiteral) interface{}        { return nil }
func (v *WhitespaceVisitor) VisitHash(node *ast.Hash) interface{}                   { return nil }
func (v *WhitespaceVisitor) VisitHashPair(node *ast.HashPair) interface{}           { return nil }

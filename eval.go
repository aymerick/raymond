package raymond

import (
	"fmt"
	"io"

	"github.com/aymerick/raymond/ast"
)

// Template evaluation visitor
type EvalVisitor struct {
	wr      io.Writer
	tpl     *Template
	curNode ast.Node
}

// Instanciate a new evaluation visitor
func NewEvalVisitor(wr io.Writer, tpl *Template) *EvalVisitor {
	return &EvalVisitor{
		wr:  wr,
		tpl: tpl,
	}
}

// fatal evaluation error
func (v *EvalVisitor) errPanic(err error) {
	panic(fmt.Errorf("Evaluation error: %s\nCurrente node:\n\t%s", err, v.curNode))
}

func (v *EvalVisitor) onNode(node ast.Node) {
	// log.Printf("onNode: %s", node)
	v.curNode = node
}

//
// Visitor interface
//

// Statements

func (v *EvalVisitor) VisitProgram(node *ast.Program) interface{} {
	v.onNode(node)

	for _, n := range node.Body {
		n.Accept(v)
	}

	return nil
}

func (v *EvalVisitor) VisitMustache(node *ast.MustacheStatement) interface{} {
	v.onNode(node)

	// @todo
	return nil
}

func (v *EvalVisitor) VisitBlock(node *ast.BlockStatement) interface{} {
	v.onNode(node)

	// @todo
	return nil
}

func (v *EvalVisitor) VisitPartial(node *ast.PartialStatement) interface{} {
	v.onNode(node)

	// @todo
	return nil
}

func (v *EvalVisitor) VisitContent(node *ast.ContentStatement) interface{} {
	v.onNode(node)

	if _, err := v.wr.Write([]byte(node.Value)); err != nil {
		v.errPanic(err)
	}

	return nil
}

func (v *EvalVisitor) VisitComment(node *ast.CommentStatement) interface{} {
	v.onNode(node)

	// @todo
	return nil
}

// Expressions

func (v *EvalVisitor) VisitSubExpression(node *ast.SubExpression) interface{} {
	v.onNode(node)

	// @todo
	return nil
}

func (v *EvalVisitor) VisitPath(node *ast.PathExpression) interface{} {
	v.onNode(node)

	// @todo
	return nil
}

// Literals

func (v *EvalVisitor) VisitString(node *ast.StringLiteral) interface{} {
	v.onNode(node)

	// @todo
	return nil
}

func (v *EvalVisitor) VisitBoolean(node *ast.BooleanLiteral) interface{} {
	v.onNode(node)

	// @todo
	return nil
}

func (v *EvalVisitor) VisitNumber(node *ast.NumberLiteral) interface{} {
	v.onNode(node)

	// @todo
	return nil
}

// Miscellaneous

func (v *EvalVisitor) VisitHash(node *ast.Hash) interface{} {
	v.onNode(node)

	// @todo
	return nil
}

func (v *EvalVisitor) VisitHashPair(node *ast.HashPair) interface{} {
	v.onNode(node)

	// @todo
	return nil
}

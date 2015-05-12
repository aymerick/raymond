package raymond

import (
	"fmt"
	"io"
	"reflect"

	"github.com/aymerick/raymond/ast"
)

var (
	// @note borrowed from https://github.com/golang/go/tree/master/src/text/template/exec.go
	errorType       = reflect.TypeOf((*error)(nil)).Elem()
	fmtStringerType = reflect.TypeOf((*fmt.Stringer)(nil)).Elem()
)

// Template evaluation visitor
type EvalVisitor struct {
	wr   io.Writer
	tpl  *Template
	data interface{}

	curNode ast.Node
}

// Instanciate a new evaluation visitor
func NewEvalVisitor(wr io.Writer, tpl *Template, data interface{}) *EvalVisitor {
	return &EvalVisitor{
		wr:   wr,
		tpl:  tpl,
		data: data,
	}
}

// fatal evaluation error
func (v *EvalVisitor) errPanic(err error) {
	panic(fmt.Errorf("Evaluation error: %s\nCurrent node:\n\t%s", err, v.curNode))
}

// fatal evaluation error message
func (v *EvalVisitor) errorf(format string, args ...interface{}) {
	v.errPanic(fmt.Errorf(format, args...))
}

// set current node
func (v *EvalVisitor) at(node ast.Node) {
	// log.Printf("at: %s", node)
	v.curNode = node
}

// evaluates field path
func (v *EvalVisitor) evalFieldPath(path *ast.PathExpression) string {
	if v.data == nil {
		return ""
	}

	fieldName := path.Parts[0]
	ctx := reflect.ValueOf(v.data)

	value := v.evalField(ctx, fieldName)

	return v.strValue(value)
}

func (v *EvalVisitor) evalField(ctx reflect.Value, fieldName string) reflect.Value {
	ctxType := ctx.Type()

	switch ctx.Kind() {
	case reflect.Struct:
		tField, ok := ctx.Type().FieldByName(fieldName)
		if !ok {
			v.errorf("%s is not a field of struct type %s", fieldName, ctxType)
		}

		if tField.PkgPath != "" {
			// field is unexported
			v.errorf("%s is an unexported field of struct type %s", fieldName, ctxType)
		}

		return ctx.FieldByIndex(tField.Index)
	case reflect.Map:
		nameVal := reflect.ValueOf(fieldName)
		if nameVal.Type().AssignableTo(ctx.Type().Key()) {
			return ctx.MapIndex(nameVal)
		}
	}

	// wat
	v.errorf("NOT IMPLEMENTED")
	panic("not reached")
}

func (v *EvalVisitor) strValue(value reflect.Value) string {
	val, ok := printableValue(value)
	if !ok {
		v.errorf("Can't print value")
	}

	return fmt.Sprintf("%s", val)
}

// printableValue returns the, possibly indirected, interface value inside v that
// is best for a call to formatted printer.
// @note borrowed from https://github.com/golang/go/tree/master/src/text/template/exec.go
func printableValue(v reflect.Value) (interface{}, bool) {
	if v.Kind() == reflect.Ptr {
		v, _ = indirect(v) // fmt.Fprint handles nil.
	}
	if !v.IsValid() {
		return "<no value>", true
	}

	if !v.Type().Implements(errorType) && !v.Type().Implements(fmtStringerType) {
		if v.CanAddr() && (reflect.PtrTo(v.Type()).Implements(errorType) || reflect.PtrTo(v.Type()).Implements(fmtStringerType)) {
			v = v.Addr()
		} else {
			switch v.Kind() {
			case reflect.Chan, reflect.Func:
				return nil, false
			}
		}
	}
	return v.Interface(), true
}

// indirect returns the item at the end of indirection, and a bool to indicate if it's nil.
// We indirect through pointers and empty interfaces (only) because
// non-empty interfaces have methods we might need.
// @note borrowed from https://github.com/golang/go/tree/master/src/text/template/exec.go
func indirect(v reflect.Value) (rv reflect.Value, isNil bool) {
	for ; v.Kind() == reflect.Ptr || v.Kind() == reflect.Interface; v = v.Elem() {
		if v.IsNil() {
			return v, true
		}
		if v.Kind() == reflect.Interface && v.NumMethod() > 0 {
			break
		}
	}
	return v, false
}

//
// Visitor interface
//

// Statements

func (v *EvalVisitor) VisitProgram(node *ast.Program) interface{} {
	v.at(node)

	for _, n := range node.Body {
		n.Accept(v)
	}

	return nil
}

func (v *EvalVisitor) VisitMustache(node *ast.MustacheStatement) interface{} {
	v.at(node)

	str, _ := node.Expression.Accept(v).(string)
	if _, err := v.wr.Write([]byte(str)); err != nil {
		v.errPanic(err)
	}

	return nil
}

func (v *EvalVisitor) VisitBlock(node *ast.BlockStatement) interface{} {
	v.at(node)

	// @todo
	return nil
}

func (v *EvalVisitor) VisitPartial(node *ast.PartialStatement) interface{} {
	v.at(node)

	// @todo
	return nil
}

func (v *EvalVisitor) VisitContent(node *ast.ContentStatement) interface{} {
	v.at(node)

	if _, err := v.wr.Write([]byte(node.Value)); err != nil {
		v.errPanic(err)
	}

	return nil
}

func (v *EvalVisitor) VisitComment(node *ast.CommentStatement) interface{} {
	v.at(node)

	// @todo
	return nil
}

// Expressions

func (v *EvalVisitor) VisitExpression(node *ast.Expression) interface{} {
	v.at(node)

	// @todo Check if this is an helper

	// so this must be a field
	path := node.FieldPath()
	if path == nil {
		v.errorf("Invalid expression or helper not found.")
	}

	// evaluate field path
	return v.evalFieldPath(path)
}

func (v *EvalVisitor) VisitSubExpression(node *ast.SubExpression) interface{} {
	v.at(node)

	// @todo
	return nil
}

func (v *EvalVisitor) VisitPath(node *ast.PathExpression) interface{} {
	v.at(node)

	// @todo
	return nil
}

// Literals

func (v *EvalVisitor) VisitString(node *ast.StringLiteral) interface{} {
	v.at(node)

	// @todo
	return nil
}

func (v *EvalVisitor) VisitBoolean(node *ast.BooleanLiteral) interface{} {
	v.at(node)

	// @todo
	return nil
}

func (v *EvalVisitor) VisitNumber(node *ast.NumberLiteral) interface{} {
	v.at(node)

	// @todo
	return nil
}

// Miscellaneous

func (v *EvalVisitor) VisitHash(node *ast.Hash) interface{} {
	v.at(node)

	// @todo
	return nil
}

func (v *EvalVisitor) VisitHashPair(node *ast.HashPair) interface{} {
	v.at(node)

	// @todo
	return nil
}

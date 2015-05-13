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

	zero reflect.Value
)

// Template evaluation visitor
type EvalVisitor struct {
	wr   io.Writer
	tpl  *Template
	data interface{}
	ctx  []reflect.Value

	curNode ast.Node
}

// Instanciate a new evaluation visitor
func NewEvalVisitor(wr io.Writer, tpl *Template, data interface{}) *EvalVisitor {
	return &EvalVisitor{
		wr:   wr,
		tpl:  tpl,
		data: data,
		ctx:  []reflect.Value{reflect.ValueOf(data)},
	}
}

func (v *EvalVisitor) pushCtx(ctx reflect.Value) {
	v.ctx = append(v.ctx, ctx)
}

func (v *EvalVisitor) popCtx() reflect.Value {
	var result reflect.Value

	result, v.ctx = v.ctx[len(v.ctx)-1], v.ctx[:len(v.ctx)-1]
	return result
}

func (v *EvalVisitor) curCtx() reflect.Value {
	return v.ctx[len(v.ctx)-1]
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

// evaluates field in given context
func (v *EvalVisitor) evalField(ctx reflect.Value, fieldName string) reflect.Value {
	result := zero

	if !ctx.IsValid() {
		return result
	}

	switch ctx.Kind() {
	case reflect.Struct:
		tField, ok := ctx.Type().FieldByName(fieldName)
		if !ok {
			v.errorf("%s is not a field of struct type %s", fieldName, ctx.Type())
		}

		if tField.PkgPath != "" {
			// field is unexported
			v.errorf("%s is an unexported field of struct type %s", fieldName, ctx.Type())
		}

		// struct field
		result = ctx.FieldByIndex(tField.Index)
	case reflect.Map:
		nameVal := reflect.ValueOf(fieldName)
		if nameVal.Type().AssignableTo(ctx.Type().Key()) {
			// map key
			result = ctx.MapIndex(nameVal)
		}
	}

	return result
}

func (v *EvalVisitor) strValue(value reflect.Value) string {
	result := ""

	ival, ok := printableValue(value)
	if !ok {
		v.errorf("Can't print value")
	}

	val := reflect.ValueOf(ival)

	switch val.Kind() {
	case reflect.Array, reflect.Slice:
		for i := 0; i < val.Len(); i++ {
			result += val.Index(i).String()
		}
	case reflect.Bool:
		s := "false"
		if val.Bool() {
			s = "true"
		}
		result = fmt.Sprintf("%s", s)
	case reflect.Int, reflect.Int8, reflect.Int16, reflect.Int32, reflect.Int64, reflect.Uint, reflect.Uint8, reflect.Uint16, reflect.Uint32, reflect.Uint64, reflect.Uintptr:
		result = fmt.Sprintf("%d", ival)
	case reflect.Float32, reflect.Float64:
		result = fmt.Sprintf("%f", ival)
	case reflect.Invalid:
		result = ""
	default:
		result = fmt.Sprintf("%s", ival)
	}

	// log.Printf("strValue(%q) => %s", ival, result)

	return result
}

// printableValue returns the, possibly indirected, interface value inside v that
// is best for a call to formatted printer.
//
// @note borrowed from https://github.com/golang/go/tree/master/src/text/template/exec.go
func printableValue(v reflect.Value) (interface{}, bool) {
	if v.Kind() == reflect.Ptr {
		v, _ = indirect(v) // fmt.Fprint handles nil.
	}
	if !v.IsValid() {
		return "", true
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
//
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

// isTrue reports whether the value is 'true', in the sense of not the zero of its type,
// and whether the value has a meaningful truth value.
//
// @note borrowed from https://github.com/golang/go/tree/master/src/text/template/exec.go
func isTrue(val reflect.Value) (truth, ok bool) {
	if !val.IsValid() {
		// Something like var x interface{}, never set. It's a form of nil.
		return false, true
	}
	switch val.Kind() {
	case reflect.Array, reflect.Map, reflect.Slice, reflect.String:
		truth = val.Len() > 0
	case reflect.Bool:
		truth = val.Bool()
	case reflect.Complex64, reflect.Complex128:
		truth = val.Complex() != 0
	case reflect.Chan, reflect.Func, reflect.Ptr, reflect.Interface:
		truth = !val.IsNil()
	case reflect.Int, reflect.Int8, reflect.Int16, reflect.Int32, reflect.Int64:
		truth = val.Int() != 0
	case reflect.Float32, reflect.Float64:
		truth = val.Float() != 0
	case reflect.Uint, reflect.Uint8, reflect.Uint16, reflect.Uint32, reflect.Uint64, reflect.Uintptr:
		truth = val.Uint() != 0
	case reflect.Struct:
		truth = true // Struct values are always true.
	default:
		return
	}
	return truth, true
}

// Finds given helper
func (v *EvalVisitor) findHelper(name string) Helper {
	// check template helpers
	if v.tpl.helpers[name] != nil {
		return v.tpl.helpers[name]
	}

	// check global helpers
	return FindHelper(name)
}

// Computes helper parameters from an expression
func (v *EvalVisitor) helperParams(node *ast.Expression) *HelperParams {
	var params []interface{}
	var hash map[string]interface{}

	for _, paramNode := range node.Params {
		param := paramNode.Accept(v)
		params = append(params, param)
	}

	// @todo Fill hash

	return NewHelperParams(params, hash)
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

	// evaluate expression
	val := reflect.ValueOf(node.Expression.Accept(v))

	str := v.strValue(val)

	// write result
	if _, err := v.wr.Write([]byte(str)); err != nil {
		v.errPanic(err)
	}

	return nil
}

func (v *EvalVisitor) VisitBlock(node *ast.BlockStatement) interface{} {
	v.at(node)

	// evaluate expression
	val := reflect.ValueOf(node.Expression.Accept(v))

	v.pushCtx(val)

	truth, _ := isTrue(val)
	if truth && (node.Program != nil) {
		node.Program.Accept(v)
	} else if node.Inverse != nil {
		node.Inverse.Accept(v)
	}

	v.popCtx()

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

	return ""
}

// Expressions

func (v *EvalVisitor) VisitExpression(node *ast.Expression) interface{} {
	v.at(node)

	// check if this is an helper
	if helperName := node.HelperName(); helperName != "" {
		if helper := v.findHelper(helperName); helper != nil {
			// call helper function
			return helper(v.helperParams(node))
		}
	}

	// field path
	if path := node.FieldPath(); path != nil {
		if val := path.Accept(v); val != nil {
			return val
		}

		return nil
	}

	// literal
	if literal, ok := node.LiteralStr(); ok {
		if val := v.evalField(v.curCtx(), literal); val.IsValid() {
			return val.Interface()
		}

		return nil
	}

	return nil
}

func (v *EvalVisitor) VisitSubExpression(node *ast.SubExpression) interface{} {
	v.at(node)

	// @todo
	return nil
}

func (v *EvalVisitor) VisitPath(node *ast.PathExpression) interface{} {
	v.at(node)

	ctx := v.curCtx()

	for i := 0; i < len(node.Parts); i++ {
		ctx = v.evalField(ctx, node.Parts[i])
		if !ctx.IsValid() {
			break
		}
	}

	if !ctx.IsValid() {
		return nil
	}

	return ctx.Interface()
}

// Literals

func (v *EvalVisitor) VisitString(node *ast.StringLiteral) interface{} {
	v.at(node)

	return node.Value
}

func (v *EvalVisitor) VisitBoolean(node *ast.BooleanLiteral) interface{} {
	v.at(node)

	return node.Value
}

func (v *EvalVisitor) VisitNumber(node *ast.NumberLiteral) interface{} {
	v.at(node)

	return node.Number()
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

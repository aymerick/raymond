package raymond

import (
	"bytes"
	"fmt"
	"html"
	"log"
	"reflect"
	"strconv"

	"github.com/aymerick/raymond/ast"
)

var (
	// @note borrowed from https://github.com/golang/go/tree/master/src/text/template/exec.go
	errorType       = reflect.TypeOf((*error)(nil)).Elem()
	fmtStringerType = reflect.TypeOf((*fmt.Stringer)(nil)).Elem()

	zero reflect.Value

	// @todo remove that debug stuff once all tests pass
	VERBOSE_EVAL = false
)

// Template evaluation visitor
type EvalVisitor struct {
	tpl  *Template
	data interface{}
	ctx  []reflect.Value

	curNode ast.Node
	blocks  []*ast.BlockStatement
}

// Instanciate a new evaluation visitor
func NewEvalVisitor(tpl *Template, data interface{}) *EvalVisitor {
	return &EvalVisitor{
		tpl:  tpl,
		data: data,
		ctx:  []reflect.Value{reflect.ValueOf(data)},
	}
}

func (v *EvalVisitor) pushCtx(ctx reflect.Value) {
	if VERBOSE_EVAL {
		log.Printf("Push context: %s", StrValue(ctx))
	}

	v.ctx = append(v.ctx, ctx)
}

func (v *EvalVisitor) popCtx() reflect.Value {
	if len(v.ctx) == 0 {
		return zero
	}

	var result reflect.Value

	result, v.ctx = v.ctx[len(v.ctx)-1], v.ctx[:len(v.ctx)-1]

	if VERBOSE_EVAL {
		log.Printf("Pop context, back to: %s", StrValue(v.curCtx()))
	}

	return result
}

func (v *EvalVisitor) curCtx() reflect.Value {
	if len(v.ctx) == 0 {
		return zero
	}

	return v.ctx[len(v.ctx)-1]
}

func (v *EvalVisitor) pushBlock(block *ast.BlockStatement) {
	v.blocks = append(v.blocks, block)
}

func (v *EvalVisitor) popBlock() *ast.BlockStatement {
	if len(v.blocks) == 0 {
		return nil
	}

	var result *ast.BlockStatement
	result, v.blocks = v.blocks[len(v.blocks)-1], v.blocks[:len(v.blocks)-1]
	return result
}

func (v *EvalVisitor) curBlock() *ast.BlockStatement {
	if len(v.blocks) == 0 {
		return nil
	}

	return v.blocks[len(v.blocks)-1]
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
	if VERBOSE_EVAL {
		log.Printf("at node: %s", node)
	}

	v.curNode = node
}

// Evaluate node with given context and returns string result
func (v *EvalVisitor) evalNodeWith(node ast.Node, ctx reflect.Value) string {
	v.pushCtx(ctx)
	result, _ := node.Accept(v).(string)
	v.popCtx()
	return result
}

// evaluates field in given context
func (v *EvalVisitor) evalField(ctx reflect.Value, fieldName string) reflect.Value {
	result := zero

	ctx, _ = indirect(ctx)
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

	if VERBOSE_EVAL {
		log.Printf("evalField(): '%s' with context %s => %s", fieldName, StrValue(ctx), StrValue(result))
	}

	return result
}

// returns string representation of a `interface{}`
func Str(value interface{}) string {
	return StrValue(reflect.ValueOf(value))
}

// returns string representation of a `reflect.Value`
func StrValue(value reflect.Value) string {
	result := ""

	ival, ok := printableValue(value)
	if !ok {
		panic("Can't print value")
	}

	val := reflect.ValueOf(ival)

	switch val.Kind() {
	case reflect.Array, reflect.Slice:
		for i := 0; i < val.Len(); i++ {
			result += StrValue(val.Index(i))
		}
	case reflect.Bool:
		result = "false"
		if val.Bool() {
			result = "true"
		}
	case reflect.Int, reflect.Int8, reflect.Int16, reflect.Int32, reflect.Int64, reflect.Uint, reflect.Uint8, reflect.Uint16, reflect.Uint32, reflect.Uint64, reflect.Uintptr:
		result = fmt.Sprintf("%d", ival)
	case reflect.Float32, reflect.Float64:
		result = strconv.FormatFloat(val.Float(), 'f', -1, 64)
	case reflect.Invalid:
		result = ""
	default:
		result = fmt.Sprintf("%s", ival)
	}

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

// IsTruth reports whether the value is 'true', in the sense of not the zero of its type,
// and whether the value has a meaningful truth value.
//
// @note borrowed from https://github.com/golang/go/tree/master/src/text/template/exec.go
func IsTruth(val reflect.Value) (truth, ok bool) {
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

// Returns true if given expression is a helper call
func (v *EvalVisitor) isHelperCall(node *ast.Expression) bool {
	if helperName := node.HelperName(); helperName != "" {
		return v.findHelper(helperName) != nil
	}
	return false
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

	if node.Hash != nil {
		hash, _ = node.Hash.Accept(v).(map[string]interface{})
	}

	return NewHelperParams(v, params, hash)
}

//
// Visitor interface
//

// Statements

func (v *EvalVisitor) VisitProgram(node *ast.Program) interface{} {
	v.at(node)

	buf := new(bytes.Buffer)

	for _, n := range node.Body {
		if str := Str(n.Accept(v)); str != "" {
			if _, err := buf.Write([]byte(str)); err != nil {
				v.errPanic(err)
			}
		}
	}

	return buf.String()
}

func (v *EvalVisitor) VisitMustache(node *ast.MustacheStatement) interface{} {
	v.at(node)

	// evaluate expression
	expr := node.Expression.Accept(v)

	// get string value
	str := Str(expr)
	if !node.Unescaped {
		// escape html
		str = html.EscapeString(str)
	}

	return str
}

func (v *EvalVisitor) VisitBlock(node *ast.BlockStatement) interface{} {
	v.at(node)
	v.pushBlock(node)

	result := ""

	// evaluate expression
	expr := node.Expression.Accept(v)

	if v.isHelperCall(node.Expression) {
		result, _ = expr.(string)
	} else {
		val := reflect.ValueOf(expr)

		truth, _ := IsTruth(val)
		if truth {
			if node.Program != nil {
				if VERBOSE_EVAL {
					log.Printf("VisitBlock(): Truthy, visiting Program")
				}

				switch val.Kind() {
				case reflect.Array, reflect.Slice:
					// Array context
					for i := 0; i < val.Len(); i++ {
						result += v.evalNodeWith(node.Program, val.Index(i))
					}
				default:
					// NOT array
					result = v.evalNodeWith(node.Program, val)
				}
			}
		} else if node.Inverse != nil {
			if VERBOSE_EVAL {
				log.Printf("VisitBlock(): Falsy, visiting Inverse")
			}

			result, _ = node.Inverse.Accept(v).(string)
		}
	}

	v.popBlock()

	return result
}

func (v *EvalVisitor) VisitPartial(node *ast.PartialStatement) interface{} {
	v.at(node)

	// @todo
	return nil
}

func (v *EvalVisitor) VisitContent(node *ast.ContentStatement) interface{} {
	v.at(node)

	// write content as is
	return node.Value
}

func (v *EvalVisitor) VisitComment(node *ast.CommentStatement) interface{} {
	v.at(node)

	// ignore comments
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

		// invalid field path
		return nil
	}

	// literal
	if literal, ok := node.LiteralStr(); ok {
		if val := v.evalField(v.curCtx(), literal); val.IsValid() {
			return val.Interface()
		}

		return nil
	}

	// wat
	return nil
}

func (v *EvalVisitor) VisitSubExpression(node *ast.SubExpression) interface{} {
	v.at(node)

	// @todo
	return nil
}

func (v *EvalVisitor) VisitPath(node *ast.PathExpression) interface{} {
	v.at(node)

	var result interface{}

	if node.Depth > len(v.ctx) {
		return nil
	}

	// go back to parent context
	var prevCtxList []reflect.Value
	for i := node.Depth; i > 0; i-- {
		prevCtxList = append(prevCtxList, v.popCtx())
	}

	// get current context
	ctx := v.curCtx()

	if VERBOSE_EVAL {
		log.Printf("VisitPath(): %s with context '%s'", node.Original, StrValue(ctx))
	}

	switch ctx.Kind() {
	case reflect.Array, reflect.Slice:
		// Array context
		var results []interface{}

		for i := 0; i < ctx.Len(); i++ {
			value := v.evalPathParts(ctx.Index(i), node.Parts)
			if value.IsValid() {
				results = append(results, value.Interface())
			}
			// else raise ?
		}

		result = results
	default:
		// NOT array context
		value := v.evalPathParts(ctx, node.Parts)
		if value.IsValid() {
			result = value.Interface()
		}
	}

	// set back contexts
	for i := len(prevCtxList); i > 0; i-- {
		var prev reflect.Value
		prev, prevCtxList = prevCtxList[len(prevCtxList)-1], prevCtxList[:len(prevCtxList)-1]
		v.pushCtx(prev)
	}

	if VERBOSE_EVAL {
		log.Printf("VisitPath(): result => %s", Str(result))
	}

	return result
}

func (v *EvalVisitor) evalPathParts(ctx reflect.Value, parts []string) reflect.Value {
	for i := 0; i < len(parts); i++ {
		part := parts[i]

		// "[foo bar]"" => "foo bar"
		if (len(part) >= 2) && (part[0] == '[') && (part[len(part)-1] == ']') {
			part = part[1 : len(part)-1]
		}

		ctx = v.evalField(ctx, part)
		if !ctx.IsValid() {
			break
		}
	}

	return ctx
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

	result := make(map[string]interface{})

	for _, pair := range node.Pairs {
		if value := pair.Accept(v); value != nil {
			result[pair.Key] = value
		}
	}

	return result
}

func (v *EvalVisitor) VisitHashPair(node *ast.HashPair) interface{} {
	v.at(node)

	return node.Val.Accept(v)
}

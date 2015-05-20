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

	ctx     []reflect.Value
	blocks  []*ast.BlockStatement
	exprs   []*ast.Expression
	exprCtx []reflect.Value

	curNode ast.Node
}

// Instanciate a new evaluation visitor
func NewEvalVisitor(tpl *Template, data interface{}) *EvalVisitor {
	return &EvalVisitor{
		tpl:  tpl,
		data: data,
		ctx:  []reflect.Value{reflect.ValueOf(data)},
	}
}

//
// Contexts stack
//

// push new context
func (v *EvalVisitor) pushCtx(ctx reflect.Value) {
	if VERBOSE_EVAL {
		log.Printf("Push context: %s", StrValue(ctx))
	}

	v.ctx = append(v.ctx, ctx)
}

// pop last context
func (v *EvalVisitor) popCtx() reflect.Value {
	if len(v.ctx) == 0 {
		return zero
	}

	var result reflect.Value

	result, v.ctx = v.ctx[len(v.ctx)-1], v.ctx[:len(v.ctx)-1]

	if VERBOSE_EVAL {
		log.Printf("Pop context, current is: %s", StrValue(v.curCtx()))
	}

	return result
}

// returns current context
func (v *EvalVisitor) curCtx() reflect.Value {
	if len(v.ctx) == 0 {
		return zero
	}

	return v.ctx[len(v.ctx)-1]
}

//
// Blocks stack
//

// push new block statement
func (v *EvalVisitor) pushBlock(block *ast.BlockStatement) {
	if VERBOSE_EVAL {
		log.Printf("Push block: %s", Str(block))
	}

	v.blocks = append(v.blocks, block)
}

// pop last block statement
func (v *EvalVisitor) popBlock() *ast.BlockStatement {
	if len(v.blocks) == 0 {
		return nil
	}

	var result *ast.BlockStatement
	result, v.blocks = v.blocks[len(v.blocks)-1], v.blocks[:len(v.blocks)-1]

	if VERBOSE_EVAL {
		log.Printf("Pop block, current is: %s", Str(v.curBlock()))
	}

	return result
}

// returns current block statement
func (v *EvalVisitor) curBlock() *ast.BlockStatement {
	if len(v.blocks) == 0 {
		return nil
	}

	return v.blocks[len(v.blocks)-1]
}

//
// Expressions stack
//

// push new expression
func (v *EvalVisitor) pushExpr(expression *ast.Expression) {
	if VERBOSE_EVAL {
		log.Printf("Push expression: %s", Str(expression))
	}

	v.exprs = append(v.exprs, expression)
}

// pop last expression
func (v *EvalVisitor) popExpr() *ast.Expression {
	if len(v.exprs) == 0 {
		return nil
	}

	var result *ast.Expression
	result, v.exprs = v.exprs[len(v.exprs)-1], v.exprs[:len(v.exprs)-1]

	if VERBOSE_EVAL {
		log.Printf("Pop expression, current is: %s", Str(v.curExpr()))
	}

	return result
}

// returns current expression
func (v *EvalVisitor) curExpr() *ast.Expression {
	if len(v.exprs) == 0 {
		return nil
	}

	return v.exprs[len(v.exprs)-1]
}

//
// Expressions context stack
//
// This the stack representing previous context for current expression
//
// This is needed to support `{{#with frank}}{{../awesome .}}{{/with}}` where '../awesome' is a function call
// When evaluating '../awesome' we go back to parent ctx, but when evaluating to '.' we must use previous ctx.
// That's that previous ctx we are storing in that stack.
//
// @todo THIS IS BUGGY ! We should use a linked list of contexts to be sure we can always access an ancestor context.
//

// push new expression context
func (v *EvalVisitor) pushExprCtx(ctx reflect.Value) {
	if VERBOSE_EVAL {
		log.Printf("Push expression context: %s", StrValue(ctx))
	}

	v.exprCtx = append(v.exprCtx, ctx)
}

// pop last expression context
func (v *EvalVisitor) popExprCtx() reflect.Value {
	if len(v.exprCtx) == 0 {
		return zero
	}

	var result reflect.Value

	result, v.exprCtx = v.exprCtx[len(v.exprCtx)-1], v.exprCtx[:len(v.exprCtx)-1]

	if VERBOSE_EVAL {
		log.Printf("Pop expression context, current is: %s", StrValue(v.curExprCtx()))
	}

	return result
}

// returns current expression context
func (v *EvalVisitor) curExprCtx() reflect.Value {
	if len(v.exprCtx) == 0 {
		return zero
	}

	return v.exprCtx[len(v.exprCtx)-1]
}

//
// Error functions
//

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

//
// Evaluation
//

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

	// check if result is a function
	result, _ = indirect(result)
	if result.Kind() == reflect.Func {
		result = v.evalFunc(result)
	}

	if VERBOSE_EVAL {
		log.Printf("evalField(): '%s' with context %s => %s | Kind: %s", fieldName, StrValue(ctx), StrValue(result), result.Kind())
	}

	return result
}

// evaluates a function
func (v *EvalVisitor) evalFunc(funcVal reflect.Value) reflect.Value {
	funcType := funcVal.Type()

	// @todo There should be a better way to get the string type
	strType := reflect.TypeOf("")

	if (funcType.NumOut() != 1) || !strType.AssignableTo(funcType.Out(0)) {
		v.errorf("Function must return a uniq string value: %q", funcVal)
	}

	if funcType.NumIn() > 1 {
		v.errorf("Function can only have a uniq argument: %q", funcVal)
	}

	args := []reflect.Value{}
	if funcType.NumIn() == 1 {
		// create helper argument
		arg := v.HelperArg(v.curExpr())

		if !reflect.TypeOf(arg).AssignableTo(funcType.In(0)) {
			v.errorf("Function argument must be a *HelperArg: %q", funcVal)
		}

		args = append(args, reflect.ValueOf(arg))
	}

	// call function
	resArr := funcVal.Call(args)

	// we already checked that func returns only one value
	return resArr[0]
}

// evaluates all path parts
func (v *EvalVisitor) evalPath(ctx reflect.Value, parts []string) reflect.Value {
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

//
// Stringification
//

// returns string representation of a `interface{}`
func Str(value interface{}) string {
	return StrValue(reflect.ValueOf(value))
}

// returns string representation of a `reflect.Value`
func StrValue(value reflect.Value) string {
	result := ""

	ival, ok := printableValue(value)
	if !ok {
		panic(fmt.Errorf("Can't print value: %q", value))
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

//
// Helpers
//

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

// Computes helper argument from an expression
func (v *EvalVisitor) HelperArg(node *ast.Expression) *HelperArg {
	var params []interface{}
	var hash map[string]interface{}

	withDepthPath := false
	if path := node.FieldPath(); path != nil && path.Depth > 0 {
		withDepthPath = true
	}

	if withDepthPath {
		v.pushCtx(v.curExprCtx())
	}

	for _, paramNode := range node.Params {
		param := paramNode.Accept(v)
		params = append(params, param)
	}

	if node.Hash != nil {
		hash, _ = node.Hash.Accept(v).(map[string]interface{})
	}

	if withDepthPath {
		v.popCtx()
	}

	return NewHelperArg(v, params, hash)
}

//
// Misc
//

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

	var result interface{}
	done := false

	v.pushExpr(node)
	v.pushExprCtx(v.curCtx())

	// check if this is an helper
	if helperName := node.HelperName(); helperName != "" {
		if helper := v.findHelper(helperName); helper != nil {
			// call helper function
			result = helper(v.HelperArg(node))
			done = true
		}
	}

	if !done {
		// field path
		if path := node.FieldPath(); path != nil {
			if val := path.Accept(v); val != nil {
				result = val
			}

			// invalid field path
			done = true
		}
	}

	if !done {
		// literal
		if literal, ok := node.LiteralStr(); ok {
			if val := v.evalField(v.curCtx(), literal); val.IsValid() {
				result = val.Interface()
			}

			done = true
		}
	}

	v.popExprCtx()
	v.popExpr()

	return result
}

func (v *EvalVisitor) VisitSubExpression(node *ast.SubExpression) interface{} {
	v.at(node)

	return node.Expression.Accept(v)
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
		log.Printf("VisitPath(): '%s' with context '%s'", node.Original, StrValue(ctx))
	}

	switch ctx.Kind() {
	case reflect.Array, reflect.Slice:
		// Array context
		var results []interface{}

		for i := 0; i < ctx.Len(); i++ {
			value := v.evalPath(ctx.Index(i), node.Parts)
			if value.IsValid() {
				results = append(results, value.Interface())
			}
			// else raise ?
		}

		result = results
	default:
		// NOT array context
		value := v.evalPath(ctx, node.Parts)
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
		log.Printf("VisitPath(): result => '%s'", Str(result))
	}

	return result
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

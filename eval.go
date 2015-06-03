package raymond

import (
	"bytes"
	"fmt"
	"reflect"
	"strconv"
	"strings"

	"github.com/aymerick/raymond/ast"
)

var (
	// @note borrowed from https://github.com/golang/go/tree/master/src/text/template/exec.go
	errorType       = reflect.TypeOf((*error)(nil)).Elem()
	fmtStringerType = reflect.TypeOf((*fmt.Stringer)(nil)).Elem()

	zero reflect.Value
)

// EvalVisitor evaluates a handlebars template with context
type EvalVisitor struct {
	tpl *Template

	// contexts stack
	ctx []reflect.Value

	// current data frame (chained with parent)
	dataFrame *DataFrame

	// block parameters stack
	blockParams []map[string]interface{}

	// block statements stack
	blocks []*ast.BlockStatement

	// expressions stack
	exprs []*ast.Expression

	// memoize expressions that were function calls
	exprFunc map[*ast.Expression]bool

	// used for info on panic
	curNode ast.Node
}

// NewEvalVisitor instanciate a new evaluation visitor with given context and initial private data frame
//
// If privData is nil, then a default data frame is created
func NewEvalVisitor(tpl *Template, ctx interface{}, privData *DataFrame) *EvalVisitor {
	frame := privData
	if frame == nil {
		frame = NewDataFrame()
	}

	return &EvalVisitor{
		tpl:       tpl,
		ctx:       []reflect.Value{reflect.ValueOf(ctx)},
		dataFrame: frame,
		exprFunc:  make(map[*ast.Expression]bool),
	}
}

// at sets current node
func (v *EvalVisitor) at(node ast.Node) {
	v.curNode = node
}

//
// Contexts stack
//

// pushCtx pushes new context to the stack
func (v *EvalVisitor) pushCtx(ctx reflect.Value) {
	v.ctx = append(v.ctx, ctx)
}

// popCtx pops last context from stack
func (v *EvalVisitor) popCtx() reflect.Value {
	if len(v.ctx) == 0 {
		return zero
	}

	var result reflect.Value
	result, v.ctx = v.ctx[len(v.ctx)-1], v.ctx[:len(v.ctx)-1]

	return result
}

// rootCtx returns root context
func (v *EvalVisitor) rootCtx() reflect.Value {
	return v.ctx[0]
}

// curCtx returns current context
func (v *EvalVisitor) curCtx() reflect.Value {
	return v.ancestorCtx(0)
}

// ancestorCtx returns ancestor context
func (v *EvalVisitor) ancestorCtx(depth int) reflect.Value {
	index := len(v.ctx) - 1 - depth
	if index < 0 {
		return zero
	}

	return v.ctx[index]
}

//
// Private data frame
//

// setDataFrame sets new data frame
func (v *EvalVisitor) setDataFrame(frame *DataFrame) {
	v.dataFrame = frame
}

// popDataFrame sets back parent data frame
func (v *EvalVisitor) popDataFrame() {
	v.dataFrame = v.dataFrame.parent
}

//
// Block Parameters stack
//

// pushBlockParams pushes new block params to the stack
func (v *EvalVisitor) pushBlockParams(params map[string]interface{}) {
	v.blockParams = append(v.blockParams, params)
}

// popBlockParams pops last block params from stack
func (v *EvalVisitor) popBlockParams() map[string]interface{} {
	var result map[string]interface{}

	if len(v.blockParams) == 0 {
		return result
	}

	result, v.blockParams = v.blockParams[len(v.blockParams)-1], v.blockParams[:len(v.blockParams)-1]
	return result
}

// blockParam iterates on stack to find given block parameter, and returns its value or nil if not founc
func (v *EvalVisitor) blockParam(name string) interface{} {
	for i := len(v.blockParams) - 1; i >= 0; i-- {
		for k, v := range v.blockParams[i] {
			if name == k {
				return v
			}
		}
	}

	return nil
}

//
// Blocks stack
//

// pushBlock pushes new block statement to stack
func (v *EvalVisitor) pushBlock(block *ast.BlockStatement) {
	v.blocks = append(v.blocks, block)
}

// popBlock pops last block statement from stack
func (v *EvalVisitor) popBlock() *ast.BlockStatement {
	if len(v.blocks) == 0 {
		return nil
	}

	var result *ast.BlockStatement
	result, v.blocks = v.blocks[len(v.blocks)-1], v.blocks[:len(v.blocks)-1]

	return result
}

// curBlock returns current block statement
func (v *EvalVisitor) curBlock() *ast.BlockStatement {
	if len(v.blocks) == 0 {
		return nil
	}

	return v.blocks[len(v.blocks)-1]
}

//
// Expressions stack
//

// pushExpr pushes new expression to stack
func (v *EvalVisitor) pushExpr(expression *ast.Expression) {
	v.exprs = append(v.exprs, expression)
}

// popExpr pops last expression from stack
func (v *EvalVisitor) popExpr() *ast.Expression {
	if len(v.exprs) == 0 {
		return nil
	}

	var result *ast.Expression
	result, v.exprs = v.exprs[len(v.exprs)-1], v.exprs[:len(v.exprs)-1]

	return result
}

// curExpr returns current expression
func (v *EvalVisitor) curExpr() *ast.Expression {
	if len(v.exprs) == 0 {
		return nil
	}

	return v.exprs[len(v.exprs)-1]
}

//
// Error functions
//

// errPanic panics
func (v *EvalVisitor) errPanic(err error) {
	panic(fmt.Errorf("Evaluation error: %s\nCurrent node:\n\t%s", err, v.curNode))
}

// errorf panics with a custom message
func (v *EvalVisitor) errorf(format string, args ...interface{}) {
	v.errPanic(fmt.Errorf(format, args...))
}


//
// Evaluation
//

// evalProgram eEvaluates program with given context and returns string result
func (v *EvalVisitor) evalProgram(program *ast.Program, ctx interface{}, data *DataFrame, key interface{}) string {
	blockParams := make(map[string]interface{})

	// compute block params
	if len(program.BlockParams) > 0 {
		blockParams[program.BlockParams[0]] = ctx
	}

	if (len(program.BlockParams) > 1) && (key != nil) {
		blockParams[program.BlockParams[1]] = key
	}

	// push contexts
	if len(blockParams) > 0 {
		v.pushBlockParams(blockParams)
	}

	ctxVal := reflect.ValueOf(ctx)
	if ctxVal.IsValid() {
		v.pushCtx(ctxVal)
	}

	if data != nil {
		v.setDataFrame(data)
	}

	// evaluate program
	result, _ := program.Accept(v).(string)

	// pop contexts
	if data != nil {
		v.popDataFrame()
	}

	if ctxVal.IsValid() {
		v.popCtx()
	}

	if len(blockParams) > 0 {
		v.popBlockParams()
	}

	return result
}

// evalPath evaluates all path parts with given context
func (v *EvalVisitor) evalPath(ctx reflect.Value, parts []string, exprRoot bool) (reflect.Value, bool) {
	partResolved := false

	for i := 0; i < len(parts); i++ {
		part := parts[i]

		// "[foo bar]"" => "foo bar"
		if (len(part) >= 2) && (part[0] == '[') && (part[len(part)-1] == ']') {
			part = part[1 : len(part)-1]
		}

		ctx = v.evalField(ctx, part, exprRoot)
		if !ctx.IsValid() {
			break
		}

		// we resolved at least one part of path
		partResolved = true
	}

	return ctx, partResolved
}

// evalField evaluates field with given context
func (v *EvalVisitor) evalField(ctx reflect.Value, fieldName string, exprRoot bool) reflect.Value {
	result := zero

	ctx, _ = indirect(ctx)
	if !ctx.IsValid() {
		return result
	}

	switch ctx.Kind() {
	case reflect.Struct:
		// check if struct have this field and that it is exported
		if tField, ok := ctx.Type().FieldByName(fieldName); ok && (tField.PkgPath == "") {
			// struct field
			result = ctx.FieldByIndex(tField.Index)
		}
	case reflect.Map:
		nameVal := reflect.ValueOf(fieldName)
		if nameVal.Type().AssignableTo(ctx.Type().Key()) {
			// map key
			result = ctx.MapIndex(nameVal)
		}
	case reflect.Array, reflect.Slice:
		if i, err := strconv.Atoi(fieldName); (err == nil) && (i < ctx.Len()) {
			result = ctx.Index(i)
		}
	}

	// check if result is a function
	result, _ = indirect(result)
	if result.Kind() == reflect.Func {
		// in that code path, we know we can't be an expression root
		result = v.evalFunc(result, exprRoot)
	}

	return result
}

// evalFunc evaluates given function
func (v *EvalVisitor) evalFunc(funcVal reflect.Value, exprRoot bool) reflect.Value {
	funcType := funcVal.Type()

	if funcType.NumOut() != 1 {
		v.errorf("Function must return a uniq value: %q", funcVal)
	}

	if funcType.NumIn() > 1 {
		v.errorf("Function can only have a uniq argument: %q", funcVal)
	}

	args := []reflect.Value{}
	if funcType.NumIn() == 1 {
		var arg *HelperArg

		if exprRoot {
			// create function arg with all params/hash
			expr := v.curExpr()
			arg = v.helperArg(expr)

			// ok, that expression was a function call
			v.exprFunc[expr] = true
		} else {
			// we are not at root of expression, so we are a parameter... and we don't like
			// infinite loops caused by trying to parse ourself forever
			arg = newEmptyHelperArg(v)
		}

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

// findBlockParam returns node's block parameter
func (v *EvalVisitor) findBlockParam(node *ast.PathExpression) (string, interface{}) {
	if len(node.Parts) > 0 {
		name := node.Parts[0]
		if value := v.blockParam(name); value != nil {
			return name, value
		}
	}

	return "", nil
}

// evalPathExpression evaluates a path expression
func (v *EvalVisitor) evalPathExpression(node *ast.PathExpression, exprRoot bool) interface{} {
	var result interface{}

	if name, value := v.findBlockParam(node); value != nil {
		// block parameter value

		// We push a new context so we can evaluate the path expression (note: this may be a bad idea).
		//
		// Example:
		//   {{#foo as |bar|}}
		//     {{bar.baz}}
		//   {{/foo}}
		//
		// With data:
		//   {"foo": {"baz": "bat"}}
		newCtx := map[string]interface{}{name: value}

		v.pushCtx(reflect.ValueOf(newCtx))
		result = v.evalCtxPathExpression(node, exprRoot)
		v.popCtx()
	} else {
		ctxTried := false

		if node.IsDataRoot() {
			// context path
			result = v.evalCtxPathExpression(node, exprRoot)

			ctxTried = true
		}

		if (result == nil) && node.Data {
			// if it is @root, then we tried to evaluate with root context but nothing was found
			// so let's try with private data

			// private data
			result = v.evalDataPathExpression(node, exprRoot)
		}

		if (result == nil) && !ctxTried {
			// context path
			result = v.evalCtxPathExpression(node, exprRoot)
		}
	}

	return result
}

// evalDataPathExpression evaluates a private data path expression
func (v *EvalVisitor) evalDataPathExpression(node *ast.PathExpression, exprRoot bool) interface{} {
	// find data frame
	frame := v.dataFrame
	for i := node.Depth; i > 0; i-- {
		if frame.parent == nil {
			return nil
		}
		frame = frame.parent
	}

	// resolve data
	// @note Can be changed to v.evalCtx() as context can't be an array
	result, _ := v.evalCtxPath(reflect.ValueOf(frame.data), node.Parts, exprRoot)
	return result
}

// evalCtxPathExpression evaluates a context path expression
func (v *EvalVisitor) evalCtxPathExpression(node *ast.PathExpression, exprRoot bool) interface{} {
	v.at(node)

	if node.IsDataRoot() {
		// `@root` - remove the first part
		parts := node.Parts[1:len(node.Parts)]

		result, _ := v.evalCtxPath(v.rootCtx(), parts, exprRoot)
		return result
	}

	return v.evalDepthPath(node.Depth, node.Parts, exprRoot)
}

// evalDepthPath iterates on contexts, starting at given depth, until there is one that resolve given path parts
func (v *EvalVisitor) evalDepthPath(depth int, parts []string, exprRoot bool) interface{} {
	var result interface{}
	partResolved := false

	ctx := v.ancestorCtx(depth)

	for (result == nil) && ctx.IsValid() && (depth <= len(v.ctx) && !partResolved) {
		// try with context
		result, partResolved = v.evalCtxPath(ctx, parts, exprRoot)

		// As soon as we find the first part of a path, we must not try to resolve with parent context if result is finally `nil`
		// Reference: "Dotted Names - Context Precedence" mustache test
		if !partResolved && (result == nil) {
			// try with previous context
			depth++
			ctx = v.ancestorCtx(depth)
		}
	}

	return result
}

// evalCtxPath evaluates path with given context
func (v *EvalVisitor) evalCtxPath(ctx reflect.Value, parts []string, exprRoot bool) (interface{}, bool) {
	var result interface{}
	partResolved := false

	switch ctx.Kind() {
	case reflect.Array, reflect.Slice:
		// Array context
		var results []interface{}

		for i := 0; i < ctx.Len(); i++ {
			value, _ := v.evalPath(ctx.Index(i), parts, exprRoot)
			if value.IsValid() {
				results = append(results, value.Interface())
			}
		}

		result = results
	default:
		// NOT array context
		var value reflect.Value

		value, partResolved = v.evalPath(ctx, parts, exprRoot)
		if value.IsValid() {
			result = value.Interface()
		}
	}

	return result, partResolved
}

//
// Helpers
//

// isHelperCall returns true if given expression is a helper call
func (v *EvalVisitor) isHelperCall(node *ast.Expression) bool {
	if helperName := node.HelperName(); helperName != "" {
		return v.findHelper(helperName) != nil
	}
	return false
}

// findHelper finds given helper
func (v *EvalVisitor) findHelper(name string) Helper {
	// check template helpers
	if v.tpl.helpers[name] != nil {
		return v.tpl.helpers[name]
	}

	// check global helpers
	return FindHelper(name)
}

// helperArg computes helper argument from an expression
func (v *EvalVisitor) helperArg(node *ast.Expression) *HelperArg {
	var params []interface{}
	var hash map[string]interface{}

	for _, paramNode := range node.Params {
		param := paramNode.Accept(v)
		params = append(params, param)
	}

	if node.Hash != nil {
		hash, _ = node.Hash.Accept(v).(map[string]interface{})
	}

	return newHelperArg(v, params, hash)
}

//
// Partials
//

// findPartial finds given partial
func (v *EvalVisitor) findPartial(name string) *Partial {
	// check template partials
	if v.tpl.partials[name] != nil {
		return v.tpl.partials[name]
	}

	// check global partials
	return FindPartial(name)
}

// partialContext computes partial context
func (v *EvalVisitor) partialContext(node *ast.PartialStatement) reflect.Value {
	if nb := len(node.Params); nb > 1 {
		v.errorf("Unsupported number of partial arguments: %d", nb)
	}

	if (len(node.Params) > 0) && (node.Hash != nil) {
		v.errorf("Passing both context and named parameters to a partial is not allowed")
	}

	if len(node.Params) == 1 {
		return reflect.ValueOf(node.Params[0].Accept(v))
	}

	if node.Hash != nil {
		hash, _ := node.Hash.Accept(v).(map[string]interface{})
		return reflect.ValueOf(hash)
	}

	return zero
}

// evalPartial evaluates a partial
func (v *EvalVisitor) evalPartial(partial *Partial, node *ast.PartialStatement) string {
	// get partial template
	partialTpl, err := partial.Template()
	if err != nil {
		v.errPanic(err)
	}

	// push partial context
	ctx := v.partialContext(node)
	if ctx.IsValid() {
		v.pushCtx(ctx)
	}

	// evaluate partial template
	result, _ := partialTpl.program.Accept(v).(string)

	// ident partial
	result = indentLines(result, node.Indent)

	if ctx.IsValid() {
		v.popCtx()
	}

	return result
}

// indentLines indents all lines of given string
func indentLines(str string, indent string) string {
	if indent == "" {
		return str
	}

	var indented []string

	lines := strings.Split(str, "\n")
	for i, line := range lines {
		if (i == (len(lines) - 1)) && (line == "") {
			// input string ends with a new line
			indented = append(indented, line)
		} else {
			indented = append(indented, indent+line)
		}
	}

	return strings.Join(indented, "\n")
}

//
// Functions
//

// wasFuncCall returns true if given expression was a function call
func (v *EvalVisitor) wasFuncCall(node *ast.Expression) bool {
	// check if expression was tagged as a function call
	return v.exprFunc[node]
}

//
// Misc
//

// indirect returns the item at the end of indirection, and a bool to indicate if it's nil.
// We indirect through pointers and empty interfaces (only) because
// non-empty interfaces have methods we might need.
//
// NOTE: borrowed from https://github.com/golang/go/tree/master/src/text/template/exec.go
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
// NOTE: borrowed from https://github.com/golang/go/tree/master/src/text/template/exec.go
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

	// check if this is a safe string
	isSafe := IsSafeString(expr)

	// get string value
	str := Str(expr)
	if !isSafe && !node.Unescaped {
		// escape html
		str = EscapeString(str)
	}

	return str
}

func (v *EvalVisitor) VisitBlock(node *ast.BlockStatement) interface{} {
	v.at(node)

	v.pushBlock(node)

	result := ""

	// evaluate expression
	expr := node.Expression.Accept(v)

	if v.isHelperCall(node.Expression) || v.wasFuncCall(node.Expression) {
		// it is the responsability of the helper/function to evaluate block
		result, _ = expr.(string)
	} else {
		val := reflect.ValueOf(expr)

		truth, _ := IsTruth(val)
		if truth {
			if node.Program != nil {
				switch val.Kind() {
				case reflect.Array, reflect.Slice:
					// Array context
					for i := 0; i < val.Len(); i++ {
						// Computes new private data frame
						frame := v.dataFrame.NewIterDataFrame(val.Len(), i, nil)

						// Evaluate program
						result += v.evalProgram(node.Program, val.Index(i).Interface(), frame, i)
					}
				default:
					// NOT array
					result = v.evalProgram(node.Program, expr, nil, nil)
				}
			}
		} else if node.Inverse != nil {
			result, _ = node.Inverse.Accept(v).(string)
		}
	}

	v.popBlock()

	return result
}

func (v *EvalVisitor) VisitPartial(node *ast.PartialStatement) interface{} {
	v.at(node)

	// partialName: helperName | sexpr
	name, ok := ast.HelperNameStr(node.Name)
	if !ok {
		if subExpr, ok := node.Name.(*ast.SubExpression); ok {
			name, _ = subExpr.Accept(v).(string)
		}
	}

	if name == "" {
		v.errorf("Unexpected partial name: %q", node.Name)
	}

	partial := v.findPartial(name)
	if partial == nil {
		v.errorf("Partial not found: %s", name)
	}

	return v.evalPartial(partial, node)
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

	// helper call
	if helperName := node.HelperName(); helperName != "" {
		if helper := v.findHelper(helperName); helper != nil {
			result = helper(v.helperArg(node))
			done = true
		}
	}

	if !done {
		// literal
		if literal, ok := node.LiteralStr(); ok {
			if val := v.evalField(v.curCtx(), literal, true); val.IsValid() {
				result = val.Interface()
				done = true
			}
		}
	}

	if !done {
		// field path
		if path := node.FieldPath(); path != nil {
			// @todo Find a cleaner way ! Don't break the pattern !
			// this is an exception to visitor pattern, because we need to pass the info
			// that this path is at root of current expression
			if val := v.evalPathExpression(path, true); val != nil {
				result = val
			}
		}
	}

	v.popExpr()

	return result
}

func (v *EvalVisitor) VisitSubExpression(node *ast.SubExpression) interface{} {
	v.at(node)

	return node.Expression.Accept(v)
}

func (v *EvalVisitor) VisitPath(node *ast.PathExpression) interface{} {
	return v.evalPathExpression(node, false)
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

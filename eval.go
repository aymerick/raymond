package raymond

import (
	"bytes"
	"fmt"
	"log"
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

	// @todo remove that debug stuff once all tests pass
	VERBOSE_EVAL = false
)

// Template evaluation visitor
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

// Instanciate a new evaluation visitor
func NewEvalVisitor(tpl *Template, data interface{}) *EvalVisitor {
	return &EvalVisitor{
		tpl:       tpl,
		ctx:       []reflect.Value{reflect.ValueOf(data)},
		dataFrame: NewDataFrame(),
		exprFunc:  make(map[*ast.Expression]bool),
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
	return v.ancestorCtx(0)
}

// get ancestor context
func (v *EvalVisitor) ancestorCtx(depth int) reflect.Value {
	index := len(v.ctx) - 1 - depth
	if index < 0 {
		return zero
	}

	return v.ctx[index]
}

//
// Block Parameters stack
//

// push new block params
func (v *EvalVisitor) pushBlockParams(params map[string]interface{}) {
	v.blockParams = append(v.blockParams, params)
}

// pop last block params
func (v *EvalVisitor) popBlockParams() map[string]interface{} {
	var result map[string]interface{}

	if len(v.blockParams) == 0 {
		return result
	}

	result, v.blockParams = v.blockParams[len(v.blockParams)-1], v.blockParams[:len(v.blockParams)-1]
	return result
}

// // returns current block params
// func (v *EvalVisitor) curBlockParams() map[string]interface{} {
// 	return v.ancestorBlockParams(0)
// }

// find block parameter value
func (v *EvalVisitor) findBlockParam(name string) interface{} {
	for i := len(v.blockParams) - 1; i >= 0; i-- {
		for k, v := range v.blockParams[i] {
			if name == k {
				return v
			}
		}
	}

	return nil
}

// // get ancestor block params
// func (v *EvalVisitor) ancestorBlockParams(depth int) map[string]interface{} {
// 	index := len(v.blockParams) - 1 - depth
// 	if index < 0 {
// 		return map[string]interface{}{}
// 	}

// 	return v.blockParams[index]
// }

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

// Evaluate program with given context and returns string result
func (v *EvalVisitor) evalProgramWith(program *ast.Program, ctx reflect.Value, index int) string {
	blockParams := make(map[string]interface{})

	// compute block params
	if len(program.BlockParams) > 0 {
		blockParams[program.BlockParams[0]] = ctx.Interface()
	}

	if (len(program.BlockParams) > 1) && (index >= 0) {
		blockParams[program.BlockParams[1]] = index
	}

	// evaluate program
	if len(blockParams) > 0 {
		v.pushBlockParams(blockParams)
	}

	v.pushCtx(ctx)

	result, _ := program.Accept(v).(string)

	v.popCtx()

	if len(blockParams) > 0 {
		v.popBlockParams()
	}

	return result
}

// evaluates all path parts
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

// evaluates field in given context
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
	}

	// check if result is a function
	result, _ = indirect(result)
	if result.Kind() == reflect.Func {
		// in that code path, we know we can't be an expression root
		result = v.evalFunc(result, exprRoot)
	}

	if VERBOSE_EVAL {
		log.Printf("evalField(): '%s' with context %s => %s | Kind: %s", fieldName, StrValue(ctx), StrValue(result), result.Kind())
	}

	return result
}

// evaluates a function
func (v *EvalVisitor) evalFunc(funcVal reflect.Value, exprRoot bool) reflect.Value {
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
			arg = NewEmptyHelperArg(v)
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

	return NewHelperArg(v, params, hash)
}

//
// Partials
//

// Finds given partial
func (v *EvalVisitor) findPartial(name string) *Partial {
	// check template partials
	if v.tpl.partials[name] != nil {
		return v.tpl.partials[name]
	}

	// check global partials
	return FindPartial(name)
}

// Computes partial context
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

// Evaluates a partial
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
	result = v.indentLines(result, node.Indent)

	if ctx.IsValid() {
		v.popCtx()
	}

	return result
}

func (v *EvalVisitor) indentLines(str string, indent string) string {
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

// Returns true if given expression was a function call
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
				if VERBOSE_EVAL {
					log.Printf("VisitBlock(): Truthy, visiting Program")
				}

				switch val.Kind() {
				case reflect.Array, reflect.Slice:
					// Array context
					for i := 0; i < val.Len(); i++ {
						// Evaluate program
						result += v.evalProgramWith(node.Program, val.Index(i), i)
					}
				default:
					// NOT array
					result = v.evalProgramWith(node.Program, val, -1)
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

// Evaluate a path expression
func (v *EvalVisitor) evalPathExpression(node *ast.PathExpression, exprRoot bool) interface{} {
	if value := v.findBlockParam(node.Str()); value != nil {
		// block parameter value
		return value
	} else if node.Data {
		// private data
		return v.evalDataPathExpression(node)
	}

	return v.evalCtxPathExpression(node, exprRoot)
}

// Evaluate a private data path expression
func (v *EvalVisitor) evalDataPathExpression(node *ast.PathExpression) interface{} {
	// find data frame
	frame := v.dataFrame
	for i := node.Depth; i > 0; i-- {
		if frame.parent == nil {
			return nil
		}
		frame = frame.parent
	}

	return frame.Find(node.Parts)
}

// Evaluate a context path expression
func (v *EvalVisitor) evalCtxPathExpression(node *ast.PathExpression, exprRoot bool) interface{} {
	v.at(node)

	var result interface{}

	depth := node.Depth
	ctx := v.ancestorCtx(depth)
	stopDeep := false

	for (result == nil) && ctx.IsValid() && (depth <= len(v.ctx) && !stopDeep) {
		if VERBOSE_EVAL {
			log.Printf("VisitPath(): '%s' with context '%s' (depth: %d)", node.Original, StrValue(ctx), depth)
		}

		switch ctx.Kind() {
		case reflect.Array, reflect.Slice:
			// Array context
			var results []interface{}

			for i := 0; i < ctx.Len(); i++ {
				value, _ := v.evalPath(ctx.Index(i), node.Parts, exprRoot)
				if value.IsValid() {
					results = append(results, value.Interface())
				}
			}

			result = results
		default:
			// NOT array context
			value, partResolved := v.evalPath(ctx, node.Parts, exprRoot)
			if value.IsValid() {
				result = value.Interface()
			}

			if partResolved {
				// As soon as we find the first part of a path, we must not try to resolve with parent context if result is finally `nil`
				// Reference: "Dotted Names - Context Precedence" mustache test
				stopDeep = true
			}
		}

		if VERBOSE_EVAL {
			log.Printf("VisitPath(): result => '%s'", Str(result))
		}

		if result == nil {
			// check ancestor
			depth++
			ctx = v.ancestorCtx(depth)
		}
	}

	return result
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

package raymond

import (
	"fmt"
	"reflect"
)

// Arguments provided to helpers
type HelperParams struct {
	// evaluation visitor
	eval *EvalVisitor

	// params
	params []interface{}
	hash   map[string]interface{}
}

// Helper function
type Helper func(p *HelperParams) string

// All registered helpers
var helpers map[string]Helper

func init() {
	helpers = make(map[string]Helper)

	// register builtin helpers
	RegisterHelper("if", ifHelper)
	RegisterHelper("unless", unlessHelper)
	RegisterHelper("with", withHelper)
	RegisterHelper("each", eachHelper)
}

// Registers a new helper function
func RegisterHelper(name string, helper Helper) {
	if helpers[name] != nil {
		panic(fmt.Errorf("Helper already registered: %s", name))
	}

	helpers[name] = helper
}

// Find a registered helper function
func FindHelper(name string) Helper {
	return helpers[name]
}

func NewHelperParams(eval *EvalVisitor, params []interface{}, hash map[string]interface{}) *HelperParams {
	return &HelperParams{
		eval:   eval,
		params: params,
		hash:   hash,
	}
}

func NewEmptyHelperParams(eval *EvalVisitor) *HelperParams {
	return &HelperParams{
		eval: eval,
		hash: make(map[string]interface{}),
	}
}

// Returns all parameters
func (p *HelperParams) Params() []interface{} {
	return p.params
}

// Returns parameter at given position
func (p *HelperParams) Param(pos int) interface{} {
	if len(p.params) > pos {
		return p.params[pos]
	} else {
		return nil
	}
}

// Get string representation of parameter at given position
func (p *HelperParams) ParamStr(pos int) string {
	return Str(p.Param(pos))
}

// Returns hash value by name
func (p *HelperParams) Option(name string) interface{} {
	return p.hash[name]
}

// Returns string representation of hash value by name
func (p *HelperParams) OptionStr(name string) string {
	return Str(p.hash[name])
}

// Returns input data by name
func (p *HelperParams) Data(name string) interface{} {
	value := p.eval.evalField(p.eval.curCtx(), name)
	if !value.IsValid() {
		return nil
	}

	return value.Interface()
}

// Get string representation of input data by name
func (p *HelperParams) DataStr(name string) string {
	return Str(p.Data(name))
}

// Returns true if first param is truthy
func (p *HelperParams) TruthFirstParam() bool {
	val := p.Param(0)
	if val == nil {
		return false
	}

	thruth, ok := IsTruth(reflect.ValueOf(val))
	if !ok {
		return false
	}

	return thruth
}

// Returns true if 'includeZero' option is set and first param is the number 0
func (p *HelperParams) IsIncludableZero() bool {
	b, ok := p.Option("includeZero").(bool)
	if ok && b {
		nb, ok := p.Param(0).(int)
		if ok && nb == 0 {
			return true
		}
	}

	return false
}

// Evaluate block
func (p *HelperParams) Block() string {
	result := ""
	if block := p.eval.curBlock(); (block != nil) && (block.Program != nil) {
		result, _ = block.Program.Accept(p.eval).(string)
	}

	return result
}

// Evaluate inverse
func (p *HelperParams) Inverse() string {
	result := ""
	if block := p.eval.curBlock(); (block != nil) && (block.Inverse != nil) {
		result, _ = block.Inverse.Accept(p.eval).(string)
	}

	return result
}

// Evaluate block with given context
func (p *HelperParams) BlockWith(ctx interface{}) string {
	p.PushCtx(ctx)
	result := p.Block()
	p.PopCtx()

	return result
}

// Push context
func (p *HelperParams) PushCtx(ctx interface{}) {
	p.eval.pushCtx(reflect.ValueOf(ctx))
}

// Pop context
func (p *HelperParams) PopCtx() interface{} {
	var value reflect.Value

	value = p.eval.popCtx()
	if !value.IsValid() {
		return value
	}

	return value.Interface()
}

//
// Builtin helpers
//

func ifHelper(p *HelperParams) string {
	if p.IsIncludableZero() || p.TruthFirstParam() {
		return p.Block()
	}

	return p.Inverse()
}

func unlessHelper(p *HelperParams) string {
	if p.IsIncludableZero() || p.TruthFirstParam() {
		return p.Inverse()
	}

	return p.Block()
}

func withHelper(p *HelperParams) string {
	if p.TruthFirstParam() {
		return p.BlockWith(p.Param(0))
	}

	return p.Inverse()
}

func eachHelper(p *HelperParams) string {
	if !p.TruthFirstParam() {
		p.Inverse()
		return ""
	}

	result := ""

	val := reflect.ValueOf(p.Param(0))
	switch val.Kind() {
	case reflect.Array, reflect.Slice:
		for i := 0; i < val.Len(); i++ {
			result += p.BlockWith(val.Index(i).Interface())
		}
	case reflect.Map:
		// note: a go hash is not ordered, so result may vary, this behaviour differs from the JS implementation
		keys := val.MapKeys()
		for i := 0; i < len(keys); i++ {
			result += p.BlockWith(val.MapIndex(keys[i]).Interface())
		}
	case reflect.Struct:
		// @todo !!!
	}

	return result
}

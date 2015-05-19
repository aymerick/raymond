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

// Get parameter at given position
func (p *HelperParams) At(pos int) interface{} {
	if len(p.params) > pos {
		return p.params[pos]
	} else {
		return nil
	}
}

// Get hash option by name
func (p *HelperParams) Option(name string) interface{} {
	return p.hash[name]
}

// Returns true if first param is truthy
func (p *HelperParams) TruthFirstParam() bool {
	val := p.At(0)
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
		nb, ok := p.At(0).(int)
		if ok && nb == 0 {
			return true
		}
	}

	return false
}

// Evaluate block
func (p *HelperParams) EvaluateBlock() {
	if block := p.eval.curBlock(); (block != nil) && (block.Program != nil) {
		block.Program.Accept(p.eval)
	}
}

// Evaluate inverse
func (p *HelperParams) EvaluateInverse() {
	if block := p.eval.curBlock(); (block != nil) && (block.Inverse != nil) {
		block.Inverse.Accept(p.eval)
	}
}

// Push context
func (p *HelperParams) PushCtx(ctx interface{}) {
	p.eval.PushCtx(reflect.ValueOf(ctx))
}

// Pop context
func (p *HelperParams) PopCtx() interface{} {
	var value reflect.Value

	value = p.eval.PopCtx()
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
		p.EvaluateBlock()
	} else {
		p.EvaluateInverse()
	}

	// irrelevant
	return ""
}

func unlessHelper(p *HelperParams) string {
	if p.IsIncludableZero() || p.TruthFirstParam() {
		p.EvaluateInverse()
	} else {
		p.EvaluateBlock()
	}

	// irrelevant
	return ""
}

func withHelper(p *HelperParams) string {
	if p.TruthFirstParam() {
		p.PushCtx(p.At(0))

		p.EvaluateBlock()

		p.PopCtx()
	} else {
		p.EvaluateInverse()
	}

	// irrelevant
	return ""
}

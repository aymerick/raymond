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

	// register default helpers
	RegisterHelper("if", ifHelper)
	RegisterHelper("unless", unlessHelper)
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

// Evaluate block
func (p *HelperParams) EvaluateBlock() {
	if p.eval.curBlock != nil && p.eval.curBlock.Program != nil {
		p.eval.curBlock.Program.Accept(p.eval)
	}
}

// Evaluate inverse
func (p *HelperParams) EvaluateInverse() {
	if p.eval.curBlock != nil && p.eval.curBlock.Inverse != nil {
		p.eval.curBlock.Inverse.Accept(p.eval)
	}
}

//
// Default helpers
//

func ifHelper(p *HelperParams) string {
	if p.TruthFirstParam() {
		p.EvaluateBlock()
	} else {
		p.EvaluateInverse()
	}

	// irrelevant
	return ""
}

func unlessHelper(p *HelperParams) string {
	if p.TruthFirstParam() {
		p.EvaluateInverse()
	} else {
		p.EvaluateBlock()
	}

	// irrelevant
	return ""
}

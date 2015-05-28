package raymond

import (
	"fmt"
	"reflect"
)

// Arguments provided to helpers
type HelperArg struct {
	// evaluation visitor
	eval *EvalVisitor

	// params
	params []interface{}
	hash   map[string]interface{}
}

// Helper function
type Helper func(h *HelperArg) string

// All global helpers
var helpers map[string]Helper

func init() {
	helpers = make(map[string]Helper)

	// register builtin helpers
	RegisterHelper("if", ifHelper)
	RegisterHelper("unless", unlessHelper)
	RegisterHelper("with", withHelper)
	RegisterHelper("each", eachHelper)
}

// Registers a new global helper function
func RegisterHelper(name string, helper Helper) {
	if helpers[name] != nil {
		panic(fmt.Errorf("Helper already registered: %s", name))
	}

	helpers[name] = helper
}

// Find a registered global helper function
func FindHelper(name string) Helper {
	return helpers[name]
}

// Instanciates a new HelperArg
func NewHelperArg(eval *EvalVisitor, params []interface{}, hash map[string]interface{}) *HelperArg {
	return &HelperArg{
		eval:   eval,
		params: params,
		hash:   hash,
	}
}

// Instanciates a new empty HelperArg
func NewEmptyHelperArg(eval *EvalVisitor) *HelperArg {
	return &HelperArg{
		eval: eval,
		hash: make(map[string]interface{}),
	}
}

// Returns all parameters
func (h *HelperArg) Params() []interface{} {
	return h.params
}

// Returns parameter at given position
func (h *HelperArg) Param(pos int) interface{} {
	if len(h.params) > pos {
		return h.params[pos]
	} else {
		return nil
	}
}

// Get string representation of parameter at given position
func (h *HelperArg) ParamStr(pos int) string {
	return Str(h.Param(pos))
}

// Returns hash value by name
func (h *HelperArg) Option(name string) interface{} {
	return h.hash[name]
}

// Returns string representation of hash value by name
func (h *HelperArg) OptionStr(name string) string {
	return Str(h.hash[name])
}

// Returns input data
func (h *HelperArg) Data() interface{} {
	return h.eval.curCtx()
}

// Returns input data by name
func (h *HelperArg) DataField(name string) interface{} {
	value := h.eval.evalField(h.eval.curCtx(), name, false)
	if !value.IsValid() {
		return nil
	}

	return value.Interface()
}

// Get string representation of input data by name
func (h *HelperArg) DataStr(name string) string {
	return Str(h.DataField(name))
}

// Returns true if first param is truthy
func (h *HelperArg) TruthFirstParam() bool {
	val := h.Param(0)
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
func (h *HelperArg) IsIncludableZero() bool {
	b, ok := h.Option("includeZero").(bool)
	if ok && b {
		nb, ok := h.Param(0).(int)
		if ok && nb == 0 {
			return true
		}
	}

	return false
}

// Evaluate block
func (h *HelperArg) Block() string {
	result := ""
	if block := h.eval.curBlock(); (block != nil) && (block.Program != nil) {
		result, _ = block.Program.Accept(h.eval).(string)
	}

	return result
}

// Evaluate inverse
func (h *HelperArg) Inverse() string {
	result := ""
	if block := h.eval.curBlock(); (block != nil) && (block.Inverse != nil) {
		result, _ = block.Inverse.Accept(h.eval).(string)
	}

	return result
}

// Evaluate block with given context
func (h *HelperArg) BlockWithCtx(ctx interface{}) string {
	h.PushCtx(ctx)
	result := h.Block()
	h.PopCtx()

	return result
}

// Evaluate block with given context and private data
func (h *HelperArg) BlockWith(ctx interface{}, data *DataFrame) string {
	h.SetDataFrame(data)
	result := h.BlockWithCtx(ctx)
	h.SetDataFrame(data.parent)

	return result
}

// Push context
func (h *HelperArg) PushCtx(ctx interface{}) {
	h.eval.pushCtx(reflect.ValueOf(ctx))
}

// Pop context
func (h *HelperArg) PopCtx() interface{} {
	var value reflect.Value

	value = h.eval.popCtx()
	if !value.IsValid() {
		return value
	}

	return value.Interface()
}

func (h *HelperArg) NewDataFrame() *DataFrame {
	return h.eval.dataFrame.Copy()
}

// Set current data frame
func (h *HelperArg) SetDataFrame(data *DataFrame) {
	h.eval.dataFrame = data
}

//
// Builtin helpers
//

func ifHelper(h *HelperArg) string {
	if h.IsIncludableZero() || h.TruthFirstParam() {
		return h.Block()
	}

	return h.Inverse()
}

func unlessHelper(h *HelperArg) string {
	if h.IsIncludableZero() || h.TruthFirstParam() {
		return h.Inverse()
	}

	return h.Block()
}

func withHelper(h *HelperArg) string {
	if h.TruthFirstParam() {
		return h.BlockWithCtx(h.Param(0))
	}

	return h.Inverse()
}

func eachHelper(h *HelperArg) string {
	if !h.TruthFirstParam() {
		h.Inverse()
		return ""
	}

	result := ""

	val := reflect.ValueOf(h.Param(0))
	switch val.Kind() {
	case reflect.Array, reflect.Slice:
		for i := 0; i < val.Len(); i++ {
			// computes private data
			data := h.NewDataFrame()
			data.Set("index", i)

			if i == 0 {
				data.Set("first", true)
			}

			if i == val.Len()-1 {
				data.Set("last", true)
			}

			result += h.BlockWith(val.Index(i).Interface(), data)
		}
	case reflect.Map:
		// note: a go hash is not ordered, so result may vary, this behaviour differs from the JS implementation
		keys := val.MapKeys()
		for i := 0; i < len(keys); i++ {
			// computes private data
			data := h.NewDataFrame()
			data.Set("key", keys[i].Interface())

			if i == 0 {
				data.Set("first", true)
			}

			if i == len(keys)-1 {
				data.Set("last", true)
			}

			result += h.BlockWith(val.MapIndex(keys[i]).Interface(), data)
		}
	case reflect.Struct:
		// @todo !!!
	}

	return result
}

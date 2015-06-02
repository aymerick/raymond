package raymond

import (
	"fmt"
	"log"
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

// Helper function: can return a string or a SafeString
type Helper func(h *HelperArg) interface{}

// All global helpers
var helpers map[string]Helper

func init() {
	helpers = make(map[string]Helper)

	// register builtin helpers
	RegisterHelper("if", ifHelper)
	RegisterHelper("unless", unlessHelper)
	RegisterHelper("with", withHelper)
	RegisterHelper("each", eachHelper)
	RegisterHelper("log", logHelper)
	RegisterHelper("lookup", lookupHelper)
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

//
// Getters
//

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
func (h *HelperArg) Hash(name string) interface{} {
	return h.hash[name]
}

// Returns string representation of hash value by name
func (h *HelperArg) HashStr(name string) string {
	return Str(h.hash[name])
}

// Returns current evaluation context
func (h *HelperArg) Ctx() interface{} {
	return h.eval.curCtx()
}

// Returns current context field value
func (h *HelperArg) Field(name string) interface{} {
	value := h.eval.evalField(h.eval.curCtx(), name, false)
	if !value.IsValid() {
		return nil
	}

	return value.Interface()
}

// Get string representation of current context field value
func (h *HelperArg) FieldStr(name string) string {
	return Str(h.Field(name))
}

//
// Private data frame
//

// Get current private data frame
func (h *HelperArg) DataFrame() *DataFrame {
	return h.eval.dataFrame
}

// Instanciates a new data frame that is a copy of current one
func (h *HelperArg) NewDataFrame() *DataFrame {
	return h.eval.dataFrame.Copy()
}

// Instanciates a new data frame and set iteration specific vars
func (h *HelperArg) NewIterDataFrame(length int, i int, key interface{}) *DataFrame {
	return h.eval.dataFrame.NewIterDataFrame(length, i, key)
}

// Set current data frame
func (h *HelperArg) SetDataFrame(data *DataFrame) {
	h.eval.setDataFrame(data)
}

// Set back parent data frame
func (h *HelperArg) PopDataFrame() {
	h.eval.popDataFrame()
}

//
// Evaluation
//

// Evaluate block with given context, private data and iteration key
func (h *HelperArg) BlockWith(ctx interface{}, data *DataFrame, key interface{}) string {
	result := ""

	if block := h.eval.curBlock(); (block != nil) && (block.Program != nil) {
		result = h.eval.evalProgram(block.Program, ctx, data, key)
	}

	return result
}

// Evaluate block with given context
func (h *HelperArg) BlockWithCtx(ctx interface{}) string {
	return h.BlockWith(ctx, nil, nil)
}

// Evaluate block
func (h *HelperArg) Block() string {
	return h.BlockWith(nil, nil, nil)
}

// Evaluate inverse
func (h *HelperArg) Inverse() string {
	result := ""
	if block := h.eval.curBlock(); (block != nil) && (block.Inverse != nil) {
		result, _ = block.Inverse.Accept(h.eval).(string)
	}

	return result
}

// Evaluate field for given context
func (h *HelperArg) Eval(ctx interface{}, field string) interface{} {
	if ctx == nil {
		return nil
	}

	if field == "" {
		return nil
	}

	val := h.eval.evalField(reflect.ValueOf(ctx), field, false)
	if !val.IsValid() {
		return nil
	}

	return val.Interface()
}

//
// Misc
//

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
	b, ok := h.Hash("includeZero").(bool)
	if ok && b {
		nb, ok := h.Param(0).(int)
		if ok && nb == 0 {
			return true
		}
	}

	return false
}

//
// Builtin helpers
//

func ifHelper(h *HelperArg) interface{} {
	if h.IsIncludableZero() || h.TruthFirstParam() {
		return h.Block()
	}

	return h.Inverse()
}

func unlessHelper(h *HelperArg) interface{} {
	if h.IsIncludableZero() || h.TruthFirstParam() {
		return h.Inverse()
	}

	return h.Block()
}

func withHelper(h *HelperArg) interface{} {
	if h.TruthFirstParam() {
		return h.BlockWithCtx(h.Param(0))
	}

	return h.Inverse()
}

func eachHelper(h *HelperArg) interface{} {
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
			data := h.NewIterDataFrame(val.Len(), i, nil)

			// evaluates block
			result += h.BlockWith(val.Index(i).Interface(), data, i)
		}
	case reflect.Map:
		// note: a go hash is not ordered, so result may vary, this behaviour differs from the JS implementation
		keys := val.MapKeys()
		for i := 0; i < len(keys); i++ {
			key := keys[i].Interface()
			ctx := val.MapIndex(keys[i]).Interface()

			// computes private data
			data := h.NewIterDataFrame(len(keys), i, key)

			// evaluates block
			result += h.BlockWith(ctx, data, key)
		}
	case reflect.Struct:
		var exportedFields []int

		// collect exported fields only
		for i := 0; i < val.NumField(); i++ {
			if tField := val.Type().Field(i); tField.PkgPath == "" {
				exportedFields = append(exportedFields, i)
			}
		}

		for i, fieldIndex := range exportedFields {
			key := val.Type().Field(fieldIndex).Name
			ctx := val.Field(fieldIndex).Interface()

			// computes private data
			data := h.NewIterDataFrame(len(exportedFields), i, key)

			// evaluates block
			result += h.BlockWith(ctx, data, key)
		}
	}

	return result
}

func logHelper(h *HelperArg) interface{} {
	log.Print(h.ParamStr(0))
	return ""
}

func lookupHelper(h *HelperArg) interface{} {
	return Str(h.Eval(h.Param(0), h.ParamStr(1)))
}

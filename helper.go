package raymond

import (
	"fmt"
	"log"
	"reflect"
)

// HelperArg represents the argument provided to helpers and context functions.
type HelperArg struct {
	// evaluation visitor
	eval *evalVisitor

	// params
	params []interface{}
	hash   map[string]interface{}
}

// Helper represents a helper function. It returns a string or a SafeString.
type Helper func(h *HelperArg) interface{}

// helpers stores all globally registered helpers
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

// RegisterHelper registers a global helper.
func RegisterHelper(name string, helper Helper) {
	if helpers[name] != nil {
		panic(fmt.Errorf("Helper already registered: %s", name))
	}

	helpers[name] = helper
}

// RegisterHelpers registers several global helpers.
func RegisterHelpers(helpers map[string]Helper) {
	for name, helper := range helpers {
		RegisterHelper(name, helper)
	}
}

// FindHelper finds a globally registered helper.
func FindHelper(name string) Helper {
	return helpers[name]
}

// newHelperArg instanciates a new HelperArg
func newHelperArg(eval *evalVisitor, params []interface{}, hash map[string]interface{}) *HelperArg {
	return &HelperArg{
		eval:   eval,
		params: params,
		hash:   hash,
	}
}

// newEmptyHelperArg instanciates a new empty HelperArg
func newEmptyHelperArg(eval *evalVisitor) *HelperArg {
	return &HelperArg{
		eval: eval,
		hash: make(map[string]interface{}),
	}
}

//
// Getters
//

// Params returns all parameters.
func (h *HelperArg) Params() []interface{} {
	return h.params
}

// Paramreturns parameter at given position.
func (h *HelperArg) Param(pos int) interface{} {
	if len(h.params) > pos {
		return h.params[pos]
	} else {
		return nil
	}
}

// ParamStr returns string representation of parameter at given position.
func (h *HelperArg) ParamStr(pos int) string {
	return Str(h.Param(pos))
}

// Hash returns entire hash.
func (h *HelperArg) Hash() map[string]interface{} {
	return h.hash
}

// HashProp returns hash property.
func (h *HelperArg) HashProp(name string) interface{} {
	return h.hash[name]
}

// HashStr returns string representation of hash property.
func (h *HelperArg) HashStr(name string) string {
	return Str(h.hash[name])
}

// Ctx returns current evaluation context.
func (h *HelperArg) Ctx() interface{} {
	return h.eval.curCtx()
}

// Field returns current context field value.
func (h *HelperArg) Field(name string) interface{} {
	value := h.eval.evalField(h.eval.curCtx(), name, false)
	if !value.IsValid() {
		return nil
	}

	return value.Interface()
}

// FieldStr returns string representation of current context field value.
func (h *HelperArg) FieldStr(name string) string {
	return Str(h.Field(name))
}

// Data returns private data value by name.
func (h *HelperArg) Data(name string) interface{} {
	return h.eval.dataFrame.Get(name)
}

// DataStr returns string representation of private data value by name.
func (h *HelperArg) DataStr(name string) string {
	return Str(h.eval.dataFrame.Get(name))
}

//
// Private data frame
//

// DataFrame returns current private data frame.
func (h *HelperArg) DataFrame() *DataFrame {
	return h.eval.dataFrame
}

// NewDataFrame instanciates a new data frame that is a copy of current one.
func (h *HelperArg) NewDataFrame() *DataFrame {
	return h.eval.dataFrame.Copy()
}

// newIterDataFrame instanciates a new data frame and set iteration specific vars.
func (h *HelperArg) NewIterDataFrame(length int, i int, key interface{}) *DataFrame {
	return h.eval.dataFrame.newIterDataFrame(length, i, key)
}

//
// Evaluation
//

// BlockWith evaluates block with given context, private data and iteration key.
func (h *HelperArg) BlockWith(ctx interface{}, data *DataFrame, key interface{}) string {
	result := ""

	if block := h.eval.curBlock(); (block != nil) && (block.Program != nil) {
		result = h.eval.evalProgram(block.Program, ctx, data, key)
	}

	return result
}

// BlockWithCtx evaluates block with given context.
func (h *HelperArg) BlockWithCtx(ctx interface{}) string {
	return h.BlockWith(ctx, nil, nil)
}

// BlockWithData evaluates block with given private data.
func (h *HelperArg) BlockWithData(data *DataFrame) string {
	return h.BlockWith(nil, data, nil)
}

// Block evaluates block.
func (h *HelperArg) Block() string {
	return h.BlockWith(nil, nil, nil)
}

// Inverse evaluates block inverse.
func (h *HelperArg) Inverse() string {
	result := ""
	if block := h.eval.curBlock(); (block != nil) && (block.Inverse != nil) {
		result, _ = block.Inverse.Accept(h.eval).(string)
	}

	return result
}

// Eval evaluates field for given context.
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

// truthFirstParam returns true if first param is truthy
func (h *HelperArg) truthFirstParam() bool {
	val := h.Param(0)
	if val == nil {
		return false
	}

	thruth, ok := isTruth(reflect.ValueOf(val))
	if !ok {
		return false
	}

	return thruth
}

// isIncludableZero returns true if 'includeZero' option is set and first param is the number 0
func (h *HelperArg) isIncludableZero() bool {
	b, ok := h.HashProp("includeZero").(bool)
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

// #if block helper
func ifHelper(h *HelperArg) interface{} {
	if h.isIncludableZero() || h.truthFirstParam() {
		return h.Block()
	}

	return h.Inverse()
}

// #unless block helper
func unlessHelper(h *HelperArg) interface{} {
	if h.isIncludableZero() || h.truthFirstParam() {
		return h.Inverse()
	}

	return h.Block()
}

// #with block helper
func withHelper(h *HelperArg) interface{} {
	if h.truthFirstParam() {
		return h.BlockWithCtx(h.Param(0))
	}

	return h.Inverse()
}

// #each block helper
func eachHelper(h *HelperArg) interface{} {
	if !h.truthFirstParam() {
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

// #log helper
func logHelper(h *HelperArg) interface{} {
	log.Print(h.ParamStr(0))
	return ""
}

// #lookup helper
func lookupHelper(h *HelperArg) interface{} {
	return Str(h.Eval(h.Param(0), h.ParamStr(1)))
}

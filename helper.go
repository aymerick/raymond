package raymond

import "fmt"

// Arguments provided to helpers
type HelperParams struct {
	params []interface{}
	hash   map[string]interface{}
}

// Helper function
type Helper func(p *HelperParams) string

// All registered helpers
var helpers map[string]Helper

func init() {
	helpers = make(map[string]Helper)
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

func NewHelperParams(params []interface{}, hash map[string]interface{}) *HelperParams {
	return &HelperParams{
		params: params,
		hash:   hash,
	}
}

// Get parameter at given position
func (p *HelperParams) at(pos int) interface{} {
	if len(p.params) > pos {
		return p.params[pos]
	} else {
		return nil
	}
}

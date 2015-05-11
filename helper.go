package raymond

import "fmt"

// Context argument provided to helpers
type HelperContext struct {
}

// Helper function
type HelperFunc func(ctx *HelperContext) string

// All registered helpers
var helpers map[string]HelperFunc

func init() {
	helpers = make(map[string]HelperFunc)
}

// Registers a new helper function
func RegisterHelper(name string, helper HelperFunc) {
	if helpers[name] != nil {
		panic(fmt.Errorf("Helper already registered: %s", name))
	}

	helpers[name] = helper
}

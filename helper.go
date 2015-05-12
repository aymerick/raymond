package raymond

import "fmt"

// Context argument provided to helpers
type HelperContext struct {
}

// Helper function
type Helper func(ctx *HelperContext) string

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

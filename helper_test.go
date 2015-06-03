package raymond

import "testing"

const (
	VERBOSE = false
)

//
// Helpers
//

func barHelper(h *HelperArg) interface{} { return "bar" }

func echoHelper(h *HelperArg) interface{} {
	str, _ := h.Param(0).(string)
	nb, ok := h.Param(1).(int)
	if !ok {
		nb = 1
	}

	result := ""
	for i := 0; i < nb; i++ {
		result += str
	}

	return result
}

func boolHelper(h *HelperArg) interface{} {
	b, _ := h.Param(0).(bool)
	if b {
		return "yes it is"
	}

	return "absolutely not"
}

func gnakHelper(h *HelperArg) interface{} {
	nb, ok := h.Param(0).(int)
	if !ok {
		nb = 1
	}

	result := ""
	for i := 0; i < nb; i++ {
		result += "GnAK!"
	}

	return result
}

//
// Tests
//

var helperTests = []Test{
	{
		"simple helper",
		`{{foo}}`,
		nil, nil,
		map[string]Helper{"foo": barHelper},
		nil,
		`bar`,
	},
	{
		"helper with literal string param",
		`{{echo "foo"}}`,
		nil, nil,
		map[string]Helper{"echo": echoHelper},
		nil,
		`foo`,
	},
	{
		"helper with identifier param",
		`{{echo foo}}`,
		map[string]interface{}{"foo": "bar"},
		nil,
		map[string]Helper{"echo": echoHelper},
		nil,
		`bar`,
	},
	{
		"helper with literal boolean param",
		`{{bool true}}`,
		nil, nil,
		map[string]Helper{"bool": boolHelper},
		nil,
		`yes it is`,
	},
	{
		"helper with literal boolean param",
		`{{bool false}}`,
		nil, nil,
		map[string]Helper{"bool": boolHelper},
		nil,
		`absolutely not`,
	},
	{
		"helper with literal boolean param",
		`{{gnak 5}}`,
		nil, nil,
		map[string]Helper{"gnak": gnakHelper},
		nil,
		`GnAK!GnAK!GnAK!GnAK!GnAK!`,
	},
	{
		"helper with several parameters",
		`{{echo "GnAK!" 3}}`,
		nil, nil,
		map[string]Helper{"echo": echoHelper},
		nil,
		`GnAK!GnAK!GnAK!`,
	},
	{
		"#if helper with true literal",
		`{{#if true}}YES MAN{{/if}}`,
		nil, nil, nil, nil,
		`YES MAN`,
	},
	{
		"#if helper with false literal",
		`{{#if false}}YES MAN{{/if}}`,
		nil, nil, nil, nil,
		``,
	},
	{
		"#if helper with truthy identifier",
		`{{#if ok}}YES MAN{{/if}}`,
		map[string]interface{}{"ok": true},
		nil, nil, nil,
		`YES MAN`,
	},
	{
		"#if helper with falsy identifier",
		`{{#if ok}}YES MAN{{/if}}`,
		map[string]interface{}{"ok": false},
		nil, nil, nil,
		``,
	},
	{
		"#unless helper with true literal",
		`{{#unless true}}YES MAN{{/unless}}`,
		nil, nil, nil, nil,
		``,
	},
	{
		"#unless helper with false literal",
		`{{#unless false}}YES MAN{{/unless}}`,
		nil, nil, nil, nil,
		`YES MAN`,
	},
	{
		"#unless helper with truthy identifier",
		`{{#unless ok}}YES MAN{{/unless}}`,
		map[string]interface{}{"ok": true},
		nil, nil, nil,
		``,
	},
	{
		"#unless helper with falsy identifier",
		`{{#unless ok}}YES MAN{{/unless}}`,
		map[string]interface{}{"ok": false},
		nil, nil, nil,
		`YES MAN`,
	},
}

//
// Let's go
//

func TestHelper(t *testing.T) {
	launchTests(t, helperTests)
}

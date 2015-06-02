package raymond

import (
	"fmt"
	"reflect"
	"testing"
)

const (
	VERBOSE = false
)

//
// Helpers
//

func barHelper(h *HelperArg) interface{} { return "bar" }

func barSuffixHelper(h *HelperArg) interface{} {
	str, _ := h.Param(0).(string)
	return "bar " + str
}

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

func linkHelper(h *HelperArg) interface{} {
	prefix, _ := h.Param(0).(string)

	return fmt.Sprintf(`<a href="%s/%s">%s</a>`, prefix, h.FieldStr("url"), h.FieldStr("text"))
}

func rawHelper(h *HelperArg) interface{} {
	result := h.Block()

	for _, param := range h.Params() {
		result += Str(param)
	}

	return result
}

func formHelper(h *HelperArg) interface{} {
	return "<form>" + h.Block() + "</form>"
}

func formCtxHelper(h *HelperArg) interface{} {
	return "<form>" + h.BlockWithCtx(h.Param(0)) + "</form>"
}

func listHelper(h *HelperArg) interface{} {
	ctx := h.Param(0)

	val := reflect.ValueOf(ctx)
	switch val.Kind() {
	case reflect.Array, reflect.Slice:
		if val.Len() > 0 {
			result := "<ul>"
			for i := 0; i < val.Len(); i++ {
				result += "<li>"
				result += h.BlockWithCtx(val.Index(i).Interface())
				result += "</li>"
			}
			result += "</ul>"

			return result
		}
	}

	return "<p>" + h.Inverse() + "</p>"
}

func blogHelper(h *HelperArg) interface{} {
	return "val is " + h.ParamStr(0)
}

func equalHelper(h *HelperArg) interface{} {
	return Str(h.ParamStr(0) == h.ParamStr(1))
}

func dashHelper(h *HelperArg) interface{} {
	return h.ParamStr(0) + "-" + h.ParamStr(1)
}

func concatHelper(h *HelperArg) interface{} {
	return h.ParamStr(0) + h.ParamStr(1)
}

func detectDataHelper(h *HelperArg) interface{} {
	if val, ok := h.DataFrame().Get("exclaim").(string); ok {
		return val
	}

	return ""
}

//
// Tests
//

var helperTests = []raymondTest{
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
	launchRaymondTests(t, helperTests)
}

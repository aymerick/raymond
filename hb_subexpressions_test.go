package raymond

import "testing"

//
// Those tests come from:
//   https://github.com/wycats/handlebars.js/blob/master/spec/subexpression.js
//
var hbSubexpressionsTests = []raymondTest{
	{
		"arg-less helper",
		"{{foo (bar)}}!",
		map[string]interface{}{},
		map[string]Helper{
			"foo": func(h *HelperArg) string {
				return h.ParamStr(0) + h.ParamStr(0)
			},
			"bar": func(h *HelperArg) string {
				return "LOL"
			},
		},
		nil,
		"LOLLOL!",
	},
	{
		"helper w args",
		"{{blog (equal a b)}}",
		map[string]interface{}{"bar": "LOL"},
		map[string]Helper{
			"blog":  blogHelper,
			"equal": equalHelper,
		},
		nil,
		"val is true",
	},
	{
		"mixed paths and helpers",
		"{{blog baz.bat (equal a b) baz.bar}}",
		map[string]interface{}{"bar": "LOL", "baz": map[string]string{"bat": "foo!", "bar": "bar!"}},
		map[string]Helper{
			"blog": func(h *HelperArg) string {
				return "val is " + h.ParamStr(0) + ", " + h.ParamStr(1) + " and " + h.ParamStr(2)
			},
			"equal": equalHelper,
		},
		nil,
		"val is foo!, true and bar!",
	},
	{
		"supports much nesting",
		"{{blog (equal (equal true true) true)}}",
		map[string]interface{}{"bar": "LOL"},
		map[string]Helper{
			"blog":  blogHelper,
			"equal": equalHelper,
		},
		nil,
		"val is true",
	},

	{
		"GH-800 : Complex subexpressions (1)",
		"{{dash 'abc' (concat a b)}}",
		map[string]interface{}{"a": "a", "b": "b", "c": map[string]string{"c": "c"}, "d": "d", "e": map[string]string{"e": "e"}},
		map[string]Helper{"dash": dashHelper, "concat": concatHelper},
		nil,
		"abc-ab",
	},
	{
		"GH-800 : Complex subexpressions (2)",
		"{{dash d (concat a b)}}",
		map[string]interface{}{"a": "a", "b": "b", "c": map[string]string{"c": "c"}, "d": "d", "e": map[string]string{"e": "e"}},
		map[string]Helper{"dash": dashHelper, "concat": concatHelper},
		nil,
		"d-ab",
	},
	{
		"GH-800 : Complex subexpressions (3)",
		"{{dash c.c (concat a b)}}",
		map[string]interface{}{"a": "a", "b": "b", "c": map[string]string{"c": "c"}, "d": "d", "e": map[string]string{"e": "e"}},
		map[string]Helper{"dash": dashHelper, "concat": concatHelper},
		nil,
		"c-ab",
	},
	{
		"GH-800 : Complex subexpressions (4)",
		"{{dash (concat a b) c.c}}",
		map[string]interface{}{"a": "a", "b": "b", "c": map[string]string{"c": "c"}, "d": "d", "e": map[string]string{"e": "e"}},
		map[string]Helper{"dash": dashHelper, "concat": concatHelper},
		nil,
		"ab-c",
	},
	{
		"GH-800 : Complex subexpressions (5)",
		"{{dash (concat a e.e) c.c}}",
		map[string]interface{}{"a": "a", "b": "b", "c": map[string]string{"c": "c"}, "d": "d", "e": map[string]string{"e": "e"}},
		map[string]Helper{"dash": dashHelper, "concat": concatHelper},
		nil,
		"ae-c",
	},

	{
		// note: test not relevant
		"provides each nested helper invocation its own options hash",
		"{{equal (equal true true) true}}",
		map[string]interface{}{},
		map[string]Helper{
			"equal": equalHelper,
		},
		nil,
		"true",
	},
	{
		"with hashes",
		"{{blog (equal (equal true true) true fun='yes')}}",
		map[string]interface{}{"bar": "LOL"},
		map[string]Helper{
			"blog":  blogHelper,
			"equal": equalHelper,
		},
		nil,
		"val is true",
	},
	{
		"as hashes",
		"{{blog fun=(equal (blog fun=1) 'val is 1')}}",
		map[string]interface{}{},
		map[string]Helper{
			"blog": func(h *HelperArg) string {
				return "val is " + h.OptionStr("fun")
			},
			"equal": equalHelper,
		},
		nil,
		"val is true",
	},
	{
		"multiple subexpressions in a hash",
		// @todo Do not use unescaping mustaches
		`{{{input aria-label=(t "Name") placeholder=(t "Example User")}}}`,
		map[string]interface{}{},
		map[string]Helper{
			"input": func(h *HelperArg) string {
				// @todo Escape values and return a SafeString
				return `<input aria-label="` + h.OptionStr("aria-label") + `" placeholder="` + h.OptionStr("placeholder") + `" />`
			},
			"t": func(h *HelperArg) string {
				// @todo Return a SafeString
				return h.ParamStr(0)
			},
		},
		nil,
		`<input aria-label="Name" placeholder="Example User" />`,
	},
	{
		"multiple subexpressions in a hash with context",
		// @todo Do not use unescaping mustaches
		`{{{input aria-label=(t item.field) placeholder=(t item.placeholder)}}}`,
		map[string]map[string]string{"item": {"field": "Name", "placeholder": "Example User"}},
		map[string]Helper{
			"input": func(h *HelperArg) string {
				// @todo Escape values and return a SafeString
				return `<input aria-label="` + h.OptionStr("aria-label") + `" placeholder="` + h.OptionStr("placeholder") + `" />`
			},
			"t": func(h *HelperArg) string {
				// @todo Return a SafeString
				return h.ParamStr(0)
			},
		},
		nil,
		`<input aria-label="Name" placeholder="Example User" />`,
	},

	// @todo "in string params mode" if relevant
	// @todo "as hashes in string params mode" if relevant

	// @todo "subexpression functions on the context"
	// @todo "subexpressions can't just be property lookups"
}

func TestHandlebarsSubexpressions(t *testing.T) {
	launchHandlebarsTests(t, hbSubexpressionsTests)
}

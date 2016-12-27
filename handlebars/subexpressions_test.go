package handlebars

import (
	"testing"

	"github.com/gobuffalo/ray"
)

//
// Those tests come from:
//   https://github.com/wycats/handlebars.js/blob/master/spec/subexpression.js
//
var subexpressionsTests = []Test{
	{
		"arg-less helper",
		"{{foo (bar)}}!",
		map[string]interface{}{},
		nil,
		map[string]interface{}{
			"foo": func(val string) string {
				return val + val
			},
			"bar": func() string {
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
		nil,
		map[string]interface{}{
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
		nil,
		map[string]interface{}{
			"blog": func(p, p2, p3 string) string {
				return "val is " + p + ", " + p2 + " and " + p3
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
		nil,
		map[string]interface{}{
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
		nil,
		map[string]interface{}{"dash": dashHelper, "concat": concatHelper},
		nil,
		"abc-ab",
	},
	{
		"GH-800 : Complex subexpressions (2)",
		"{{dash d (concat a b)}}",
		map[string]interface{}{"a": "a", "b": "b", "c": map[string]string{"c": "c"}, "d": "d", "e": map[string]string{"e": "e"}},
		nil,
		map[string]interface{}{"dash": dashHelper, "concat": concatHelper},
		nil,
		"d-ab",
	},
	{
		"GH-800 : Complex subexpressions (3)",
		"{{dash c.c (concat a b)}}",
		map[string]interface{}{"a": "a", "b": "b", "c": map[string]string{"c": "c"}, "d": "d", "e": map[string]string{"e": "e"}},
		nil,
		map[string]interface{}{"dash": dashHelper, "concat": concatHelper},
		nil,
		"c-ab",
	},
	{
		"GH-800 : Complex subexpressions (4)",
		"{{dash (concat a b) c.c}}",
		map[string]interface{}{"a": "a", "b": "b", "c": map[string]string{"c": "c"}, "d": "d", "e": map[string]string{"e": "e"}},
		nil,
		map[string]interface{}{"dash": dashHelper, "concat": concatHelper},
		nil,
		"ab-c",
	},
	{
		"GH-800 : Complex subexpressions (5)",
		"{{dash (concat a e.e) c.c}}",
		map[string]interface{}{"a": "a", "b": "b", "c": map[string]string{"c": "c"}, "d": "d", "e": map[string]string{"e": "e"}},
		nil,
		map[string]interface{}{"dash": dashHelper, "concat": concatHelper},
		nil,
		"ae-c",
	},

	{
		// note: test not relevant
		"provides each nested helper invocation its own options hash",
		"{{equal (equal true true) true}}",
		map[string]interface{}{},
		nil,
		map[string]interface{}{
			"equal": equalHelper,
		},
		nil,
		"true",
	},
	{
		"with hashes",
		"{{blog (equal (equal true true) true fun='yes')}}",
		map[string]interface{}{"bar": "LOL"},
		nil,
		map[string]interface{}{
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
		nil,
		map[string]interface{}{
			"blog": func(options *ray.Options) string {
				return "val is " + options.HashStr("fun")
			},
			"equal": equalHelper,
		},
		nil,
		"val is true",
	},
	{
		"multiple subexpressions in a hash",
		`{{input aria-label=(t "Name") placeholder=(t "Example User")}}`,
		map[string]interface{}{},
		nil,
		map[string]interface{}{
			"input": func(options *ray.Options) ray.SafeString {
				return ray.SafeString(`<input aria-label="` + options.HashStr("aria-label") + `" placeholder="` + options.HashStr("placeholder") + `" />`)
			},
			"t": func(param string) ray.SafeString {
				return ray.SafeString(param)
			},
		},
		nil,
		`<input aria-label="Name" placeholder="Example User" />`,
	},
	{
		"multiple subexpressions in a hash with context",
		`{{input aria-label=(t item.field) placeholder=(t item.placeholder)}}`,
		map[string]map[string]string{"item": {"field": "Name", "placeholder": "Example User"}},
		nil,
		map[string]interface{}{
			"input": func(options *ray.Options) ray.SafeString {
				return ray.SafeString(`<input aria-label="` + options.HashStr("aria-label") + `" placeholder="` + options.HashStr("placeholder") + `" />`)
			},
			"t": func(param string) ray.SafeString {
				return ray.SafeString(param)
			},
		},
		nil,
		`<input aria-label="Name" placeholder="Example User" />`,
	},

	// @todo "in string params mode"

	// @todo "as hashes in string params mode"

	{
		"subexpression functions on the context",
		"{{foo (bar)}}!",
		map[string]interface{}{"bar": func() string { return "LOL" }},
		nil,
		map[string]interface{}{
			"foo": func(val string) string {
				return val + val
			},
		},
		nil,
		"LOLLOL!",
	},

	// @todo "subexpressions can't just be property lookups" should raise error
}

func TestSubexpressions(t *testing.T) {
	launchTests(t, subexpressionsTests)
}

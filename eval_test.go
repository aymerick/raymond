package raymond

import "testing"

var evalTests = []Test{
	{
		"only content",
		"this is content",
		nil, nil, nil, nil,
		"this is content",
	},
	{
		"checks path in parent contexts",
		"{{#a}}{{one}}{{#b}}{{one}}{{two}}{{one}}{{/b}}{{/a}}",
		map[string]interface{}{"a": map[string]int{"one": 1}, "b": map[string]int{"two": 2}},
		nil, nil, nil,
		"1121",
	},
	{
		"block params",
		"{{#foo as |bar|}}{{bar}}{{/foo}}{{bar}}",
		map[string]string{"foo": "baz", "bar": "bat"},
		nil, nil, nil,
		"bazbat",
	},
	{
		"block params on array",
		"{{#foo as |bar i|}}{{i}}.{{bar}} {{/foo}}",
		map[string][]string{"foo": {"baz", "bar", "bat"}},
		nil, nil, nil,
		"0.baz 1.bar 2.bat ",
	},
	{
		"nested block params",
		"{{#foos as |foo iFoo|}}{{#wats as |wat iWat|}}{{iFoo}}.{{iWat}}.{{foo}}-{{wat}} {{/wats}}{{/foos}}",
		map[string][]string{"foos": {"baz", "bar"}, "wats": {"the", "phoque"}},
		nil, nil, nil,
		"0.0.baz-the 0.1.baz-phoque 1.0.bar-the 1.1.bar-phoque ",
	},
	{
		"block params with path reference",
		"{{#foo as |bar|}}{{bar.baz}}{{/foo}}",
		map[string]map[string]string{"foo": {"baz": "bat"}},
		nil, nil, nil,
		"bat",
	},

	// @todo Test with a struct for data
	// @todo Test with a "../../path" (depth 2 path) while context is only depth 1
}

func TestEval(t *testing.T) {
	launchTests(t, evalTests)
}

var evalErrors = []Test{
	{
		"functions with wrong number of arguments",
		"{{foo}}",
		map[string]interface{}{"foo": func(a, b *HelperArg) string { return "foo" }},
		nil, nil, nil,
		"Function can only have a uniq argument",
	},
	{
		"functions with wrong argument type",
		"{{foo}}",
		map[string]interface{}{"foo": func(a string) string { return "foo" }},
		nil, nil, nil,
		"Function argument must be a *HelperArg",
	},
	{
		"functions with wrong number of returned values",
		"{{foo}}",
		map[string]interface{}{"foo": func() (string, string) { return "foo", "bar" }},
		nil, nil, nil,
		"Function must return a uniq value",
	},
}

func TestEvalErrors(t *testing.T) {
	launchErrorTests(t, evalErrors)
}

//
// strValue() / Str() tests
//

type strTest struct {
	name   string
	input  interface{}
	output string
}

var strTests = []strTest{
	{"String", "foo", "foo"},
	{"Boolean true", true, "true"},
	{"Boolean false", false, "false"},
	{"Integer", 25, "25"},
	{"Float", 25.75, "25.75"},
	{"Nil", nil, ""},
	{"[]string", []string{"foo", "bar"}, "foobar"},
	{"[]interface{} (strings)", []interface{}{"foo", "bar"}, "foobar"},
	{"[]Boolean", []bool{true, false}, "truefalse"},
}

func TestStr(t *testing.T) {
	for _, test := range strTests {
		if res := Str(test.input); res != test.output {
			t.Errorf("Failed to stringify: %s\nexpected:\n\t'%s'got:\n\t%q", test.name, test.output, res)
		}
	}
}

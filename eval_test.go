package raymond

import "testing"

var evalTests = []raymondTest{
	{
		"only content",
		"this is content",
		nil,
		nil,
		"this is content",
	},
	// @todo Test with a struct for data

	// @todo Test with a "../../path" (depth 2 path) while context is only depth 1
}

func TestEval(t *testing.T) {
	launchRaymondTests(t, evalTests)
}

var evalErrors = []raymondTest{
	{
		"functions with wrong number of arguments",
		"{{foo}}",
		map[string]interface{}{"foo": func(a, b *HelperArg) string { return "foo" }},
		nil,
		"Function can only have a uniq argument",
	},
	{
		"functions with wrong argument type",
		"{{foo}}",
		map[string]interface{}{"foo": func(a string) string { return "foo" }},
		nil,
		"Function argument must be a *HelperArg",
	},
	{
		"functions with wrong number of returned values",
		"{{foo}}",
		map[string]interface{}{"foo": func() (string, string) { return "foo", "bar" }},
		nil,
		"Function must return a uniq string value",
	},
	{
		"functions with wrong returned value type",
		"{{foo}}",
		map[string]interface{}{"foo": func() bool { return true }},
		nil,
		"Function must return a uniq string value",
	},
}

func TestEvalErrors(t *testing.T) {
	launchRaymondErrorTests(t, evalErrors)
}

//
// StrValue() / Str() tests
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
	stats.tests(len(strTests))

	for _, test := range strTests {
		if res := Str(test.input); res != test.output {
			t.Errorf("Failed to stringify: %s\nexpected:\n\t'%s'got:\n\t%q", test.name, test.output, res)
			stats.failed()
		}
	}

	stats.output()
}

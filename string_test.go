package raymond

import (
	"fmt"
	"testing"
)

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

func ExampleStr() {
	output := Str(3) + " foos are " + Str(true) + " and " + Str(-1.25) + " bars are " + Str(false) + "\n"
	output += "But you know '" + Str(nil) + "' John Snow\n"
	output += "map: " + Str(map[string]string{"foo": "bar"}) + "\n"
	output += "array: " + Str([]interface{}{true, 10, "foo", 5, "bar"})

	fmt.Println(output)
	// Output: 3 foos are true and -1.25 bars are false
	// But you know '' John Snow
	// map: map[foo:bar]
	// array: true10foo5bar
}

func ExampleSafeString() {
	RegisterHelper("em", func() SafeString {
		return SafeString("<em>FOO BAR</em>")
	})

	tpl := MustParse("{{em}}")

	result := tpl.MustExec(nil)
	fmt.Print(result)
	// Output: <em>FOO BAR</em>
}

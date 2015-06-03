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

func ExampleStr_bool() {
	output := "foo is " + Str(true) + " but bar is " + Str(false)

	fmt.Println(output)
	// Output: foo is true but bar is false
}

func ExampleStr_numbers() {
	output := "I saw " + Str(3) + " foo with " + Str(-1.25) + " bar"

	fmt.Println(output)
	// Output: I saw 3 foo with -1.25 bar
}

func ExampleStr_nil() {
	output := "You know '" + Str(nil) + "' John Snow"

	fmt.Println(output)
	// Output: You know '' John Snow
}

func ExampleStr_map() {
	output := Str(map[string]string{"foo": "bar"})

	fmt.Println(output)
	// Output: map[foo:bar]
}

func ExampleStr_array() {
	output := Str([]interface{}{true, 10, "foo", 5, "bar"})

	fmt.Println(output)
	// Output: true10foo5bar
}

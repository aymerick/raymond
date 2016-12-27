package ray

import "fmt"

func ExampleIsTrue() {
	output := "Empty array: " + Str(IsTrue([0]string{})) + "\n"
	output += "Non empty array: " + Str(IsTrue([1]string{"foo"})) + "\n"

	output += "Empty slice: " + Str(IsTrue([]string{})) + "\n"
	output += "Non empty slice: " + Str(IsTrue([]string{"foo"})) + "\n"

	output += "Empty map: " + Str(IsTrue(map[string]string{})) + "\n"
	output += "Non empty map: " + Str(IsTrue(map[string]string{"foo": "bar"})) + "\n"

	output += "Empty string: " + Str(IsTrue("")) + "\n"
	output += "Non empty string: " + Str(IsTrue("foo")) + "\n"

	output += "true bool: " + Str(IsTrue(true)) + "\n"
	output += "false bool: " + Str(IsTrue(false)) + "\n"

	output += "0 integer: " + Str(IsTrue(0)) + "\n"
	output += "positive integer: " + Str(IsTrue(10)) + "\n"
	output += "negative integer: " + Str(IsTrue(-10)) + "\n"

	output += "0 float: " + Str(IsTrue(0.0)) + "\n"
	output += "positive float: " + Str(IsTrue(10.0)) + "\n"
	output += "negative integer: " + Str(IsTrue(-10.0)) + "\n"

	output += "struct: " + Str(IsTrue(struct{}{})) + "\n"
	output += "nil: " + Str(IsTrue(nil)) + "\n"

	fmt.Println(output)
	// Output: Empty array: false
	// Non empty array: true
	// Empty slice: false
	// Non empty slice: true
	// Empty map: false
	// Non empty map: true
	// Empty string: false
	// Non empty string: true
	// true bool: true
	// false bool: false
	// 0 integer: false
	// positive integer: true
	// negative integer: true
	// 0 float: false
	// positive float: true
	// negative integer: true
	// struct: true
	// nil: false
}

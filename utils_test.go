package raymond

import "fmt"

func ExampleIsTruth() {
	output := "Empty array: " + Str(IsTruth([0]string{})) + "\n"
	output += "Non empty array: " + Str(IsTruth([1]string{"foo"})) + "\n"

	output += "Empty slice: " + Str(IsTruth([]string{})) + "\n"
	output += "Non empty slice: " + Str(IsTruth([]string{"foo"})) + "\n"

	output += "Empty map: " + Str(IsTruth(map[string]string{})) + "\n"
	output += "Non empty map: " + Str(IsTruth(map[string]string{"foo": "bar"})) + "\n"

	output += "Empty string: " + Str(IsTruth("")) + "\n"
	output += "Non empty string: " + Str(IsTruth("foo")) + "\n"

	output += "true bool: " + Str(IsTruth(true)) + "\n"
	output += "false bool: " + Str(IsTruth(false)) + "\n"

	output += "0 integer: " + Str(IsTruth(0)) + "\n"
	output += "positive integer: " + Str(IsTruth(10)) + "\n"
	output += "negative integer: " + Str(IsTruth(-10)) + "\n"

	output += "0 float: " + Str(IsTruth(0.0)) + "\n"
	output += "positive float: " + Str(IsTruth(10.0)) + "\n"
	output += "negative integer: " + Str(IsTruth(-10.0)) + "\n"

	output += "struct: " + Str(IsTruth(struct{}{})) + "\n"
	output += "nil: " + Str(IsTruth(nil)) + "\n"

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

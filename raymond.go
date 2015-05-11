package raymond

import "bytes"

// Renders a template with input data and returns result.
//
// Panics on error.
//
// Note that this function call is not optimal as your template is parsed
// everytime you call it. You should use `Parse()` function instead.
func Render(source string, data interface{}) string {
	buf := new(bytes.Buffer)

	// parse template
	tpl, err := Parse(source)
	if err != nil {
		panic(err)
	}

	// renders template
	if err := tpl.Exec(buf, data); err != nil {
		panic(err)
	}

	return buf.String()
}

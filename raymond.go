package raymond

// Renders a template with input data and returns result.
//
// Note that this function call is not optimal as your template is parsed
// everytime you call it. You should use `Parse()` function instead.
func Render(source string, data interface{}) (string, error) {
	// parse template
	tpl, err := Parse(source)
	if err != nil {
		return "", err
	}

	// renders template
	str, err := tpl.Exec(data)
	if err != nil {
		return "", err
	}

	return str, nil
}

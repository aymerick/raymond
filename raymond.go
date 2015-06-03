package raymond

// Render parses a template and evaluates it with given context
//
// Note that this function call is not optimal as your template is parsed everytime you call it. You should use Parse() function instead.
func Render(source string, ctx interface{}) (string, error) {
	// parse template
	tpl, err := Parse(source)
	if err != nil {
		return "", err
	}

	// renders template
	str, err := tpl.Exec(ctx)
	if err != nil {
		return "", err
	}

	return str, nil
}

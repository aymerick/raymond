// Package raymond provides handlebars evaluation
//
// Basic usage:
//
//  source := "<h1>{{title}}</h1><div class="body">{{body}}</div>"
//
//  ctx := map[string]string{
//      "title": "My New Post",
//      "body":  "This is my first post!",
//  }
//
//  tpl := raymond.MustParse(source)
//  result := tpl.MustExec(ctx)
//
//  // alternatively, for one shot templating:
//  result := raymond.MustRender(source, ctx)
//
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

// MustRender parses a template and evaluates it with given context. Panics on error
//
// Note that this function call is not optimal as your template is parsed everytime you call it. You should use Parse() function instead.
func MustRender(source string, ctx interface{}) string {
	return MustParse(source).MustExec(ctx)
}

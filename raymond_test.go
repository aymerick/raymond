package raymond

import "fmt"

func Example() {
	source := "<h1>{{title}}</h1><p>{{body.content}}</p>"

	ctx := map[string]interface{}{
		"title": "foo",
		"body":  map[string]string{"content": "bar"},
	}

	// parse template
	tpl := MustParse(source)

	// evaluate template with context
	output := tpl.MustExec(ctx)

	// alternatively, for one shots:
	// output :=  MustRender(source, ctx)

	fmt.Print(output)
	// Output: <h1>foo</h1><p>bar</p>
}

func ExampleRender() {
	tpl := "<h1>{{title}}</h1><p>{{body.content}}</p>"

	ctx := map[string]interface{}{
		"title": "foo",
		"body":  map[string]string{"content": "bar"},
	}

	// render template with context
	output, err := Render(tpl, ctx)
	if err != nil {
		panic(err)
	}

	fmt.Print(output)
	// Output: <h1>foo</h1><p>bar</p>
}

func ExampleMustRender() {
	tpl := "<h1>{{title}}</h1><p>{{body.content}}</p>"

	ctx := map[string]interface{}{
		"title": "foo",
		"body":  map[string]string{"content": "bar"},
	}

	// render template with context
	output := MustRender(tpl, ctx)

	fmt.Print(output)
	// Output: <h1>foo</h1><p>bar</p>
}

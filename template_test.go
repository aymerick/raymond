package ray

import (
	"fmt"
	"testing"
)

var sourceBasic = `<div class="entry">
  <h1>{{title}}</h1>
  <div class="body">
    {{body}}
  </div>
</div>`

var basicAST = `CONTENT[ '<div class="entry">
  <h1>' ]
{{ PATH:title [] }}
CONTENT[ '</h1>
  <div class="body">
    ' ]
{{ PATH:body [] }}
CONTENT[ '
  </div>
</div>' ]
`

func TestNewTemplate(t *testing.T) {
	t.Parallel()

	tpl := newTemplate(sourceBasic)
	if tpl.source != sourceBasic {
		t.Errorf("Failed to instantiate template")
	}
}

func TestParse(t *testing.T) {
	t.Parallel()

	tpl, err := Parse(sourceBasic)
	if err != nil || (tpl.source != sourceBasic) {
		t.Errorf("Failed to parse template")
	}

	if str := tpl.PrintAST(); str != basicAST {
		t.Errorf("Template parsing incorrect: %s", str)
	}
}

func TestClone(t *testing.T) {
	t.Parallel()

	sourcePartial := `I am a {{wat}} partial`
	sourcePartial2 := `Partial for the {{wat}}`

	tpl := MustParse(sourceBasic)
	tpl.RegisterPartial("p", sourcePartial)

	if (len(tpl.partials) != 1) || (tpl.partials["p"] == nil) {
		t.Errorf("What?")
	}

	cloned := tpl.Clone()

	if (len(cloned.partials) != 1) || (cloned.partials["p"] == nil) {
		t.Errorf("Template partials must be cloned")
	}

	cloned.RegisterPartial("p2", sourcePartial2)

	if (len(cloned.partials) != 2) || (cloned.partials["p"] == nil) || (cloned.partials["p2"] == nil) {
		t.Errorf("Failed to register a partial on cloned template")
	}

	if (len(tpl.partials) != 1) || (tpl.partials["p"] == nil) {
		t.Errorf("Modification of a cloned template MUST NOT affect original template")
	}
}

func ExampleTemplate_Exec() {
	source := "<h1>{{title}}</h1><p>{{body.content}}</p>"

	ctx := map[string]interface{}{
		"title": "foo",
		"body":  map[string]string{"content": "bar"},
	}

	// parse template
	tpl := MustParse(source)

	// evaluate template with context
	output, err := tpl.Exec(ctx)
	if err != nil {
		panic(err)
	}

	fmt.Print(output)
	// Output: <h1>foo</h1><p>bar</p>
}

func ExampleTemplate_MustExec() {
	source := "<h1>{{title}}</h1><p>{{body.content}}</p>"

	ctx := map[string]interface{}{
		"title": "foo",
		"body":  map[string]string{"content": "bar"},
	}

	// parse template
	tpl := MustParse(source)

	// evaluate template with context
	output := tpl.MustExec(ctx)

	fmt.Print(output)
	// Output: <h1>foo</h1><p>bar</p>
}

func ExampleTemplate_ExecWith() {
	source := "<h1>{{title}}</h1><p>{{#body}}{{content}} and {{@baz.bat}}{{/body}}</p>"

	ctx := map[string]interface{}{
		"title": "foo",
		"body":  map[string]string{"content": "bar"},
	}

	// parse template
	tpl := MustParse(source)

	// computes private data frame
	frame := NewDataFrame()
	frame.Set("baz", map[string]string{"bat": "unicorns"})

	// evaluate template
	output, err := tpl.ExecWith(ctx, frame)
	if err != nil {
		panic(err)
	}

	fmt.Print(output)
	// Output: <h1>foo</h1><p>bar and unicorns</p>
}

func ExampleTemplate_PrintAST() {
	source := "<h1>{{title}}</h1><p>{{#body}}{{content}} and {{@baz.bat}}{{/body}}</p>"

	// parse template
	tpl := MustParse(source)

	// print AST
	output := tpl.PrintAST()

	fmt.Print(output)
	// Output: CONTENT[ '<h1>' ]
	// {{ PATH:title [] }}
	// CONTENT[ '</h1><p>' ]
	// BLOCK:
	//   PATH:body []
	//   PROGRAM:
	//     {{     PATH:content []
	//  }}
	//     CONTENT[ ' and ' ]
	//     {{     @PATH:baz/bat []
	//  }}
	//   CONTENT[ '</p>' ]
	//
}

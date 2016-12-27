package ray

import (
	"fmt"
	"testing"

	"github.com/stretchr/testify/require"
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
	r := require.New(t)

	tpl, err := Parse(sourceBasic)
	r.NoError(err)
	r.Equal(sourceBasic, tpl.source)
	str := tpl.PrintAST()
	r.Equal(basicAST, str)
}

func TestClone(t *testing.T) {
	t.Parallel()
	r := require.New(t)

	sourcePartial := `I am a {{wat}} partial`
	sourcePartial2 := `Partial for the {{wat}}`

	tpl := MustParse(sourceBasic)
	tpl.RegisterPartial("p", sourcePartial)

	r.Len(tpl.partials, 1)
	r.NotNil(tpl.partials["p"])

	cloned := tpl.Clone()

	r.Len(cloned.partials, 1)
	r.NotNil(cloned.partials["p"])

	cloned.RegisterPartial("p2", sourcePartial2)

	r.Len(cloned.partials, 2)
	r.NotNil(cloned.partials["p"])
	r.NotNil(cloned.partials["p2"])

	r.Len(tpl.partials, 1)
	r.NotNil(tpl.partials["p"])
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

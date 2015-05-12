package raymond

import (
	"bytes"
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
	tpl := NewTemplate(sourceBasic)
	if tpl.source != sourceBasic {
		t.Errorf("Faild to instantiate template")
	}
}

func TestParseTemplate(t *testing.T) {
	tpl, err := Parse(sourceBasic)
	if err != nil || (tpl.source != sourceBasic) {
		t.Errorf("Faild to parse template")
	}

	if str := tpl.PrintAST(); str != basicAST {
		t.Errorf("Template parsing incorrect: %s", str)
	}
}

type tplTest struct {
	name   string
	input  string
	data   interface{}
	output string
}

var tplTests = []tplTest{
	{"only content", "this is content", nil, "this is content"},

	//
	// Next tests come from:
	//   https://github.com/wycats/handlebars.js/blob/master/spec/basic.js
	//
	{"most basic", "{{foo}}", map[string]string{"foo": "foo"}, "foo"},
}

func TestRenderTemplate(t *testing.T) {
	for _, test := range tplTests {
		var err error
		var tpl *Template

		buf := new(bytes.Buffer)

		// parse template
		tpl, err = Parse(test.input)
		if err != nil {
			t.Errorf("Test '%s' failed - Failed to parse template\ninput:\n\t'%s'\nerror:\n\t%s", test.name, test.input, err)
		} else {
			// render template
			err = tpl.Exec(buf, test.data)
			if err != nil {
				t.Errorf("Test '%s' failed\ninput:\n\t'%s'\nerror:\n\t%s\nAST:\n\t%s", test.name, test.input, err, tpl.PrintAST())
			} else {
				// check output
				output := buf.String()
				if test.output != output {
					t.Errorf("Test '%s' failed\ninput:\n\t'%s'\nexpected\n\t%q\ngot\n\t%q\nAST:\n\t%s", test.name, test.input, test.output, output, tpl.PrintAST())
				}
			}
		}
	}
}

package raymond

import (
	"bytes"
	"testing"
)

//
// Basic rendering test
//

var testInput = `<div class="entry">
  <h1>{{title}}</h1>
  <div class="body">
    {{body}}
  </div>
</div>`

var testOutput = `<div class="entry">
  <h1>foo</h1>
  <div class="body">
    bar
  </div>
</div>`

func TestRender(t *testing.T) {
	output := Render(testInput, map[string]string{"title": "foo", "body": "bar"})
	if output != testOutput {
		t.Errorf("Failed to render template\ninput:\n\n'%s'\n\nexpected:\n\n%s\n\ngot:\n\n%s", testInput, testOutput, output)
	}
}

//
// Generic test
//

type raymondTest struct {
	name    string
	input   string
	data    interface{}
	helpers map[string]Helper
	output  string
}

func launchRaymondTests(t *testing.T, tests []raymondTest) {
	for _, test := range tests {
		var err error
		var tpl *Template

		buf := new(bytes.Buffer)

		// parse template
		tpl, err = Parse(test.input)
		if err != nil {
			t.Errorf("Test '%s' failed - Failed to parse template\ninput:\n\t'%s'\nerror:\n\t%s", test.name, test.input, err)
		} else {
			if len(test.helpers) > 0 {
				// register helpers
				tpl.RegisterHelpers(test.helpers)
			}

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

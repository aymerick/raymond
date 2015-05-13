package raymond

import (
	"bytes"
	"log"
	"testing"
)

const (
	VERBOSE = false
)

type helperTest struct {
	name    string
	input   string
	data    interface{}
	helpers map[string]Helper
	output  string
}

var helperTests = []helperTest{
	{
		"simple helper",
		`{{foo}}`,
		nil,
		map[string]Helper{"foo": barHelper},
		`bar`,
	},
	{
		"helper with literal string param",
		`{{echo "foo"}}`,
		nil,
		map[string]Helper{"echo": echoHelper},
		`foo`,
	},
	{
		"helper with identifier param",
		`{{echo foo}}`,
		map[string]interface{}{"foo": "bar"},
		map[string]Helper{"echo": echoHelper},
		`bar`,
	},
	{
		"helper with literal boolean param",
		`{{bool true}}`,
		nil,
		map[string]Helper{"bool": boolHelper},
		`yes it is`,
	},
	{
		"helper with literal boolean param",
		`{{bool false}}`,
		nil,
		map[string]Helper{"bool": boolHelper},
		`absolutely not`,
	},
	{
		"helper with literal boolean param",
		`{{gnak 5}}`,
		nil,
		map[string]Helper{"gnak": gnakHelper},
		`GnAK!GnAK!GnAK!GnAK!GnAK!`,
	},
	{
		"helper with several parameters",
		`{{echo "GnAK!" 3}}`,
		nil,
		map[string]Helper{"echo": echoHelper},
		`GnAK!GnAK!GnAK!`,
	},

	//
	// Next tests come from:
	//   https://github.com/wycats/handlebars.js/blob/master/spec/helper.js
	//

	// {
	// 	"helper with complex lookup",
	// 	"{{#goodbyes}}{{{link ../prefix}}}{{/goodbyes}}",
	// 	map[string]interface{}{"prefix": "/root", "goodbyes": []map[string]string{{"text": "Goodbye", "url": "goodbye"}}},
	// 	map[string]Helper{"link": linkHelper},
	// 	`<a href="/root/goodbye">Goodbye</a>`,
	// },
}

//
// Helpers
//

func barHelper(p *HelperParams) string { return "bar" }

func echoHelper(p *HelperParams) string {
	str, _ := p.at(0).(string)
	nb, ok := p.at(1).(int)
	if !ok {
		nb = 1
	}

	result := ""
	for i := 0; i < nb; i++ {
		result += str
	}

	return result
}

func boolHelper(p *HelperParams) string {
	b, _ := p.at(0).(bool)
	if b {
		return "yes it is"
	}

	return "absolutely not"
}

func gnakHelper(p *HelperParams) string {
	nb, ok := p.at(0).(int)
	if !ok {
		nb = 1
	}

	result := ""
	for i := 0; i < nb; i++ {
		result += "GnAK!"
	}

	return result
}

//
// Let's go
//

func TestHelper(t *testing.T) {
	for _, test := range helperTests {
		if VERBOSE {
			log.Printf("\n\n**********************************")
			log.Printf("Testing: %s", test.name)
		}

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

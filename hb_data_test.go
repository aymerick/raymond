package raymond

import "testing"

//
// Those tests come from:
//   https://github.com/wycats/handlebars.js/blob/master/spec/data.js
//
var hbDataTests = []raymondTest{
	{
		"passing in data to a compiled function that expects data - works with helpers",
		"{{hello}}",
		map[string]string{"noun": "cat"},
		map[string]interface{}{"adjective": "happy"},
		map[string]Helper{"hello": func(h *HelperArg) interface{} {
			return h.DataStr("adjective") + " " + h.FieldStr("noun")
		}},
		nil,
		"happy cat",
	},
}

func TestHandlebarsData(t *testing.T) {
	launchHandlebarsTests(t, hbDataTests)
}

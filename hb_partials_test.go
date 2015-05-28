package raymond

import "testing"

//
// Those tests come from:
//   https://github.com/wycats/handlebars.js/blob/master/spec/partials.js
//
var hbPartialsTests = []raymondTest{
	{
		"basic partials",
		"Dudes: {{#dudes}}{{> dude}}{{/dudes}}",
		map[string]interface{}{"dudes": []map[string]string{{"name": "Yehuda", "url": "http://yehuda"}, {"name": "Alan", "url": "http://alan"}}},
		nil,
		map[string]string{"dude": "{{name}} ({{url}}) "},
		"Dudes: Yehuda (http://yehuda) Alan (http://alan) ",
	},
}

func TestHandlebarsPartials(t *testing.T) {
	launchHandlebarsTests(t, hbPartialsTests)
}

package raymond

import "testing"

//
// Those tests come from:
//   https://github.com/wycats/handlebars.js/blob/master/spec/helper.js
//
var hbHelpersTests = []raymondTest{
	{
		"helper with complex lookup",
		"{{#goodbyes}}{{{link ../prefix}}}{{/goodbyes}}",
		map[string]interface{}{"prefix": "/root", "goodbyes": []map[string]string{{"text": "Goodbye", "url": "goodbye"}}},
		map[string]Helper{"link": linkHelper},
		"<a href=\"/root/goodbye\">Goodbye</a>",
	},
	{
		"helper for raw block gets raw content",
		"{{{{raw}}}} {{test}} {{{{/raw}}}}",
		map[string]interface{}{"test": "hello"},
		map[string]Helper{"raw": rawHelper},
		" {{test}} ",
	},
	{
		"helper for raw block gets parameters",
		"{{{{raw 1 2 3}}}} {{test}} {{{{/raw}}}}",
		map[string]interface{}{"test": "hello"},
		map[string]Helper{"raw": rawHelper},
		" {{test}} 123",
	},
	{
		"helper block with complex lookup expression",
		"{{#goodbyes}}{{../name}}{{/goodbyes}}",
		map[string]interface{}{"name": "Alan"},
		map[string]Helper{"goodbyes": func(p *HelperParams) string {
			// @todo This is ugly, compared to the JS implementation API.
			//       We should capture the result of block evaluation.
			for _, str := range []string{"Goodbye", "goodbye", "GOODBYE"} {
				p.Write(str + " ")
				p.EvaluateBlockWith(str)
				p.Write("! ")
			}
			return ""
		}},
		"Goodbye Alan! goodbye Alan! GOODBYE Alan! ",
	},

	// @todo Add remaining tests
}

func TestHandlebarsHelpers(t *testing.T) {
	launchRaymondTests(t, hbHelpersTests)
}

package raymond

import "testing"

//
// Those tests come from:
//   https://github.com/wycats/handlebars.js/blob/master/spec/helper.js
//
var hbHelpersTests = []raymondTest{
	// @todo !!
	{},
	// {
	// 	"helper with complex lookup",
	// 	"{{#goodbyes}}{{{link ../prefix}}}{{/goodbyes}}",
	// 	map[string]interface{}{"prefix": "/root", "goodbyes": []map[string]string{{"text": "Goodbye", "url": "goodbye"}}},
	// 	map[string]Helper{"link": linkHelper},
	// 	`<a href="/root/goodbye">Goodbye</a>`,
	// },
}

func TestHandlebarsHelpers(t *testing.T) {
	launchRaymondTests(t, hbHelpersTests)
}

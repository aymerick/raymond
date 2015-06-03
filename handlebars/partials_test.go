package handlebars

import (
	"testing"

	"github.com/aymerick/raymond"
)

//
// Those tests come from:
//   https://github.com/wycats/handlebars.js/blob/master/spec/partials.js
//
var partialsTests = []Test{
	{
		"basic partials",
		"Dudes: {{#dudes}}{{> dude}}{{/dudes}}",
		map[string]interface{}{"dudes": []map[string]string{{"name": "Yehuda", "url": "http://yehuda"}, {"name": "Alan", "url": "http://alan"}}},
		nil, nil,
		map[string]string{"dude": "{{name}} ({{url}}) "},
		"Dudes: Yehuda (http://yehuda) Alan (http://alan) ",
	},
	{
		"dynamic partials",
		"Dudes: {{#dudes}}{{> (partial)}}{{/dudes}}",
		map[string]interface{}{"dudes": []map[string]string{{"name": "Yehuda", "url": "http://yehuda"}, {"name": "Alan", "url": "http://alan"}}},
		nil,
		map[string]raymond.Helper{"partial": func(h *raymond.HelperArg) interface{} {
			return "dude"
		}},
		map[string]string{"dude": "{{name}} ({{url}}) "},
		"Dudes: Yehuda (http://yehuda) Alan (http://alan) ",
	},

	// @todo "failing dynamic partials"

	{
		"partials with context",
		"Dudes: {{>dude dudes}}",
		map[string]interface{}{"dudes": []map[string]string{{"name": "Yehuda", "url": "http://yehuda"}, {"name": "Alan", "url": "http://alan"}}},
		nil, nil,
		map[string]string{"dude": "{{#this}}{{name}} ({{url}}) {{/this}}"},
		"Dudes: Yehuda (http://yehuda) Alan (http://alan) ",
	},
	{
		"partials with undefined context",
		"Dudes: {{>dude dudes}}",
		map[string]interface{}{},
		nil, nil,
		map[string]string{"dude": "{{foo}} Empty"},
		"Dudes:  Empty",
	},

	// @todo "partials with duplicate parameters"

	{
		"partials with parameters",
		"Dudes: {{#dudes}}{{> dude others=..}}{{/dudes}}",
		map[string]interface{}{"foo": "bar", "dudes": []map[string]string{{"name": "Yehuda", "url": "http://yehuda"}, {"name": "Alan", "url": "http://alan"}}},
		nil, nil,
		map[string]string{"dude": "{{others.foo}}{{name}} ({{url}}) "},
		"Dudes: barYehuda (http://yehuda) barAlan (http://alan) ",
	},
	{
		"partial in a partial",
		"Dudes: {{#dudes}}{{>dude}}{{/dudes}}",
		map[string]interface{}{"dudes": []map[string]string{{"name": "Yehuda", "url": "http://yehuda"}, {"name": "Alan", "url": "http://alan"}}},
		nil, nil,
		map[string]string{"dude": "{{name}} {{> url}} ", "url": `<a href="{{url}}">{{url}}</a>`},
		`Dudes: Yehuda <a href="http://yehuda">http://yehuda</a> Alan <a href="http://alan">http://alan</a> `,
	},

	// @todo "rendering undefined partial throws an exception"

	// @todo "registering undefined partial throws an exception"

	// SKIP: "rendering template partial in vm mode throws an exception"
	// SKIP: "rendering function partial in vm mode"

	{
		"GH-14: a partial preceding a selector",
		"Dudes: {{>dude}} {{anotherDude}}",
		map[string]string{"name": "Jeepers", "anotherDude": "Creepers"},
		nil, nil,
		map[string]string{"dude": "{{name}}"},
		"Dudes: Jeepers Creepers",
	},
	{
		"Partials with slash paths",
		"Dudes: {{> shared/dude}}",
		map[string]string{"name": "Jeepers", "anotherDude": "Creepers"},
		nil, nil,
		map[string]string{"shared/dude": "{{name}}"},
		"Dudes: Jeepers",
	},
	{
		"Partials with slash and point paths",
		"Dudes: {{> shared/dude.thing}}",
		map[string]string{"name": "Jeepers", "anotherDude": "Creepers"},
		nil, nil,
		map[string]string{"shared/dude.thing": "{{name}}"},
		"Dudes: Jeepers",
	},

	// @todo "Global Partials"

	// @todo "Multiple partial registration"

	{
		"Partials with integer path",
		"Dudes: {{> 404}}",
		map[string]string{"name": "Jeepers", "anotherDude": "Creepers"},
		nil, nil,
		map[string]string{"404": "{{name}}"}, // @note Difference with JS test: partial name is a string
		"Dudes: Jeepers",
	},
	// @note This is not supported by our implementation. But really... who cares ?
	// {
	// 	"Partials with complex path",
	// 	"Dudes: {{> 404/asdf?.bar}}",
	// 	map[string]string{"name": "Jeepers", "anotherDude": "Creepers"},
	// 	nil, nil,
	// 	map[string]string{"404/asdf?.bar": "{{name}}"},
	// 	"Dudes: Jeepers",
	// },
	{
		"Partials with escaped",
		"Dudes: {{> [+404/asdf?.bar]}}",
		map[string]string{"name": "Jeepers", "anotherDude": "Creepers"},
		nil, nil,
		map[string]string{"+404/asdf?.bar": "{{name}}"},
		"Dudes: Jeepers",
	},
	{
		"Partials with string",
		"Dudes: {{> '+404/asdf?.bar'}}",
		map[string]string{"name": "Jeepers", "anotherDude": "Creepers"},
		nil, nil,
		map[string]string{"+404/asdf?.bar": "{{name}}"},
		"Dudes: Jeepers",
	},
	{
		"should handle empty partial",
		"Dudes: {{#dudes}}{{> dude}}{{/dudes}}",
		map[string]interface{}{"dudes": []map[string]string{{"name": "Yehuda", "url": "http://yehuda"}, {"name": "Alan", "url": "http://alan"}}},
		nil, nil,
		map[string]string{"dude": ""},
		"Dudes: ",
	},

	// @todo "throw on missing partial"

	// SKIP: "should pass compiler flags"

	{
		"standalone partials (1) - indented partials",
		"Dudes:\n{{#dudes}}\n  {{>dude}}\n{{/dudes}}",
		map[string]interface{}{"dudes": []map[string]string{{"name": "Yehuda", "url": "http://yehuda"}, {"name": "Alan", "url": "http://alan"}}},
		nil, nil,
		map[string]string{"dude": "{{name}}\n"},
		"Dudes:\n  Yehuda\n  Alan\n",
	},
	{
		"standalone partials (2) - nested indented partials",
		"Dudes:\n{{#dudes}}\n  {{>dude}}\n{{/dudes}}",
		map[string]interface{}{"dudes": []map[string]string{{"name": "Yehuda", "url": "http://yehuda"}, {"name": "Alan", "url": "http://alan"}}},
		nil, nil,
		map[string]string{"dude": "{{name}}\n {{> url}}", "url": "{{url}}!\n"},
		"Dudes:\n  Yehuda\n   http://yehuda!\n  Alan\n   http://alan!\n",
	},

	// // @todo preventIndent option
	// {
	// 	"standalone partials (3) - prevent nested indented partials",
	// 	"Dudes:\n{{#dudes}}\n  {{>dude}}\n{{/dudes}}",
	// 	map[string]interface{}{"dudes": []map[string]string{{"name": "Yehuda", "url": "http://yehuda"}, {"name": "Alan", "url": "http://alan"}}},
	// 	nil, nil,
	// 	map[string]string{"dude": "{{name}}\n {{> url}}", "url": "{{url}}!\n"},
	// 	"Dudes:\n  Yehuda\n http://yehuda!\n  Alan\n http://alan!\n",
	// },

	// @todo "compat mode"
}

func TestPartials(t *testing.T) {
	launchTests(t, partialsTests)
}

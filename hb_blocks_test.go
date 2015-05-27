package raymond

import "testing"

//
// Those tests come from:
//   https://github.com/wycats/handlebars.js/blob/master/spec/blocks.js
//
var hbBlocksTests = []raymondTest{
	{
		"array (1) - Arrays iterate over the contents when not empty",
		"{{#goodbyes}}{{text}}! {{/goodbyes}}cruel {{world}}!",
		map[string]interface{}{"goodbyes": []map[string]string{{"text": "goodbye"}, {"text": "Goodbye"}, {"text": "GOODBYE"}}, "world": "world"},
		nil,
		"goodbye! Goodbye! GOODBYE! cruel world!",
	},
	{
		"array (2) - Arrays ignore the contents when empty",
		"{{#goodbyes}}{{text}}! {{/goodbyes}}cruel {{world}}!",
		map[string]interface{}{"goodbyes": []map[string]string{}, "world": "world"},
		nil,
		"cruel world!",
	},
	{
		"array without data",
		"{{#goodbyes}}{{text}}{{/goodbyes}} {{#goodbyes}}{{text}}{{/goodbyes}}",
		map[string]interface{}{"goodbyes": []map[string]string{{"text": "goodbye"}, {"text": "Goodbye"}, {"text": "GOODBYE"}}, "world": "world"},
		nil,
		"goodbyeGoodbyeGOODBYE goodbyeGoodbyeGOODBYE",
	},
	// {
	// 	"array with @index - The @index variable is used",
	// 	"{{#goodbyes}}{{@index}}. {{text}}! {{/goodbyes}}cruel {{world}}!",
	// 	map[string]interface{}{"goodbyes": []map[string]string{{"text": "goodbye"}, {"text": "Goodbye"}, {"text": "GOODBYE"}}, "world": "world"},
	// 	nil,
	// 	"0. goodbye! 1. Goodbye! 2. GOODBYE! cruel world!",
	// },
	{
		"empty block (1) - Arrays iterate over the contents when not empty",
		"{{#goodbyes}}{{/goodbyes}}cruel {{world}}!",
		map[string]interface{}{"goodbyes": []map[string]string{{"text": "goodbye"}, {"text": "Goodbye"}, {"text": "GOODBYE"}}, "world": "world"},
		nil,
		"cruel world!",
	},
	{
		"empty block (1) - Arrays ignore the contents when empty",
		"{{#goodbyes}}{{/goodbyes}}cruel {{world}}!",
		map[string]interface{}{"goodbyes": []map[string]string{}, "world": "world"},
		nil,
		"cruel world!",
	},
	{
		"block with complex lookup - Templates can access variables in contexts up the stack with relative path syntax",
		"{{#goodbyes}}{{text}} cruel {{../name}}! {{/goodbyes}}",
		map[string]interface{}{"goodbyes": []map[string]string{{"text": "goodbye"}, {"text": "Goodbye"}, {"text": "GOODBYE"}}, "name": "Alan"},
		nil,
		"goodbye cruel Alan! Goodbye cruel Alan! GOODBYE cruel Alan! ",
	},
	{
		"multiple blocks with complex lookup",
		"{{#goodbyes}}{{../name}}{{../name}}{{/goodbyes}}",
		map[string]interface{}{"goodbyes": []map[string]string{{"text": "goodbye"}, {"text": "Goodbye"}, {"text": "GOODBYE"}}, "name": "Alan"},
		nil,
		"AlanAlanAlanAlanAlanAlan",
	},

	// @todo "{{#goodbyes}}{{text}} cruel {{foo/../name}}! {{/goodbyes}}" should throw error

	{
		"block with deep nested complex lookup",
		"{{#outer}}Goodbye {{#inner}}cruel {{../sibling}} {{../../omg}}{{/inner}}{{/outer}}",
		map[string]interface{}{"omg": "OMG!", "outer": []map[string]interface{}{{"sibling": "sad", "inner": []map[string]string{{"text": "goodbye"}}}}},
		nil,
		"Goodbye cruel sad OMG!",
	},
	{
		"inverted sections with unset value - Inverted section rendered when value isn't set.",
		"{{#goodbyes}}{{this}}{{/goodbyes}}{{^goodbyes}}Right On!{{/goodbyes}}",
		map[string]interface{}{},
		nil,
		"Right On!",
	},
	{
		"inverted sections with false value - Inverted section rendered when value is false.",
		"{{#goodbyes}}{{this}}{{/goodbyes}}{{^goodbyes}}Right On!{{/goodbyes}}",
		map[string]interface{}{"goodbyes": false},
		nil,
		"Right On!",
	},
	{
		"inverted section with empty set - Inverted section rendered when value is empty set.",
		"{{#goodbyes}}{{this}}{{/goodbyes}}{{^goodbyes}}Right On!{{/goodbyes}}",
		map[string]interface{}{"goodbyes": []interface{}{}},
		nil,
		"Right On!",
	},
	{
		"block inverted sections",
		"{{#people}}{{name}}{{^}}{{none}}{{/people}}",
		map[string]interface{}{"none": "No people"},
		nil,
		"No people",
	},
	{
		"chained inverted sections (1)",
		"{{#people}}{{name}}{{else if none}}{{none}}{{/people}}",
		map[string]interface{}{"none": "No people"},
		nil,
		"No people",
	},
	{
		"chained inverted sections (2)",
		"{{#people}}{{name}}{{else if nothere}}fail{{else unless nothere}}{{none}}{{/people}}",
		map[string]interface{}{"none": "No people"},
		nil,
		"No people",
	},
	{
		"chained inverted sections (3)",
		"{{#people}}{{name}}{{else if none}}{{none}}{{else}}fail{{/people}}",
		map[string]interface{}{"none": "No people"},
		nil,
		"No people",
	},

	// @todo "{{#people}}{{name}}{{else if none}}{{none}}{{/if}}" should throw error

	{
		"block inverted sections with empty arrays",
		"{{#people}}{{name}}{{^}}{{none}}{{/people}}",
		map[string]interface{}{"none": "No people", "people": map[string]interface{}{}},
		nil,
		"No people",
	},
	{
		"block standalone else sections (1)",
		"{{#people}}\n{{name}}\n{{^}}\n{{none}}\n{{/people}}\n",
		map[string]interface{}{"none": "No people"},
		nil,
		"No people\n",
	},
	{
		"block standalone else sections (2)",
		"{{#none}}\n{{.}}\n{{^}}\n{{none}}\n{{/none}}\n",
		map[string]interface{}{"none": "No people"},
		nil,
		"No people\n",
	},
	{
		"block standalone else sections (3)",
		"{{#people}}\n{{name}}\n{{^}}\n{{none}}\n{{/people}}\n",
		map[string]interface{}{"none": "No people"},
		nil,
		"No people\n",
	},
	{
		"block standalone chained else sections (1)",
		"{{#people}}\n{{name}}\n{{else if none}}\n{{none}}\n{{/people}}\n",
		map[string]interface{}{"none": "No people"},
		nil,
		"No people\n",
	},
	{
		"block standalone chained else sections (2)",
		"{{#people}}\n{{name}}\n{{else if none}}\n{{none}}\n{{^}}\n{{/people}}\n",
		map[string]interface{}{"none": "No people"},
		nil,
		"No people\n",
	},
	{
		"should handle nesting",
		"{{#data}}\n{{#if true}}\n{{.}}\n{{/if}}\n{{/data}}\nOK.",
		map[string]interface{}{"data": []int{1, 3, 5}},
		nil,
		"1\n3\n5\nOK.",
	},
	// // @todo compat mode
	// {
	// 	"block with deep recursive lookup lookup",
	// 	"{{#outer}}Goodbye {{#inner}}cruel {{omg}}{{/inner}}{{/outer}}",
	// 	map[string]interface{}{"omg": "OMG!", "outer": []map[string]interface{}{{"inner": []map[string]string{{"text": "goodbye"}}}}},
	// 	nil,
	// 	"Goodbye cruel OMG!",
	// },
	// // @todo compat mode
	// {
	// 	"block with deep recursive pathed lookup",
	// 	"{{#outer}}Goodbye {{#inner}}cruel {{omg.yes}}{{/inner}}{{/outer}}",
	// 	map[string]interface{}{"omg": map[string]string{"yes": "OMG!"}, "outer": []map[string]interface{}{{"inner": []map[string]string{{"yes": "no", "text": "goodbye"}}}}},
	// 	nil,
	// 	"Goodbye cruel OMG!",
	// },
	{
		"block with missed recursive lookup",
		"{{#outer}}Goodbye {{#inner}}cruel {{omg.yes}}{{/inner}}{{/outer}}",
		map[string]interface{}{"omg": map[string]string{"no": "OMG!"}, "outer": []map[string]interface{}{{"inner": []map[string]string{{"yes": "no", "text": "goodbye"}}}}},
		nil,
		"Goodbye cruel ",
	},
}

func TestHandlebarsBlocks(t *testing.T) {
	launchHandlebarsTests(t, hbBlocksTests)
}

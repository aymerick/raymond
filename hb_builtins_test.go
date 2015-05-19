package raymond

import "testing"

//
// Those tests come from:
//   https://github.com/wycats/handlebars.js/blob/master/spec/builtin.js
//
var hbBuiltinsTests = []raymondTest{
	{
		"#if - if with boolean argument shows the contents when true",
		"{{#if goodbye}}GOODBYE {{/if}}cruel {{world}}!",
		map[string]interface{}{"goodbye": true, "world": "world"},
		nil,
		"GOODBYE cruel world!",
	},
	{
		"#if - if with string argument shows the contents",
		"{{#if goodbye}}GOODBYE {{/if}}cruel {{world}}!",
		map[string]interface{}{"goodbye": "dummy", "world": "world"},
		nil,
		"GOODBYE cruel world!",
	},
	{
		"#if - if with boolean argument does not show the contents when false",
		"{{#if goodbye}}GOODBYE {{/if}}cruel {{world}}!",
		map[string]interface{}{"goodbye": false, "world": "world"},
		nil,
		"cruel world!",
	},
	{
		"#if - if with undefined does not show the contents",
		"{{#if goodbye}}GOODBYE {{/if}}cruel {{world}}!",
		map[string]interface{}{"world": "world"},
		nil,
		"cruel world!",
	},
	{
		"#if - if with non-empty array shows the contents",
		"{{#if goodbye}}GOODBYE {{/if}}cruel {{world}}!",
		map[string]interface{}{"goodbye": []string{"foo"}, "world": "world"},
		nil,
		"GOODBYE cruel world!",
	},
	{
		"#if - if with empty array does not show the contents",
		"{{#if goodbye}}GOODBYE {{/if}}cruel {{world}}!",
		map[string]interface{}{"goodbye": []string{}, "world": "world"},
		nil,
		"cruel world!",
	},
	{
		"#if - if with zero does not show the contents",
		"{{#if goodbye}}GOODBYE {{/if}}cruel {{world}}!",
		map[string]interface{}{"goodbye": 0, "world": "world"},
		nil,
		"cruel world!",
	},
	{
		"#if - if with zero and includeZero option shows the contents",
		"{{#if goodbye includeZero=true}}GOODBYE {{/if}}cruel {{world}}!",
		map[string]interface{}{"goodbye": 0, "world": "world"},
		nil,
		"GOODBYE cruel world!",
	},

	// {
	// 	"#if - if with function shows the contents when function returns true",
	// 	"{{#if goodbye}}GOODBYE {{/if}}cruel {{world}}!",
	// 	map[string]interface{}{},
	// 	nil,
	// 	"GOODBYE cruel world!",
	// },
	// {
	// 	"#if - if with function shows the contents when function returns string",
	// 	"{{#if goodbye}}GOODBYE {{/if}}cruel {{world}}!",
	// 	map[string]interface{}{},
	// 	nil,
	// 	"GOODBYE cruel world!",
	// },
	// {
	// 	"#if - if with function does not show the contents when returns false",
	// 	"{{#if goodbye}}GOODBYE {{/if}}cruel {{world}}!",
	// 	map[string]interface{}{},
	// 	nil,
	// 	"cruel world!",
	// },
	//    {
	//        "#if - if with function does not show the contents when returns undefined",
	//        "{{#if goodbye}}GOODBYE {{/if}}cruel {{world}}!",
	//        map[string]interface{}{},
	//        nil,
	//        "cruel world!",
	//    },

	{
		"#with",
		"{{#with person}}{{first}} {{last}}{{/with}}",
		map[string]interface{}{"person": map[string]string{"first": "Alan", "last": "Johnson"}},
		nil,
		"Alan Johnson",
	},
	// {
	//     "#with - with with function argument",
	//     "{{#with person}}{{first}} {{last}}{{/with}}",
	//     map[string]interface{}{},
	//     nil,
	//     "Alan Johnson",
	// },
	{
		"#with - with with else",
		"{{#with person}}Person is present{{else}}Person is not present{{/with}}",
		map[string]interface{}{},
		nil,
		"Person is not present",
	},

	{
		"#each - each with array argument iterates over the contents when not empty",
		"{{#each goodbyes}}{{text}}! {{/each}}cruel {{world}}!",
		map[string]interface{}{"goodbyes": []map[string]string{{"text": "goodbye"}, {"text": "Goodbye"}, {"text": "GOODBYE"}}, "world": "world"},
		nil,
		"goodbye! Goodbye! GOODBYE! cruel world!",
	},
	{
		"#each - each with array argument ignores the contents when empty",
		"{{#each goodbyes}}{{text}}! {{/each}}cruel {{world}}!",
		map[string]interface{}{"goodbyes": []map[string]string{}, "world": "world"},
		nil,
		"cruel world!",
	},
	{
		"#each - each without data (1)",
		"{{#each goodbyes}}{{text}}! {{/each}}cruel {{world}}!",
		map[string]interface{}{"goodbyes": []map[string]string{{"text": "goodbye"}, {"text": "Goodbye"}, {"text": "GOODBYE"}}, "world": "world"},
		nil,
		"goodbye! Goodbye! GOODBYE! cruel world!",
	},
	{
		"#each - each without data (2)",
		"{{#each .}}{{.}}{{/each}}",
		map[string]interface{}{"goodbyes": "cruel", "world": "world"},
		nil,
		// note: a go hash is not ordered, so result may vary, this behaviour differs from the JS implementation
		[]string{"cruelworld", "worldcruel"},
	},
	{
		"#each - each without context",
		"{{#each goodbyes}}{{text}}! {{/each}}cruel {{world}}!",
		nil,
		nil,
		"cruel !",
	},

	// @todo "#each - each with an object and @key"

	// {
	// 	"#each - each with @index",
	// 	"{{#each goodbyes}}{{@index}}. {{text}}! {{/each}}cruel {{world}}!",
	// 	map[string]interface{}{"goodbyes": []map[string]string{{"text": "goodbye"}, {"text": "Goodbye"}, {"text": "GOODBYE"}}, "world": "world"},
	// 	nil,
	// 	"0. goodbye! 1. Goodbye! 2. GOODBYE! cruel world!",
	// },
	// {
	// 	"#each - each with nested @index",
	// 	"{{#each goodbyes}}{{@index}}. {{text}}! {{#each ../goodbyes}}{{@index}} {{/each}}After {{@index}} {{/each}}{{@index}}cruel {{world}}!",
	// 	map[string]interface{}{"goodbyes": []map[string]string{{"text": "goodbye"}, {"text": "Goodbye"}, {"text": "GOODBYE"}}, "world": "world"},
	// 	nil,
	// 	"0. goodbye! 0 1 2 After 0 1. Goodbye! 0 1 2 After 1 2. GOODBYE! 0 1 2 After 2 cruel world!",
	// },

	// {
	// 	"#each - each with block params",
	// 	"{{#each goodbyes as |value index|}}{{index}}. {{value.text}}! {{#each ../goodbyes as |childValue childIndex|}} {{index}} {{childIndex}}{{/each}} After {{index}} {{/each}}{{index}}cruel {{world}}!",
	// 	map[string]interface{}{"goodbyes": []map[string]string{{"text": "goodbye"}, {"text": "Goodbye"}}, "world": "world"},
	// 	nil,
	// 	"0. goodbye!  0 0 0 1 After 0 1. Goodbye!  1 0 1 1 After 1 cruel world!",
	// },

	// @todo Add remaining tests
}

func TestHandlebarsBuiltins(t *testing.T) {
	launchRaymondTests(t, hbBuiltinsTests)
}

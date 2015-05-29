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
		nil, nil, nil,
		"GOODBYE cruel world!",
	},
	{
		"#if - if with string argument shows the contents",
		"{{#if goodbye}}GOODBYE {{/if}}cruel {{world}}!",
		map[string]interface{}{"goodbye": "dummy", "world": "world"},
		nil, nil, nil,
		"GOODBYE cruel world!",
	},
	{
		"#if - if with boolean argument does not show the contents when false",
		"{{#if goodbye}}GOODBYE {{/if}}cruel {{world}}!",
		map[string]interface{}{"goodbye": false, "world": "world"},
		nil, nil, nil,
		"cruel world!",
	},
	{
		"#if - if with undefined does not show the contents",
		"{{#if goodbye}}GOODBYE {{/if}}cruel {{world}}!",
		map[string]interface{}{"world": "world"},
		nil, nil, nil,
		"cruel world!",
	},
	{
		"#if - if with non-empty array shows the contents",
		"{{#if goodbye}}GOODBYE {{/if}}cruel {{world}}!",
		map[string]interface{}{"goodbye": []string{"foo"}, "world": "world"},
		nil, nil, nil,
		"GOODBYE cruel world!",
	},
	{
		"#if - if with empty array does not show the contents",
		"{{#if goodbye}}GOODBYE {{/if}}cruel {{world}}!",
		map[string]interface{}{"goodbye": []string{}, "world": "world"},
		nil, nil, nil,
		"cruel world!",
	},
	{
		"#if - if with zero does not show the contents",
		"{{#if goodbye}}GOODBYE {{/if}}cruel {{world}}!",
		map[string]interface{}{"goodbye": 0, "world": "world"},
		nil, nil, nil,
		"cruel world!",
	},
	{
		"#if - if with zero and includeZero option shows the contents",
		"{{#if goodbye includeZero=true}}GOODBYE {{/if}}cruel {{world}}!",
		map[string]interface{}{"goodbye": 0, "world": "world"},
		nil, nil, nil,
		"GOODBYE cruel world!",
	},
	// {
	// 	"#if - if with function shows the contents when function returns true",
	// 	"{{#if goodbye}}GOODBYE {{/if}}cruel {{world}}!",
	// 	map[string]interface{}{
	// 		"goodbye": func() bool { return true },
	// 		"world":   "world",
	// 	},
	// 	nil,
	// 	nil,
	//  nil,
	// 	"GOODBYE cruel world!",
	// },
	{
		"#if - if with function shows the contents when function returns string",
		"{{#if goodbye}}GOODBYE {{/if}}cruel {{world}}!",
		map[string]interface{}{
			"goodbye": func(h *HelperArg) string { return "world" },
			"world":   "world",
		},
		nil, nil, nil,
		"GOODBYE cruel world!",
	},

	// {
	// 	"#if - if with function does not show the contents when returns false",
	// 	"{{#if goodbye}}GOODBYE {{/if}}cruel {{world}}!",
	// 	map[string]interface{}{},
	// 	nil,
	// 	nil,
	//  nil,
	// 	"cruel world!",
	// },
	// {
	// 	"#if - if with function does not show the contents when returns undefined",
	// 	"{{#if goodbye}}GOODBYE {{/if}}cruel {{world}}!",
	// 	map[string]interface{}{},
	// 	nil,
	// 	nil,
	//  nil,
	// 	"cruel world!",
	// },

	{
		"#with",
		"{{#with person}}{{first}} {{last}}{{/with}}",
		map[string]interface{}{"person": map[string]string{"first": "Alan", "last": "Johnson"}},
		nil, nil, nil,
		"Alan Johnson",
	},
	// {
	// 	"#with - with with function argument",
	// 	"{{#with person}}{{first}} {{last}}{{/with}}",
	// 	map[string]interface{}{},
	// 	nil,
	// 	nil,
	//  nil,
	// 	"Alan Johnson",
	// },
	{
		"#with - with with else",
		"{{#with person}}Person is present{{else}}Person is not present{{/with}}",
		map[string]interface{}{},
		nil, nil, nil,
		"Person is not present",
	},

	{
		"#each - each with array argument iterates over the contents when not empty",
		"{{#each goodbyes}}{{text}}! {{/each}}cruel {{world}}!",
		map[string]interface{}{"goodbyes": []map[string]string{{"text": "goodbye"}, {"text": "Goodbye"}, {"text": "GOODBYE"}}, "world": "world"},
		nil, nil, nil,
		"goodbye! Goodbye! GOODBYE! cruel world!",
	},
	{
		"#each - each with array argument ignores the contents when empty",
		"{{#each goodbyes}}{{text}}! {{/each}}cruel {{world}}!",
		map[string]interface{}{"goodbyes": []map[string]string{}, "world": "world"},
		nil, nil, nil,
		"cruel world!",
	},
	{
		"#each - each without data (1)",
		"{{#each goodbyes}}{{text}}! {{/each}}cruel {{world}}!",
		map[string]interface{}{"goodbyes": []map[string]string{{"text": "goodbye"}, {"text": "Goodbye"}, {"text": "GOODBYE"}}, "world": "world"},
		nil, nil, nil,
		"goodbye! Goodbye! GOODBYE! cruel world!",
	},
	{
		"#each - each without data (2)",
		"{{#each .}}{{.}}{{/each}}",
		map[string]interface{}{"goodbyes": "cruel", "world": "world"},
		nil, nil, nil,
		// note: a go hash is not ordered, so result may vary, this behaviour differs from the JS implementation
		[]string{"cruelworld", "worldcruel"},
	},
	{
		"#each - each without context",
		"{{#each goodbyes}}{{text}}! {{/each}}cruel {{world}}!",
		nil, nil, nil, nil,
		"cruel !",
	},

	// NOTE: we test with a map instead of an object
	{
		"#each - each with an object and @key (map)",
		"{{#each goodbyes}}{{@key}}. {{text}}! {{/each}}cruel {{world}}!",
		map[string]interface{}{"goodbyes": map[interface{}]map[string]string{"<b>#1</b>": {"text": "goodbye"}, 2: {"text": "GOODBYE"}}, "world": "world"},
		nil, nil, nil,
		[]string{"&lt;b&gt;#1&lt;/b&gt;. goodbye! 2. GOODBYE! cruel world!", "2. GOODBYE! &lt;b&gt;#1&lt;/b&gt;. goodbye! cruel world!"},
	},
	// NOTE: An additional test with a struct, but without an html stuff for the key, because it is impossible
	{
		"#each - each with an object and @key (struct)",
		"{{#each goodbyes}}{{@key}}. {{text}}! {{/each}}cruel {{world}}!",
		map[string]interface{}{
			"goodbyes": struct {
				Foo map[string]string
				Bar map[string]int
			}{map[string]string{"text": "baz"}, map[string]int{"text": 10}},
			"world": "world",
		},
		nil, nil, nil,
		[]string{"Foo. baz! Bar. 10! cruel world!", "Bar. 10! Foo. baz! cruel world!"},
	},
	{
		"#each - each with @index",
		"{{#each goodbyes}}{{@index}}. {{text}}! {{/each}}cruel {{world}}!",
		map[string]interface{}{"goodbyes": []map[string]string{{"text": "goodbye"}, {"text": "Goodbye"}, {"text": "GOODBYE"}}, "world": "world"},
		nil, nil, nil,
		"0. goodbye! 1. Goodbye! 2. GOODBYE! cruel world!",
	},
	{
		"#each - each with nested @index",
		"{{#each goodbyes}}{{@index}}. {{text}}! {{#each ../goodbyes}}{{@index}} {{/each}}After {{@index}} {{/each}}{{@index}}cruel {{world}}!",
		map[string]interface{}{"goodbyes": []map[string]string{{"text": "goodbye"}, {"text": "Goodbye"}, {"text": "GOODBYE"}}, "world": "world"},
		nil, nil, nil,
		"0. goodbye! 0 1 2 After 0 1. Goodbye! 0 1 2 After 1 2. GOODBYE! 0 1 2 After 2 cruel world!",
	},
	{
		"#each - each with block params",
		"{{#each goodbyes as |value index|}}{{index}}. {{value.text}}! {{#each ../goodbyes as |childValue childIndex|}} {{index}} {{childIndex}}{{/each}} After {{index}} {{/each}}{{index}}cruel {{world}}!",
		map[string]interface{}{"goodbyes": []map[string]string{{"text": "goodbye"}, {"text": "Goodbye"}}, "world": "world"},
		nil, nil, nil,
		"0. goodbye!  0 0 0 1 After 0 1. Goodbye!  1 0 1 1 After 1 cruel world!",
	},
	// @note: That test differs from JS impl because maps and structs are not ordered in go
	{
		"#each - each object with @index",
		"{{#each goodbyes}}{{@index}}. {{text}}! {{/each}}cruel {{world}}!",
		map[string]interface{}{"goodbyes": map[string]map[string]string{"a": {"text": "goodbye"}, "b": {"text": "Goodbye"}}, "world": "world"},
		nil, nil, nil,
		[]string{"0. goodbye! 1. Goodbye! cruel world!", "0. Goodbye! 1. goodbye! cruel world!"},
	},
	{
		"#each - each with nested @first",
		"{{#each goodbyes}}({{#if @first}}{{text}}! {{/if}}{{#each ../goodbyes}}{{#if @first}}{{text}}!{{/if}}{{/each}}{{#if @first}} {{text}}!{{/if}}) {{/each}}cruel {{world}}!",
		map[string]interface{}{"goodbyes": []map[string]string{{"text": "goodbye"}, {"text": "Goodbye"}, {"text": "GOODBYE"}}, "world": "world"},
		nil, nil, nil,
		"(goodbye! goodbye! goodbye!) (goodbye!) (goodbye!) cruel world!",
	},
	// @note: That test differs from JS impl because maps and structs are not ordered in go
	{
		"#each - each object with @first",
		"{{#each goodbyes}}{{#if @first}}{{text}}! {{/if}}{{/each}}cruel {{world}}!",
		map[string]interface{}{"goodbyes": map[string]map[string]string{"foo": {"text": "goodbye"}, "bar": {"text": "Goodbye"}}, "world": "world"},
		nil, nil, nil,
		[]string{"goodbye! cruel world!", "Goodbye! cruel world!"},
	},
	{
		"#each - each with @last",
		"{{#each goodbyes}}{{#if @last}}{{text}}! {{/if}}{{/each}}cruel {{world}}!",
		map[string]interface{}{"goodbyes": []map[string]string{{"text": "goodbye"}, {"text": "Goodbye"}, {"text": "GOODBYE"}}, "world": "world"},
		nil, nil, nil,
		"GOODBYE! cruel world!",
	},
	// @note: That test differs from JS impl because maps and structs are not ordered in go
	{
		"#each - each object with @last",
		"{{#each goodbyes}}{{#if @last}}{{text}}! {{/if}}{{/each}}cruel {{world}}!",
		map[string]interface{}{"goodbyes": map[string]map[string]string{"foo": {"text": "goodbye"}, "bar": {"text": "Goodbye"}}, "world": "world"},
		nil, nil, nil,
		[]string{"goodbye! cruel world!", "Goodbye! cruel world!"},
	},
	{
		"#each - each with nested @last",
		"{{#each goodbyes}}({{#if @last}}{{text}}! {{/if}}{{#each ../goodbyes}}{{#if @last}}{{text}}!{{/if}}{{/each}}{{#if @last}} {{text}}!{{/if}}) {{/each}}cruel {{world}}!",
		map[string]interface{}{"goodbyes": []map[string]string{{"text": "goodbye"}, {"text": "Goodbye"}, {"text": "GOODBYE"}}, "world": "world"},
		nil, nil, nil,
		"(GOODBYE!) (GOODBYE!) (GOODBYE! GOODBYE! GOODBYE!) cruel world!",
	},

	{
		"#each - each with function argument (1)",
		"{{#each goodbyes}}{{text}}! {{/each}}cruel {{world}}!",
		map[string]interface{}{"goodbyes": func() []map[string]string {
			return []map[string]string{{"text": "goodbye"}, {"text": "Goodbye"}, {"text": "GOODBYE"}}
		}, "world": "world"},
		nil, nil, nil,
		"goodbye! Goodbye! GOODBYE! cruel world!",
	},
	{
		"#each - each with function argument (2)",
		"{{#each goodbyes}}{{text}}! {{/each}}cruel {{world}}!",
		map[string]interface{}{"goodbyes": []map[string]string{}, "world": "world"},
		nil, nil, nil,
		"cruel world!",
	},
	// {
	// 	"#each - data passed to helpers",
	// 	"{{#each letters}}{{this}}{{detectDataInsideEach}}{{/each}}",
	// 	map[string][]string{"letters": {"a", "b", "c"}},
	// 	map[string]interface{}{"exclaim": "!"},
	// 	map[string]Helper{"detectDataInsideEach": detectDataHelper},
	// 	nil,
	// 	"a!b!c!",
	// },

	// @todo Add remaining tests
}

func TestHandlebarsBuiltins(t *testing.T) {
	launchHandlebarsTests(t, hbBuiltinsTests)
}

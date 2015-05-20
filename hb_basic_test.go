package raymond

import "testing"

//
// @todo Adds tests from:
//   https://github.com/wycats/handlebars.js/blob/master/spec/data.js
//   https://github.com/wycats/handlebars.js/blob/master/spec/partials.js
//   https://github.com/wycats/handlebars.js/blob/master/spec/regressions.js
//   https://github.com/wycats/handlebars.js/blob/master/spec/whitespace-control.js
//
//   https://github.com/wycats/handlebars.js/blob/master/spec/mustache/
//

//
// Those tests come from:
//   https://github.com/wycats/handlebars.js/blob/master/spec/basic.js
//
var hbBasicTests = []raymondTest{
	{
		"most basic",
		"{{foo}}",
		map[string]string{"foo": "foo"},
		nil,
		"foo",
	},
	{
		"escaping (1)",
		"\\{{foo}}",
		map[string]string{"foo": "food"},
		nil,
		"{{foo}}",
	},
	{
		"escaping (2)",
		"content \\{{foo}}",
		map[string]string{},
		nil,
		"content {{foo}}",
	},
	{
		"escaping (3)",
		"\\\\{{foo}}",
		map[string]string{"foo": "food"},
		nil,
		"\\food",
	},
	{
		"escaping (4)",
		"content \\\\{{foo}}",
		map[string]string{"foo": "food"},
		nil,
		"content \\food",
	},
	{
		"escaping (5)",
		"\\\\ {{foo}}",
		map[string]string{"foo": "food"},
		nil,
		"\\\\ food",
	},
	{
		"compiling with a basic context",
		"Goodbye\n{{cruel}}\n{{world}}!",
		map[string]string{"cruel": "cruel", "world": "world"},
		nil,
		"Goodbye\ncruel\nworld!",
	},
	{
		"compiling with an undefined context (1)",
		"Goodbye\n{{cruel}}\n{{world.bar}}!",
		nil,
		nil,
		"Goodbye\n\n!",
	},
	{
		"compiling with an undefined context (2)",
		"{{#unless foo}}Goodbye{{../test}}{{test2}}{{/unless}}",
		nil,
		nil,
		"Goodbye",
	},
	{
		"comments (1)",
		"{{! Goodbye}}Goodbye\n{{cruel}}\n{{world}}!",
		map[string]string{"cruel": "cruel", "world": "world"},
		nil,
		"Goodbye\ncruel\nworld!",
	},
	// {"comments (2)", "    {{~! comment ~}}      blah", nil, "blah"},
	// {"comments (3)", "    {{~!-- long-comment --~}}      blah", nil, "blah"},
	// {"comments (4)", "    {{! comment ~}}      blah", nil, "    blah"},
	// {"comments (5)", "    {{!-- long-comment --~}}      blah", nil, "    blah"},
	// {"comments (6)", "    {{~! comment}}      blah", nil, "      blah"},
	// {"comments (7)", "    {{~!-- long-comment --}}      blah", nil, "      blah"},
	{
		"boolean (1)",
		"{{#goodbye}}GOODBYE {{/goodbye}}cruel {{world}}!",
		map[string]interface{}{"goodbye": true, "world": "world"},
		nil,
		"GOODBYE cruel world!",
	},
	{
		"boolean (2)",
		"{{#goodbye}}GOODBYE {{/goodbye}}cruel {{world}}!",
		map[string]interface{}{"goodbye": false, "world": "world"},
		nil,
		"cruel world!",
	},
	{
		"zeros (1)",
		"num1: {{num1}}, num2: {{num2}}",
		map[string]interface{}{"num1": 42, "num2": 0},
		nil,
		"num1: 42, num2: 0",
	},
	{
		"zeros (2)",
		"num: {{.}}",
		0,
		nil,
		"num: 0",
	},
	{
		"zeros (3)",
		"num: {{num1/num2}}",
		map[string]map[string]interface{}{"num1": {"num2": 0}},
		nil,
		"num: 0",
	},
	{
		"false (1)",
		"val1: {{val1}}, val2: {{val2}}",
		map[string]interface{}{"val1": false, "val2": false},
		nil,
		"val1: false, val2: false",
	},
	{
		"false (2)",
		"val: {{.}}",
		false,
		nil,
		"val: false",
	},
	{
		"false (3)",
		"val: {{val1/val2}}",
		map[string]map[string]interface{}{"val1": {"val2": false}},
		nil,
		"val: false",
	},
	{
		"false (4)",
		"val1: {{{val1}}}, val2: {{{val2}}}",
		map[string]interface{}{"val1": false, "val2": false},
		nil,
		"val1: false, val2: false",
	},
	{
		"false (5)",
		"val: {{{val1/val2}}}",
		map[string]map[string]interface{}{"val1": {"val2": false}},
		nil,
		"val: false",
	},
	{
		"newlines (1)",
		"Alan's\nTest",
		nil,
		nil,
		"Alan's\nTest",
	},
	{
		"newlines (2)",
		"Alan's\rTest",
		nil,
		nil,
		"Alan's\rTest",
	},
	{
		"escaping text (1)",
		"Awesome's",
		map[string]string{},
		nil,
		"Awesome's",
	},
	{
		"escaping text (2)",
		"Awesome\\",
		map[string]string{},
		nil,
		"Awesome\\",
	},
	{
		"escaping text (3)",
		"Awesome\\\\ foo",
		map[string]string{},
		nil,
		"Awesome\\\\ foo",
	},
	{
		"escaping text (4)",
		"Awesome {{foo}}",
		map[string]string{"foo": "\\"},
		nil,
		"Awesome \\",
	},
	{
		"escaping text (5)",
		" ' ' ",
		map[string]string{},
		nil,
		" ' ' ",
	},
	{
		"escaping expressions (6)",
		"{{{awesome}}}",
		map[string]string{"awesome": "&'\\<>"},
		nil,
		"&'\\<>",
	},
	{
		"escaping expressions (7)",
		"{{&awesome}}",
		map[string]string{"awesome": "&'\\<>"},
		nil,
		"&'\\<>",
	},
	{
		"escaping expressions (8)",
		"{{awesome}}",
		map[string]string{"awesome": "&\"'`\\<>"},
		nil,
		"&amp;&#34;&#39;`\\&lt;&gt;",
	},
	{
		"escaping expressions (9)",
		"{{awesome}}",
		map[string]string{"awesome": "Escaped, <b> looks like: &lt;b&gt;"},
		nil,
		"Escaped, &lt;b&gt; looks like: &amp;lt;b&amp;gt;",
	},

	// @todo "functions returning safestrings shouldn't be escaped"

	{
		"functions (1)",
		"{{awesome}}",
		map[string]interface{}{"awesome": func() string { return "Awesome" }},
		nil,
		"Awesome",
	},
	{
		"functions (2)",
		"{{awesome}}",
		map[string]interface{}{"awesome": func(h *HelperArg) string { return h.DataStr("more") }, "more": "More awesome"},
		nil,
		"More awesome",
	},
	{
		"functions with context argument",
		"{{awesome frank}}",
		map[string]interface{}{"awesome": func(h *HelperArg) string { return h.ParamStr(0) }, "frank": "Frank"},
		nil,
		"Frank",
	},

	// @todo "functions with context argument"
	// @todo "pathed functions with context argument"
	// @todo "depthed functions with context argument"
	// @todo "block functions with context argument"
	// @todo "depthed block functions with context argument"
	// @todo "block functions without context argument"
	// @todo "pathed block functions without context argument"
	// @todo "depthed block functions without context argument"

	{
		"paths with hyphens (1)",
		"{{foo-bar}}",
		map[string]string{"foo-bar": "baz"},
		nil,
		"baz",
	},
	{
		"paths with hyphens (2)",
		"{{foo.foo-bar}}",
		map[string]map[string]string{"foo": {"foo-bar": "baz"}},
		nil,
		"baz",
	},
	{
		"paths with hyphens (3)",
		"{{foo/foo-bar}}",
		map[string]map[string]string{"foo": {"foo-bar": "baz"}},
		nil,
		"baz",
	},
	{
		"nested paths",
		"Goodbye {{alan/expression}} world!",
		map[string]map[string]string{"alan": {"expression": "beautiful"}},
		nil,
		"Goodbye beautiful world!",
	},
	{
		"nested paths with empty string value",
		"Goodbye {{alan/expression}} world!",
		map[string]map[string]string{"alan": {"expression": ""}},
		nil,
		"Goodbye  world!",
	},
	{
		"literal paths (1)",
		"Goodbye {{[@alan]/expression}} world!",
		map[string]map[string]string{"@alan": {"expression": "beautiful"}},
		nil,
		"Goodbye beautiful world!",
	},
	{
		"literal paths (2)",
		"Goodbye {{[foo bar]/expression}} world!",
		map[string]map[string]string{"foo bar": {"expression": "beautiful"}},
		nil,
		"Goodbye beautiful world!",
	},
	{
		"literal references",
		"Goodbye {{[foo bar]}} world!",
		map[string]string{"foo bar": "beautiful"},
		nil,
		"Goodbye beautiful world!",
	},

	// @todo "that current context path ({{.}}) doesn't hit helpers"

	{
		"complex but empty paths (1)",
		"{{person/name}}",
		map[string]map[string]interface{}{"person": {"name": nil}},
		nil,
		"",
	},
	{
		"complex but empty paths (2)",
		"{{person/name}}",
		map[string]map[string]string{"person": {}},
		nil,
		"",
	},
	{
		"this keyword in paths (1)",
		"{{#goodbyes}}{{this}}{{/goodbyes}}",
		map[string]interface{}{"goodbyes": []string{"goodbye", "Goodbye", "GOODBYE"}},
		nil,
		"goodbyeGoodbyeGOODBYE",
	},
	{
		"this keyword in paths (2)",
		"{{#hellos}}{{this/text}}{{/hellos}}",
		map[string]interface{}{"hellos": []interface{}{
			map[string]string{"text": "hello"},
			map[string]string{"text": "Hello"},
			map[string]string{"text": "HELLO"},
		}},
		nil,
		"helloHelloHELLO",
	},

	// @todo "{{#hellos}}{{text/this/foo}}{{/hellos}}" should throw error 'Invalid path: text/this'

	{
		"this keyword nested inside path' (1)",
		"{{[this]}}",
		map[string]string{"this": "bar"},
		nil,
		"bar",
	},
	{
		"this keyword nested inside path' (2)",
		"{{text/[this]}}",
		map[string]map[string]string{"text": {"this": "bar"}},
		nil,
		"bar",
	},
	{
		"this keyword in helpers (1)",
		"{{#goodbyes}}{{foo this}}{{/goodbyes}}",
		map[string]interface{}{"goodbyes": []string{"goodbye", "Goodbye", "GOODBYE"}},
		map[string]Helper{"foo": barSuffixHelper},
		"bar goodbyebar Goodbyebar GOODBYE",
	},
	{
		"this keyword in helpers (2)",
		"{{#hellos}}{{foo this/text}}{{/hellos}}",
		map[string]interface{}{"hellos": []map[string]string{{"text": "hello"}, {"text": "Hello"}, {"text": "HELLO"}}},
		map[string]Helper{"foo": barSuffixHelper},
		"bar hellobar Hellobar HELLO",
	},

	// @todo "{{#hellos}}{{foo text/this/foo}}{{/hellos}}" should throw error 'Invalid path: text/this'

	{
		"this keyword nested inside helpers param (1)",
		"{{foo [this]}}",
		map[string]interface{}{"this": "bar"},
		map[string]Helper{"foo": echoHelper},
		"bar",
	},
	{
		"this keyword nested inside helpers param (2)",
		"{{foo text/[this]}}",
		map[string]map[string]string{"text": {"this": "bar"}},
		map[string]Helper{"foo": echoHelper},
		"bar",
	},
	{
		"pass string literals (1)",
		`{{"foo"}}`,
		map[string]string{},
		nil,
		"",
	},
	{
		"pass string literals (2)",
		`{{"foo"}}`,
		map[string]string{"foo": "bar"},
		nil,
		"bar",
	},
	{
		"pass string literals (3)",
		`{{#"foo"}}{{.}}{{/"foo"}}`,
		map[string]interface{}{"foo": []string{"bar", "baz"}},
		nil,
		"barbaz",
	},
	{
		"pass number literals (1)",
		"{{12}}",
		map[string]string{},
		nil,
		"",
	},
	{
		"pass number literals (2)",
		"{{12}}",
		map[string]string{"12": "bar"},
		nil,
		"bar",
	},
	{
		"pass number literals (3)",
		"{{12.34}}",
		map[string]string{},
		nil,
		"",
	},
	{
		"pass number literals (4)",
		"{{12.34}}",
		map[string]string{"12.34": "bar"},
		nil,
		"bar",
	},

	// @todo {"pass number literals (5)", "{{12.34 1}}", ...function..., "bar1"},

	{
		"pass boolean literals (1)",
		"{{true}}",
		map[string]string{},
		nil,
		"",
	},
	{
		"pass boolean literals (2)",
		"{{true}}",
		map[string]string{"": "foo"},
		nil,
		"",
	},
	{
		"pass boolean literals (3)",
		"{{false}}",
		map[string]string{"false": "foo"},
		nil,
		"foo",
	},

	// @todo
	// {
	//  "should handle literals in subexpression",
	//  "{{foo (false)}}",
	//  ...,
	//  ...,
	//  "bar",
	// },
}

func TestHandlebarsBasic(t *testing.T) {
	launchRaymondTests(t, hbBasicTests)
}

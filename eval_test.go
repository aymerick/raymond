package raymond

import (
	"bytes"
	"testing"
)

type evalTest struct {
	name   string
	input  string
	data   interface{}
	output string
}

var evalTests = []evalTest{
	{
		"only content",
		"this is content",
		nil,
		"this is content",
	},
	// @todo Test with a struct for data

	//
	// Next tests come from:
	//   https://github.com/wycats/handlebars.js/blob/master/spec/basic.js
	//

	{
		"most basic",
		"{{foo}}",
		map[string]string{"foo": "foo"},
		"foo",
	},
	{
		"escaping (1)",
		"\\{{foo}}",
		map[string]string{"foo": "food"},
		"{{foo}}",
	},
	{
		"escaping (2)",
		"content \\{{foo}}",
		map[string]string{},
		"content {{foo}}",
	},
	{
		"escaping (3)",
		"\\\\{{foo}}",
		map[string]string{"foo": "food"},
		"\\food",
	},
	{
		"escaping (4)",
		"content \\\\{{foo}}",
		map[string]string{"foo": "food"},
		"content \\food",
	},
	{
		"escaping (5)",
		"\\\\ {{foo}}",
		map[string]string{"foo": "food"},
		"\\\\ food",
	},
	{
		"compiling with a basic context",
		"Goodbye\n{{cruel}}\n{{world}}!",
		map[string]string{"cruel": "cruel", "world": "world"},
		"Goodbye\ncruel\nworld!",
	},
	{
		"compiling with an undefined context (1)",
		"Goodbye\n{{cruel}}\n{{world.bar}}!",
		nil, "Goodbye\n\n!",
	},
	{
		"compiling with an undefined context (2)",
		"{{#unless foo}}Goodbye{{../test}}{{test2}}{{/unless}}",
		nil, "Goodbye",
	},
	{
		"comments (1)",
		"{{! Goodbye}}Goodbye\n{{cruel}}\n{{world}}!",
		map[string]string{"cruel": "cruel", "world": "world"},
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
		"GOODBYE cruel world!",
	},
	{
		"boolean (2)",
		"{{#goodbye}}GOODBYE {{/goodbye}}cruel {{world}}!",
		map[string]interface{}{"goodbye": false, "world": "world"},
		"cruel world!",
	},
	{
		"zeros (1)",
		"num1: {{num1}}, num2: {{num2}}",
		map[string]interface{}{"num1": 42, "num2": 0},
		"num1: 42, num2: 0",
	},
	{
		"zeros (2)",
		"num: {{.}}",
		0,
		"num: 0",
	},
	{
		"zeros (3)",
		"num: {{num1/num2}}",
		map[string]map[string]interface{}{"num1": {"num2": 0}},
		"num: 0",
	},
	{
		"false (1)",
		"val1: {{val1}}, val2: {{val2}}",
		map[string]interface{}{"val1": false, "val2": false},
		"val1: false, val2: false",
	},
	{
		"false (2)",
		"val: {{.}}",
		false, "val: false",
	},
	{
		"false (3)",
		"val: {{val1/val2}}",
		map[string]map[string]interface{}{"val1": {"val2": false}},
		"val: false",
	},
	{
		"false (4)",
		"val1: {{{val1}}}, val2: {{{val2}}}",
		map[string]interface{}{"val1": false, "val2": false},
		"val1: false, val2: false",
	},
	{
		"false (5)",
		"val: {{{val1/val2}}}",
		map[string]map[string]interface{}{"val1": {"val2": false}},
		"val: false",
	},
	{
		"newlines (1)",
		"Alan's\nTest",
		nil,
		"Alan's\nTest",
	},
	{
		"newlines (2)",
		"Alan's\rTest",
		nil,
		"Alan's\rTest",
	},
	{
		"escaping text (1)",
		"Awesome's",
		map[string]string{},
		"Awesome's",
	},
	{
		"escaping text (2)",
		"Awesome\\",
		map[string]string{},
		"Awesome\\",
	},
	{
		"escaping text (3)",
		"Awesome\\\\ foo",
		map[string]string{},
		"Awesome\\\\ foo",
	},
	{
		"escaping text (4)",
		"Awesome {{foo}}",
		map[string]string{"foo": "\\"},
		"Awesome \\",
	},
	{
		"escaping text (5)",
		" ' ' ",
		map[string]string{},
		" ' ' ",
	},
	{
		"escaping expressions (6)",
		"{{{awesome}}}",
		map[string]string{"awesome": "&'\\<>"},
		"&'\\<>",
	},
	{
		"escaping expressions (7)",
		"{{&awesome}}",
		map[string]string{"awesome": "&'\\<>"},
		"&'\\<>",
	},
	{
		"escaping expressions (8)",
		"{{awesome}}",
		map[string]string{"awesome": "&\"'`\\<>"},
		"&amp;&#34;&#39;`\\&lt;&gt;",
	},
	{
		"escaping expressions (9)",
		"{{awesome}}",
		map[string]string{"awesome": "Escaped, <b> looks like: &lt;b&gt;"},
		"Escaped, &lt;b&gt; looks like: &amp;lt;b&amp;gt;",
	},

	// @todo "functions returning safestrings shouldn't be escaped"
	// @todo "functions"
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
		"baz",
	},
	{
		"paths with hyphens (2)",
		"{{foo.foo-bar}}",
		map[string]map[string]string{"foo": {"foo-bar": "baz"}},
		"baz",
	},
	{
		"paths with hyphens (3)",
		"{{foo/foo-bar}}",
		map[string]map[string]string{"foo": {"foo-bar": "baz"}},
		"baz",
	},
	{
		"nested paths",
		"Goodbye {{alan/expression}} world!",
		map[string]map[string]string{"alan": {"expression": "beautiful"}},
		"Goodbye beautiful world!",
	},
	{
		"nested paths with empty string value",
		"Goodbye {{alan/expression}} world!",
		map[string]map[string]string{"alan": {"expression": ""}},
		"Goodbye  world!",
	},
	{
		"literal paths (1)",
		"Goodbye {{[@alan]/expression}} world!",
		map[string]map[string]string{"@alan": {"expression": "beautiful"}},
		"Goodbye beautiful world!",
	},
	{
		"literal paths (2)",
		"Goodbye {{[foo bar]/expression}} world!",
		map[string]map[string]string{"foo bar": {"expression": "beautiful"}},
		"Goodbye beautiful world!",
	},
	{
		"literal references",
		"Goodbye {{[foo bar]}} world!",
		map[string]string{"foo bar": "beautiful"},
		"Goodbye beautiful world!",
	},

	// @todo "that current context path ({{.}}) doesn't hit helpers"

	{
		"complex but empty paths (1)",
		"{{person/name}}",
		map[string]map[string]interface{}{"person": {"name": nil}},
		"",
	},
	{
		"complex but empty paths (2)",
		"{{person/name}}",
		map[string]map[string]string{"person": {}},
		"",
	},

	// {"this keyword in paths (1)", "{{#goodbyes}}{{this}}{{/goodbyes}}", map[string]interface{}{"goodbyes": []string{"goodbye", "Goodbye", "GOODBYE"}}, "goodbyeGoodbyeGOODBYE"},
	// {"this keyword in paths (2)", "{{#hellos}}{{this/text}}{{/hellos}}", map[string]interface{}{"hellos": []interface{}{map[string]string{"text": "hello"}, map[string]string{"text": "Hello"}, map[string]string{"text": "HELLO"}}}, "helloHelloHELLO"},

	// @todo "{{#hellos}}{{text/this/foo}}{{/hellos}}" should throw error 'Invalid path: text/this'

	// {"this keyword nested inside path' (1)", "{{[this]}}", map[string]string{"this": "bar"}, "bar"},
	// {"this keyword nested inside path' (2)", "{{text/[this]}}", map[string]map[string]string{"text": {"this": "bar"}}, "bar"},

	// @todo {"this keyword in helpers (1)", "{{#goodbyes}}{{foo this}}{{/goodbyes}}", ..., "bar goodbyebar Goodbyebar GOODBYE"},
	// @todo {"this keyword in helpers (2)", "{{#hellos}}{{foo this/text}}{{/hellos}}", ..., "bar hellobar Hellobar HELLO', 'This keyword evaluates in more complex paths"},

	// @todo "this keyword nested inside helpers param"

	{
		"pass string literals (1)",
		`{{"foo"}}`,
		map[string]string{},
		"",
	},
	{
		"pass string literals (2)",
		`{{"foo"}}`,
		map[string]string{"foo": "bar"},
		"bar",
	},
	{
		"pass string literals (3)",
		`{{#"foo"}}{{.}}{{/"foo"}}`,
		map[string]interface{}{"foo": []string{"bar", "baz"}},
		"barbaz",
	},
	{
		"pass number literals (1)",
		"{{12}}",
		map[string]string{},
		"",
	},
	{
		"pass number literals (2)",
		"{{12}}",
		map[string]string{"12": "bar"},
		"bar",
	},
	{
		"pass number literals (3)",
		"{{12.34}}",
		map[string]string{},
		"",
	},
	{
		"pass number literals (4)",
		"{{12.34}}",
		map[string]string{"12.34": "bar"},
		"bar",
	},

	// @todo {"pass number literals (5)", "{{12.34 1}}", ...function..., "bar1"},

	{
		"pass boolean literals (1)",
		"{{true}}",
		map[string]string{},
		"",
	},
	{
		"pass boolean literals (2)",
		"{{true}}",
		map[string]string{"": "foo"},
		"",
	},
	{
		"pass boolean literals (3)",
		"{{false}}",
		map[string]string{"false": "foo"},
		"foo",
	},

	// @todo "should handle literals in subexpression"

	//
	// @todo Adds tests from:
	//   https://github.com/wycats/handlebars.js/blob/master/spec/blocks.js
	//   https://github.com/wycats/handlebars.js/blob/master/spec/builtin.js
	//   https://github.com/wycats/handlebars.js/blob/master/spec/data.js
	//   https://github.com/wycats/handlebars.js/blob/master/spec/partials.js
	//   https://github.com/wycats/handlebars.js/blob/master/spec/regressions.js
	//   https://github.com/wycats/handlebars.js/blob/master/spec/strict.js
	//   https://github.com/wycats/handlebars.js/blob/master/spec/string-params.js
	//   https://github.com/wycats/handlebars.js/blob/master/spec/subexpression.js
	//   https://github.com/wycats/handlebars.js/blob/master/spec/whitespace-control.js
	//
	//   https://github.com/wycats/handlebars.js/blob/master/spec/mustache/
	//
}

func TestEval(t *testing.T) {
	for _, test := range evalTests {
		var err error
		var tpl *Template

		buf := new(bytes.Buffer)

		// parse template
		tpl, err = Parse(test.input)
		if err != nil {
			t.Errorf("Test '%s' failed - Failed to parse template\ninput:\n\t'%s'\nerror:\n\t%s", test.name, test.input, err)
		} else {
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

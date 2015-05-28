package raymond

import (
	"fmt"
	"regexp"
	"testing"
)

//
// @todo Adds tests from:
//   https://github.com/wycats/handlebars.js/blob/master/spec/data.js
//   https://github.com/wycats/handlebars.js/blob/master/spec/regressions.js
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
		nil,
		"foo",
	},
	{
		"escaping (1)",
		"\\{{foo}}",
		map[string]string{"foo": "food"},
		nil,
		nil,
		"{{foo}}",
	},
	{
		"escaping (2)",
		"content \\{{foo}}",
		map[string]string{},
		nil,
		nil,
		"content {{foo}}",
	},
	{
		"escaping (3)",
		"\\\\{{foo}}",
		map[string]string{"foo": "food"},
		nil,
		nil,
		"\\food",
	},
	{
		"escaping (4)",
		"content \\\\{{foo}}",
		map[string]string{"foo": "food"},
		nil,
		nil,
		"content \\food",
	},
	{
		"escaping (5)",
		"\\\\ {{foo}}",
		map[string]string{"foo": "food"},
		nil,
		nil,
		"\\\\ food",
	},
	{
		"compiling with a basic context",
		"Goodbye\n{{cruel}}\n{{world}}!",
		map[string]string{"cruel": "cruel", "world": "world"},
		nil,
		nil,
		"Goodbye\ncruel\nworld!",
	},
	{
		"compiling with an undefined context (1)",
		"Goodbye\n{{cruel}}\n{{world.bar}}!",
		nil,
		nil,
		nil,
		"Goodbye\n\n!",
	},
	{
		"compiling with an undefined context (2)",
		"{{#unless foo}}Goodbye{{../test}}{{test2}}{{/unless}}",
		nil,
		nil,
		nil,
		"Goodbye",
	},
	{
		"comments (1)",
		"{{! Goodbye}}Goodbye\n{{cruel}}\n{{world}}!",
		map[string]string{"cruel": "cruel", "world": "world"},
		nil,
		nil,
		"Goodbye\ncruel\nworld!",
	},
	{
		"comments (2)",
		"    {{~! comment ~}}      blah",
		nil,
		nil,
		nil,
		"blah",
	},
	{
		"comments (3)",
		"    {{~!-- long-comment --~}}      blah",
		nil,
		nil,
		nil,
		"blah",
	},
	{
		"comments (4)",
		"    {{! comment ~}}      blah",
		nil,
		nil,
		nil,
		"    blah",
	},
	{
		"comments (5)",
		"    {{!-- long-comment --~}}      blah",
		nil,
		nil,
		nil,
		"    blah",
	},
	{
		"comments (6)",
		"    {{~! comment}}      blah",
		nil,
		nil,
		nil,
		"      blah",
	},
	{
		"comments (7)",
		"    {{~!-- long-comment --}}      blah",
		nil,
		nil,
		nil,
		"      blah",
	},
	{
		"boolean (1)",
		"{{#goodbye}}GOODBYE {{/goodbye}}cruel {{world}}!",
		map[string]interface{}{"goodbye": true, "world": "world"},
		nil,
		nil,
		"GOODBYE cruel world!",
	},
	{
		"boolean (2)",
		"{{#goodbye}}GOODBYE {{/goodbye}}cruel {{world}}!",
		map[string]interface{}{"goodbye": false, "world": "world"},
		nil,
		nil,
		"cruel world!",
	},
	{
		"zeros (1)",
		"num1: {{num1}}, num2: {{num2}}",
		map[string]interface{}{"num1": 42, "num2": 0},
		nil,
		nil,
		"num1: 42, num2: 0",
	},
	{
		"zeros (2)",
		"num: {{.}}",
		0,
		nil,
		nil,
		"num: 0",
	},
	{
		"zeros (3)",
		"num: {{num1/num2}}",
		map[string]map[string]interface{}{"num1": {"num2": 0}},
		nil,
		nil,
		"num: 0",
	},
	{
		"false (1)",
		"val1: {{val1}}, val2: {{val2}}",
		map[string]interface{}{"val1": false, "val2": false},
		nil,
		nil,
		"val1: false, val2: false",
	},
	{
		"false (2)",
		"val: {{.}}",
		false,
		nil,
		nil,
		"val: false",
	},
	{
		"false (3)",
		"val: {{val1/val2}}",
		map[string]map[string]interface{}{"val1": {"val2": false}},
		nil,
		nil,
		"val: false",
	},
	{
		"false (4)",
		"val1: {{{val1}}}, val2: {{{val2}}}",
		map[string]interface{}{"val1": false, "val2": false},
		nil,
		nil,
		"val1: false, val2: false",
	},
	{
		"false (5)",
		"val: {{{val1/val2}}}",
		map[string]map[string]interface{}{"val1": {"val2": false}},
		nil,
		nil,
		"val: false",
	},
	{
		"newlines (1)",
		"Alan's\nTest",
		nil,
		nil,
		nil,
		"Alan's\nTest",
	},
	{
		"newlines (2)",
		"Alan's\rTest",
		nil,
		nil,
		nil,
		"Alan's\rTest",
	},
	{
		"escaping text (1)",
		"Awesome's",
		map[string]string{},
		nil,
		nil,
		"Awesome's",
	},
	{
		"escaping text (2)",
		"Awesome\\",
		map[string]string{},
		nil,
		nil,
		"Awesome\\",
	},
	{
		"escaping text (3)",
		"Awesome\\\\ foo",
		map[string]string{},
		nil,
		nil,
		"Awesome\\\\ foo",
	},
	{
		"escaping text (4)",
		"Awesome {{foo}}",
		map[string]string{"foo": "\\"},
		nil,
		nil,
		"Awesome \\",
	},
	{
		"escaping text (5)",
		" ' ' ",
		map[string]string{},
		nil,
		nil,
		" ' ' ",
	},
	{
		"escaping expressions (6)",
		"{{{awesome}}}",
		map[string]string{"awesome": "&'\\<>"},
		nil,
		nil,
		"&'\\<>",
	},
	{
		"escaping expressions (7)",
		"{{&awesome}}",
		map[string]string{"awesome": "&'\\<>"},
		nil,
		nil,
		"&'\\<>",
	},
	{
		"escaping expressions (8)",
		"{{awesome}}",
		map[string]string{"awesome": "&\"'`\\<>"},
		nil,
		nil,
		"&amp;&quot;&apos;`\\&lt;&gt;",
	},
	{
		"escaping expressions (9)",
		"{{awesome}}",
		map[string]string{"awesome": "Escaped, <b> looks like: &lt;b&gt;"},
		nil,
		nil,
		"Escaped, &lt;b&gt; looks like: &amp;lt;b&amp;gt;",
	},

	// @todo "functions returning safestrings shouldn't be escaped"

	{
		"functions (1)",
		"{{awesome}}",
		map[string]interface{}{"awesome": func() string { return "Awesome" }},
		nil,
		nil,
		"Awesome",
	},
	{
		"functions (2)",
		"{{awesome}}",
		map[string]interface{}{"awesome": func(h *HelperArg) string {
			return h.DataStr("more")
		}, "more": "More awesome"},
		nil,
		nil,
		"More awesome",
	},
	{
		"functions with context argument",
		"{{awesome frank}}",
		map[string]interface{}{"awesome": func(h *HelperArg) string {
			return h.ParamStr(0)
		}, "frank": "Frank"},
		nil,
		nil,
		"Frank",
	},
	{
		"pathed functions with context argument",
		"{{bar.awesome frank}}",
		map[string]interface{}{"bar": map[string]interface{}{"awesome": func(h *HelperArg) string {
			return h.ParamStr(0)
		}}, "frank": "Frank"},
		nil,
		nil,
		"Frank",
	},
	{
		"depthed functions with context argument",
		"{{#with frank}}{{../awesome .}}{{/with}}",
		map[string]interface{}{"awesome": func(h *HelperArg) string {
			return h.ParamStr(0)
		}, "frank": "Frank"},
		nil,
		nil,
		"Frank",
	},
	{
		"block functions with context argument",
		"{{#awesome 1}}inner {{.}}{{/awesome}}",
		map[string]interface{}{"awesome": func(h *HelperArg) string {
			return h.BlockWith(h.Param(0))
		}},
		nil,
		nil,
		"inner 1",
	},
	{
		"depthed block functions with context argument",
		"{{#with value}}{{#../awesome 1}}inner {{.}}{{/../awesome}}{{/with}}",
		map[string]interface{}{
			"awesome": func(h *HelperArg) string {
				return h.BlockWith(h.Param(0))
			},
			"value": true,
		},
		nil,
		nil,
		"inner 1",
	},
	{
		"block functions without context argument",
		"{{#awesome}}inner{{/awesome}}",
		map[string]interface{}{
			"awesome": func(h *HelperArg) string {
				return h.Block()
			},
		},
		nil,
		nil,
		"inner",
	},
	// @note I don't even understand how this test passes with the JS implementation
	// {
	// 	"pathed block functions without context argument",
	// 	"{{#foo.awesome}}inner{{/foo.awesome}}",
	// 	map[string]map[string]interface{}{
	// 		"foo": {
	// 			"awesome": func(h *HelperArg) string {
	// 				return h.Data()
	// 			},
	// 		},
	// 	},
	// 	nil,
	// 	nil,
	// 	"inner",
	// },
	// @note I don't even understand how this test passes with the JS implementation
	// {
	// 	"depthed block functions without context argument",
	// 	"{{#with value}}{{#../awesome}}inner{{/../awesome}}{{/with}}",
	// 	map[string]interface{}{
	// 		"value": true,
	// 		"awesome": func(h *HelperArg) string {
	// 			return h.Data()
	// 		},
	// 	},
	// 	nil,
	// 	nil,
	// 	"inner",
	// },
	{
		"paths with hyphens (1)",
		"{{foo-bar}}",
		map[string]string{"foo-bar": "baz"},
		nil,
		nil,
		"baz",
	},
	{
		"paths with hyphens (2)",
		"{{foo.foo-bar}}",
		map[string]map[string]string{"foo": {"foo-bar": "baz"}},
		nil,
		nil,
		"baz",
	},
	{
		"paths with hyphens (3)",
		"{{foo/foo-bar}}",
		map[string]map[string]string{"foo": {"foo-bar": "baz"}},
		nil,
		nil,
		"baz",
	},
	{
		"nested paths",
		"Goodbye {{alan/expression}} world!",
		map[string]map[string]string{"alan": {"expression": "beautiful"}},
		nil,
		nil,
		"Goodbye beautiful world!",
	},
	{
		"nested paths with empty string value",
		"Goodbye {{alan/expression}} world!",
		map[string]map[string]string{"alan": {"expression": ""}},
		nil,
		nil,
		"Goodbye  world!",
	},
	{
		"literal paths (1)",
		"Goodbye {{[@alan]/expression}} world!",
		map[string]map[string]string{"@alan": {"expression": "beautiful"}},
		nil,
		nil,
		"Goodbye beautiful world!",
	},
	{
		"literal paths (2)",
		"Goodbye {{[foo bar]/expression}} world!",
		map[string]map[string]string{"foo bar": {"expression": "beautiful"}},
		nil,
		nil,
		"Goodbye beautiful world!",
	},
	{
		"literal references",
		"Goodbye {{[foo bar]}} world!",
		map[string]string{"foo bar": "beautiful"},
		nil,
		nil,
		"Goodbye beautiful world!",
	},
	// {
	// 	"that current context path ({{.}}) doesn't hit helpers",
	// 	"test: {{.}}",
	// 	map[string]string{"helper": "awesome"},
	// 	nil,
	// 	nil,
	// 	"test: ",
	// },
	{
		"complex but empty paths (1)",
		"{{person/name}}",
		map[string]map[string]interface{}{"person": {"name": nil}},
		nil,
		nil,
		"",
	},
	{
		"complex but empty paths (2)",
		"{{person/name}}",
		map[string]map[string]string{"person": {}},
		nil,
		nil,
		"",
	},
	{
		"this keyword in paths (1)",
		"{{#goodbyes}}{{this}}{{/goodbyes}}",
		map[string]interface{}{"goodbyes": []string{"goodbye", "Goodbye", "GOODBYE"}},
		nil,
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
		nil,
		"helloHelloHELLO",
	},
	{
		"this keyword nested inside path' (1)",
		"{{[this]}}",
		map[string]string{"this": "bar"},
		nil,
		nil,
		"bar",
	},
	{
		"this keyword nested inside path' (2)",
		"{{text/[this]}}",
		map[string]map[string]string{"text": {"this": "bar"}},
		nil,
		nil,
		"bar",
	},
	{
		"this keyword in helpers (1)",
		"{{#goodbyes}}{{foo this}}{{/goodbyes}}",
		map[string]interface{}{"goodbyes": []string{"goodbye", "Goodbye", "GOODBYE"}},
		map[string]Helper{"foo": barSuffixHelper},
		nil,
		"bar goodbyebar Goodbyebar GOODBYE",
	},
	{
		"this keyword in helpers (2)",
		"{{#hellos}}{{foo this/text}}{{/hellos}}",
		map[string]interface{}{"hellos": []map[string]string{{"text": "hello"}, {"text": "Hello"}, {"text": "HELLO"}}},
		map[string]Helper{"foo": barSuffixHelper},
		nil,
		"bar hellobar Hellobar HELLO",
	},
	{
		"this keyword nested inside helpers param (1)",
		"{{foo [this]}}",
		map[string]interface{}{"this": "bar"},
		map[string]Helper{"foo": echoHelper},
		nil,
		"bar",
	},
	{
		"this keyword nested inside helpers param (2)",
		"{{foo text/[this]}}",
		map[string]map[string]string{"text": {"this": "bar"}},
		map[string]Helper{"foo": echoHelper},
		nil,
		"bar",
	},
	{
		"pass string literals (1)",
		`{{"foo"}}`,
		map[string]string{},
		nil,
		nil,
		"",
	},
	{
		"pass string literals (2)",
		`{{"foo"}}`,
		map[string]string{"foo": "bar"},
		nil,
		nil,
		"bar",
	},
	{
		"pass string literals (3)",
		`{{#"foo"}}{{.}}{{/"foo"}}`,
		map[string]interface{}{"foo": []string{"bar", "baz"}},
		nil,
		nil,
		"barbaz",
	},
	{
		"pass number literals (1)",
		"{{12}}",
		map[string]string{},
		nil,
		nil,
		"",
	},
	{
		"pass number literals (2)",
		"{{12}}",
		map[string]string{"12": "bar"},
		nil,
		nil,
		"bar",
	},
	{
		"pass number literals (3)",
		"{{12.34}}",
		map[string]string{},
		nil,
		nil,
		"",
	},
	{
		"pass number literals (4)",
		"{{12.34}}",
		map[string]string{"12.34": "bar"},
		nil,
		nil,
		"bar",
	},
	{
		"pass number literals (5)",
		"{{12.34 1}}",
		map[string]interface{}{"12.34": func(h *HelperArg) string {
			return "bar" + h.ParamStr(0)
		}},
		nil,
		nil,
		"bar1",
	},
	{
		"pass boolean literals (1)",
		"{{true}}",
		map[string]string{},
		nil,
		nil,
		"",
	},
	{
		"pass boolean literals (2)",
		"{{true}}",
		map[string]string{"": "foo"},
		nil,
		nil,
		"",
	},
	{
		"pass boolean literals (3)",
		"{{false}}",
		map[string]string{"false": "foo"},
		nil,
		nil,
		"foo",
	},
	{
		"should handle literals in subexpression",
		"{{foo (false)}}",
		map[string]interface{}{"false": func() string { return "bar" }},
		map[string]Helper{"foo": func(h *HelperArg) string {
			return h.ParamStr(0)
		}},
		nil,
		"bar",
	},
}

func TestHandlebarsBasic(t *testing.T) {
	launchHandlebarsTests(t, hbBasicTests)
}

func TestHandlebarsBasicErrors(t *testing.T) {
	var err error

	inputs := []string{
		// this keyword nested inside path
		"{{#hellos}}{{text/this/foo}}{{/hellos}}",
		// this keyword nested inside helpers param
		"{{#hellos}}{{foo text/this/foo}}{{/hellos}}",
	}

	expectedError := regexp.QuoteMeta("Invalid path: text/this")

	stats.handlebarsTests(len(inputs))

	for _, input := range inputs {
		_, err = Parse(input)
		if err == nil {
			t.Errorf("Test failed - Error expected")
			stats.failed()
		}

		match, errMatch := regexp.MatchString(expectedError, fmt.Sprint(err))
		if errMatch != nil {
			panic("Failed to match regexp")
		}

		if !match {
			t.Errorf("Test failed - Expected error:\n\t%s\n\nGot:\n\t%s", expectedError, err)
			stats.failed()
		}
	}

	stats.output()
}

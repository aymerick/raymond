package handlebars

import (
	"fmt"
	"reflect"
	"strings"
	"testing"

	"github.com/aymerick/raymond"
)

//
// Helpers
//

func barSuffixHelper(h *raymond.HelperArg) interface{} {
	str, _ := h.Param(0).(string)
	return "bar " + str
}

func echoHelper(h *raymond.HelperArg) interface{} {
	str, _ := h.Param(0).(string)
	nb, ok := h.Param(1).(int)
	if !ok {
		nb = 1
	}

	result := ""
	for i := 0; i < nb; i++ {
		result += str
	}

	return result
}

func linkHelper(h *raymond.HelperArg) interface{} {
	prefix, _ := h.Param(0).(string)

	return fmt.Sprintf(`<a href="%s/%s">%s</a>`, prefix, h.FieldStr("url"), h.FieldStr("text"))
}

func rawHelper(h *raymond.HelperArg) interface{} {
	result := h.Block()

	for _, param := range h.Params() {
		result += raymond.Str(param)
	}

	return result
}

func formHelper(h *raymond.HelperArg) interface{} {
	return "<form>" + h.Block() + "</form>"
}

func formCtxHelper(h *raymond.HelperArg) interface{} {
	return "<form>" + h.BlockWithCtx(h.Param(0)) + "</form>"
}

func listHelper(h *raymond.HelperArg) interface{} {
	ctx := h.Param(0)

	val := reflect.ValueOf(ctx)
	switch val.Kind() {
	case reflect.Array, reflect.Slice:
		if val.Len() > 0 {
			result := "<ul>"
			for i := 0; i < val.Len(); i++ {
				result += "<li>"
				result += h.BlockWithCtx(val.Index(i).Interface())
				result += "</li>"
			}
			result += "</ul>"

			return result
		}
	}

	return "<p>" + h.Inverse() + "</p>"
}

func blogHelper(h *raymond.HelperArg) interface{} {
	return "val is " + h.ParamStr(0)
}

func equalHelper(h *raymond.HelperArg) interface{} {
	return raymond.Str(h.ParamStr(0) == h.ParamStr(1))
}

func dashHelper(h *raymond.HelperArg) interface{} {
	return h.ParamStr(0) + "-" + h.ParamStr(1)
}

func concatHelper(h *raymond.HelperArg) interface{} {
	return h.ParamStr(0) + h.ParamStr(1)
}

func detectDataHelper(h *raymond.HelperArg) interface{} {
	if val, ok := h.DataFrame().Get("exclaim").(string); ok {
		return val
	}

	return ""
}

//
// Those tests come from:
//   https://github.com/wycats/handlebars.js/blob/master/spec/helper.js
//
var helpersTests = []Test{
	{
		"helper with complex lookup",
		"{{#goodbyes}}{{{link ../prefix}}}{{/goodbyes}}",
		map[string]interface{}{"prefix": "/root", "goodbyes": []map[string]string{{"text": "Goodbye", "url": "goodbye"}}},
		nil,
		map[string]raymond.Helper{"link": linkHelper},
		nil,
		`<a href="/root/goodbye">Goodbye</a>`,
	},
	{
		"helper for raw block gets raw content",
		"{{{{raw}}}} {{test}} {{{{/raw}}}}",
		map[string]interface{}{"test": "hello"},
		nil,
		map[string]raymond.Helper{"raw": rawHelper},
		nil,
		" {{test}} ",
	},
	{
		"helper for raw block gets parameters",
		"{{{{raw 1 2 3}}}} {{test}} {{{{/raw}}}}",
		map[string]interface{}{"test": "hello"},
		nil,
		map[string]raymond.Helper{"raw": rawHelper},
		nil,
		" {{test}} 123",
	},
	{
		"helper block with complex lookup expression",
		"{{#goodbyes}}{{../name}}{{/goodbyes}}",
		map[string]interface{}{"name": "Alan"},
		nil,
		map[string]raymond.Helper{"goodbyes": func(h *raymond.HelperArg) interface{} {
			out := ""
			for _, str := range []string{"Goodbye", "goodbye", "GOODBYE"} {
				out += str + " " + h.BlockWithCtx(str) + "! "
			}
			return out
		}},
		nil,
		"Goodbye Alan! goodbye Alan! GOODBYE Alan! ",
	},
	{
		"helper with complex lookup and nested template",
		"{{#goodbyes}}{{#link ../prefix}}{{text}}{{/link}}{{/goodbyes}}",
		map[string]interface{}{"prefix": "/root", "goodbyes": []map[string]string{{"text": "Goodbye", "url": "goodbye"}}},
		nil,
		map[string]raymond.Helper{"link": linkHelper},
		nil,
		`<a href="/root/goodbye">Goodbye</a>`,
	},
	{
		// note: The JS implementation returns undefined, we return empty string
		"helper returning undefined value (1)",
		" {{nothere}}",
		map[string]interface{}{},
		nil,
		map[string]raymond.Helper{"nothere": func(h *raymond.HelperArg) interface{} {
			return ""
		}},
		nil,
		" ",
	},
	{
		// note: The JS implementation returns undefined, we return empty string
		"helper returning undefined value (2)",
		" {{#nothere}}{{/nothere}}",
		map[string]interface{}{},
		nil,
		map[string]raymond.Helper{"nothere": func(h *raymond.HelperArg) interface{} {
			return ""
		}},
		nil,
		" ",
	},
	{
		"block helper",
		"{{#goodbyes}}{{text}}! {{/goodbyes}}cruel {{world}}!",
		map[string]interface{}{"world": "world"},
		nil,
		map[string]raymond.Helper{"goodbyes": func(h *raymond.HelperArg) interface{} {
			return h.BlockWithCtx(map[string]string{"text": "GOODBYE"})
		}},
		nil,
		"GOODBYE! cruel world!",
	},
	{
		"block helper staying in the same context",
		"{{#form}}<p>{{name}}</p>{{/form}}",
		map[string]interface{}{"name": "Yehuda"},
		nil,
		map[string]raymond.Helper{"form": formHelper},
		nil,
		"<form><p>Yehuda</p></form>",
	},
	{
		"block helper should have context in this",
		"<ul>{{#people}}<li>{{#link}}{{name}}{{/link}}</li>{{/people}}</ul>",
		map[string]interface{}{"people": []map[string]interface{}{{"name": "Alan", "id": 1}, {"name": "Yehuda", "id": 2}}},
		nil,
		map[string]raymond.Helper{"link": func(h *raymond.HelperArg) interface{} {
			return fmt.Sprintf("<a href=\"/people/%s\">%s</a>", h.FieldStr("id"), h.Block())
		}},
		nil,
		`<ul><li><a href="/people/1">Alan</a></li><li><a href="/people/2">Yehuda</a></li></ul>`,
	},
	{
		"block helper for undefined value",
		"{{#empty}}shouldn't render{{/empty}}",
		nil, nil, nil, nil,
		"",
	},
	{
		"block helper passing a new context",
		"{{#form yehuda}}<p>{{name}}</p>{{/form}}",
		map[string]map[string]string{"yehuda": {"name": "Yehuda"}},
		nil,
		map[string]raymond.Helper{"form": formCtxHelper},
		nil,
		"<form><p>Yehuda</p></form>",
	},
	{
		"block helper passing a complex path context",
		"{{#form yehuda/cat}}<p>{{name}}</p>{{/form}}",
		map[string]map[string]interface{}{"yehuda": {"name": "Yehuda", "cat": map[string]string{"name": "Harold"}}},
		nil,
		map[string]raymond.Helper{"form": formCtxHelper},
		nil,
		"<form><p>Harold</p></form>",
	},
	{
		"nested block helpers",
		"{{#form yehuda}}<p>{{name}}</p>{{#link}}Hello{{/link}}{{/form}}",
		map[string]map[string]string{"yehuda": {"name": "Yehuda"}},
		nil,
		map[string]raymond.Helper{"link": func(h *raymond.HelperArg) interface{} {
			return fmt.Sprintf("<a href=\"%s\">%s</a>", h.FieldStr("name"), h.Block())
		}, "form": formCtxHelper},
		nil,
		`<form><p>Yehuda</p><a href="Yehuda">Hello</a></form>`,
	},
	{
		"block helper inverted sections (1) - an inverse wrapper is passed in as a new context",
		"{{#list people}}{{name}}{{^}}<em>Nobody's here</em>{{/list}}",
		map[string][]map[string]string{"people": {{"name": "Alan"}, {"name": "Yehuda"}}},
		nil,
		map[string]raymond.Helper{"list": listHelper},
		nil,
		`<ul><li>Alan</li><li>Yehuda</li></ul>`,
	},
	{
		"block helper inverted sections (2) - an inverse wrapper can be optionally called",
		"{{#list people}}{{name}}{{^}}<em>Nobody's here</em>{{/list}}",
		map[string][]map[string]string{"people": {}},
		nil,
		map[string]raymond.Helper{"list": listHelper},
		nil,
		`<p><em>Nobody's here</em></p>`,
	},
	{
		"block helper inverted sections (3) - the context of an inverse is the parent of the block",
		"{{#list people}}Hello{{^}}{{message}}{{/list}}",
		map[string]interface{}{"people": []interface{}{}, "message": "Nobody's here"},
		nil,
		map[string]raymond.Helper{"list": listHelper},
		nil,
		`<p>Nobody&apos;s here</p>`,
	},

	{
		"pathed lambdas with parameters (1)",
		"{{./helper 1}}",
		map[string]interface{}{
			"helper": func() string { return "winning" },
			"hash": map[string]interface{}{
				"helper": func() string { return "winning" },
			}},
		nil,
		map[string]raymond.Helper{"./helper": func(h *raymond.HelperArg) interface{} { return "fail" }},
		nil,
		"winning",
	},
	{
		"pathed lambdas with parameters (2)",
		"{{hash/helper 1}}",
		map[string]interface{}{
			"helper": func() string { return "winning" },
			"hash": map[string]interface{}{
				"helper": func() string { return "winning" },
			}},
		nil,
		map[string]raymond.Helper{"./helper": func(h *raymond.HelperArg) interface{} { return "fail" }},
		nil,
		"winning",
	},

	{
		"helpers hash - providing a helpers hash (1)",
		"Goodbye {{cruel}} {{world}}!",
		map[string]interface{}{"cruel": "cruel"},
		nil,
		map[string]raymond.Helper{"world": func(h *raymond.HelperArg) interface{} { return "world" }},
		nil,
		"Goodbye cruel world!",
	},
	{
		"helpers hash - providing a helpers hash (2)",
		"Goodbye {{#iter}}{{cruel}} {{world}}{{/iter}}!",
		map[string]interface{}{"iter": []map[string]string{{"cruel": "cruel"}}},
		nil,
		map[string]raymond.Helper{"world": func(h *raymond.HelperArg) interface{} { return "world" }},
		nil,
		"Goodbye cruel world!",
	},
	{
		"helpers hash - in cases of conflict, helpers win (1)",
		"{{{lookup}}}",
		map[string]interface{}{"lookup": "Explicit"},
		nil,
		map[string]raymond.Helper{"lookup": func(h *raymond.HelperArg) interface{} { return "helpers" }},
		nil,
		"helpers",
	},
	{
		"helpers hash - in cases of conflict, helpers win (2)",
		"{{lookup}}",
		map[string]interface{}{"lookup": "Explicit"},
		nil,
		map[string]raymond.Helper{"lookup": func(h *raymond.HelperArg) interface{} { return "helpers" }},
		nil,
		"helpers",
	},
	{
		"helpers hash - the helpers hash is available is nested contexts",
		"{{#outer}}{{#inner}}{{helper}}{{/inner}}{{/outer}}",
		map[string]interface{}{"outer": map[string]interface{}{"inner": map[string]interface{}{"unused": []string{}}}},
		nil,
		map[string]raymond.Helper{"helper": func(h *raymond.HelperArg) interface{} { return "helper" }},
		nil,
		"helper",
	},

	// @todo "helpers hash - the helper hash should augment the global hash"

	// @todo "registration"

	{
		"decimal number literals work",
		"Message: {{hello -1.2 1.2}}",
		nil, nil,
		map[string]raymond.Helper{"hello": func(h *raymond.HelperArg) interface{} {
			ts, t2s := "NaN", "NaN"

			if v, ok := h.Param(0).(float64); ok {
				ts = raymond.Str(v)
			}

			if v, ok := h.Param(1).(float64); ok {
				t2s = raymond.Str(v)
			}

			return "Hello " + ts + " " + t2s + " times"
		}},
		nil,
		"Message: Hello -1.2 1.2 times",
	},
	{
		"negative number literals work",
		"Message: {{hello -12}}",
		nil, nil,
		map[string]raymond.Helper{"hello": func(h *raymond.HelperArg) interface{} {
			times := "NaN"

			if v, ok := h.Param(0).(int); ok {
				times = raymond.Str(v)
			}

			return "Hello " + times + " times"
		}},
		nil,
		"Message: Hello -12 times",
	},

	{
		"String literal parameters - simple literals work",
		`Message: {{hello "world" 12 true false}}`,
		nil, nil,
		map[string]raymond.Helper{"hello": func(h *raymond.HelperArg) interface{} {
			times, bool1, bool2 := "NaN", "NaB", "NaB"

			param, ok := h.Param(0).(string)
			if !ok {
				param = "NaN"
			}

			if v, ok := h.Param(1).(int); ok {
				times = raymond.Str(v)
			}

			if v, ok := h.Param(2).(bool); ok {
				bool1 = raymond.Str(v)
			}

			if v, ok := h.Param(3).(bool); ok {
				bool2 = raymond.Str(v)
			}

			return "Hello " + param + " " + times + " times: " + bool1 + " " + bool2
		}},
		nil,
		"Message: Hello world 12 times: true false",
	},

	// @todo "using a quote in the middle of a parameter raises an error"

	{
		"String literal parameters - escaping a String is possible",
		"Message: {{{hello \"\\\"world\\\"\"}}}",
		nil, nil,
		map[string]raymond.Helper{"hello": func(h *raymond.HelperArg) interface{} {
			return "Hello " + h.ParamStr(0)
		}},
		nil,
		`Message: Hello "world"`,
	},
	{
		"String literal parameters - it works with ' marks",
		"Message: {{{hello \"Alan's world\"}}}",
		nil, nil,
		map[string]raymond.Helper{"hello": func(h *raymond.HelperArg) interface{} {
			return "Hello " + h.ParamStr(0)
		}},
		nil,
		`Message: Hello Alan's world`,
	},

	{
		"multiple parameters - simple multi-params work",
		"Message: {{goodbye cruel world}}",
		map[string]string{"cruel": "cruel", "world": "world"},
		nil,
		map[string]raymond.Helper{"goodbye": func(h *raymond.HelperArg) interface{} {
			return "Goodbye " + h.ParamStr(0) + " " + h.ParamStr(1)
		}},
		nil,
		"Message: Goodbye cruel world",
	},
	{
		"multiple parameters - block multi-params work",
		"Message: {{#goodbye cruel world}}{{greeting}} {{adj}} {{noun}}{{/goodbye}}",
		map[string]string{"cruel": "cruel", "world": "world"},
		nil,
		map[string]raymond.Helper{"goodbye": func(h *raymond.HelperArg) interface{} {
			return h.BlockWithCtx(map[string]interface{}{"greeting": "Goodbye", "adj": h.Param(0), "noun": h.Param(1)})
		}},
		nil,
		"Message: Goodbye cruel world",
	},

	{
		"hash - helpers can take an optional hash",
		`{{goodbye cruel="CRUEL" world="WORLD" times=12}}`,
		nil, nil,
		map[string]raymond.Helper{"goodbye": func(h *raymond.HelperArg) interface{} {
			return "GOODBYE " + h.HashStr("cruel") + " " + h.HashStr("world") + " " + h.HashStr("times") + " TIMES"
		}},
		nil,
		"GOODBYE CRUEL WORLD 12 TIMES",
	},
	{
		"hash - helpers can take an optional hash with booleans (1)",
		`{{goodbye cruel="CRUEL" world="WORLD" print=true}}`,
		nil, nil,
		map[string]raymond.Helper{"goodbye": func(h *raymond.HelperArg) interface{} {
			p, ok := h.HashProp("print").(bool)
			if ok {
				if p {
					return "GOODBYE " + h.HashStr("cruel") + " " + h.HashStr("world")
				} else {
					return "NOT PRINTING"
				}
			}

			return "THIS SHOULD NOT HAPPEN"
		}},
		nil,
		"GOODBYE CRUEL WORLD",
	},
	{
		"hash - helpers can take an optional hash with booleans (2)",
		`{{goodbye cruel="CRUEL" world="WORLD" print=false}}`,
		nil, nil,
		map[string]raymond.Helper{"goodbye": func(h *raymond.HelperArg) interface{} {
			p, ok := h.HashProp("print").(bool)
			if ok {
				if p {
					return "GOODBYE " + h.HashStr("cruel") + " " + h.HashStr("world")
				} else {
					return "NOT PRINTING"
				}
			}

			return "THIS SHOULD NOT HAPPEN"
		}},
		nil,
		"NOT PRINTING",
	},
	{
		"block helpers can take an optional hash",
		`{{#goodbye cruel="CRUEL" times=12}}world{{/goodbye}}`,
		nil, nil,
		map[string]raymond.Helper{"goodbye": func(h *raymond.HelperArg) interface{} {
			return "GOODBYE " + h.HashStr("cruel") + " " + h.Block() + " " + h.HashStr("times") + " TIMES"
		}},
		nil,
		"GOODBYE CRUEL world 12 TIMES",
	},
	{
		"block helpers can take an optional hash with single quoted stings",
		`{{#goodbye cruel='CRUEL' times=12}}world{{/goodbye}}`,
		nil, nil,
		map[string]raymond.Helper{"goodbye": func(h *raymond.HelperArg) interface{} {
			return "GOODBYE " + h.HashStr("cruel") + " " + h.Block() + " " + h.HashStr("times") + " TIMES"
		}},
		nil,
		"GOODBYE CRUEL world 12 TIMES",
	},
	{
		"block helpers can take an optional hash with booleans (1)",
		`{{#goodbye cruel="CRUEL" print=true}}world{{/goodbye}}`,
		nil, nil,
		map[string]raymond.Helper{"goodbye": func(h *raymond.HelperArg) interface{} {
			p, ok := h.HashProp("print").(bool)
			if ok {
				if p {
					return "GOODBYE " + h.HashStr("cruel") + " " + h.Block()
				} else {
					return "NOT PRINTING"
				}
			}

			return "THIS SHOULD NOT HAPPEN"
		}},
		nil,
		"GOODBYE CRUEL world",
	},
	{
		"block helpers can take an optional hash with booleans (1)",
		`{{#goodbye cruel="CRUEL" print=false}}world{{/goodbye}}`,
		nil, nil,
		map[string]raymond.Helper{"goodbye": func(h *raymond.HelperArg) interface{} {
			p, ok := h.HashProp("print").(bool)
			if ok {
				if p {
					return "GOODBYE " + h.HashStr("cruel") + " " + h.Block()
				} else {
					return "NOT PRINTING"
				}
			}

			return "THIS SHOULD NOT HAPPEN"
		}},
		nil,
		"NOT PRINTING",
	},

	// @todo "helperMissing - if a context is not found, helperMissing is used" throw error

	// @todo "helperMissing - if a context is not found, custom helperMissing is used"

	// @todo "helperMissing - if a value is not found, custom helperMissing is used"

	{
		"block helpers can take an optional hash with booleans (1)",
		`{{#goodbye cruel="CRUEL" print=false}}world{{/goodbye}}`,
		nil, nil,
		map[string]raymond.Helper{"goodbye": func(h *raymond.HelperArg) interface{} {
			p, ok := h.HashProp("print").(bool)
			if ok {
				if p {
					return "GOODBYE " + h.HashStr("cruel") + " " + h.Block()
				} else {
					return "NOT PRINTING"
				}
			}

			return "THIS SHOULD NOT HAPPEN"
		}},
		nil,
		"NOT PRINTING",
	},

	// @todo "knownHelpers/knownHelpersOnly" tests

	// @todo "blockHelperMissing" tests

	// @todo "name field" tests

	{
		"name conflicts - helpers take precedence over same-named context properties",
		`{{goodbye}} {{cruel world}}`,
		map[string]string{"goodbye": "goodbye", "world": "world"},
		nil,
		map[string]raymond.Helper{
			"goodbye": func(h *raymond.HelperArg) interface{} {
				return strings.ToUpper(h.FieldStr("goodbye"))
			},
			"cruel": func(h *raymond.HelperArg) interface{} {
				return "cruel " + strings.ToUpper(h.ParamStr(0))
			},
		},
		nil,
		"GOODBYE cruel WORLD",
	},
	{
		"name conflicts - helpers take precedence over same-named context properties",
		`{{#goodbye}} {{cruel world}}{{/goodbye}}`,
		map[string]string{"goodbye": "goodbye", "world": "world"},
		nil,
		map[string]raymond.Helper{
			"goodbye": func(h *raymond.HelperArg) interface{} {
				return strings.ToUpper(h.FieldStr("goodbye")) + h.Block()
			},
			"cruel": func(h *raymond.HelperArg) interface{} {
				return "cruel " + strings.ToUpper(h.ParamStr(0))
			},
		},
		nil,
		"GOODBYE cruel WORLD",
	},
	{
		"name conflicts - Scoped names take precedence over helpers",
		`{{this.goodbye}} {{cruel world}} {{cruel this.goodbye}}`,
		map[string]string{"goodbye": "goodbye", "world": "world"},
		nil,
		map[string]raymond.Helper{
			"goodbye": func(h *raymond.HelperArg) interface{} {
				return strings.ToUpper(h.FieldStr("goodbye"))
			},
			"cruel": func(h *raymond.HelperArg) interface{} {
				return "cruel " + strings.ToUpper(h.ParamStr(0))
			},
		},
		nil,
		"goodbye cruel WORLD cruel GOODBYE",
	},
	{
		"name conflicts - Scoped names take precedence over block helpers",
		`{{#goodbye}} {{cruel world}}{{/goodbye}} {{this.goodbye}}`,
		map[string]string{"goodbye": "goodbye", "world": "world"},
		nil,
		map[string]raymond.Helper{
			"goodbye": func(h *raymond.HelperArg) interface{} {
				return strings.ToUpper(h.FieldStr("goodbye")) + h.Block()
			},
			"cruel": func(h *raymond.HelperArg) interface{} {
				return "cruel " + strings.ToUpper(h.ParamStr(0))
			},
		},
		nil,
		"GOODBYE cruel WORLD goodbye",
	},

	// @todo "block params" tests
}

func TestHelpers(t *testing.T) {
	launchTests(t, helpersTests)
}

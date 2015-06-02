package raymond

import (
	"fmt"
	"strings"
	"testing"
)

//
// Those tests come from:
//   https://github.com/wycats/handlebars.js/blob/master/spec/helper.js
//
var hbHelpersTests = []raymondTest{
	{
		"helper with complex lookup",
		"{{#goodbyes}}{{{link ../prefix}}}{{/goodbyes}}",
		map[string]interface{}{"prefix": "/root", "goodbyes": []map[string]string{{"text": "Goodbye", "url": "goodbye"}}},
		nil,
		map[string]Helper{"link": linkHelper},
		nil,
		`<a href="/root/goodbye">Goodbye</a>`,
	},
	{
		"helper for raw block gets raw content",
		"{{{{raw}}}} {{test}} {{{{/raw}}}}",
		map[string]interface{}{"test": "hello"},
		nil,
		map[string]Helper{"raw": rawHelper},
		nil,
		" {{test}} ",
	},
	{
		"helper for raw block gets parameters",
		"{{{{raw 1 2 3}}}} {{test}} {{{{/raw}}}}",
		map[string]interface{}{"test": "hello"},
		nil,
		map[string]Helper{"raw": rawHelper},
		nil,
		" {{test}} 123",
	},
	{
		"helper block with complex lookup expression",
		"{{#goodbyes}}{{../name}}{{/goodbyes}}",
		map[string]interface{}{"name": "Alan"},
		nil,
		map[string]Helper{"goodbyes": func(h *HelperArg) string {
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
		map[string]Helper{"link": linkHelper},
		nil,
		`<a href="/root/goodbye">Goodbye</a>`,
	},
	{
		// note: The JS implementation returns undefined, we return empty string
		"helper returning undefined value (1)",
		" {{nothere}}",
		map[string]interface{}{},
		nil,
		map[string]Helper{"nothere": func(h *HelperArg) string {
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
		map[string]Helper{"nothere": func(h *HelperArg) string {
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
		map[string]Helper{"goodbyes": func(h *HelperArg) string {
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
		map[string]Helper{"form": formHelper},
		nil,
		"<form><p>Yehuda</p></form>",
	},
	{
		"block helper should have context in this",
		"<ul>{{#people}}<li>{{#link}}{{name}}{{/link}}</li>{{/people}}</ul>",
		map[string]interface{}{"people": []map[string]interface{}{{"name": "Alan", "id": 1}, {"name": "Yehuda", "id": 2}}},
		nil,
		map[string]Helper{"link": func(h *HelperArg) string {
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
		map[string]Helper{"form": formCtxHelper},
		nil,
		"<form><p>Yehuda</p></form>",
	},
	{
		"block helper passing a complex path context",
		"{{#form yehuda/cat}}<p>{{name}}</p>{{/form}}",
		map[string]map[string]interface{}{"yehuda": {"name": "Yehuda", "cat": map[string]string{"name": "Harold"}}},
		nil,
		map[string]Helper{"form": formCtxHelper},
		nil,
		"<form><p>Harold</p></form>",
	},
	{
		"nested block helpers",
		"{{#form yehuda}}<p>{{name}}</p>{{#link}}Hello{{/link}}{{/form}}",
		map[string]map[string]string{"yehuda": {"name": "Yehuda"}},
		nil,
		map[string]Helper{"link": func(h *HelperArg) string {
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
		map[string]Helper{"list": listHelper},
		nil,
		`<ul><li>Alan</li><li>Yehuda</li></ul>`,
	},
	{
		"block helper inverted sections (2) - an inverse wrapper can be optionally called",
		"{{#list people}}{{name}}{{^}}<em>Nobody's here</em>{{/list}}",
		map[string][]map[string]string{"people": {}},
		nil,
		map[string]Helper{"list": listHelper},
		nil,
		`<p><em>Nobody's here</em></p>`,
	},
	{
		"block helper inverted sections (3) - the context of an inverse is the parent of the block",
		"{{#list people}}Hello{{^}}{{message}}{{/list}}",
		map[string]interface{}{"people": []interface{}{}, "message": "Nobody's here"},
		nil,
		map[string]Helper{"list": listHelper},
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
		map[string]Helper{"./helper": func(h *HelperArg) string { return "fail" }},
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
		map[string]Helper{"./helper": func(h *HelperArg) string { return "fail" }},
		nil,
		"winning",
	},

	{
		"helpers hash - providing a helpers hash (1)",
		"Goodbye {{cruel}} {{world}}!",
		map[string]interface{}{"cruel": "cruel"},
		nil,
		map[string]Helper{"world": func(h *HelperArg) string { return "world" }},
		nil,
		"Goodbye cruel world!",
	},
	{
		"helpers hash - providing a helpers hash (2)",
		"Goodbye {{#iter}}{{cruel}} {{world}}{{/iter}}!",
		map[string]interface{}{"iter": []map[string]string{{"cruel": "cruel"}}},
		nil,
		map[string]Helper{"world": func(h *HelperArg) string { return "world" }},
		nil,
		"Goodbye cruel world!",
	},
	{
		"helpers hash - in cases of conflict, helpers win (1)",
		"{{{lookup}}}",
		map[string]interface{}{"lookup": "Explicit"},
		nil,
		map[string]Helper{"lookup": func(h *HelperArg) string { return "helpers" }},
		nil,
		"helpers",
	},
	{
		"helpers hash - in cases of conflict, helpers win (2)",
		"{{lookup}}",
		map[string]interface{}{"lookup": "Explicit"},
		nil,
		map[string]Helper{"lookup": func(h *HelperArg) string { return "helpers" }},
		nil,
		"helpers",
	},
	{
		"helpers hash - the helpers hash is available is nested contexts",
		"{{#outer}}{{#inner}}{{helper}}{{/inner}}{{/outer}}",
		map[string]interface{}{"outer": map[string]interface{}{"inner": map[string]interface{}{"unused": []string{}}}},
		nil,
		map[string]Helper{"helper": func(h *HelperArg) string { return "helper" }},
		nil,
		"helper",
	},

	// @todo "helpers hash - the helper hash should augment the global hash"

	// @todo "registration"

	{
		"decimal number literals work",
		"Message: {{hello -1.2 1.2}}",
		nil, nil,
		map[string]Helper{"hello": func(h *HelperArg) string {
			ts, t2s := "NaN", "NaN"

			if v, ok := h.Param(0).(float64); ok {
				ts = Str(v)
			}

			if v, ok := h.Param(1).(float64); ok {
				t2s = Str(v)
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
		map[string]Helper{"hello": func(h *HelperArg) string {
			times := "NaN"

			if v, ok := h.Param(0).(int); ok {
				times = Str(v)
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
		map[string]Helper{"hello": func(h *HelperArg) string {
			times, bool1, bool2 := "NaN", "NaB", "NaB"

			param, ok := h.Param(0).(string)
			if !ok {
				param = "NaN"
			}

			if v, ok := h.Param(1).(int); ok {
				times = Str(v)
			}

			if v, ok := h.Param(2).(bool); ok {
				bool1 = Str(v)
			}

			if v, ok := h.Param(3).(bool); ok {
				bool2 = Str(v)
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
		map[string]Helper{"hello": func(h *HelperArg) string {
			return "Hello " + h.ParamStr(0)
		}},
		nil,
		`Message: Hello "world"`,
	},
	{
		"String literal parameters - it works with ' marks",
		"Message: {{{hello \"Alan's world\"}}}",
		nil, nil,
		map[string]Helper{"hello": func(h *HelperArg) string {
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
		map[string]Helper{"goodbye": func(h *HelperArg) string {
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
		map[string]Helper{"goodbye": func(h *HelperArg) string {
			return h.BlockWithCtx(map[string]interface{}{"greeting": "Goodbye", "adj": h.Param(0), "noun": h.Param(1)})
		}},
		nil,
		"Message: Goodbye cruel world",
	},

	{
		"hash - helpers can take an optional hash",
		`{{goodbye cruel="CRUEL" world="WORLD" times=12}}`,
		nil, nil,
		map[string]Helper{"goodbye": func(h *HelperArg) string {
			return "GOODBYE " + h.HashStr("cruel") + " " + h.HashStr("world") + " " + h.HashStr("times") + " TIMES"
		}},
		nil,
		"GOODBYE CRUEL WORLD 12 TIMES",
	},
	{
		"hash - helpers can take an optional hash with booleans (1)",
		`{{goodbye cruel="CRUEL" world="WORLD" print=true}}`,
		nil, nil,
		map[string]Helper{"goodbye": func(h *HelperArg) string {
			p, ok := h.Hash("print").(bool)
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
		map[string]Helper{"goodbye": func(h *HelperArg) string {
			p, ok := h.Hash("print").(bool)
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
		map[string]Helper{"goodbye": func(h *HelperArg) string {
			return "GOODBYE " + h.HashStr("cruel") + " " + h.Block() + " " + h.HashStr("times") + " TIMES"
		}},
		nil,
		"GOODBYE CRUEL world 12 TIMES",
	},
	{
		"block helpers can take an optional hash with single quoted stings",
		`{{#goodbye cruel='CRUEL' times=12}}world{{/goodbye}}`,
		nil, nil,
		map[string]Helper{"goodbye": func(h *HelperArg) string {
			return "GOODBYE " + h.HashStr("cruel") + " " + h.Block() + " " + h.HashStr("times") + " TIMES"
		}},
		nil,
		"GOODBYE CRUEL world 12 TIMES",
	},
	{
		"block helpers can take an optional hash with booleans (1)",
		`{{#goodbye cruel="CRUEL" print=true}}world{{/goodbye}}`,
		nil, nil,
		map[string]Helper{"goodbye": func(h *HelperArg) string {
			p, ok := h.Hash("print").(bool)
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
		map[string]Helper{"goodbye": func(h *HelperArg) string {
			p, ok := h.Hash("print").(bool)
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
		map[string]Helper{"goodbye": func(h *HelperArg) string {
			p, ok := h.Hash("print").(bool)
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
		map[string]Helper{
			"goodbye": func(h *HelperArg) string {
				return strings.ToUpper(h.FieldStr("goodbye"))
			},
			"cruel": func(h *HelperArg) string {
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
		map[string]Helper{
			"goodbye": func(h *HelperArg) string {
				return strings.ToUpper(h.FieldStr("goodbye")) + h.Block()
			},
			"cruel": func(h *HelperArg) string {
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
		map[string]Helper{
			"goodbye": func(h *HelperArg) string {
				return strings.ToUpper(h.FieldStr("goodbye"))
			},
			"cruel": func(h *HelperArg) string {
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
		map[string]Helper{
			"goodbye": func(h *HelperArg) string {
				return strings.ToUpper(h.FieldStr("goodbye")) + h.Block()
			},
			"cruel": func(h *HelperArg) string {
				return "cruel " + strings.ToUpper(h.ParamStr(0))
			},
		},
		nil,
		"GOODBYE cruel WORLD goodbye",
	},

	// @todo "block params" tests
}

func TestHandlebarsHelpers(t *testing.T) {
	launchHandlebarsTests(t, hbHelpersTests)
}

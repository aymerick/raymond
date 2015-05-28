package raymond

import (
	"fmt"
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
		map[string]Helper{"link": linkHelper},
		nil,
		`<a href="/root/goodbye">Goodbye</a>`,
	},
	{
		"helper for raw block gets raw content",
		"{{{{raw}}}} {{test}} {{{{/raw}}}}",
		map[string]interface{}{"test": "hello"},
		map[string]Helper{"raw": rawHelper},
		nil,
		" {{test}} ",
	},
	{
		"helper for raw block gets parameters",
		"{{{{raw 1 2 3}}}} {{test}} {{{{/raw}}}}",
		map[string]interface{}{"test": "hello"},
		map[string]Helper{"raw": rawHelper},
		nil,
		" {{test}} 123",
	},
	{
		"helper block with complex lookup expression",
		"{{#goodbyes}}{{../name}}{{/goodbyes}}",
		map[string]interface{}{"name": "Alan"},
		map[string]Helper{"goodbyes": func(h *HelperArg) string {
			out := ""
			for _, str := range []string{"Goodbye", "goodbye", "GOODBYE"} {
				out += str + " " + h.BlockWith(str) + "! "
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
		map[string]Helper{"link": linkHelper},
		nil,
		`<a href="/root/goodbye">Goodbye</a>`,
	},
	{
		// note: The JS implementation returns undefined, we return empty string
		"helper returning undefined value (1)",
		" {{nothere}}",
		map[string]interface{}{},
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
		map[string]Helper{"goodbyes": func(h *HelperArg) string {
			return h.BlockWith(map[string]string{"text": "GOODBYE"})
		}},
		nil,
		"GOODBYE! cruel world!",
	},
	{
		"block helper staying in the same context",
		"{{#form}}<p>{{name}}</p>{{/form}}",
		map[string]interface{}{"name": "Yehuda"},
		map[string]Helper{"form": formHelper},
		nil,
		"<form><p>Yehuda</p></form>",
	},
	{
		"block helper should have context in this",
		"<ul>{{#people}}<li>{{#link}}{{name}}{{/link}}</li>{{/people}}</ul>",
		map[string]interface{}{"people": []map[string]interface{}{{"name": "Alan", "id": 1}, {"name": "Yehuda", "id": 2}}},
		map[string]Helper{"link": func(h *HelperArg) string {
			return fmt.Sprintf("<a href=\"/people/%s\">%s</a>", h.DataStr("id"), h.Block())
		}},
		nil,
		`<ul><li><a href="/people/1">Alan</a></li><li><a href="/people/2">Yehuda</a></li></ul>`,
	},
	{
		"block helper for undefined value",
		"{{#empty}}shouldn't render{{/empty}}",
		nil,
		nil,
		nil,
		"",
	},
	{
		"block helper passing a new context",
		"{{#form yehuda}}<p>{{name}}</p>{{/form}}",
		map[string]map[string]string{"yehuda": {"name": "Yehuda"}},
		map[string]Helper{"form": formCtxHelper},
		nil,
		"<form><p>Yehuda</p></form>",
	},
	{
		"block helper passing a complex path context",
		"{{#form yehuda/cat}}<p>{{name}}</p>{{/form}}",
		map[string]map[string]interface{}{"yehuda": {"name": "Yehuda", "cat": map[string]string{"name": "Harold"}}},
		map[string]Helper{"form": formCtxHelper},
		nil,
		"<form><p>Harold</p></form>",
	},
	{
		"nested block helpers",
		"{{#form yehuda}}<p>{{name}}</p>{{#link}}Hello{{/link}}{{/form}}",
		map[string]map[string]string{"yehuda": {"name": "Yehuda"}},
		map[string]Helper{"link": func(h *HelperArg) string {
			return fmt.Sprintf("<a href=\"%s\">%s</a>", h.DataStr("name"), h.Block())
		}, "form": formCtxHelper},
		nil,
		`<form><p>Yehuda</p><a href="Yehuda">Hello</a></form>`,
	},
	{
		"block helper inverted sections (1) - an inverse wrapper is passed in as a new context",
		"{{#list people}}{{name}}{{^}}<em>Nobody's here</em>{{/list}}",
		map[string][]map[string]string{"people": {{"name": "Alan"}, {"name": "Yehuda"}}},
		map[string]Helper{"list": listHelper},
		nil,
		`<ul><li>Alan</li><li>Yehuda</li></ul>`,
	},
	{
		"block helper inverted sections (2) - an inverse wrapper can be optionally called",
		"{{#list people}}{{name}}{{^}}<em>Nobody's here</em>{{/list}}",
		map[string][]map[string]string{"people": {}},
		map[string]Helper{"list": listHelper},
		nil,
		`<p><em>Nobody's here</em></p>`,
	},
	{
		"block helper inverted sections (3) - the context of an inverse is the parent of the block",
		"{{#list people}}Hello{{^}}{{message}}{{/list}}",
		map[string]interface{}{"people": []interface{}{}, "message": "Nobody's here"},
		map[string]Helper{"list": listHelper},
		nil,
		`<p>Nobody&apos;s here</p>`,
	},

	// @todo "pathed lambas with parameters"

	// {
	// 	"",
	// 	"",
	// 	map[string]interface{}{},
	// 	nil,
	// 	nil,
	// 	"",
	// },

	// @todo Add remaining tests
}

func TestHandlebarsHelpers(t *testing.T) {
	launchHandlebarsTests(t, hbHelpersTests)
}

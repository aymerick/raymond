package handlebars

import (
	"testing"

	"github.com/aymerick/raymond"
)

//
// Those tests come from:
//   https://github.com/wycats/handlebars.js/blob/master/spec/data.js
//
var dataTests = []Test{
	{
		"passing in data to a compiled function that expects data - works with helpers",
		"{{hello}}",
		map[string]string{"noun": "cat"},
		map[string]interface{}{"adjective": "happy"},
		map[string]raymond.Helper{"hello": func(h *raymond.HelperArg) interface{} {
			return h.DataStr("adjective") + " " + h.FieldStr("noun")
		}},
		nil,
		"happy cat",
	},
	{
		"data can be looked up via @foo",
		"{{@hello}}",
		nil,
		map[string]interface{}{"hello": "hello"},
		nil, nil,
		"hello",
	},
	{
		"deep @foo triggers automatic top-level data",
		`{{#let world="world"}}{{#if foo}}{{#if foo}}Hello {{@world}}{{/if}}{{/if}}{{/let}}`,
		map[string]bool{"foo": true},
		map[string]interface{}{"hello": "hello"},
		map[string]raymond.Helper{"let": func(h *raymond.HelperArg) interface{} {
			frame := h.NewDataFrame()

			for k, v := range h.Hash() {
				frame.Set(k, v)
			}

			return h.BlockWithData(frame)
		}},
		nil,
		"Hello world",
	},
	{
		"parameter data can be looked up via @foo",
		`{{hello @world}}`,
		nil,
		map[string]interface{}{"world": "world"},
		map[string]raymond.Helper{"hello": func(h *raymond.HelperArg) interface{} {
			return "Hello " + h.ParamStr(0)
		}},
		nil,
		"Hello world",
	},
	{
		"hash values can be looked up via @foo",
		`{{hello noun=@world}}`,
		nil,
		map[string]interface{}{"world": "world"},
		map[string]raymond.Helper{"hello": func(h *raymond.HelperArg) interface{} {
			return "Hello " + h.HashStr("noun")
		}},
		nil,
		"Hello world",
	},
	{
		"nested parameter data can be looked up via @foo.bar",
		`{{hello @world.bar}}`,
		nil,
		map[string]interface{}{"world": map[string]string{"bar": "world"}},
		map[string]raymond.Helper{"hello": func(h *raymond.HelperArg) interface{} {
			return "Hello " + h.ParamStr(0)
		}},
		nil,
		"Hello world",
	},
	{
		"nested parameter data does not fail with @world.bar",
		`{{hello @world.bar}}`,
		nil,
		map[string]interface{}{"foo": map[string]string{"bar": "world"}},
		map[string]raymond.Helper{"hello": func(h *raymond.HelperArg) interface{} {
			return "Hello " + h.ParamStr(0)
		}},
		nil,
		// @todo Test differs with JS implementation: we don't output `undefined`
		"Hello ",
	},

	// @todo "parameter data throws when using complex scope references",

	{
		"data can be functions",
		`{{@hello}}`,
		nil,
		map[string]interface{}{"hello": func() string { return "hello" }},
		nil, nil,
		"hello",
	},
	{
		"data can be functions with params",
		`{{@hello "hello"}}`,
		nil,
		map[string]interface{}{"hello": func(h *raymond.HelperArg) string { return h.ParamStr(0) }},
		nil, nil,
		"hello",
	},

	{
		"data is inherited downstream",
		`{{#let foo=1 bar=2}}{{#let foo=bar.baz}}{{@bar}}{{@foo}}{{/let}}{{@foo}}{{/let}}`,
		map[string]map[string]string{"bar": {"baz": "hello world"}},
		nil,
		map[string]raymond.Helper{"let": func(h *raymond.HelperArg) interface{} {
			frame := h.NewDataFrame()

			for k, v := range h.Hash() {
				frame.Set(k, v)
			}

			return h.BlockWithData(frame)
		}},
		nil,
		"2hello world1",
	},
	{
		"passing in data to a compiled function that expects data - works with helpers in partials",
		`{{>myPartial}}`,
		map[string]string{"noun": "cat"},
		map[string]interface{}{"adjective": "happy"},
		map[string]raymond.Helper{"hello": func(h *raymond.HelperArg) interface{} {
			return h.DataStr("adjective") + " " + h.FieldStr("noun")
		}},
		map[string]string{
			"myPartial": "{{hello}}",
		},
		"happy cat",
	},
	{
		"passing in data to a compiled function that expects data - works with helpers and parameters",
		`{{hello world}}`,
		map[string]interface{}{"exclaim": true, "world": "world"},
		map[string]interface{}{"adjective": "happy"},
		map[string]raymond.Helper{"hello": func(h *raymond.HelperArg) interface{} {
			str := "error"
			if b, ok := h.Field("exclaim").(bool); ok {
				if b {
					str = "!"
				} else {
					str = ""
				}
			}

			return h.DataStr("adjective") + " " + h.ParamStr(0) + str
		}},
		nil,
		"happy world!",
	},
	{
		"passing in data to a compiled function that expects data - works with block helpers",
		`{{#hello}}{{world}}{{/hello}}`,
		map[string]bool{"exclaim": true},
		map[string]interface{}{"adjective": "happy"},
		map[string]raymond.Helper{
			"hello": func(h *raymond.HelperArg) interface{} {
				return h.Block()
			},
			"world": func(h *raymond.HelperArg) interface{} {
				str := "error"
				if b, ok := h.Field("exclaim").(bool); ok {
					if b {
						str = "!"
					} else {
						str = ""
					}
				}

				return h.DataStr("adjective") + " world" + str
			},
		},
		nil,
		"happy world!",
	},
	{
		"passing in data to a compiled function that expects data - works with block helpers that use ..",
		`{{#hello}}{{world ../zomg}}{{/hello}}`,
		map[string]interface{}{"exclaim": true, "zomg": "world"},
		map[string]interface{}{"adjective": "happy"},
		map[string]raymond.Helper{
			"hello": func(h *raymond.HelperArg) interface{} {
				return h.BlockWithCtx(map[string]string{"exclaim": "?"})
			},
			"world": func(h *raymond.HelperArg) interface{} {
				return h.DataStr("adjective") + " " + h.ParamStr(0) + h.FieldStr("exclaim")
			},
		},
		nil,
		"happy world?",
	},
	{
		"passing in data to a compiled function that expects data - data is passed to with block helpers where children use ..",
		`{{#hello}}{{world ../zomg}}{{/hello}}`,
		map[string]interface{}{"exclaim": true, "zomg": "world"},
		map[string]interface{}{"adjective": "happy", "accessData": "#win"},
		map[string]raymond.Helper{
			"hello": func(h *raymond.HelperArg) interface{} {
				return h.DataStr("accessData") + " " + h.BlockWithCtx(map[string]string{"exclaim": "?"})
			},
			"world": func(h *raymond.HelperArg) interface{} {
				return h.DataStr("adjective") + " " + h.ParamStr(0) + h.FieldStr("exclaim")
			},
		},
		nil,
		"#win happy world?",
	},
	{
		"you can override inherited data when invoking a helper",
		`{{#hello}}{{world zomg}}{{/hello}}`,
		map[string]interface{}{"exclaim": true, "zomg": "planet"},
		map[string]interface{}{"adjective": "happy"},
		map[string]raymond.Helper{
			"hello": func(h *raymond.HelperArg) interface{} {
				ctx := map[string]string{"exclaim": "?", "zomg": "world"}
				data := h.NewDataFrame()
				data.Set("adjective", "sad")

				return h.BlockWith(ctx, data)
			},
			"world": func(h *raymond.HelperArg) interface{} {
				return h.DataStr("adjective") + " " + h.ParamStr(0) + h.FieldStr("exclaim")
			},
		},
		nil,
		"sad world?",
	},
	{
		"you can override inherited data when invoking a helper with depth",
		`{{#hello}}{{world ../zomg}}{{/hello}}`,
		map[string]interface{}{"exclaim": true, "zomg": "world"},
		map[string]interface{}{"adjective": "happy"},
		map[string]raymond.Helper{
			"hello": func(h *raymond.HelperArg) interface{} {
				ctx := map[string]string{"exclaim": "?"}
				data := h.NewDataFrame()
				data.Set("adjective", "sad")

				return h.BlockWith(ctx, data)
			},
			"world": func(h *raymond.HelperArg) interface{} {
				return h.DataStr("adjective") + " " + h.ParamStr(0) + h.FieldStr("exclaim")
			},
		},
		nil,
		"sad world?",
	},
	{
		"@root - the root context can be looked up via @root",
		`{{@root.foo}}`,
		map[string]interface{}{"foo": "hello"},
		nil, nil, nil,
		"hello",
	},
	{
		"@root - passed root values take priority",
		`{{@root.foo}}`,
		nil,
		map[string]interface{}{"root": map[string]string{"foo": "hello"}},
		nil, nil,
		"hello",
	},
	{
		"nesting - the root context can be looked up via @root",
		`{{#helper}}{{#helper}}{{@./depth}} {{@../depth}} {{@../../depth}}{{/helper}}{{/helper}}`,
		map[string]interface{}{"foo": "hello"},
		map[string]interface{}{"depth": 0},
		map[string]raymond.Helper{
			"helper": func(h *raymond.HelperArg) interface{} {
				data := h.NewDataFrame()

				if depth, ok := h.Data("depth").(int); ok {
					data.Set("depth", depth+1)
				}

				return h.BlockWithData(data)
			},
		},
		nil,
		"2 1 0",
	},
}

func TestData(t *testing.T) {
	launchTests(t, dataTests)
}

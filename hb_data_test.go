package raymond

import "testing"

//
// Those tests come from:
//   https://github.com/wycats/handlebars.js/blob/master/spec/data.js
//
var hbDataTests = []raymondTest{
	{
		"passing in data to a compiled function that expects data - works with helpers",
		"{{hello}}",
		map[string]string{"noun": "cat"},
		map[string]interface{}{"adjective": "happy"},
		map[string]Helper{"hello": func(h *HelperArg) interface{} {
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
		map[string]Helper{"let": func(h *HelperArg) interface{} {
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
		map[string]Helper{"hello": func(h *HelperArg) interface{} {
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
		map[string]Helper{"hello": func(h *HelperArg) interface{} {
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
		map[string]Helper{"hello": func(h *HelperArg) interface{} {
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
		map[string]Helper{"hello": func(h *HelperArg) interface{} {
			return "Hello " + h.ParamStr(0)
		}},
		nil,
		// @todo Test differs with JS implementation: we don't output `undefined`
		"Hello ",
	},

	// @todo "parameter data throws when using complex scope references",

	// // @todo Implements data as function
	// {
	// 	"data can be functions",
	// 	`{{@hello}}`,
	// 	nil,
	// 	map[string]interface{}{"hello": func() string { return "hello" }},
	// 	nil, nil,
	// 	"hello",
	// },
	// // @todo Implements data as function
	// {
	//  "data can be functions with params",
	//  `{{@hello "hello"}}`,
	//  nil,
	//  map[string]interface{}{"hello": func(h *HelperArg) string { return h.ParamStr(0) }},
	//  nil, nil,
	//  "hello",
	// },

	{
		"data is inherited downstream",
		`{{#let foo=1 bar=2}}{{#let foo=bar.baz}}{{@bar}}{{@foo}}{{/let}}{{@foo}}{{/let}}`,
		map[string]map[string]string{"bar": {"baz": "hello world"}},
		nil,
		map[string]Helper{"let": func(h *HelperArg) interface{} {
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
		map[string]Helper{"hello": func(h *HelperArg) interface{} {
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
		map[string]Helper{"hello": func(h *HelperArg) interface{} {
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
		map[string]Helper{
			"hello": func(h *HelperArg) interface{} {
				return h.Block()
			},
			"world": func(h *HelperArg) interface{} {
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
		map[string]Helper{
			"hello": func(h *HelperArg) interface{} {
				return h.BlockWithCtx(map[string]string{"exclaim": "?"})
			},
			"world": func(h *HelperArg) interface{} {
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
		map[string]Helper{
			"hello": func(h *HelperArg) interface{} {
				return h.DataStr("accessData") + " " + h.BlockWithCtx(map[string]string{"exclaim": "?"})
			},
			"world": func(h *HelperArg) interface{} {
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
		map[string]Helper{
			"hello": func(h *HelperArg) interface{} {
				ctx := map[string]string{"exclaim": "?", "zomg": "world"}
				data := h.NewDataFrame()
				data.Set("adjective", "sad")

				return h.BlockWith(ctx, data, nil)
			},
			"world": func(h *HelperArg) interface{} {
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
		map[string]Helper{
			"hello": func(h *HelperArg) interface{} {
				ctx := map[string]string{"exclaim": "?"}
				data := h.NewDataFrame()
				data.Set("adjective", "sad")

				return h.BlockWith(ctx, data, nil)
			},
			"world": func(h *HelperArg) interface{} {
				return h.DataStr("adjective") + " " + h.ParamStr(0) + h.FieldStr("exclaim")
			},
		},
		nil,
		"sad world?",
	},

	// @todo Add remaining tests
}

func TestHandlebarsData(t *testing.T) {
	launchHandlebarsTests(t, hbDataTests)
}

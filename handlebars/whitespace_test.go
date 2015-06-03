package handlebars

import "testing"

//
// Those tests come from:
//   https://github.com/wycats/handlebars.js/blob/master/spec/whitespace-control.js
//
var whitespaceControlTests = []Test{
	{
		"should strip whitespace around mustache calls (1)",
		" {{~foo~}} ",
		map[string]string{"foo": "bar<"},
		nil, nil, nil,
		"bar&lt;",
	},
	{
		"should strip whitespace around mustache calls (2)",
		" {{~foo}} ",
		map[string]string{"foo": "bar<"},
		nil, nil, nil,
		"bar&lt; ",
	},
	{
		"should strip whitespace around mustache calls (3)",
		" {{foo~}} ",
		map[string]string{"foo": "bar<"},
		nil, nil, nil,
		" bar&lt;",
	},
	{
		"should strip whitespace around mustache calls (4)",
		" {{~&foo~}} ",
		map[string]string{"foo": "bar<"},
		nil, nil, nil,
		"bar<",
	},
	{
		"should strip whitespace around mustache calls (5)",
		" {{~{foo}~}} ",
		map[string]string{"foo": "bar<"},
		nil, nil, nil,
		"bar<",
	},
	{
		"should strip whitespace around mustache calls (6)",
		"1\n{{foo~}} \n\n 23\n{{bar}}4",
		nil, nil, nil, nil,
		"1\n23\n4",
	},

	{
		"blocks - should strip whitespace around simple block calls (1)",
		" {{~#if foo~}} bar {{~/if~}} ",
		map[string]string{"foo": "bar<"},
		nil, nil, nil,
		"bar",
	},
	{
		"blocks - should strip whitespace around simple block calls (2)",
		" {{#if foo~}} bar {{/if~}} ",
		map[string]string{"foo": "bar<"},
		nil, nil, nil,
		" bar ",
	},
	{
		"blocks - should strip whitespace around simple block calls (3)",
		" {{~#if foo}} bar {{~/if}} ",
		map[string]string{"foo": "bar<"},
		nil, nil, nil,
		" bar ",
	},
	{
		"blocks - should strip whitespace around simple block calls (4)",
		" {{#if foo}} bar {{/if}} ",
		map[string]string{"foo": "bar<"},
		nil, nil, nil,
		"  bar  ",
	},
	{
		"blocks - should strip whitespace around simple block calls (5)",
		" \n\n{{~#if foo~}} \n\nbar \n\n{{~/if~}}\n\n ",
		map[string]string{"foo": "bar<"},
		nil, nil, nil,
		"bar",
	},
	{
		"blocks - should strip whitespace around simple block calls (6)",
		" a\n\n{{~#if foo~}} \n\nbar \n\n{{~/if~}}\n\na ",
		map[string]string{"foo": "bar<"},
		nil, nil, nil,
		" abara ",
	},

	{
		"should strip whitespace around inverse block calls (1)",
		" {{~^if foo~}} bar {{~/if~}} ",
		nil, nil, nil, nil,
		"bar",
	},
	{
		"should strip whitespace around inverse block calls (2)",
		" {{^if foo~}} bar {{/if~}} ",
		nil, nil, nil, nil,
		" bar ",
	},
	{
		"should strip whitespace around inverse block calls (3)",
		" {{~^if foo}} bar {{~/if}} ",
		nil, nil, nil, nil,
		" bar ",
	},
	{
		"should strip whitespace around inverse block calls (4)",
		" {{^if foo}} bar {{/if}} ",
		nil, nil, nil, nil,
		"  bar  ",
	},
	{
		"should strip whitespace around inverse block calls (5)",
		" \n\n{{~^if foo~}} \n\nbar \n\n{{~/if~}}\n\n ",
		nil, nil, nil, nil,
		"bar",
	},

	{
		"should strip whitespace around complex block calls (1)",
		"{{#if foo~}} bar {{~^~}} baz {{~/if}}",
		map[string]string{"foo": "bar<"},
		nil, nil, nil,
		"bar",
	},
	{
		"should strip whitespace around complex block calls (2)",
		"{{#if foo~}} bar {{^~}} baz {{/if}}",
		map[string]string{"foo": "bar<"},
		nil, nil, nil,
		"bar ",
	},
	{
		"should strip whitespace around complex block calls (3)",
		"{{#if foo}} bar {{~^~}} baz {{~/if}}",
		map[string]string{"foo": "bar<"},
		nil, nil, nil,
		" bar",
	},
	{
		"should strip whitespace around complex block calls (4)",
		"{{#if foo}} bar {{^~}} baz {{/if}}",
		map[string]string{"foo": "bar<"},
		nil, nil, nil,
		" bar ",
	},
	{
		"should strip whitespace around complex block calls (5)",
		"{{#if foo~}} bar {{~else~}} baz {{~/if}}",
		map[string]string{"foo": "bar<"},
		nil, nil, nil,
		"bar",
	},
	{
		"should strip whitespace around complex block calls (6)",
		"\n\n{{~#if foo~}} \n\nbar \n\n{{~^~}} \n\nbaz \n\n{{~/if~}}\n\n",
		map[string]string{"foo": "bar<"},
		nil, nil, nil,
		"bar",
	},
	{
		"should strip whitespace around complex block calls (7)",
		"\n\n{{~#if foo~}} \n\n{{{foo}}} \n\n{{~^~}} \n\nbaz \n\n{{~/if~}}\n\n",
		map[string]string{"foo": "bar<"},
		nil, nil, nil,
		"bar<",
	},
	{
		"should strip whitespace around complex block calls (8)",
		"{{#if foo~}} bar {{~^~}} baz {{~/if}}",
		nil, nil, nil, nil,
		"baz",
	},
	{
		"should strip whitespace around complex block calls (9)",
		"{{#if foo}} bar {{~^~}} baz {{/if}}",
		nil, nil, nil, nil,
		"baz ",
	},
	{
		"should strip whitespace around complex block calls (10)",
		"{{#if foo~}} bar {{~^}} baz {{~/if}}",
		nil, nil, nil, nil,
		" baz",
	},
	{
		"should strip whitespace around complex block calls (11)",
		"{{#if foo~}} bar {{~^}} baz {{/if}}",
		nil, nil, nil, nil,
		" baz ",
	},
	{
		"should strip whitespace around complex block calls (12)",
		"{{#if foo~}} bar {{~else~}} baz {{~/if}}",
		nil, nil, nil, nil,
		"baz",
	},
	{
		"should strip whitespace around complex block calls (13)",
		"\n\n{{~#if foo~}} \n\nbar \n\n{{~^~}} \n\nbaz \n\n{{~/if~}}\n\n",
		nil, nil, nil, nil,
		"baz",
	},

	{
		"should strip whitespace around partials (1)",
		"foo {{~> dude~}} ",
		nil, nil, nil,
		map[string]string{"dude": "bar"},
		"foobar",
	},
	{
		"should strip whitespace around partials (2)",
		"foo {{> dude~}} ",
		nil, nil, nil,
		map[string]string{"dude": "bar"},
		"foo bar",
	},
	{
		"should strip whitespace around partials (3)",
		"foo {{> dude}} ",
		nil, nil, nil,
		map[string]string{"dude": "bar"},
		"foo bar ",
	},
	{
		"should strip whitespace around partials (4)",
		"foo\n {{~> dude}} ",
		nil, nil, nil,
		map[string]string{"dude": "bar"},
		"foobar",
	},
	{
		"should strip whitespace around partials (5)",
		"foo\n {{> dude}} ",
		nil, nil, nil,
		map[string]string{"dude": "bar"},
		"foo\n bar",
	},

	{
		"should only strip whitespace once",
		" {{~foo~}} {{foo}} {{foo}} ",
		map[string]string{"foo": "bar"},
		nil, nil, nil,
		"barbar bar ",
	},
}

func TestWhitespaceControl(t *testing.T) {
	launchTests(t, whitespaceControlTests)
}

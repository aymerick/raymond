package raymond

import "testing"

//
// Those tests come from:
//   https://github.com/wycats/handlebars.js/blob/master/spec/whitespace-control.js
//
var hbWhitespaceControlTests = []raymondTest{
	{
		"should strip whitespace around mustache calls (1)",
		" {{~foo~}} ",
		map[string]string{"foo": "bar<"},
		nil,
		"bar&lt;",
	},
	{
		"should strip whitespace around mustache calls (2)",
		" {{~foo}} ",
		map[string]string{"foo": "bar<"},
		nil,
		"bar&lt; ",
	},
	{
		"should strip whitespace around mustache calls (3)",
		" {{foo~}} ",
		map[string]string{"foo": "bar<"},
		nil,
		" bar&lt;",
	},
	{
		"should strip whitespace around mustache calls (4)",
		" {{~&foo~}} ",
		map[string]string{"foo": "bar<"},
		nil,
		"bar<",
	},
	{
		"should strip whitespace around mustache calls (5)",
		" {{~{foo}~}} ",
		map[string]string{"foo": "bar<"},
		nil,
		"bar<",
	},
	{
		"should strip whitespace around mustache calls (6)",
		"1\n{{foo~}} \n\n 23\n{{bar}}4",
		nil,
		nil,
		"1\n23\n4",
	},
}

func TestHandlebarsWhitespaceControl(t *testing.T) {
	launchHandlebarsTests(t, hbWhitespaceControlTests)
}

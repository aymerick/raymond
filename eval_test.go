package raymond

import (
	"testing"
)

var evalTests = []Test{
	{
		"only content",
		"this is content",
		nil, nil, nil, nil,
		"this is content",
	},
	{
		"checks path in parent contexts",
		"{{#a}}{{one}}{{#b}}{{one}}{{two}}{{one}}{{/b}}{{/a}}",
		map[string]interface{}{"a": map[string]int{"one": 1}, "b": map[string]int{"two": 2}},
		nil, nil, nil,
		"1121",
	},
	{
		"block params",
		"{{#foo as |bar|}}{{bar}}{{/foo}}{{bar}}",
		map[string]string{"foo": "baz", "bar": "bat"},
		nil, nil, nil,
		"bazbat",
	},
	{
		"block params on array",
		"{{#foo as |bar i|}}{{i}}.{{bar}} {{/foo}}",
		map[string][]string{"foo": {"baz", "bar", "bat"}},
		nil, nil, nil,
		"0.baz 1.bar 2.bat ",
	},
	{
		"nested block params",
		"{{#foos as |foo iFoo|}}{{#wats as |wat iWat|}}{{iFoo}}.{{iWat}}.{{foo}}-{{wat}} {{/wats}}{{/foos}}",
		map[string][]string{"foos": {"baz", "bar"}, "wats": {"the", "phoque"}},
		nil, nil, nil,
		"0.0.baz-the 0.1.baz-phoque 1.0.bar-the 1.1.bar-phoque ",
	},
	{
		"block params with path reference",
		"{{#foo as |bar|}}{{bar.baz}}{{/foo}}",
		map[string]map[string]string{"foo": {"baz": "bat"}},
		nil, nil, nil,
		"bat",
	},
	{
		"falsy block evaluation",
		"{{#foo}}bar{{/foo}} baz",
		map[string]interface{}{"foo": false},
		nil, nil, nil,
		" baz",
	},
	{
		"block helper returns a SafeString",
		"{{title}} - {{#bold}}{{body}}{{/bold}}",
		map[string]string{
			"title": "My new blog post",
			"body":  "I have so many things to say!",
		},
		nil,
		map[string]interface{}{"bold": func(options *Options) SafeString {
			return SafeString(`<div class="mybold">` + options.Fn() + "</div>")
		}},
		nil,
		`My new blog post - <div class="mybold">I have so many things to say!</div>`,
	},
	{
		"chained blocks",
		"{{#if a}}A{{else if b}}B{{else}}C{{/if}}",
		map[string]interface{}{"b": false},
		nil, nil, nil,
		"C",
	},

	// @todo Test with a "../../path" (depth 2 path) while context is only depth 1
}

func TestEval(t *testing.T) {
	t.Parallel()

	launchTests(t, evalTests)
}

var evalErrors = []Test{
	{
		"functions with wrong number of arguments",
		`{{foo "bar"}}`,
		map[string]interface{}{"foo": func(a string, b string) string { return "foo" }},
		nil, nil, nil,
		"Helper 'foo' called with wrong number of arguments, needed 2 but got 1",
	},
	{
		"functions with wrong number of returned values (1)",
		"{{foo}}",
		map[string]interface{}{"foo": func() {}},
		nil, nil, nil,
		"Helper function must return a string or a SafeString",
	},
	{
		"functions with wrong number of returned values (2)",
		"{{foo}}",
		map[string]interface{}{"foo": func() (string, bool, string) { return "foo", true, "bar" }},
		nil, nil, nil,
		"Helper function must return a string or a SafeString",
	},
}

func TestEvalErrors(t *testing.T) {
	launchErrorTests(t, evalErrors)
}

func TestEvalStruct(t *testing.T) {
	t.Parallel()

	source := `<div class="post">
  <h1>By {{author.FirstName}} {{Author.lastName}}</h1>
  <div class="body">{{Body}}</div>

  <h1>Comments</h1>

  {{#each comments}}
  <h2>By {{Author.FirstName}} {{author.LastName}}</h2>
  <div class="body">{{body}}</div>
  {{/each}}
</div>`

	expected := `<div class="post">
  <h1>By Jean Valjean</h1>
  <div class="body">Life is difficult</div>

  <h1>Comments</h1>

  <h2>By Marcel Beliveau</h2>
  <div class="body">LOL!</div>
</div>`

	type Person struct {
		FirstName string
		LastName  string
	}

	type Comment struct {
		Author Person
		Body   string
	}

	type Post struct {
		Author   Person
		Body     string
		Comments []Comment
	}

	ctx := Post{
		Person{"Jean", "Valjean"},
		"Life is difficult",
		[]Comment{
			Comment{
				Person{"Marcel", "Beliveau"},
				"LOL!",
			},
		},
	}

	output := MustRender(source, ctx)
	if output != expected {
		t.Errorf("Failed to evaluate with struct context")
	}
}

func TestEvalStructTag(t *testing.T) {
	t.Parallel()

	source := `<div class="person">
	<h1>{{real-name}}</h1>
	<ul>
	  <li>City: {{info.location}}</li>
	  <li>Rug: {{info.[r.u.g]}}</li>
	  <li>Activity: {{info.activity}}</li>
	</ul>
	{{#each other-names}}
	<p>{{alias-name}}</p>
	{{/each}}
</div>`

	expected := `<div class="person">
	<h1>Lebowski</h1>
	<ul>
	  <li>City: Venice</li>
	  <li>Rug: Tied The Room Together</li>
	  <li>Activity: Bowling</li>
	</ul>
	<p>his dudeness</p>
	<p>el duderino</p>
</div>`

	type Alias struct {
		Name string `handlebars:"alias-name"`
	}

	type CharacterInfo struct {
		City     string `handlebars:"location"`
		Rug      string `handlebars:"r.u.g"`
		Activity string `handlebars:"not-activity"`
	}

	type Character struct {
		RealName string `handlebars:"real-name"`
		Info     CharacterInfo
		Aliases  []Alias `handlebars:"other-names"`
	}

	ctx := Character{
		"Lebowski",
		CharacterInfo{"Venice", "Tied The Room Together", "Bowling"},
		[]Alias{
			{"his dudeness"},
			{"el duderino"},
		},
	}

	output := MustRender(source, ctx)
	if output != expected {
		t.Errorf("Failed to evaluate with struct tag context")
	}
}

type TestFoo struct {
}

func (t *TestFoo) Subject() string {
	return "foo"
}

func TestEvalMethod(t *testing.T) {
	t.Parallel()

	source := `Subject is {{subject}}! YES I SAID {{Subject}}!`
	expected := `Subject is foo! YES I SAID foo!`

	ctx := &TestFoo{}

	output := MustRender(source, ctx)
	if output != expected {
		t.Errorf("Failed to evaluate struct method: %s", output)
	}
}

type TestBar struct {
}

func (t *TestBar) Subject() interface{} {
	return testBar
}

func testBar() string {
	return "bar"
}

func TestEvalMethodReturningFunc(t *testing.T) {
	t.Parallel()

	source := `Subject is {{subject}}! YES I SAID {{Subject}}!`
	expected := `Subject is bar! YES I SAID bar!`

	ctx := &TestBar{}

	output := MustRender(source, ctx)
	if output != expected {
		t.Errorf("Failed to evaluate struct method: %s", output)
	}
}

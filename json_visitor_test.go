package raymond

import (
	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/require"
	"testing"
)

func TestJSONVisitor(t *testing.T) {
	for _, tt := range []struct {
		name   string
		source string
		want   map[string]interface{}
	}{{
		name:   "basic",
		source: sourceBasic,
		want: map[string]interface{}{
			"title": "test_title",
			"body":  "test_body",
		},
	}, {
		name: "nested vars",
		source: `<div class="entry">
  <h1>{{title.name}}</h1>
  <div class="body">
    {{body.content}}
  </div>
</div>`,
		want: map[string]interface{}{
			"title": map[string]interface{}{"name": "test_name"},
			"body":  map[string]interface{}{"content": "test_content"},
		},
	}, {
		name: "block params",
		source: `{{#foo as |bar|}}
{{bar.baz}}
{{/foo}}`,
		want: map[string]interface{}{
			"bar": map[string]interface{}{"baz": "test_baz"},
		},
	}, {
		name: "with block",
		source: `{{#with people.[0].[0]}}
  {{name}}
{{/with}}`,
		want: map[string]interface{}{"people": newList(newList(map[string]interface{}{"name": "test_name"}))},
	}, {
		name:   "if block",
		source: `{{#if people.name}} {{people.name}}{{/if}}`,
		want: map[string]interface{}{
			"people": map[string]interface{}{"name": "test_name"},
		},
	}, {
		name:   "if block with incomplete foo path and complete foo path",
		source: `{{#if foo}} {{foo.bar.baz}} {{/if}}`,
		want:   map[string]interface{}{"foo": map[string]interface{}{"bar": map[string]interface{}{"baz": "test_baz"}}},
	}, {
		name:   "accesses multiple elements of a map in multiple paths",
		source: `{{bar.baz}} {{name.first}}{{name.last}}`,
		want: map[string]interface{}{
			"bar":  map[string]interface{}{"baz": "test_baz"},
			"name": map[string]interface{}{"first": "test_first", "last": "test_last"}},
	}, {
		name:   "large template",
		source: largeTemplate,
		want:   map[string]interface{}{"bar": "test_bar", "foo": "test_foo", "name": "test_name", "phone": "test_phone"},
	}, {
		name:   "multi with",
		source: "{{#with foo}}{{#with bar}}{{baz}}{{/with}}{{/with}}",
		want:   map[string]interface{}{"foo": map[string]interface{}{"bar": map[string]interface{}{"baz": "test_baz"}}},
	}, {
		name:   "multi as",
		source: "{{#with foo as |bee|}}{{#with bee.bar as |bazinga|}}{{bazinga.baz}}{{/with}}{{/with}}",
		want:   map[string]interface{}{"foo": map[string]interface{}{"bar": map[string]interface{}{"baz": "test_baz"}}},
	}, {
		name:   "multi multi with",
		source: "{{#with fizz}}{{#with foo}}{{#with bar}}{{baz}}{{bop}}{{/with}}{{/with}}{{/with}}",
		want:   map[string]interface{}{"fizz": map[string]interface{}{"foo": map[string]interface{}{"bar": map[string]interface{}{"baz": "test_baz", "bop": "test_bop"}}}},
	}, {
		name:   "multi with same names",
		source: "{{#with foo}}{{#with foo}}{{baz}}{{/with}}{{/with}}",
		want:   map[string]interface{}{"foo": map[string]interface{}{"foo": map[string]interface{}{"baz": "test_baz"}}},
	}, {
		name:   "up a level",
		source: "{{#with foo}}{{#with foo}}{{../baz}}{{/with}}{{/with}}",
		want:   map[string]interface{}{"foo": map[string]interface{}{"baz": "test_baz"}},
	}, {
		name:   "each lookup",
		source: "{{#each people}} {{.}} lives in {{lookup ../cities @index}}{{/each}}",
		want:   map[string]interface{}{"people": newList("test_people"), "cities": newList("test_cities")},
	}, {
		name:   "each lookup complex",
		source: "{{#each people}} {{./foo/bar/baz}} lives in {{lookup ../cities @index}}{{/each}}",
		want:   map[string]interface{}{"people": newList(map[string]interface{}{"foo": map[string]interface{}{"bar": map[string]interface{}{"baz": "test_baz"}}}), "cities": newList("test_cities")},
	}, {
		name:   "each",
		source: "{{#with foo}}{{#each foo}}{{baz}}{{/each}}{{/with}}",
		want:   map[string]interface{}{"foo": map[string]interface{}{"foo": newList(map[string]interface{}{"baz": "test_baz"})}},
	}, {
		name:   "multiple paths in a non-block helper block",
		source: `{{#foo bar baz}} {{name.first name.last}} {{/foo}}`,
		want: map[string]interface{}{
			"bar": "test_bar",
			"baz": "test_baz",
			"name": map[string]interface{}{
				"first": "test_first",
				"last":  "test_last"}},
	}} {
		t.Run(tt.name, func(t *testing.T) {
			tpl, err := Parse(tt.source)
			require.NoError(t, err)
			require.Equal(t, tpl.source, tt.source)

			//fmt.Println(tpl.PrintAST())

			vars, err := tpl.ExtractTemplateVars()
			require.NoError(t, err)
			assert.Equal(t, tt.want, vars)
		})
	}
}

var largeTemplate = `<html>
    {{#if name}}
    <div>Hello {{name}}!</dev>
    {{else}}
    <div>Hello there!</div>
    {{/if}}
    
    {{#ifGt foo bar}}
    <br><br><div>foo is greater than bar</div>
    {{/ifGt}}
    
    {{#ifGt foo 10}}
    <br><br><div>foo is greater than 10</div>
    {{else}}
    <br><br><div>foo is not greater than 10</div>
    {{/ifGt}}
    
    
    {{#ifLt foo bar}}
    <br><br><div>foo is less than bar</div><br><br>
    {{/ifLt}}
    
    {{#ifLt foo 10}}
    <div>foo is less than 10</div>
    {{else}}
    <div>foo is not less than 10</div>
    {{/ifLt}}
    
    {{#ifEq foo bar}}
    <br><br><div>foo is equal to bar</div><br><br>
    {{/ifEq}}
    
    {{#ifEq foo 10}}
    <div>foo is equal to 10</div>
    {{else}}
    <div>foo is not equal to 10</div>
    {{/ifEq}}
    
    {{#ifMatchesRegexStr "^(\+\d{1,2}\s)?\(?\d{3}\)?[\s.-]\d{3}[\s.-]\d{4}$" phone}}
    <br><div>phone var is a phone number</div>
    {{else}}
    <br><div>phone var is not a phone number</div>
    {{/ifMatchesRegexStr}}
</html>`

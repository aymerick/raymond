package raymond

import "testing"

var testInput = `<div class="entry">
  <h1>{{title}}</h1>
  <div class="body">
    {{body}}
  </div>
</div>`

var testOutput = `<div class="entry">
  <h1>foo</h1>
  <div class="body">
    bar
  </div>
</div>`

func TestRender(t *testing.T) {
	output := Render(testInput, map[string]string{"title": "foo", "body": "bar"})
	if output != testOutput {
		t.Errorf("Failed to render template\ninput:\n\n'%s'\n\nexpected:\n\n%s\n\ngot:\n\n%s", testInput, testOutput, output)
	}
}

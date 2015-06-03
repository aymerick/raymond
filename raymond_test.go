package raymond

import (
	"testing"
)

//
// Basic rendering test
//

var testInput = `<div class="entry">
  <h1>{{title}}</h1>
  <div class="body">
    {{body}}
  </div>
</div>`

var testCtx = map[string]string{"title": "foo", "body": "bar"}

var testOutput = `<div class="entry">
  <h1>foo</h1>
  <div class="body">
    bar
  </div>
</div>`

func TestRender(t *testing.T) {
	output, err := Render(testInput, testCtx)
	if err != nil || (output != testOutput) {
		t.Errorf("Failed to render template\ninput:\n\n'%s'\n\nexpected:\n\n%s\n\ngot:\n\n%serror:\n\n%s", testInput, testOutput, output, err)
	}
}

func TestMustRender(t *testing.T) {
	output := MustRender(testInput, testCtx)
	if (output != testOutput) {
		t.Errorf("Failed to render template\ninput:\n\n'%s'\n\nexpected:\n\n%s\n\ngot:\n\n%s", testInput, testOutput, output)
	}
}

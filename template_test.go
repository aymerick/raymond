package raymond

import "testing"

var sourceBasic = `<div class="entry">
  <h1>{{title}}</h1>
  <div class="body">
    {{body}}
  </div>
</div>`

var basicAST = `CONTENT[ '<div class="entry">
  <h1>' ]
{{ PATH:title [] }}
CONTENT[ '</h1>
  <div class="body">
    ' ]
{{ PATH:body [] }}
CONTENT[ '
  </div>
</div>' ]
`

func TestNewTemplate(t *testing.T) {
	stats.test()

	tpl := NewTemplate(sourceBasic)
	if tpl.source != sourceBasic {
		t.Errorf("Faild to instantiate template")
		stats.failed()
	}

	stats.output()
}

func TestParseTemplate(t *testing.T) {
	stats.test()

	tpl, err := Parse(sourceBasic)
	if err != nil || (tpl.source != sourceBasic) {
		t.Errorf("Faild to parse template")
		stats.failed()
	}

	if str := tpl.PrintAST(); str != basicAST {
		t.Errorf("Template parsing incorrect: %s", str)
		stats.failed()
	}

	stats.output()
}

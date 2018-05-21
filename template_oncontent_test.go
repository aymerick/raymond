package raymond

import (
	"github.com/cmaster11/raymond/ast"
	"fmt"
	"github.com/cmaster11/raymond/parser"
	"testing"
)

type Map = map[string]interface{}
type Array = []interface{}

var serializeSource = `
Hello! My name is {{name}}.
{{#brum}}I have a brum named {{brum.name}}!{{else}}I don't have a brum{{/brum}}
I think {{> mmm}}
`

func TestTemplate_Serialize(t *testing.T) {

	onParserContent := func(text string) string {
		return "((" + text + "))"
	}

	parserOptions := &parser.ParserOptions{
		OnContent: onParserContent,
	}

	tpl, err := ParseWithOptions(serializeSource, parserOptions)
	if err != nil {
		t.Fatal(err)
	}

	tpl.RegisterPartial("mmm", "{{#each thoughts as |th idx|}}{{idx}}: {{th}}\n{{/each}}")

	fmt.Println(tpl.Serialize())

	tpl.OnContent = func(nodeType ast.NodeType, text string) string {
		switch nodeType {
		case ast.NodePartial:
			return "{" + text + "}"
		case ast.NodeMustache:
			return "<<" + text + ">>"
		default:
			panic("invalid node")
		}
	}

	ctx := Map{
		"name": "lol",
		"brum": Map{
			"name": "chuchu",
		},
		"thoughts": Array{
			"mmm", "omm", "mnaaa",
		},
	}

	output := tpl.MustExec(ctx)

	fmt.Println(output)
}

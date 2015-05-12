package raymond

import (
	"io"
	"runtime"

	"github.com/aymerick/raymond/ast"
	"github.com/aymerick/raymond/parser"
)

// Template
type Template struct {
	source  string
	program *ast.Program
}

// Instanciate a template an parse it
func Parse(source string) (*Template, error) {
	tpl := NewTemplate(source)

	// parse template
	if err := tpl.Parse(); err != nil {
		return nil, err
	}

	return tpl, nil
}

// Instanciate a new template
func NewTemplate(source string) *Template {
	return &Template{
		source: source,
	}
}

// Parse template
func (tpl *Template) Parse() error {
	if tpl.program == nil {
		var err error

		tpl.program, err = parser.Parse(tpl.source)
		if err != nil {
			return err
		}
	}

	return nil
}

// Returns string version of parsed template
func (tpl *Template) PrintAST() string {
	return ast.PrintNode(tpl.program)
}

// Renders a template with input data
func (tpl *Template) Exec(wr io.Writer, data interface{}) (err error) {
	defer errRecover(&err)

	// parses template if necessary
	err = tpl.Parse()
	if err != nil {
		return
	}

	// setup visitor
	v := NewEvalVisitor(wr, tpl, data)

	// visit AST
	tpl.program.Accept(v)

	return
}

// recovers exec panic
func errRecover(errp *error) {
	e := recover()
	if e != nil {
		switch err := e.(type) {
		case runtime.Error:
			panic(e)
		case error:
			*errp = err
		default:
			panic(e)
		}
	}
}

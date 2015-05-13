package raymond

import (
	"fmt"
	"io"
	"runtime"

	"github.com/aymerick/raymond/ast"
	"github.com/aymerick/raymond/parser"
)

// Template
type Template struct {
	source  string
	program *ast.Program
	helpers map[string]Helper
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
		source:  source,
		helpers: make(map[string]Helper),
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

// Register several helpers
func (tpl *Template) RegisterHelpers(helpers map[string]Helper) {
	for name, helper := range helpers {
		tpl.RegisterHelper(name, helper)
	}
}

// Register an helper
func (tpl *Template) RegisterHelper(name string, helper Helper) {
	if tpl.helpers[name] != nil {
		panic(fmt.Sprintf("Helper %s already registered", name))
	}

	tpl.helpers[name] = helper
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

// Returns string version of parsed template
func (tpl *Template) PrintAST() string {
	return ast.PrintNode(tpl.program)
}

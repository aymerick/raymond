package raymond

import (
	"fmt"
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

// Instanciate a template an parse it. Panics on error.
func MustParse(source string) *Template {
	result, err := Parse(source)
	if err != nil {
		panic(err)
	}
	return result
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
func (tpl *Template) Exec(data interface{}) (result string, err error) {
	defer errRecover(&err)

	// parses template if necessary
	err = tpl.Parse()
	if err != nil {
		return
	}

	// setup visitor
	v := NewEvalVisitor(tpl, data)

	// visit AST
	resEval := tpl.program.Accept(v)

	// get result
	result, _ = resEval.(string)

	// named return values
	return
}

// Renders a template with input data. Panics on error.
func (tpl *Template) MustExec(data interface{}) string {
	result, err := tpl.Exec(data)
	if err != nil {
		panic(err)
	}
	return result
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

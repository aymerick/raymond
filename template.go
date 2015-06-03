package raymond

import (
	"fmt"
	"runtime"

	"github.com/aymerick/raymond/ast"
	"github.com/aymerick/raymond/parser"
)

// Template represents a handlebars template.
type Template struct {
	source   string
	program  *ast.Program
	helpers  map[string]Helper
	partials map[string]*partial
}

// newTemplate instanciate a new template without parsing it
func newTemplate(source string) *Template {
	return &Template{
		source:   source,
		helpers:  make(map[string]Helper),
		partials: make(map[string]*partial),
	}
}

// Parse instanciates a template by parsing given source.
func Parse(source string) (*Template, error) {
	tpl := newTemplate(source)

	// parse template
	if err := tpl.Parse(); err != nil {
		return nil, err
	}

	return tpl, nil
}

// MustParse instanciates a template by parsing given source. It panics on error.
func MustParse(source string) *Template {
	result, err := Parse(source)
	if err != nil {
		panic(err)
	}
	return result
}

// Parse parses the template.
//
// It can be called several times, the parsing will be done only once.
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

// RegisterHelper registers a helper.
func (tpl *Template) RegisterHelper(name string, helper Helper) {
	if tpl.helpers[name] != nil {
		panic(fmt.Sprintf("Helper %s already registered", name))
	}

	tpl.helpers[name] = helper
}

// RegisterHelpers register several helpers.
func (tpl *Template) RegisterHelpers(helpers map[string]Helper) {
	for name, helper := range helpers {
		tpl.RegisterHelper(name, helper)
	}
}

// RegisterPartial registers a partial.
func (tpl *Template) RegisterPartial(name string, partial string) {
	if tpl.partials[name] != nil {
		panic(fmt.Sprintf("Partial %s already registered", name))
	}

	tpl.partials[name] = newPartial(name, partial)
}

// RegisterPartials registers several partials.
func (tpl *Template) RegisterPartials(partials map[string]string) {
	for name, partial := range partials {
		tpl.RegisterPartial(name, partial)
	}
}

// Exec renders template with given context.
func (tpl *Template) Exec(ctx interface{}) (result string, err error) {
	return tpl.ExecWith(ctx, nil)
}

// MustExec renders a template with given context. It panics on error.
func (tpl *Template) MustExec(ctx interface{}) string {
	result, err := tpl.Exec(ctx)
	if err != nil {
		panic(err)
	}
	return result
}

// ExecWith renders a template with given context and private data frame.
func (tpl *Template) ExecWith(ctx interface{}, privData *DataFrame) (result string, err error) {
	defer errRecover(&err)

	// parses template if necessary
	err = tpl.Parse()
	if err != nil {
		return
	}

	// setup visitor
	v := newEvalVisitor(tpl, ctx, privData)

	// visit AST
	result, _ = tpl.program.Accept(v).(string)

	// named return values
	return
}

// errRecover recovers evaluation panic
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

// PrintAST returns string representation of parsed template.
func (tpl *Template) PrintAST() string {
	if err := tpl.Parse(); err != nil {
		return fmt.Sprintf("PARSER ERROR: %s", err)
	}

	return ast.Print(tpl.program)
}

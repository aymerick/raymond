package raymond

import "fmt"

// partial represents a partial template
type partial struct {
	name   string
	source string
	tpl    *Template
}

// partials storres all global partials
var partials map[string]*partial

func init() {
	partials = make(map[string]*partial)
}

// newPartial instanciates a new partial
func newPartial(name string, source string) *partial {
	return &partial{
		name:   name,
		source: source,
	}
}

// RegisterPartial registers a global partial.
func RegisterPartial(name string, source string) {
	if partials[name] != nil {
		panic(fmt.Errorf("Partial already registered: %s", name))
	}

	partials[name] = newPartial(name, source)
}

// RegisterPartials registers several global partials.
func RegisterPartials(partials map[string]string) {
	for name, p := range partials {
		RegisterPartial(name, p)
	}
}

// findPartial finds a registered global partial
func findPartial(name string) *partial {
	return partials[name]
}

// Template returns parsed partial template
func (p *partial) template() (*Template, error) {
	if p.tpl == nil {
		var err error

		p.tpl, err = Parse(p.source)
		if err != nil {
			return nil, err
		}
	}

	return p.tpl, nil
}

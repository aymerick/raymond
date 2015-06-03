package raymond

import "fmt"

// Partial represents a partial template
type Partial struct {
	name   string
	source string
	tpl    *Template
}

// partials storres all global partials
var partials map[string]*Partial

func init() {
	partials = make(map[string]*Partial)
}

// newPartial instanciates a new partial
func NewPartial(name string, source string) *Partial {
	return &Partial{
		name:   name,
		source: source,
	}
}

// RegisterPartial registers a global partial
func RegisterPartial(name string, source string) {
	if partials[name] != nil {
		panic(fmt.Errorf("Partial already registered: %s", name))
	}

	partials[name] = NewPartial(name, source)
}

// RegisterPartials registers several global partials
func RegisterPartials(partials map[string]string) {
	for name, partial := range partials {
		RegisterPartial(name, partial)
	}
}

// FindPartial finds a registered global partial
func FindPartial(name string) *Partial {
	return partials[name]
}

// Template returns parsed partial template
func (p *Partial) Template() (*Template, error) {
	if p.tpl == nil {
		var err error

		p.tpl, err = Parse(p.source)
		if err != nil {
			return nil, err
		}
	}

	return p.tpl, nil
}

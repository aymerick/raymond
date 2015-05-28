package raymond

import "fmt"

// A partial template
type Partial struct {
	name   string
	source string
	tpl    *Template
}

// All global partials
var partials map[string]*Partial

func init() {
	partials = make(map[string]*Partial)
}

// Instanciate a new partial
func NewPartial(name string, source string) *Partial {
	return &Partial{
		name:   name,
		source: source,
	}
}

// Registers a new global partial
func RegisterPartial(name string, source string) {
	if partials[name] != nil {
		panic(fmt.Errorf("Partial already registered: %s", name))
	}

	partials[name] = NewPartial(name, source)
}

// Find a registered global partial
func FindPartial(name string) *Partial {
	return partials[name]
}

// Return partial templae
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

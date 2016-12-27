package ray

import (
	"fmt"
	"regexp"
	"testing"

	"github.com/stretchr/testify/require"
)

type Test struct {
	name     string
	input    string
	data     interface{}
	privData map[string]interface{}
	helpers  map[string]interface{}
	partials map[string]string
	output   interface{}
}

func launchTests(t *testing.T, tests []Test) {
	// NOTE: TestMustache() makes Parallel testing fail
	// t.Parallel()
	r := require.New(t)

	for _, test := range tests {
		var err error
		var tpl *Template

		// parse template
		tpl, err = Parse(test.input)
		r.NoError(err)

		if len(test.helpers) > 0 {
			// register helpers
			tpl.RegisterHelpers(test.helpers)
		}

		if len(test.partials) > 0 {
			// register partials
			tpl.RegisterPartials(test.partials)
		}

		// setup private data frame
		var privData *DataFrame
		if test.privData != nil {
			privData = NewDataFrame()
			for k, v := range test.privData {
				privData.Set(k, v)
			}
		}

		// render template
		output, err := tpl.ExecWith(test.data, privData)
		r.NoError(err)
		// check output
		var expectedArr []string
		expectedArr, ok := test.output.([]string)
		if ok {
			match := false
			for _, expectedStr := range expectedArr {
				if expectedStr == output {
					match = true
					break
				}
			}

			r.True(match)
		} else {
			expectedStr, ok := test.output.(string)
			r.True(ok)
			r.Equal(expectedStr, output)
		}
	}
}

func launchErrorTests(t *testing.T, tests []Test) {
	t.Parallel()
	r := require.New(t)

	for _, test := range tests {
		var err error
		var tpl *Template

		// parse template
		tpl, err = Parse(test.input)
		r.NoError(err)
		if len(test.helpers) > 0 {
			// register helpers
			tpl.RegisterHelpers(test.helpers)
		}

		if len(test.partials) > 0 {
			// register partials
			tpl.RegisterPartials(test.partials)
		}

		// setup private data frame
		var privData *DataFrame
		if test.privData != nil {
			privData := NewDataFrame()
			for k, v := range test.privData {
				privData.Set(k, v)
			}
		}

		// render template
		output, err := tpl.ExecWith(test.data, privData)
		if err == nil {
			t.Fatalf("Test '%s' failed - Error expected\ninput:\n\t'%s'\ngot\n\t%q\nAST:\n%q", test.name, test.input, output, tpl.PrintAST())
		} else {
			var errMatch error
			match := false

			// check output
			var expectedArr []string
			expectedArr, ok := test.output.([]string)
			if ok {
				if len(expectedArr) > 0 {
					for _, expectedStr := range expectedArr {
						match, errMatch = regexp.MatchString(regexp.QuoteMeta(expectedStr), fmt.Sprint(err))
						r.NoError(errMatch)

						if match {
							break
						}
					}
				} else {
					// nothing to test
					match = true
				}
			} else {
				expectedStr, ok := test.output.(string)
				r.True(ok)

				if expectedStr != "" {
					match, errMatch = regexp.MatchString(regexp.QuoteMeta(expectedStr), fmt.Sprint(err))
					r.NoError(errMatch)
				} else {
					// nothing to test
					match = true
				}
			}

			r.True(match)
		}
	}
}

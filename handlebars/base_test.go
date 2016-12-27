package handlebars

import (
	"io/ioutil"
	"path"
	"strconv"
	"testing"

	"github.com/gobuffalo/ray"
	"github.com/stretchr/testify/require"
)

// cf. https://github.com/aymerick/go-fuzz-tests/ray
const dumpTpl = false

var dumpTplNb = 0

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
	t.Parallel()
	r := require.New(t)

	for _, test := range tests {
		var err error
		var tpl *ray.Template

		if dumpTpl {
			filename := strconv.Itoa(dumpTplNb)
			err := ioutil.WriteFile(path.Join(".", "dump_tpl", filename), []byte(test.input), 0644)
			r.NoError(err)
			dumpTplNb++
		}

		// parse template
		tpl, err = ray.Parse(test.input)
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
		var privData *ray.DataFrame
		if test.privData != nil {
			privData = ray.NewDataFrame()
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

package raymond

import (
	"fmt"
	"regexp"
	"sync"
	"testing"
)

//
// Basic rendering test
//

var testInput = `<div class="entry">
  <h1>{{title}}</h1>
  <div class="body">
    {{body}}
  </div>
</div>`

var testOutput = `<div class="entry">
  <h1>foo</h1>
  <div class="body">
    bar
  </div>
</div>`

func TestRender(t *testing.T) {
	stats.test()

	output, err := Render(testInput, map[string]string{"title": "foo", "body": "bar"})
	if err != nil || (output != testOutput) {
		t.Errorf("Failed to render template\ninput:\n\n'%s'\n\nexpected:\n\n%s\n\ngot:\n\n%serror:\n\n%s", testInput, testOutput, output, err)
		stats.failed()
	}

	stats.output()
}

//
// Test stats
//

type raymondStats struct {
	sync.Mutex
	total      int
	fail       int
	handlebars int
	mustache   int
}

var (
	stats = &raymondStats{}
)

func (stats *raymondStats) tests(nb int) {
	stats.Lock()
	stats.total += nb
	stats.Unlock()
}

func (stats *raymondStats) test() {
	stats.tests(1)
}

func (stats *raymondStats) handlebarsTests(nb int) {
	stats.Lock()
	stats.handlebars += nb
	stats.Unlock()
}

func (stats *raymondStats) mustacheTests(nb int) {
	stats.Lock()
	stats.mustache += nb
	stats.Unlock()
}

func (stats *raymondStats) failed() {
	stats.Lock()
	stats.fail += 1
	stats.Unlock()
}

func (stats *raymondStats) output() {
	fmt.Printf("[stats] Total: %d (handlebars: %d / mustache: %d) - Failed: %d\n", stats.total, stats.handlebars, stats.mustache, stats.fail)
}

//
// Generic test
//

type raymondTest struct {
	name    string
	input   string
	data    interface{}
	helpers map[string]Helper
	output  interface{}
}

// launch an array of tests
func launchRaymondTests(t *testing.T, tests []raymondTest) {
	stats.tests(len(tests))

	for _, test := range tests {
		var err error
		var tpl *Template

		// log.Printf("****************************************")
		// log.Printf("* TEST: '%s'", test.name)

		// parse template
		tpl, err = Parse(test.input)
		if err != nil {
			t.Errorf("Test '%s' failed - Failed to parse template\ninput:\n\t'%s'\nerror:\n\t%s", test.name, test.input, err)
			stats.failed()
		} else {
			if len(test.helpers) > 0 {
				// register helpers
				tpl.RegisterHelpers(test.helpers)
			}

			// render template
			output, err := tpl.Exec(test.data)
			if err != nil {
				t.Errorf("Test '%s' failed\ninput:\n\t'%s'\ndata:\n\t%s\nerror:\n\t%s\nAST:\n\t%s", test.name, test.input, Str(test.data), err, tpl.PrintAST())
				stats.failed()
			} else {
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

					if !match {
						t.Errorf("Test '%s' failed\ninput:\n\t'%s'\ndata:\n\t%s\nexpected\n\t%q\ngot\n\t%q\nAST:\n%s", test.name, test.input, Str(test.data), expectedArr, output, tpl.PrintAST())
						stats.failed()
					}
				} else {
					expectedStr, ok := test.output.(string)
					if !ok {
						panic(fmt.Errorf("Erroneous test output description: %q", test.output))
					}

					if expectedStr != output {
						t.Errorf("Test '%s' failed\ninput:\n\t'%s'\ndata:\n\t%s\nexpected\n\t%q\ngot\n\t%q\nAST:\n%s", test.name, test.input, Str(test.data), expectedStr, output, tpl.PrintAST())
						stats.failed()
					}
				}
			}
		}
	}

	stats.output()
}

// launch handlebars tests with stats
func launchHandlebarsTests(t *testing.T, tests []raymondTest) {
	stats.handlebarsTests(len(tests))

	launchRaymondTests(t, tests)
}

// launch mustache tests with stats
func launchMustacheTests(t *testing.T, tests []raymondTest) {
	stats.mustacheTests(len(tests))

	launchRaymondTests(t, tests)
}

// launch an array of error tests
func launchRaymondErrorTests(t *testing.T, tests []raymondTest) {
	stats.tests(len(tests))

	for _, test := range tests {
		var err error
		var tpl *Template

		// log.Printf("****************************************")
		// log.Printf("* TEST: '%s'", test.name)

		// parse template
		tpl, err = Parse(test.input)
		if err != nil {
			t.Errorf("Test '%s' failed - Failed to parse template\ninput:\n\t'%s'\nerror:\n\t%s", test.name, test.input, err)
			stats.failed()
		} else {
			if len(test.helpers) > 0 {
				// register helpers
				tpl.RegisterHelpers(test.helpers)
			}

			// render template
			output, err := tpl.Exec(test.data)
			if err == nil {
				t.Errorf("Test '%s' failed - Error expected\ninput:\n\t'%s'\ngot\n\t%q\nAST:\n%q", test.name, test.input, output, tpl.PrintAST())
				stats.failed()
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
							if errMatch != nil {
								panic("Failed to match regexp")
							}

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
					if !ok {
						panic(fmt.Errorf("Erroneous test output description: %q", test.output))
					}

					if expectedStr != "" {
						match, errMatch = regexp.MatchString(regexp.QuoteMeta(expectedStr), fmt.Sprint(err))
						if errMatch != nil {
							panic("Failed to match regexp")
						}
					} else {
						// nothing to test
						match = true
					}
				}

				if !match {
					t.Errorf("Test '%s' failed - Incorrect error returned\ninput:\n\t'%s'\ndata:\n\t%s\nexpected\n\t%q\ngot\n\t%q", test.name, test.input, Str(test.data), test.output, err)
					stats.failed()
				}
			}
		}
	}

	stats.output()
}

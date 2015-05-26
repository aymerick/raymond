package raymond

import (
	"io/ioutil"
	"path"
	"regexp"
	"testing"

	"gopkg.in/yaml.v2"
)

//
// Note that, as the JS implementation, we do not support:
//   - support alternative delimeters
//   - the mustache lambda spec
//

type mustacheTest struct {
	Name     string
	Desc     string
	Data     interface{}
	Template string
	Expected string
	Partials map[string]string
}

type mustacheTestFile struct {
	Overview string
	Tests    []mustacheTest
}

var (
	rAltDelim = regexp.MustCompile(regexp.QuoteMeta("{{="))
)

func testsFromMustacheFile(fileName string) []raymondTest {
	result := []raymondTest{}

	fileData, err := ioutil.ReadFile(path.Join("mustache", "specs", fileName))
	if err != nil {
		panic(err)
	}

	var testFile mustacheTestFile
	if err := yaml.Unmarshal(fileData, &testFile); err != nil {
		panic(err)
	}

	for _, mustacheTest := range testFile.Tests {
		if mustBeSkipped(mustacheTest) {
			continue
		}

		test := raymondTest{
			name:   mustacheTest.Name,
			input:  mustacheTest.Template,
			data:   mustacheTest.Data,
			output: mustacheTest.Expected,
		}

		result = append(result, test)
	}

	return result
}

// returns true if test must be skipped
func mustBeSkipped(test mustacheTest) bool {
	// @todo Skip partials tests "Failed Lookup" and "Standalone Indentation"
	return haveAltDelimiter(test)
}

// returns true if test have alternative delimeter in template or in partials
func haveAltDelimiter(test mustacheTest) bool {
	// check template
	if rAltDelim.MatchString(test.Template) {
		return true
	}

	// check partials
	for _, partial := range test.Partials {
		if rAltDelim.MatchString(partial) {
			return true
		}
	}

	return false
}

// func TestMustacheComments(t *testing.T) {
// 	launchRaymondTests(t, testsFromMustacheFile("comments.yml"))
// }

func TestMustacheDelimiters(t *testing.T) {
	launchRaymondTests(t, testsFromMustacheFile("delimiters.yml"))
}

func TestMustacheInterpolation(t *testing.T) {
	launchRaymondTests(t, testsFromMustacheFile("interpolation.yml"))
}

// func TestMustacheInverted(t *testing.T) {
// 	launchRaymondTests(t, testsFromMustacheFile("inverted.yml"))
// }

// func TestMustachePartials(t *testing.T) {
// 	launchRaymondTests(t, testsFromMustacheFile("partials.yml"))
// }

// func TestMustacheSections(t *testing.T) {
// 	launchRaymondTests(t, testsFromMustacheFile("sections.yml"))
// }

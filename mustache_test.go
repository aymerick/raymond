package raymond

import (
	"io/ioutil"
	"path"
	"regexp"
	"strings"
	"testing"

	"gopkg.in/yaml.v2"
)

//
// Note that, as the JS implementation, we do not support:
//   - alternative delimeters
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

func TestMustache(t *testing.T) {
	skipFiles := map[string]bool{
		// mustache lambdas differ from handlebars lambdas
		"~lambdas.yml": true,
	}

	for _, fileName := range mustacheTestFiles() {
		if skipFiles[fileName] {
			// fmt.Printf("Skipped file: %s\n", fileName)
			continue
		}

		launchMustacheTests(t, testsFromMustacheFile(fileName))
	}
}

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
		if mustBeSkipped(mustacheTest, fileName) {
			// fmt.Printf("Skipped test: %s\n", mustacheTest.Name)
			continue
		}

		test := raymondTest{
			name:     mustacheTest.Name,
			input:    mustacheTest.Template,
			data:     mustacheTest.Data,
			partials: mustacheTest.Partials,
			output:   mustacheTest.Expected,
		}

		result = append(result, test)
	}

	return result
}

// returns true if test must be skipped
func mustBeSkipped(test mustacheTest, fileName string) bool {
	// handlebars does not support alternative delimiters
	return haveAltDelimiter(test) ||
		// the JS implementation skips those tests
		fileName == "partials.yml" && (test.Name == "Failed Lookup" || test.Name == "Standalone Indentation")
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

func mustacheTestFiles() []string {
	var result []string

	files, err := ioutil.ReadDir(path.Join("mustache", "specs"))
	if err != nil {
		panic(err)
	}

	for _, file := range files {
		fileName := file.Name()

		if !file.IsDir() && strings.HasSuffix(fileName, ".yml") {
			result = append(result, fileName)
		}
	}

	return result
}

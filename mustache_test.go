package raymond

import (
	"io/ioutil"
	"path"
	"testing"

	"gopkg.in/yaml.v2"
)

// @todo Replace that by adding yaml tags in raymondTest struct
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

// func TestMustacheComments(t *testing.T) {
// 	launchRaymondTests(t, testsFromMustacheFile("comments.yml"))
// }

// func TestMustacheDelimiters(t *testing.T) {
// 	launchRaymondTests(t, testsFromMustacheFile("delimiters.yml"))
// }

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

// func TestMustacheLambdas(t *testing.T) {
// 	launchRaymondTests(t, testsFromMustacheFile("~lambdas.yml"))
// }

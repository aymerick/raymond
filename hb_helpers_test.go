package raymond

import "testing"

//
// Those tests come from:
//   https://github.com/wycats/handlebars.js/blob/master/spec/helper.js
//
var hbHelpersTests = []raymondTest{}

func TestHandlebarsHelpers(t *testing.T) {
	launchRaymondTests(t, hbHelpersTests)
}

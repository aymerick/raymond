package raymond

import "testing"

//
// Those tests come from:
//   https://github.com/wycats/handlebars.js/blob/master/spec/blocks.js
//
var hbBlocksTests = []raymondTest{}

func TestHandlebarsBlocks(t *testing.T) {
	launchRaymondTests(t, hbBlocksTests)
}

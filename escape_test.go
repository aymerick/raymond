package raymond

import "fmt"

func ExampleEscape() {
	tpl := MustParse("{{{link url text}}}")

	tpl.RegisterHelper("link", func(h *HelperArg) interface{} {
		url := Escape(h.ParamStr(0))
		text := Escape(h.ParamStr(1))

		return SafeString("<a href='" + url + "'>" + text + "</a>")
	})

	ctx := map[string]string{
		"url":  "http://www.aymerick.com/",
		"text": "This is a <em>cool</em> website",
	}

	result := tpl.MustExec(ctx)
	fmt.Print(result)
	// Output: <a href='http://www.aymerick.com/'>This is a &lt;em&gt;cool&lt;/em&gt; website</a>
}

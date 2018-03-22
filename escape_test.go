package raymond

import (
	"strings"
	"testing"
)

func TestDefaultEscape(t *testing.T) {
	expected := "<a href='http://www.aymerick.com/'>This is a &lt;em&gt;cool&lt;/em&gt; website</a>"
	tpl := MustParse("{{link url text}}", nil)
	esc := &HTMLEscaper{}

	tpl.RegisterHelper("link", func(url string, text string) SafeString {
		return SafeString("<a href='" + esc.Escape(url) + "'>" + esc.Escape(text) + "</a>")
	})

	ctx := map[string]string{
		"url":  "http://www.aymerick.com/",
		"text": "This is a <em>cool</em> website",
	}

	result := tpl.MustExec(ctx)
	if result != expected {
		t.Errorf("expected %s, got %s", expected, result)
	}
}

type TestEscaper struct{}

func (t TestEscaper) Escape(s string) string {
	return strings.Replace(s, "em", "EM", -1)
}

func TestCustomEscape(t *testing.T) {
	expected := "This is a <EM>cool</EM> website"
	opts := &TemplateOptions{
		Escaper: &TestEscaper{},
	}
	tpl := MustParse("{{ text }}", opts)

	ctx := map[string]string{
		"text": "This is a <em>cool</em> website",
	}

	result := tpl.MustExec(ctx)
	if result != expected {
		t.Errorf("expected %s, got %s", expected, result)
	}
}

# raymond

[Handlebars](http://handlebarsjs.com) for [golang](https://golang.org).


## Todo

- [ ] documentation
- [ ] test with <https://github.com/dvyukov/go-fuzz>
- [ ] check performances


## Quick start

    $ go get github.com/aymerick/raymond

The quick and dirty way of rendering a handlebars template:

```go
package main

import (
    "fmt"

    "github.com/aymerick/raymond"
)

func main() {
    tpl := `<div class="entry">
  <h1>{{title}}</h1>
  <div class="body">
    {{body}}
  </div>
</div>
`

    data := map[string]string{
        "title": "My New Post",
        "body":  "This is my first post!",
    }

    result, err := raymond.Render(tpl, data)
    if err != nil {
        panic("Please fill a bug :)")
    }

    fmt.Print(result)
}
```

Displays:

```html
<div class="entry">
  <h1>My New Post</h1>
  <div class="body">
    This is my first post!
  </div>
</div>
```

Please note that the template will be parsed everytime you call `Render()` function. So you probably want to read the next section.


## Correct usage

To avoid parsing a template several times, use the `Parse()` and `Exec()` function:

```go
package main

import (
    "fmt"

    "github.com/aymerick/raymond"
)

func main() {
    source := `<div class="entry">
  <h1>{{title}}</h1>
  <div class="body">
    {{body}}
  </div>
</div>
`

    dataList := []map[string]string{
        {
            "title": "My New Post",
            "body":  "This is my first post!",
        },
        {
            "title": "Here is another post",
            "body":  "This is my second post!",
        },
    }

    // parse template
    tpl, err := raymond.Parse(source)
    if err != nil {
        panic(err)
    }

    for _, data := range dataList {
        // render template
        result, err := tpl.Exec(data)
        if err != nil {
            panic(err)
        }

        fmt.Print(result)
    }
}

```

Displays:

```html
<div class="entry">
  <h1>My New Post</h1>
  <div class="body">
    This is my first post!
  </div>
</div>
<div class="entry">
  <h1>Here is another post</h1>
  <div class="body">
    This is my second post!
  </div>
</div>
```

You can use `MustParse()` and `MustExec()` functions if you don't want to deal with errors:

```go
    // parse template
    tpl := raymond.MustParse(source)

    // render template
    result := tpl.MustExec(data)
```


## HTML escaping

By default, the result of a mustache expression is HTML escaped. Use the triple mustache `{{{` to output unescaped values.

```go
  source := `<div class="entry">
  <h1>{{title}}</h1>
  <div class="body">
    {{{body}}}
  </div>
</div>
`

    data := map[string]string{
        "title": "All about <p> Tags",
        "body":  "<p>This is a post about &lt;p&gt; tags</p>",
    }

    tpl := raymond.MustParse(source)
    result := tpl.MustExec(data)

    fmt.Print(result)
```

```html
<div class="entry">
  <h1>All about &lt;p&gt; Tags</h1>
  <div class="body">
    <p>This is a post about &lt;p&gt; tags</p>
  </div>
</div>
```

When returning HTML from a helper, you should return a `SafeString` if you don't want it to be escaped by default. When using `SafeStrin`g all unknown or unsafe data should be manually escaped with the `EscapeString` method.

```go
  source := `{{{link text url}}}`

  data := map[string]string{
    "text": "This is a <em>cool</em> website",
    "url":  "http://www.aymerick.com/",
  }

  tpl := raymond.MustParse(source)

  tpl.RegisterHelper("link", func(h *raymond.HelperArg) interface{} {
    text := raymond.EscapeString(h.ParamStr(0))
    url := raymond.EscapeString(h.ParamStr(1))

    return raymond.SafeString("<a href='" + url + "'>" + text + "</a>")
  })

  result := tpl.MustExec(data)
  fmt.Print(result)
```

```html
<a href='http://www.aymerick.com/'>This is a &lt;em&gt;cool&lt;/em&gt; website</a>
```


## Block Expressions

@todo doc


## Handlebars Paths

@todo doc


## Helpers

@todo doc


## Block helpers

@todo doc


## Built-In Helpers

### The `if` block helper

@todo doc

### The `unless` block helper

@todo doc

### The `with` block helper

@todo doc

### The `each` block helper

@todo doc


## Partials

@todo doc


## Mustache

Handlebars is a superset of [mustache](https://mustache.github.io) but it differs on those points:

- Alternative delimiters are not supported
- There is no recursive lookup


## Limitations

These handlebars features are currently not implemented:

- `strict` mode: errors on missing lookup
- `stringParams` mode: resolves a parameter to it's name if the value isn't present in the context stack
- `compat` mode : enables recursive lookup
- `preventIndent` mode: disables indentation of nested partials
- `knownHelpersOnly` mode: allows only known builtin helpers
- `trackIds` mode: informs helpers about the paths that were used to lookup an argument for a given value
- `blockHelperMissing` helper: helper called when a helper can not be directly resolved
- `helperMissing` helper: helper called when a potential helper expression was not found


## Test

    $ go test ./...

    $ go test -run="HandlebarsBasic"


## References

  - <http://handlebarsjs.com/>
  - <https://mustache.github.io/mustache.5.html>
  - <https://github.com/golang/go/tree/master/src/text/template>
  - <https://www.youtube.com/watch?v=HxaD_trXwRE>

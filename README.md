# raymond

Handlebars for [golang](https://golang.org), supporting the same features as [handlebars.js](http://handlebarsjs.com) `3.0`.


## Quick Start

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


## Correct Usage

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


## HTML Escaping

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

Output:

```html
<div class="entry">
  <h1>All about &lt;p&gt; Tags</h1>
  <div class="body">
    <p>This is a post about &lt;p&gt; tags</p>
  </div>
</div>
```

When returning HTML from a helper, you should return a `SafeString` if you don't want it to be escaped by default. When using `SafeString` all unknown or unsafe data should be manually escaped with the `EscapeString` method.

```go
  tpl := raymond.MustParse("{{{link text url}}}")

  tpl.RegisterHelper("link", func(h *raymond.HelperArg) interface{} {
    text := raymond.EscapeString(h.ParamStr(0))
    url := raymond.EscapeString(h.ParamStr(1))

    return raymond.SafeString("<a href='" + url + "'>" + text + "</a>")
  })

  data := map[string]string{
    "text": "This is a <em>cool</em> website",
    "url":  "http://www.aymerick.com/",
  }

  result := tpl.MustExec(data)
  fmt.Print(result)
```

Output:

```html
<a href='http://www.aymerick.com/'>This is a &lt;em&gt;cool&lt;/em&gt; website</a>
```


## Helpers

@todo doc

### Private data

@todo doc


## Partials

### Template partials

You can register template partials before execution:

```go
  tpl := raymond.MustParse("{{> foo}} baz")
  tpl.RegisterPartial("foo", "<span>bar</span>")

  result := tpl.MustExec(nil)
  fmt.Print(result)
```

Output:

```html
<span>bar</span> baz
```

You can register several partials at once:

```go
tpl := raymond.MustParse("{{> foo}} and {{> baz}}")
tpl.RegisterPartials(map[string]string{
  "foo": "<span>bar</span>",
  "baz": "<span>bat</span>",
})

result := tpl.MustExec(nil)
fmt.Print(result)
```

Output:

```html
<span>bar</span> and <span>bat</span>
```

### Global Partials

You can registers global partials that will be accessible by all templates:

```go
  raymond.RegisterPartial("foo", "<span>bar</span>")

  tpl := raymond.MustParse("{{> foo}} baz")
  result := tpl.MustExec(nil)
  fmt.Print(result)
```

Or:

```go
  raymond.RegisterPartials(map[string]string{
    "foo": "<span>bar</span>",
    "baz": "<span>bat</span>",
  })

  tpl := raymond.MustParse("{{> foo}} and {{> baz}}")
  result := tpl.MustExec(nil)
  fmt.Print(result)
```

### Dynamic Partials

It's possible to dynamically select the partial to be executed by using sub expression syntax.

For example, that template randomly evaluates the `foo` or `baz` partial:

```go
  tpl := raymond.MustParse("{{> (whichPartial) }}")
  tpl.RegisterPartials(map[string]string{
    "foo": "<span>bar</span>",
    "baz": "<span>bat</span>",
  })

  data := map[string]interface{}{
    "whichPartial": func() string {
      rand.Seed(time.Now().UTC().UnixNano())

      names := []string{"foo", "baz"}
      return names[rand.Intn(len(names))]
    },
  }

  result := tpl.MustExec(data)
  fmt.Print(result)
```

### Partial Contexts

It's possible to execute partials on a custom context by passing in the context to the partial call.

For example:

```go
  tpl := raymond.MustParse("User: {{> userDetails user }}")
  tpl.RegisterPartial("userDetails", "{{firstname}} {{lastname}}")

  data := map[string]interface{}{
    "user": map[string]string{
      "firstname": "Jean",
      "lastname":  "Valjean",
    },
  }

  result := tpl.MustExec(data)
  fmt.Print(result)
```

Displays:

```html
User: Jean Valjean
```

### Partial Parameters

Custom data can be passed to partials through hash parameters.

For example:

```go
  tpl := raymond.MustParse("{{> myPartial name=hero }}")
  tpl.RegisterPartial("myPartial", "his name is: {{name}}")

  data := map[string]interface{}{
    "hero": "Goldorak",
  }

  result := tpl.MustExec(data)
  fmt.Print(result)
```

Displays:

```html
his name is: Goldorak
```

## Utilities

@todo doc for `IsSafeString()`

@todo doc for `Str()`


## Mustache

Handlebars is a superset of [mustache](https://mustache.github.io) but it differs on those points:

- Alternative delimiters are not supported
- There is no recursive lookup


## Limitations

These handlebars options are currently NOT implemented:

- `compat` - enables recursive field lookup
- `knownHelpers` - list of helpers that are known to exist (truthy) at template execution time
- `knownHelpersOnly` - allows further optimzations based on the known helpers list
- `trackIds` - include the id names used to resolve parameters for helpers
- `noEscape` - disables HTML escaping globally
- `strict` - templates will throw rather than silently ignore missing fields
- `assumeObjects` - removes object existence checks when traversing paths
- `preventIndent` - disables the auto-indententation of nested partials
- `stringParams` - resolves a parameter to it's name if the value isn't present in the context stack

These handlebars features are currently NOT implemented:

- Function in private data frame
- `blockHelperMissing` - helper called when a helper can not be directly resolved
- `helperMissing` - helper called when a potential helper expression was not found
- `@contextPath` - value set in `trackIds` mode that records the lookup path for the current context
- `@level` - log level


## Todo

- [ ] function in private data frame
- [ ] test with <https://github.com/dvyukov/go-fuzz>
- [ ] benchmarks


## Test

    $ go test ./...

    $ go test -run="HandlebarsBasic"


## References

  - <http://handlebarsjs.com/>
  - <https://mustache.github.io/mustache.5.html>
  - <https://github.com/golang/go/tree/master/src/text/template>
  - <https://www.youtube.com/watch?v=HxaD_trXwRE>


## Others Implementations

- [handlebars.js](http://handlebarsjs.com) - javascript
- [handlebars.java](https://github.com/jknack/handlebars.java) - java
- [handlebars.rb](https://github.com/cowboyd/handlebars.rb) - ruby
- [handlebars.php](https://github.com/XaminProject/handlebars.php) - php
- [handlebars-objc](https://github.com/Bertrand/handlebars-objc) - Objective C
- [rumblebars](https://github.com/nicolas-cherel/rumblebars) - rust

# raymond

Handlebars for golang

**WARNING: This is a work in progress, some features are missing.**


## Todo

- [ ] `strict` mode
- [ ] `stringParams` mode
- [ ] `compat` mode
- [ ] `preventIndent` mode
- [ ] permits helpers to return safe strings
- [ ] the `lookup` helper
- [ ] the `log` helper
- [ ] pass all handlebars.js tests
- [ ] documentation
- [ ] test with <https://github.com/dvyukov/go-fuzz>


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
package main

import (
    "fmt"

    "github.com/aymerick/raymond"
)

func main() {
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
}

```

```html
<div class="entry">
  <h1>All about &lt;p&gt; Tags</h1>
  <div class="body">
    <p>This is a post about &lt;p&gt; tags</p>
  </div>
</div>
```

@todo How a helper can return a safe string ?


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


## Limitations

@todo doc


## Test

    $ go test ./...

    $ go test -run="HandlebarsBasic"


## References

  - <http://handlebarsjs.com/>
  - <https://mustache.github.io/mustache.5.html>
  - <https://github.com/golang/go/tree/master/src/text/template>
  - <https://www.youtube.com/watch?v=HxaD_trXwRE>

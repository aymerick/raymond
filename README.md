# raymond [![Build Status](https://secure.travis-ci.org/aymerick/raymond.svg?branch=master)](http://travis-ci.org/aymerick/raymond) [![GoDoc](https://godoc.org/github.com/aymerick/raymond?status.svg)](http://godoc.org/github.com/aymerick/raymond)

Handlebars for [golang](https://golang.org) with the same features as [handlebars.js](http://handlebarsjs.com) `3.0`.

The full API documentation is available here: <http://godoc.org/github.com/aymerick/raymond>.


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

    ctx := map[string]string{
        "title": "My New Post",
        "body":  "This is my first post!",
    }

    result, err := raymond.Render(tpl, ctx)
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

    ctxList := []map[string]string{
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

    for _, ctx := range ctxList {
        // render template
        result, err := tpl.Exec(ctx)
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
    result := tpl.MustExec(ctx)
```


## Context

The rendering context can contain any type of objects, including `array`, `slice`, `map`, `struct` and `func`.

When using structs, be warned that only exported fields are accessible:

```go
package main

import (
  "fmt"

  "github.com/aymerick/raymond"
)

func main() {
  source := `<div class="post">
  <h1>By {{Author.FirstName}} {{Author.LastName}}</h1>
  <div class="body">{{Body}}</div>

  <h1>Comments</h1>

  {{#each Comments}}
  <h2>By {{Author.FirstName}} {{Author.LastName}}</h2>
  <div class="body">{{Body}}</div>
  {{/each}}
</div>`

  type Person struct {
    FirstName string
    LastName  string
  }

  type Comment struct {
    Author Person
    Body   string
  }

  type Post struct {
    Author   Person
    Body     string
    Comments []Comment
  }

  ctx := Post{
    Person{"Jean", "Valjean"},
    "Life is difficult",
    []Comment{
      Comment{
        Person{"Marcel", "Beliveau"},
        "LOL!",
      },
    },
  }

  output := raymond.MustRender(source, ctx)

  fmt.Print(output)
}
```

Output:

```html
<div class="post">
  <h1>By Jean Valjean</h1>
  <div class="body">Life is difficult</div>

  <h1>Comments</h1>

  <h2>By Marcel Beliveau</h2>
  <div class="body">LOL!</div>
</div>
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

  ctx := map[string]string{
    "title": "All about <p> Tags",
    "body":  "<p>This is a post about &lt;p&gt; tags</p>",
  }

  tpl := raymond.MustParse(source)
  result := tpl.MustExec(ctx)

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

When returning HTML from a helper, you should return a `SafeString` if you don't want it to be escaped by default. When using `SafeString` all unknown or unsafe data should be manually escaped with the `Escape` method.

```go
  raymond.RegisterHelper("link", func(url, text string) raymond.SafeString {
    return raymond.SafeString("<a href='" + raymond.Escape(url) + "'>" + raymond.Escape(text) + "</a>")
  })

  tpl := raymond.MustParse("{{link url text}}")

  ctx := map[string]string{
    "url":  "http://www.aymerick.com/",
    "text": "This is a <em>cool</em> website",
  }

  result := tpl.MustExec(ctx)
  fmt.Print(result)
```

Output:

```html
<a href='http://www.aymerick.com/'>This is a &lt;em&gt;cool&lt;/em&gt; website</a>
```


## Helpers

Helpers can be accessed from any context in a template. You can register a helper with the `RegisterHelper` function.

For example:

```html
<div class="post">
  <h1>By {{fullName author}}</h1>
  <div class="body">{{body}}</div>

  <h1>Comments</h1>

  {{#each comments}}
  <h2>By {{fullName author}}</h2>
  <div class="body">{{body}}</div>
  {{/each}}
</div>
```

With this context and helper:

```go
ctx := map[string]interface{}{
  "author": map[string]string{"firstName": "Jean", "lastName": "Valjean"},
  "body":   "Life is difficult",
  "comments": []map[string]interface{}{{
    "author": map[string]string{"firstName": "Marcel", "lastName": "Beliveau"},
    "body":   "LOL!",
  }},
}

raymond.RegisterHelper("fullName", func(person map[string]string) string {
  return person["firstName"] + " " + person["lastName"]
})
```

Outputs:

```html
<div class="post">
  <h1>By Jean Valjean</h1>
  <div class="body">Life is difficult</div>

  <h1>Comments</h1>

  <h2>By Marcel Beliveau</h2>
  <div class="body">LOL!</div>
</div>
```

Helper arguments can be any type. Following example uses structs instead of maps and produces the same output as the previous one:

```go
  type Person struct {
    FirstName string
    LastName  string
  }

  type Comment struct {
    Author Person
    Body   string
  }

  type Post struct {
    Author   Person
    Body     string
    Comments []Comment
  }

  ctx := Post{
    Person{"Jean", "Valjean"},
    "Life is difficult",
    []Comment{
      Comment{
        Person{"Marcel", "Beliveau"},
        "LOL!",
      },
    },
  }

  RegisterHelper("fullName", func(person Person) string {
    return person.FirstName + " " + person.LastName
  })
```


### Template Helpers

You can register a helper on a specific template, and in that case that helper will be only available to that template:

```go
  tpl := raymond.MustParse("User: {{fullName user.firstName user.lastName}}")

  tpl.RegisterHelper("fullName", func(firstName, lastName string) string {
    return firstName + " " + lastName
  })
```


### Built-In Helpers

Those built-in helpers are available to all templates.


#### The `if` block helper

You can use the `if` helper to conditionally render a block. If its argument returns `false`, `nil`, `0`, `""`, an empty array or an empty map, then raymond will not render the block.

```html
<div class="entry">
  {{#if author}}
    <h1>{{firstName}} {{lastName}}</h1>
  {{/if}}
</div>
```

When using a block expression, you can specify a template section to run if the expression returns a falsy value. That section, marked by `{{else}}` is called an "else section".

```html
<div class="entry">
  {{#if author}}
    <h1>{{firstName}} {{lastName}}</h1>
  {{else}}
    <h1>Unknown Author</h1>
  {{/if}}
</div>
```


#### The `unless` block helper

You can use the `unless` helper as the inverse of the `if` helper. Its block will be rendered if the expression returns a falsy value.

```html
<div class="entry">
  {{#unless license}}
  <h3 class="warning">WARNING: This entry does not have a license!</h3>
  {{/unless}}
</div>
```


#### The `each` block helper

You can iterate over an array, a map or a struct instance using the built-in each helper. Inside the block, you can use `this` to reference the element being iterated over.

For example:

```html
<ul class="people">
  {{#each people}}
    <li>{{this}}</li>
  {{/each}}
</ul>
```

With this context:

```go
map[string]interface{}{
  "people": []string{
    "Marcel", "Jean-Claude", "Yvette",
  },
}
```

Outputs:

```html
<ul class="people">
    <li>Marcel</li>
    <li>Jean-Claude</li>
    <li>Yvette</li>
</ul>
```

You can optionally provide an `{{else}}` section which will display only when the passed argument is empty.

```html
{{#each paragraphs}}
  <p>{{this}}</p>
{{else}}
  <p class="empty">No content</p>
{{/each}}
```

When looping through items in `each`, you can optionally reference the current loop index via `{{@index}}`.

```html
{{#each array}}
  {{@index}}: {{this}}
{{/each}}
```

Additionally for map and struct instance iteration, `{{@key}}` references the current key name:

```html
{{#each map}}
  {{@key}}: {{this}}
{{/each}}
```

The first and last steps of iteration are noted via the `@first` and `@last` variables.


#### The `with` block helper

You can shift the context for a section of a template by using the built-in `with` block helper.

```html
<div class="entry">
  <h1>{{title}}</h1>

  {{#with author}}
  <h2>By {{firstName}} {{lastName}}</h2>
  {{/with}}
</div>
```

With this context:

```go
  map[string]interface{}{
    "title": "My first post!",
    "author": map[string]string{
      "firstName": "Jean",
      "lastName":  "Valjean",
    },
  }
```

Outputs:

```html
<div class="entry">
  <h1>My first post!</h1>

  <h2>By Jean Valjean</h2>
</div>
```

You can optionally provide an `{{else}}` section which will display only when the passed argument is empty.

```html
{{#with author}}
  <p>{{name}}</p>
{{else}}
  <p class="empty">No content</p>
{{/with}}
```


#### The `lookup` helper

The `lookup` helper allows for dynamic parameter resolution using handlebars variables.

```html
{{#each bar}}
  {{lookup ../foo @index}}
{{/each}}
```


#### The `log` helper

The `log` helper allows for logging while evaluating a template.

```html
{{log "Look at me!"}}
```

Note that the handlebars.js `@level` variable is not supported.


### Block Helpers

Block helpers make it possible to define custom iterators and other functionality that can invoke the passed block with a new context.


#### Block Evaluation

As an example, let's define a block helper that adds some markup to the wrapped text.

```html
<div class="entry">
  <h1>{{title}}</h1>
  <div class="body">
    {{#bold}}{{body}}{{/bold}}
  </div>
</div>
```

The `bold` helper will add markup to make its text bold.

```go
  raymond.RegisterHelper("bold", func(options *raymond.Options) raymond.SafeString {
    return raymond.SafeString(`<div class="mybold">` + options.Fn() + "</div>")
  })
```

As you can see, an helper evaluates the block content by calling `options.Fn()`.

If you want to evaluate the block with another context, then use `options.FnWithCtx(ctx)`, like this french version of built-in `with` block helper:

```go
  raymond.RegisterHelper("avec", func(context interface{}, options *raymond.Options) string {
    return options.FnWithCtx(context)
  })
```

With that template:

```html
{{#avec obj.text}}{{this}}{{/avec}}
```


#### Conditional

Let's write a french version of `if` block helper:

```go
  source := `{{#si yep}}YEP !{{/si}}`

  ctx := map[string]interface{}{"yep": true}

  raymond.RegisterHelper("si", func(conditional bool, options *raymond.Options) string {
    if conditional {
      return options.Fn()
    }
    return ""
  })
```

Note that as the first parameter of the helper is typed as `bool` an automatic conversion is made if corresponding context value is not a boolean. So this helper works with that context too:

```go
  ctx := map[string]interface{}{"yep": "message"}
```

See `IsTruth()` function for more informations on boolean conversion.


#### Else Block Evaluation

We can enhance the `si` block helper to evaluate the `else block` by calling `options.Inverse()` if conditional is false:

```go
  source := `{{#si yep}}YEP !{{else}}NOP !{{/si}}`

  ctx := map[string]interface{}{"yep": false}

  raymond.RegisterHelper("si", func(conditional bool, options *raymond.Options) string {
    if conditional {
      return options.Fn()
    }
    return options.Inverse()
  })
```

Outputs:
```
NOP !
```


#### Block Parameters

It's possible to receive named parameters from supporting helpers

```html
  {{#each users as |user userId|}}
    Id: {{userId}} Name: {{user.name}}
  {{/each}}
```

In this particular example, `user` will have the same value as the current context and `userId` will have the index value for the iteration.

This allows for nested helpers to avoid name conflicts that can occur with private variables.

For example:

```html
{{#each users as |user userId|}}
  {{#each user.book as |book bookId|}}
    User: {{userId}} Book: {{bookId}}
  {{/each}}
{{/each}}
```

With this context:

```go
  ctx := map[string]interface{}{
    "users": map[string]interface{}{
      "marcel": map[string]interface{}{
        "book": map[string]interface{}{
          "book1": "My first book",
          "book2": "My second book",
        },
      },
      "didier": map[string]interface{}{
        "book": map[string]interface{}{
          "bookA": "Good book",
          "bookB": "Bad book",
        },
      },
    },
  }
```

Outputs:

```html
  User: marcel Book: book1
  User: marcel Book: book2
  User: didier Book: bookA
  User: didier Book: bookB
```

As you can see, the second block parameter is the map key. When using structs, it is the struct field name.

When using arrays the second parameter is element index:

```go
  ctx := map[string]interface{}{
    "users": []map[string]interface{}{
      {
        "id": "marcel",
        "book": []map[string]interface{}{
          {"id": "book1", "title": "My first book"},
          {"id": "book2", "title": "My second book"},
        },
      },
      {
        "id": "didier",
        "book": []map[string]interface{}{
          {"id": "bookA", "title": "Good book"},
          {"id": "bookB", "title": "Bad book"},
        },
      },
    },
  }
```

Outputs:

```html
    User: 0 Book: 0
    User: 0 Book: 1
    User: 1 Book: 0
    User: 1 Book: 1
```

### Helper Parameters

@todo doc

@todo doc automatique string conversion


### Options Argument

@todo doc


### Helper Hash Arguments

@todo doc


### Private Data

@todo doc


### Utilites

@todo doc for `Str()`

@todo doc for `IsTruth()`... describes boolean conversion


## Context Functions

@todo doc


## Partials

### Template Partials

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

  ctx := map[string]interface{}{
    "whichPartial": func() string {
      rand.Seed(time.Now().UTC().UnixNano())

      names := []string{"foo", "baz"}
      return names[rand.Intn(len(names))]
    },
  }

  result := tpl.MustExec(ctx)
  fmt.Print(result)
```


### Partial Contexts

It's possible to execute partials on a custom context by passing in the context to the partial call.

For example:

```go
  tpl := raymond.MustParse("User: {{> userDetails user }}")
  tpl.RegisterPartial("userDetails", "{{firstname}} {{lastname}}")

  ctx := map[string]interface{}{
    "user": map[string]string{
      "firstname": "Jean",
      "lastname":  "Valjean",
    },
  }

  result := tpl.MustExec(ctx)
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

  ctx := map[string]interface{}{
    "hero": "Goldorak",
  }

  result := tpl.MustExec(ctx)
  fmt.Print(result)
```

Displays:

```html
his name is: Goldorak
```


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

- raw block content is not passed as a parameter to helper
- `blockHelperMissing` - helper called when a helper can not be directly resolved
- `helperMissing` - helper called when a potential helper expression was not found
- `@contextPath` - value set in `trackIds` mode that records the lookup path for the current context
- `@level` - log level


## Todo

- [ ] add a test for inverse statement with the `each` helper
- [ ] test with <https://github.com/dvyukov/go-fuzz>
- [ ] benchmarks


## Handlebars Lexer

You should not use the lexer directly, but for your information here is an example:

```go
package main

import (
    "fmt"

    "github.com/aymerick/raymond/lexer"
)

func main() {
  source := "You know {{nothing}} John Snow"

  output := ""

  lex := lexer.Scan(source)
  for {
    // consume next token
    token := lex.NextToken()

    output += fmt.Sprintf(" %s", token)

    // stops when all tokens have been consumed, or on error
    if token.Kind == lexer.TokenEOF || token.Kind == lexer.TokenError {
      break
    }
  }

  fmt.Print(output)
}
```

Outputs:

```
Content{"You know "} Open{"{{"} ID{"nothing"} Close{"}}"} Content{" John Snow"} EOF
```


## Handlebars Parser

You should not use the parser directly, but for your information here is an example:

```go
package main

import (
  "fmt"

  "github.com/aymerick/raymond/ast"
  "github.com/aymerick/raymond/parser"
)

func main() {
  source := "You know {{nothing}} John Snow"

  // parse template
  program, err := parser.Parse(source)
  if err != nil {
    panic(err)
  }

  // print AST
  output := ast.Print(program)

  fmt.Print(output)
}
```

Outputs:

```
CONTENT[ 'You know ' ]
{{ PATH:nothing [] }}
CONTENT[ ' John Snow' ]
```


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

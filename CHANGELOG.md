# Raymond Changelog

### Raymond 1.1.0 _(June 15, 2015)_

- Permits templates references with lowercase versions of struct fields.
- Adds `ParseFile()` function.
- Adds `RegisterPartialFile()`, `RegisterPartialFiles()` and `Clone()` methods on `Template`.
- Helpers can now be struct methods.
- Ensures safe concurrent access to helpers and partials.

### Raymond 1.0.0 _(June 09, 2015)_

- This is the first release. Raymond supports almost all handlebars features. See https://github.com/aymerick/raymond#limitations for a list of differences with the javascript implementation.

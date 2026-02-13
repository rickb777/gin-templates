# gin-templates

[![GoDoc](https://img.shields.io/badge/api-Godoc-blue.svg)](https://pkg.go.dev/github.com/rickb777/gin-templates)
[![Go Report Card](https://goreportcard.com/badge/github.com/rickb777/gin-templates)](https://goreportcard.com/report/github.com/rickb777/gin-templates)
[![Issues](https://img.shields.io/github/issues/rickb777/gin-templates.svg)](https://github.com/rickb777/gin-templates/issues)

* Loads a tree of HTML templates into [Gin](https://github.com/gin-gonic/gin)
* Allows many templates from a directory tree

## Installation

    go get -u github.com/rickb777/gin-templates

## Usage

Given a directory `templates` containing your templates and given that their names all end in ".html", use this setup

```go
engine := gin.New()
engine.FuncMap = (... as needed ...)
engine.HTMLRender = gin_templates.LoadTemplates("templates", ".html", engine.FuncMap)
```

This scans all ".html" files in `templates` and its subdirectories. So, for example, if files `templates/foo/home.html` and `templates/foo/bar/baz.html` exist, when loaded they will have the names `foo/home.html` and `foo/bar/baz.html` respectively.

As usual, your handlers will use `Context.HTML()`, e.g

```go
c.HTML(200, "foo/home.html", dataModel)
```
will execute the template with the path `foo/home.html`.

It keeps the same  `DebugMode` and `ReleaseMode` behaviour as the core Gin rendering. So in `DebugMode`, templates are hot-reloaded every time they are used.

## Rationale

Go templates work very well and are widely used. There is, however, an unfortunate feature of the main `template.ParseFiles` and `template.ParseGlob` functions only work well for small numbers of templates. The name they store of each template is based on the filename ignoring its path; therefore any name collisions cause templates to be unusable. For a larger number of templates, this problem becomes increasingly likely.

So `gin-templates` provides a way to load templates without being affected by this issue. The template files can be organised in directories as needed and their names will reflect their path relative to the root directory of the templates.

## Status

This library is in early development.

## Credits

This package was inspired by [multitemplate](https://github.com/gin-contrib/multitemplate).
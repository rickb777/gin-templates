# gin-templates

[![GoDoc](https://img.shields.io/badge/api-Godoc-blue.svg)](https://pkg.go.dev/github.com/rickb777/gin-templates)
[![Build Status](https://api.travis-ci.org/rickb777/gin-templates.svg?branch=master)](https://travis-ci.org/rickb777/gin-templates/builds)
[![Coverage Status](https://coveralls.io/repos/rickb777/gin-templates/badge.svg?branch=master&service=github)](https://coveralls.io/github/rickb777/gin-templates?branch=master)
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
	engine.HTMLRender := gin_templates.LoadTemplates(engine, "templates", ".html")
```

It keeps the same  `DebugMode` and `ReleaseMode` behaviour as the core Gin rendering. So in `DebugMode`, templates are hot-reloaded every time they are used.

## Rationale

Go templates work very well and are widely used. There is, however, an unfortunate feature of the main `template.ParseFiles` and `template.ParseGlob` functions only work well for small numbers of templates. The name they store of each template is based on the filename ignoring its path; therefore any name collisions cause templates to be unusable.

So `gin-templates` provides a way to load templates without being affected by this issue. The template files can be organised in directories as needed and their names will reflect their path relative to the root directory of the templates.

For example, if files `foo/home.html` and `foo/bar/baz.html` exist, when loaded they will have the same names, i.e. `foo/home.html` and `foo/bar/baz.html` respectively.

## Status

This library is in early development.

## Credits

This package was inspired by [multitemplate](https://github.com/gin-contrib/multitemplate).
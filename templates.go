package templates

import (
	"fmt"
	"html/template"
	"net/http"
	"os"
	"path/filepath"
	"strings"
	"time"

	"github.com/gin-gonic/gin"
	"github.com/gin-gonic/gin/render"
	"github.com/spf13/afero"
)

var Fs = afero.NewOsFs()

// ResponseProcessor interface creates the contract for custom content negotiation.
// This matches processor.ResponseProcessor (github.com/rickb777/negotiator) without
// introducing an explicit dependency.
type ResponseProcessor interface {
	// CanProcess is the predicate that determines whether this processor
	// will handle a given request.
	CanProcess(mediaRange string, lang string) bool
	// ContentType returns the content type for this response.
	ContentType() string
	// Process renders the data model to the response writer, without setting any headers.
	// If the processor encounters an error, it should panic.
	Process(w http.ResponseWriter, template string, dataModel interface{}) error
}

type HTMLProcessor interface {
	render.HTMLRender
	ResponseProcessor
}

// LoadTemplates finds all the templates in the directory dir and its subdirectories
// that have names ending with the given suffix. The function map can be nil if not
// required.
func LoadTemplates(dir, suffix string, funcMap template.FuncMap) HTMLProcessor {
	if funcMap == nil {
		funcMap = template.FuncMap{}
	}

	rootDir := filepath.Clean(dir)

	files := findTemplates(rootDir, suffix)

	if len(files) == 0 {
		panic("No HTML files were found in " + rootDir)
	}

	root := parseTemplates(rootDir, files, funcMap)

	if gin.IsDebugging() {
		return &processor{
			HTMLRender: htmlDebug{
				root:    root,
				rootDir: rootDir,
				suffix:  suffix,
				files:   files,
				funcMap: funcMap,
			},
		}
	}

	return &processor{HTMLRender: htmlProduction{root: root}}
}

//-------------------------------------------------------------------------------------------------

func findTemplates(rootDir, suffix string) map[string]time.Time {
	cleanRoot := filepath.Clean(rootDir)
	files := make(map[string]time.Time)

	err := afero.Walk(Fs, cleanRoot, func(path string, info os.FileInfo, e1 error) error {
		if e1 != nil {
			panic(fmt.Sprintf("Cannot load templates from: %s: %v\n", rootDir, e1))
		}

		if !info.IsDir() && strings.HasSuffix(path, suffix) {
			files[path] = time.Time{}
		}

		return nil
	})

	if err != nil {
		panic(fmt.Sprintf("Cannot load templates from: %s: %v\n", rootDir, err))
	}

	return files
}

func parseTemplates(rootDir string, files map[string]time.Time, funcMap template.FuncMap) *template.Template {
	pfx := len(rootDir) + 1
	root := template.New("")

	for path := range files {
		b, e2 := afero.ReadFile(Fs, path)
		if e2 != nil {
			panic(fmt.Sprintf("Read template error: %s: %v\n", path, e2))
		}

		name := path[pfx:]
		t := root.New(name).Funcs(funcMap)
		t, e2 = t.Parse(string(b))
		if e2 != nil {
			panic(fmt.Sprintf("Parse template error: %s: %v\n", path, e2))
		}
	}

	return root
}

//-------------------------------------------------------------------------------------------------

const (
	TextHtml         = "text/html"
	ApplicationXhtml = "application/xhtml+xml"
)

// ContentType is the value returned when serving HTML files. This defaults to
// TextHtml, but set it to ApplicationXhtml if more appropriate.
var ContentType = TextHtml

// processor adds methods to a HTMLRender to allow it to take part in content negotiation.
type processor struct {
	render.HTMLRender
}

func (p processor) CanProcess(mediaRange string, lang string) bool {
	return mediaRange == TextHtml || mediaRange == ApplicationXhtml
}

func (p processor) Process(w http.ResponseWriter, template string, dataModel interface{}) error {
	return p.Instance(template, dataModel).Render(w)
}

func (p processor) ContentType() string {
	return ContentType
}

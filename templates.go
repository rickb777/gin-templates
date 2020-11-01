package templates

import (
	"fmt"
	"html/template"
	"os"
	"path/filepath"
	"strings"
	"time"

	"github.com/gin-gonic/gin"
	"github.com/gin-gonic/gin/render"
	"github.com/spf13/afero"
)

var Fs = afero.NewOsFs()

type htmlProduction struct {
	root *template.Template
}

type htmlDebug struct {
	root    *template.Template
	rootDir string
	suffix  string
	files   map[string]time.Time
	funcMap template.FuncMap
}

// LoadTemplates finds all the templates in the directory dir and its subdirectories
// that have names ending with the given suffix. The function map can be nil if not
// required.
func LoadTemplates(dir, suffix string, funcMap template.FuncMap) render.HTMLRender {
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
		return htmlDebug{
			root:    root,
			rootDir: rootDir,
			suffix:  suffix,
			files:   files,
			funcMap: funcMap,
		}
	}

	return htmlProduction{root: root}
}

func (r htmlProduction) Instance(name string, data interface{}) render.Render {
	return render.HTML{
		Template: r.root,
		Name:     name,
		Data:     data,
	}
}

func (r htmlDebug) Instance(name string, data interface{}) render.Render {
	path := r.rootDir + "/" + name
	if _, exists := r.files[path]; !exists {
		r.files = findTemplates(r.rootDir, r.suffix)
	}
	return render.HTML{
		Template: r.getCurrentTemplateTree(),
		Name:     name,
		Data:     data,
	}
}

func (r htmlDebug) getCurrentTemplateTree() *template.Template {
	changed := checkForChanges(r.files)
	if changed {
		r.root = parseTemplates(r.rootDir, r.files, r.funcMap)
	}
	return r.root
}

func checkForChanges(files map[string]time.Time) bool {
	changed := false
	for path, modTime := range files {
		fi, err := Fs.Stat(path)
		if err == nil {
			if fi.ModTime().After(modTime) {
				files[path] = fi.ModTime()
				changed = true
			}
		} else if !os.IsNotExist(err) {
			delete(files, path)
		}

	}

	return changed
}

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

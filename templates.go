package templates

import (
	"github.com/gin-gonic/gin"
	"github.com/gin-gonic/gin/render"
	"html/template"
	"io/ioutil"
	"log"
	"os"
	"path/filepath"
	"strings"
	"time"
)

type htmlProduction struct {
	root *template.Template
}

type htmlDebug struct {
	root    *template.Template
	rootDir string
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

	root := template.New("")
	parseTemplates(root, rootDir, files, funcMap)

	if gin.IsDebugging() {
		return htmlDebug{rootDir: rootDir, root: root, files: files, funcMap: funcMap}
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
	return render.HTML{
		Template: r.loadTemplate(),
		Name:     name,
		Data:     data,
	}
}

func (r htmlDebug) loadTemplate() *template.Template {
	changed := checkForChanges(r.files)
	return parseTemplates(r.root, r.rootDir, changed, r.funcMap)
}

func checkForChanges(files map[string]time.Time) map[string]time.Time {
	changed := make(map[string]time.Time)
	for path, modTime := range files {
		fi, err := os.Stat(path)
		if err == nil {
			if fi.ModTime().After(modTime) {
				files[path] = fi.ModTime()
				changed[path] = fi.ModTime()
			}
		} else {
			log.Printf("%q err %v\n", path, err)
		}
	}

	return changed
}

func findTemplates(rootDir, suffix string) map[string]time.Time {
	cleanRoot := filepath.Clean(rootDir)
	files := make(map[string]time.Time)

	err := filepath.Walk(cleanRoot, func(path string, info os.FileInfo, e1 error) error {
		if e1 != nil {
			panic("cannot load templates from " + rootDir + ": " + e1.Error())
		}

		if !info.IsDir() && strings.HasSuffix(path, suffix) {
			files[path] = time.Time{}
		}

		return nil
	})

	if err != nil {
		panic("cannot load templates from " + rootDir + ": " + err.Error())
	}

	return files
}

func parseTemplates(root *template.Template, rootDir string, files map[string]time.Time, funcMap template.FuncMap) *template.Template {
	pfx := len(rootDir) + 1

	for path := range files {
		b, e2 := ioutil.ReadFile(path)
		if e2 != nil {
			panic(path + " " + e2.Error())
		}

		name := path[pfx:]
		t := root.New(name).Funcs(funcMap)
		t, e2 = t.Parse(string(b))
		if e2 != nil {
			panic(path + " " + e2.Error())
		}
	}

	return root
}

package templates

import (
	"github.com/gin-gonic/gin/render"
	"html/template"
	"os"
	"time"
)

type htmlProduction struct {
	root *template.Template
}

func (r htmlProduction) Instance(name string, data interface{}) render.Render {
	return render.HTML{
		Template: r.root,
		Name:     name,
		Data:     data,
	}
}

//-------------------------------------------------------------------------------------------------

type htmlDebug struct {
	root    *template.Template
	rootDir string
	suffix  string
	files   map[string]time.Time
	funcMap template.FuncMap
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

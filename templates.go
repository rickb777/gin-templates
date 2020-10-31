package gin_templates

import (
	"github.com/gin-gonic/gin"
	"github.com/gin-gonic/gin/render"
	"html/template"
	"io/ioutil"
	"os"
	"path/filepath"
	"strings"
)

type htmlProduction struct {
	root *template.Template
}

type htmlDebug struct {
	rootDir string
	files   map[string]os.FileInfo
	funcMap template.FuncMap
}

func LoadTemplates(engine *gin.Engine, dir, suffix string) render.HTMLRender {
	funcMap := engine.FuncMap
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
		return htmlDebug{rootDir: rootDir, files: files, funcMap: engine.FuncMap}
	}

	return htmlProduction{root: root}
}

func (r htmlDebug) Instance(name string, data interface{}) render.Render {
	return render.HTML{
		Template: r.loadTemplate(),
		Name:     name,
		Data:     data,
	}
}

func (r htmlDebug) loadTemplate() *template.Template {
	// TODO filter out files that haven't changed
	return parseTemplates(template.New(""), r.rootDir, r.files, r.funcMap)
}

func (r htmlProduction) Instance(name string, data interface{}) render.Render {
	return render.HTML{
		Template: r.root,
		Name:     name,
		Data:     data,
	}
}

func findTemplates(rootDir, suffix string) map[string]os.FileInfo {
	cleanRoot := filepath.Clean(rootDir)
	files := make(map[string]os.FileInfo)

	err := filepath.Walk(cleanRoot, func(path string, info os.FileInfo, e1 error) error {
		if e1 != nil {
			panic("cannot load templates from " + rootDir + ": " + e1.Error())
		}

		if !info.IsDir() && strings.HasSuffix(path, suffix) {
			files[path] = info
		}

		return nil
	})

	if err != nil {
		panic("cannot load templates from " + rootDir + ": " + err.Error())
	}

	return files
}

func parseTemplates(root *template.Template, rootDir string, files map[string]os.FileInfo, funcMap template.FuncMap) *template.Template {
	pfx := len(rootDir) + 1

	for path, _ := range files {
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

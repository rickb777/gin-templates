package templates_test

import (
	"github.com/gin-gonic/gin"
	"github.com/rickb777/gin-templates"
	"net/http/httptest"
	"strings"
	"testing"
)

func TestDebugInstance(t *testing.T) {
	engine := gin.New()
	r := templates.LoadTemplates("test-data", ".html", engine.FuncMap)
	w := httptest.NewRecorder()

	err := r.Instance("foo/home.html", map[string]string{"Title": "Hello"}).Render(w)
	must(t, err)

	s := w.Body.String()
	mustContain(t, s, "Hello")
	mustContain(t, s, "Home")
}

func TestProductionInstance(t *testing.T) {
	gin.SetMode(gin.ReleaseMode)
	engine := gin.New()
	r := templates.LoadTemplates("test-data", ".html", engine.FuncMap)
	w := httptest.NewRecorder()

	err := r.Instance("foo/home.html", map[string]string{"Title": "Hello"}).Render(w)
	must(t, err)

	s := w.Body.String()
	mustContain(t, s, "Hello")
	mustContain(t, s, "Home")
}

func must(t *testing.T, err error) {
	if err != nil {
		t.Fatal(err)
	}
}

func mustContain(t *testing.T, s, wanted string) {
	if !strings.Contains(s, wanted) {
		t.Errorf("missing %q in %s", wanted, s)
	}

}

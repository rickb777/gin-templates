package templates_test

import (
	"github.com/gin-gonic/gin"
	"github.com/rickb777/gin-templates"
	"net/http/httptest"
	"strings"
	"testing"
	"time"
)

func TestDebugInstance(t *testing.T) {
	gin.SetMode(gin.DebugMode)
	engine := gin.New()
	r := templates.LoadTemplates("test-data", ".html", engine.FuncMap)
	w1 := httptest.NewRecorder()

	t0 := time.Now()
	instance1 := r.Instance("foo/home.html", map[string]string{"Title": "Hello"})
	d1 := time.Now().Sub(t0)
	err := instance1.Render(w1)
	must(t, err)

	s := w1.Body.String()
	mustContain(t, s, "Hello")
	mustContain(t, s, "Home")

	w2 := httptest.NewRecorder()

	t2 := time.Now()
	instance2 := r.Instance("foo/home.html", map[string]string{"Title": "Hello"})
	d2 := time.Now().Sub(t2)
	err = instance2.Render(w2)
	must(t, err)

	s = w2.Body.String()
	mustContain(t, s, "Hello")
	mustContain(t, s, "Home")
	if d2 >= d1 {
		t.Errorf("expected %v to be less than %v", d2, d1)
	}
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

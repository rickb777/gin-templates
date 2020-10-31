package gin_templates

import (
	"github.com/gin-gonic/gin"
	"net/http/httptest"
	"strings"
	"testing"
)

func TestDebugInstance(t *testing.T) {
	engine := gin.New()
	r := LoadTemplates(engine, "test-data", ".html")
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
	r := LoadTemplates(engine, "test-data", ".html")
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

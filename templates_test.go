package templates_test

import (
	"github.com/gin-gonic/gin"
	. "github.com/onsi/gomega"
	"github.com/rickb777/gin-templates"
	"github.com/spf13/afero"
	"net/http/httptest"
	"os"
	"testing"
	"time"
)

func TestDebugInstance_using_fakes(t *testing.T) {
	g := NewGomegaWithT(t)
	gin.SetMode(gin.DebugMode)
	rec := &recorder{fs: afero.NewMemMapFs()}
	templates.Fs = rec
	rec.fs.MkdirAll("foo/bar", 0755)
	afero.WriteFile(rec.fs, "test-data/foo/home.html", []byte("<html>Home</html>"), 0644)
	afero.WriteFile(rec.fs, "test-data/foo/bar/baz.html", []byte("<html>Baz</html>"), 0644)

	engine := gin.New()
	r := templates.LoadTemplates("test-data", ".html", engine.FuncMap)

	//---------- request 1 ----------
	g.Expect(r.CanProcess("text/plain", "")).To(BeFalse())
	g.Expect(r.CanProcess("text/html", "")).To(BeTrue())
	g.Expect(r.ContentType()).To(Equal("text/html"))

	w := httptest.NewRecorder()
	t0 := time.Now()
	instance := r.Instance("foo/home.html", map[string]string{"Title": "Hello"})
	d1 := time.Now().Sub(t0)
	err := instance.Render(w)
	g.Expect(err).NotTo(HaveOccurred())

	s := w.Body.String()
	g.Expect(s).To(ContainSubstring("Home"))
	g.Expect(rec.opened).To(ContainElements("test-data/foo/home.html", "test-data/foo/bar/baz.html"))

	//---------- request 2: no change so no parsing ----------
	rec.opened = nil
	w = httptest.NewRecorder()

	t2 := time.Now()
	instance = r.Instance("foo/home.html", map[string]string{"Title": "Hello"})
	d2 := time.Now().Sub(t2)
	err = instance.Render(w)
	g.Expect(err).NotTo(HaveOccurred())

	s = w.Body.String()
	g.Expect(s).To(ContainSubstring("Home"))
	g.Expect(rec.opened).To(BeEmpty())
	g.Expect(d2).To(BeNumerically("<", d1)) // it should be faster

	//---------- request 3: an altered file ----------
	rec.opened = nil
	w = httptest.NewRecorder()
	afero.WriteFile(rec.fs, "test-data/foo/bar/baz.html", []byte("<html>Updated</html>"), 0644)

	templates.ContentType = templates.ApplicationXhtml
	g.Expect(r.CanProcess("application/xhtml+xml", "")).To(BeTrue())
	g.Expect(r.ContentType()).To(Equal("application/xhtml+xml"))

	r.Process(w, "foo/bar/baz.html", map[string]string{"Title": "Hello"})

	s = w.Body.String()
	g.Expect(s).To(ContainSubstring("Updated"))
	g.Expect(rec.opened).To(ContainElements("test-data/foo/home.html", "test-data/foo/bar/baz.html"))

	//---------- request 4: a new file ----------
	rec.opened = nil
	w = httptest.NewRecorder()
	afero.WriteFile(rec.fs, "test-data/foo/bar/new.html", []byte("<html>New</html>"), 0644)

	instance = r.Instance("foo/bar/new.html", map[string]string{"Title": "Hello"})
	err = instance.Render(w)
	g.Expect(err).NotTo(HaveOccurred())

	s = w.Body.String()
	g.Expect(s).To(ContainSubstring("New"))
	g.Expect(rec.opened).To(ContainElements("test-data/foo/home.html", "test-data/foo/bar/baz.html", "test-data/foo/bar/new.html"))

	//---------- request 5: ok after deleting a file ----------
	rec.opened = nil
	w = httptest.NewRecorder()
	rec.fs.Remove("test-data/foo/bar/baz.html")

	instance = r.Instance("foo/bar/new.html", map[string]string{"Title": "Hello"})
	err = instance.Render(w)
	g.Expect(err).NotTo(HaveOccurred())

	s = w.Body.String()
	g.Expect(s).To(ContainSubstring("New"))
	g.Expect(rec.opened).To(ContainElements("test-data/foo/home.html", "test-data/foo/bar/new.html"))
}

func TestProductionInstance_using_files(t *testing.T) {
	g := NewGomegaWithT(t)
	gin.SetMode(gin.ReleaseMode)
	templates.Fs = afero.NewOsFs() // real test files

	r := templates.LoadTemplates("test-data", ".html", nil)
	w := httptest.NewRecorder()

	err := r.Instance("foo/home.html", map[string]string{"Title": "Hello"}).Render(w)
	g.Expect(err).NotTo(HaveOccurred())

	s := w.Body.String()
	g.Expect(s).To(ContainSubstring("Hello"))
	g.Expect(s).To(ContainSubstring("Home"))
}

func mustContain(t *testing.T, ss []string, wanted ...string) {
	for _, w := range wanted {
		found := false
		for _, s := range ss {
			if w == s {
				found = true
				break
			}
		}
		if !found {
			t.Errorf("missing %q in %v", w, ss)
		}
	}
}

//-------------------------------------------------------------------------------------------------

type recorder struct {
	fs     afero.Fs
	opened []string
}

func (r *recorder) Create(name string) (afero.File, error) {
	return r.fs.Create(name)
}

func (r *recorder) Mkdir(name string, perm os.FileMode) error {
	return r.fs.Mkdir(name, perm)
}

func (r *recorder) MkdirAll(path string, perm os.FileMode) error {
	return r.fs.MkdirAll(path, perm)
}

func (r *recorder) Open(name string) (afero.File, error) {
	r.opened = append(r.opened, name)
	return r.fs.Open(name)
}

func (r *recorder) OpenFile(name string, flag int, perm os.FileMode) (afero.File, error) {
	//r.opened = append(r.opened, name)
	return r.fs.OpenFile(name, flag, perm)
}

func (r *recorder) Remove(name string) error {
	return r.fs.Remove(name)
}

func (r *recorder) RemoveAll(path string) error {
	return r.fs.RemoveAll(path)
}

func (r *recorder) Rename(oldname, newname string) error {
	return r.fs.Rename(oldname, newname)
}

func (r *recorder) Stat(name string) (os.FileInfo, error) {
	return r.fs.Stat(name)
}

func (r *recorder) Name() string {
	return r.fs.Name()
}

func (r *recorder) Chmod(name string, mode os.FileMode) error {
	return r.fs.Chmod(name, mode)
}

func (r *recorder) Chown(name string, uid, gid int) error {
	return r.fs.Chown(name, uid, gid)
}

func (r *recorder) Chtimes(name string, atime time.Time, mtime time.Time) error {
	return r.fs.Chtimes(name, atime, mtime)
}

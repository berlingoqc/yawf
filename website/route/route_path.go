package route

import (
	"errors"
	"log"
	"net/http"
	"strings"
	"text/template"

	"github.com/gorilla/mux"
	"github.com/russross/blackfriday/v2"
)

// RoutePath est un des chemins directe dans mon path
type RoutePath struct {
	ContentTmplPath string
	Path            string
	TemplateRoot    string

	Handler func(r *http.Request) map[string]interface{}
}

func (b *RoutePath) Initialize(templateRoot string, r *mux.Router) error {
	b.TemplateRoot = templateRoot
	p := strings.Split(b.ContentTmplPath, "/")
	if len(p) == 0 {
		return errors.New("No slash in ContentTmplPath")
	}
	file := p[len(p)-1]
	b.CreateBaseRouteContent(r, b.ContentTmplPath, file, b.Path)
	return nil
}

// GetBaseTemplate ...
func (w *RoutePath) GetBaseTemplate(tmplFile string, name string) (*template.Template, error) {
	tmpl, err := template.ParseFiles(w.TemplateRoot+"/shared/layout.html", w.TemplateRoot+tmplFile, w.TemplateRoot+"/shared/footer.html")
	if err != nil {
		return nil, err
	}
	return tmpl.Lookup(name), nil
}

// CreateBaseRouteContent ...
func (w *RoutePath) CreateBaseRouteContent(r *mux.Router, tmplfile string, tmplname string, handler string) {
	slayout, err := w.GetBaseTemplate(tmplfile, tmplname)
	if err != nil {
		log.Fatal(err)
	}
	r.Path(handler).HandlerFunc(func(rw http.ResponseWriter, r *http.Request) {
		// si j'ai un handler je l'execute et je donne les données a mon template
		var data map[string]interface{}
		if w.Handler != nil {
			data = w.Handler(r)
		}
		slayout, _ = w.GetBaseTemplate(tmplfile, tmplname)
		slayout.ExecuteTemplate(rw, "layout", data)
	})
}

// GetRoutePath crée une instance d'une route de base
func GetRoutePath(tmplFile string, path string) *RoutePath {
	// Get le nom du fichier dans le tmplFile
	b := &RoutePath{
		ContentTmplPath: tmplFile,
		Path:            path,
	}
	return b
}

// GetMarkdownRoutePath retourne une route qui retour un fichier markdown
func GetMarkdownRoutePath(path string, getContent func() ([]byte, error)) *RoutePath {
	b := &RoutePath{
		ContentTmplPath: "/shared/markdown_page.html",
		Path:            path,
		Handler: func(r *http.Request) map[string]interface{} {
			m := make(map[string]interface{})
			d, e := getContent()
			if e != nil {
				// do something
				return nil
			}
			m["md"] = string(blackfriday.Run(d))
			return m
		},
	}
	return b
}

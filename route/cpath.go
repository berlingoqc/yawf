package route

import (
	"html/template"
	"net/http"
	"strings"

	"github.com/gorilla/mux"
)

var (
	// GetHandlerMap create the map that is use on each CPath handler to get global data
	GetHandlerMap = func() map[string]interface{} {
		return make(map[string]interface{})
	}
	// AssetFolder is the template folder
	AssetFolder = ""
	// LayoutTmpl is the tmpl layout file
	LayoutTmpl = ""
	// FooterTmpl is the tmpl footer file
	FooterTmpl = ""
)

// CPath represent an endpoint for a template inside my layout
type CPath struct {
	Path    string
	Tmpl    string
	Handler func(map[string]interface{}, *http.Request)

	FuncMap template.FuncMap
}

// GetPath ...
func (c *CPath) GetPath() string {
	return c.Path
}

// Initialize ...
func (c *CPath) Initialize(r *mux.Router, data map[string]interface{}) {
	p := strings.Split(c.Tmpl, "/")
	if len(p) == 0 {
		panic("No slash in Tmpl")
	}
	file := p[len(p)-1]
	c.CreateBaseRouteContent(r, c.Tmpl, file, c.Path)
}

// GetBaseTemplate ...
func (c *CPath) GetBaseTemplate(tmplFile string, name string) (*template.Template, error) {
	tmpl := template.New(name)
	if c.FuncMap != nil {
		tmpl = tmpl.Funcs(c.FuncMap)
	}
	tmpl, err := tmpl.ParseFiles(AssetFolder+LayoutTmpl, AssetFolder+tmplFile, AssetFolder+FooterTmpl)
	if err != nil {
		return nil, err
	}
	return tmpl.Lookup(name), nil
}

// CreateBaseRouteContent ...
func (c *CPath) CreateBaseRouteContent(r *mux.Router, tmplfile string, tmplname string, handler string) {
	data := GetHandlerMap()
	c.FuncMap = data["Func"].(template.FuncMap)
	slayout, err := c.GetBaseTemplate(tmplfile, tmplname)
	if err != nil {
		panic(err)
	}
	r.Path(handler).HandlerFunc(func(rw http.ResponseWriter, r *http.Request) {
		data = GetHandlerMap()
		// si j'ai un handler je l'execute et je donne les données a mon template
		if c.Handler != nil {
			c.Handler(data, r)
			// Regarde si il y a une erreur ou si on doit rediriger ailleur
			if data != nil {
				if _, ok := data["Error"]; ok {
					data["Tmpl"] = tmplfile
					tmplfile = "/error.html"
					tmplname = "error.html"
				} else {
					// Si la template request d'avoir info de l'usage courrant

				}
			}
		}
		slayout, err = c.GetBaseTemplate(tmplfile, tmplname)
		if err != nil {
			panic(err)
		}
		slayout.ExecuteTemplate(rw, "layout", data)
	})
}

// GetCPath get a cpath
func GetCPath(path string, tmpl string) *CPath {
	r := &CPath{Path: path, Tmpl: tmpl}
	return r
}

// GetCPathMarkdown get a cpath that render a markdown file in the template
func GetCPathMarkdown(path string, file string) *CPath {
	r := &CPath{Path: path, Tmpl: "/shared/markdown_page.html"}
	r.Handler = func(m map[string]interface{}, r *http.Request) {
		m["File"] = file
	}

	return r
}

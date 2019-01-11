package route

import (
	"encoding/json"
	"errors"
	"fmt"
	"html/template"
	"io/ioutil"
	"log"
	"net/http"
	"net/url"
	"reflect"
	"strconv"
	"strings"

	"github.com/gorilla/mux"
	"github.com/russross/blackfriday/v2"
)

var (
	GetHandlerMap = func() map[string]interface{} {
		return make(map[string]interface{})
	}
)

func RespondWithError(w http.ResponseWriter, code int, message string) {
	RespondWithJSON(w, code, map[string]string{"error": message})
}

func RespondWithJSON(w http.ResponseWriter, code int, payload interface{}) {
	resp, e := json.Marshal(payload)
	if e != nil {
		log.Printf("Failed to marshal message %v\n", e.Error())
		resp = []byte(e.Error())
	}

	w.Header().Set("Content-Type", "application/json")
	w.WriteHeader(code)
	w.Write(resp)
}

func GetStringURLQuery(name string, u *url.URL) string {
	n, ok := u.Query()[name]
	if !ok {
		return ""
	}
	return n[0]
}

func StructFromQuery(url *url.URL, t interface{}) error {
	data := url.Query()
	typeT := reflect.TypeOf(t)
	if typeT.Kind() != reflect.Ptr {
		return errors.New("interface{} must be ptr but is " + string(typeT.Kind()))
	}
	typeT = typeT.Elem()
	values := reflect.ValueOf(t)
	values = values.Elem()
	for i := 0; i < typeT.NumField(); i++ {
		ft := typeT.Field(i)
		fv := values.Field(i)
		for k, v := range data {
			if k == ft.Name {
				if len(v) == 1 {
					switch fv.Type().Kind() {
					case reflect.String:
						fv.SetString(v[0])
						break
					case reflect.Int:
						i, err := strconv.Atoi(v[0])
						if err != nil {
							return fmt.Errorf("Error parsing field %v with value %v error %v", ft.Name, v[0], err)
						}
						fv.SetInt(int64(i))
						break
					default:
						break
					}
				}
			}
		}
	}
	return nil
}

// Widget is a portable template that can be render on the server and
// query with jQuery
type Widget struct {
	// Name of this widget is for query in map
	Name string
	// The Path to the template
	File string
	// Stuct reprensent the wanted data struct from the request
	// when rendering the widget
	Struct interface{}
	// Render is the function call to render the widget in the response writer
	Render func(interface{}, http.ResponseWriter, *http.Request) interface{}
}

// ModuleWidgetAPI create a api to serve widget use by a module and provide to the others
func ModuleWidgetAPI(r *mux.Route, root string, widgets map[string]Widget) {
	r.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		// go find the key that tell use the name of the Widget wanted
		n := GetStringURLQuery("id_widget", r.URL)
		if n == "" {
			RespondWithError(w, http.StatusBadRequest, "No widget query")
			return
		}
		widget, ok := widgets[n]
		if !ok {
			RespondWithError(w, http.StatusBadRequest, "Widget don't exists "+n)
			return
		}
		// Si on n'a besoin d'une struct essate de la parser depuis la map
		if widget.Struct != nil {
			// Essaye de parser la struct
			err := StructFromQuery(r.URL, widget.Struct)
			if err != nil {
				RespondWithError(w, http.StatusBadRequest, err.Error())
				return
			}
		}
		// On peut render notre widget dans le writer
		data := widget.Render(widget.Struct, w, r)
		if data == nil {
			return
		}
		tmpl := template.New("")
		tmpl, err := tmpl.ParseFiles(root + "/template" + widget.File)
		tmpl = tmpl.Lookup(widget.File)
		if err != nil {
			// Respond with error widget
			RespondWithError(w, http.StatusBadRequest, err.Error())
			return
		}
		tmpl.Execute(w, data)
	})
}

// AddMarkdownFolderHandler create a handler that retourne the
// contains of the markdown file from the repository given
func AddMarkdownFolderHandler(r *mux.Route, mdDirectory string) error {
	r.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		n, ok := r.URL.Query()["name"]
		if !ok || len(n) != 1 {
			w.Write([]byte("Error"))
		}
		data, err := ioutil.ReadFile(mdDirectory + "/" + n[0])
		if err != nil {
			w.Write([]byte(err.Error()))
		}
		w.Write(blackfriday.Run(data))

	})

	return nil
}

// RoutePath est un des chemins directe dans mon path
type RoutePath struct {
	ContentTmplPath string
	Path            string
	TemplateRoot    string

	Handler func(map[string]interface{}, *http.Request)

	funcMap template.FuncMap
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
	tmpl := template.New(name)
	if w.funcMap != nil {
		tmpl = tmpl.Funcs(w.funcMap)
	}
	tmpl, err := tmpl.ParseFiles(w.TemplateRoot+"/shared/layout.html", w.TemplateRoot+tmplFile, w.TemplateRoot+"/shared/footer.html")
	if err != nil {
		return nil, err
	}
	return tmpl.Lookup(name), nil
}

// CreateBaseRouteContent ...
func (w *RoutePath) CreateBaseRouteContent(r *mux.Router, tmplfile string, tmplname string, handler string) {
	data := GetHandlerMap()
	w.funcMap = data["Func"].(template.FuncMap)
	slayout, err := w.GetBaseTemplate(tmplfile, tmplname)
	if err != nil {
		log.Fatal(err)
	}
	r.Path(handler).HandlerFunc(func(rw http.ResponseWriter, r *http.Request) {
		data = GetHandlerMap()
		// si j'ai un handler je l'execute et je donne les données a mon template
		if w.Handler != nil {
			w.Handler(data, r)
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

func GetMarkdownRouteFile(fileName string, path string) *RoutePath {
	return GetMarkdownRoutePath(path, func() ([]byte, error) {
		return ioutil.ReadFile(fileName)
	})
}

// GetMarkdownRoutePath retourne une route qui retour un fichier markdown
func GetMarkdownRoutePath(path string, getContent func() ([]byte, error)) *RoutePath {
	b := &RoutePath{
		ContentTmplPath: "/shared/markdown_page.html",
		Path:            path,
		Handler: func(m map[string]interface{}, r *http.Request) {
			d, e := getContent()
			if e != nil {
				m["Error"] = e
				return
			}
			m["md"] = string(blackfriday.Run(d))
		},
	}
	return b
}

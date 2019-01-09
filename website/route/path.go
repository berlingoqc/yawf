package route

import (
	"net/http"

	"github.com/gorilla/mux"
)

// WPath correspond a un des mes sous-path de mon site avec son router
type WPath struct {
	Name        string
	Router      *mux.Router
	Route       map[string]*RoutePath
	StaticRoute map[string]string

	Middleware []func(http.Handler) http.Handler
}

// Initialize donnes les informations au template des routes pour la creation des templates
// et initilalize les middleware
func (w *WPath) Initialize(templateRoot string) error {
	var e error
	for _, v := range w.Route {
		e = v.Initialize(templateRoot, w.Router)
		if e != nil {
			return e
		}
	}

	for k, v := range w.StaticRoute {
		w.Router.Path(k).HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
			http.ServeFile(w, r, v)
		})
	}

	for _, v := range w.Middleware {
		w.Router.Use(v)
	}

	// Initialize l'envoie de fichiers statique

	return nil
}

func (w *WPath) AddMiddleware(h func(http.Handler) http.Handler) {
	w.Middleware = append(w.Middleware, h)
}

// GetWPath crée une nouvelle instance de WPath avec les informations données
func GetWPath(name string, router *mux.Router, routes ...*RoutePath) *WPath {
	w := &WPath{
		Name:        name,
		Router:      router,
		Route:       make(map[string]*RoutePath),
		StaticRoute: make(map[string]string),
	}

	for _, r := range routes {
		w.Route[r.Path] = r
	}

	return w
}

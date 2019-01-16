package route

import (
	"net/http"

	"github.com/berlingoqc/yawf/route/security"
	"github.com/gorilla/mux"
)

// Item is a element of WPath
type Item struct {
	Path     IPath
	Security *security.PathSecurity
}

// WPath represent a point that contains multiple IPath and that have a middleware
type WPath struct {
	Name   string
	Router *mux.Router
	Route  map[string]*Item

	// Security is the default that apply when the item is null
	Security *security.PathSecurity

	Middleware []func(http.Handler) http.Handler
}

// InitializeWPath a wpath
func InitializeWPath(wp *WPath, ctx map[string]interface{}) {
	for _, v := range wp.Route {
		v.Path.Initialize(wp.Router, ctx)
	}
}

// GetWPath return a new instance of wpath
func GetWPath(path string, r *mux.Router) *WPath {
	w := &WPath{
		Name:   path,
		Router: r,
		Route:  make(map[string]*Item),
	}

	return w
}

// AddWPathItem add an item to a WPath
func AddWPathItem(w *WPath, items ...IPath) {
	for _, i := range items {
		item := &Item{
			Path: i,
		}
		w.Route[i.GetPath()] = item
	}

}

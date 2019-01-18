package route

import (
	"fmt"
	"net/http"
	"strings"

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

// HaveSecurity tell if the wpath needed to have a auth middleware
func HaveSecurity(wp *WPath) bool {
	if wp.Security == nil {
		for _, r := range wp.Route {
			if r.Security != nil {
				return true
			}
		}
	} else {
		return true
	}

	return false
}

// GetAuthMiddleware return a custom middleware for the authentification
// needed by this wpath if handler is null there is no authentification
// to make
func GetAuthMiddleware(wp *WPath) func(http.Handler) http.Handler {
	// Valide if there is no auth for the wpath
	if !HaveSecurity(wp) {
		return nil
	}
	return func(next http.Handler) http.Handler {
		return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
			// Get le subpath relative au path du WPath
			url := r.URL.String()
			subpath := strings.Replace(url, wp.Name, "", 1)
			fmt.Printf("Auth middleware WPath %v url %v item %v \n", wp.Name, url, subpath)
			if item, ok := wp.Route[subpath]; ok {
				if item.Security != nil {
					if !security.PathSecurityValidation(item.Security, w, r) {
						// Erreur de securiter c'est occuper du redirect deja
						return
					}
				} else {
					// Applique la securité global si présente
					if wp.Security != nil {
						if !security.PathSecurityValidation(wp.Security, w, r) {
							return
						}
					}
				}
			} else {
				// Essaye d'acceder dequoi qui est pas dans mes path, redirige
				security.RedirectTo(w, r, "/", nil)
				return
			}

			next.ServeHTTP(w, r)
		})
	}
}

// InitializeWPath a wpath
func InitializeWPath(wp *WPath, ctx map[string]interface{}) {
	for _, v := range wp.Route {
		v.Path.Initialize(wp.Router, ctx)
	}
	// Crée mon middleware d'authentification si au moins un element
	// et securiser dans le WPath
	authMiddleware := GetAuthMiddleware(wp)
	if authMiddleware != nil {
		fmt.Printf("Got a authentification middleware for %v \n", wp.Name)
		wp.Middleware = append(wp.Middleware, authMiddleware)
	} else {
		fmt.Printf("WPath %v does not requries authentification \n", wp.Name)
	}
	for _, m := range wp.Middleware {
		wp.Router.Use(m)
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

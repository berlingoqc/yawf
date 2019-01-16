package route

import (
	"net/http"

	"github.com/gorilla/mux"
)

// FPath is a path that serve Static file
type FPath struct {
	Path string
	Root string
}

// GetPath ...
func (f *FPath) GetPath() string { return f.Path }

// Initialize ...
func (f *FPath) Initialize(r *mux.Router, data map[string]interface{}) {
	r.PathPrefix(f.Path).Handler(http.StripPrefix(f.Path, http.FileServer(http.Dir(f.Root))))
}

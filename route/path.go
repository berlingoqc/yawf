package route

import (
	"github.com/gorilla/mux"
)

// IPath is the interface that represent all my path
type IPath interface {
	Initialize(r *mux.Router, data map[string]interface{})
	GetPath() string
}

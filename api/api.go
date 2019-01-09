package api

import (
	"github.com/gorilla/mux"
)

type WApi struct {
	Path   string
	Router *mux.Router
}

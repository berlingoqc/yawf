package route

import (
	"errors"
	"io/ioutil"
	"net/http"

	"github.com/gorilla/mux"
	"github.com/russross/blackfriday/v2"
)

// MPath represent a path that server markdown
type MPath struct {
	Path       string
	Directory  string
	GetContent func(name string) ([]byte, error)
}

// GetPath ...
func (m *MPath) GetPath() string {
	return m.Path
}

// Initialize ...
func (m *MPath) Initialize(rr *mux.Router, data map[string]interface{}) {
	rr.Path(m.Path).HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		var err error
		defer func() {
			if err != nil {
				respondWithError(w, http.StatusBadRequest, err.Error())
			}
		}()
		n, ok := r.URL.Query()["name"]
		if !ok || len(n) != 1 {
			err = errors.New("No name query")
			return
		}
		data, err := m.GetContent(n[0])
		if err != nil {
			return
		}
		w.Write(blackfriday.Run(data))
	})
}

// GetMPathFile get a MPath that get it's file
// from a local repository
func GetMPathFile(path string, directory string) *MPath {
	m := &MPath{
		Path:      path,
		Directory: directory,
		GetContent: func(name string) ([]byte, error) {
			return ioutil.ReadFile(directory + "/" + name)
		},
	}
	return m
}

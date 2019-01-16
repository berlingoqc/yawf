package route

import (
	"errors"
	"fmt"
	"html/template"
	"net/http"

	"github.com/berlingoqc/yawf/utility"

	"github.com/berlingoqc/yawf/conv"

	"github.com/berlingoqc/yawf/route/security"

	"github.com/gorilla/mux"
)

// WidgetQuery is the struct to query a widget
type WidgetQuery struct {
	// Name of the widget wanted
	ID string
	// Data is a map[string]interface{} serialize in json
	Data string
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

// WidPath represent a path that server widget
type WidPath struct {
	Path        string
	Security    *security.PathSecurity
	Widgets     map[string]*Widget
	AssetFolder string
}

// GetPath ...
func (w *WidPath) GetPath() string {
	return w.Path
}

// GetSecurity ...
func (w *WidPath) GetSecurity() *security.PathSecurity {
	return w.Security
}

// Initialize ...
func (w *WidPath) Initialize(r *mux.Router, data map[string]interface{}) error {
	r.Path("/widget").HandlerFunc(func(rw http.ResponseWriter, r *http.Request) {
		// Get la struct depuis la request
		var err error
		defer func() {
			if err != nil {
				fmt.Printf("Error rendering widget %v\n", err.Error())
				respondWithError(rw, http.StatusBadRequest, err.Error())
			}
		}()
		wq := &WidgetQuery{}
		if err = conv.QueryToStruct(r.URL.Query(), wq); err != nil {
			return
		}
		widget, ok := w.Widgets[wq.ID]
		if !ok {
			err = errors.New("Can't find widget " + wq.ID)
			return
		}
		if widget.Struct != nil {
			// Parse la struct requise pour la template
			var m map[string]interface{}
			m, err = utility.JSONToMap(wq.Data)
			if err != nil {
				return
			}
			err = conv.MapToStruct(m, widget.Struct)
			if err != nil {
				return
			}
		}
		data := widget.Render(widget.Struct, rw, r)
		if data == nil {
			err = errors.New("Error rendering widget " + wq.ID)
			return
		}
		tmpleData, err := utility.ReadFileString(w.AssetFolder + widget.File)
		if err != nil {
			return
		}
		tmpl, err := template.New("widget").Parse(tmpleData)
		if err != nil {
			return
		}
		err = tmpl.Execute(rw, data)
		if err != nil {
			return
		}
		rw.Header().Add("Content-Type", "text/plain")
	})
	return nil
}

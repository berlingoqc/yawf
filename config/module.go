package config

import (
	"errors"
	"plugin"

	"github.com/berlingoqc/yawf/website/route"
	"github.com/gorilla/mux"

	"github.com/berlingoqc/yawf/db"
)

var (
	availableModules map[string]IModule
)

func init() {
	availableModules = make(map[string]IModule)
}

// AddAvailableModule add a module to the list of available module, automaticly call when
// you laod a module dynamically
func AddAvailableModule(name string, module IModule) {
	availableModules[name] = module
}

// GetModule return of the available module if present
func GetModule(name string) IModule {
	if d, ok := availableModules[name]; ok {
		return d
	}
	return nil
}

// GetModules return the map of available module
func GetModules() map[string]IModule {
	return availableModules
}

// ModuleInfo contient l'informations sur les modules
// disponible dans le framework
type ModuleInfo struct {
	Name        string
	Package     string
	Version     string
	Description string

	Provides            []string
	Dependencie         []string
	OptionalDependencie []string
}

// IModule is the interface for the module of the website. A module correspond
// to route with pages, api, widgets, tasks and database. They can be load
// dynamicly
type IModule interface {
	// Initialize receive a map that correspond to the context of the website
	Initialize(data map[string]interface{}) error
	// GetInfo return the informations about this module
	GetInfo() ModuleInfo
	// GetDBInstance return a open connection to the database of the module
	GetDBInstance() (db.IDB, error)
	// GetNeededAssets return the list of asset required by this module
	GetNeededAssets() []string
	// GetNavigationItems return the list of items to add to the main navigation bar
	GetNavigationItems() []interface{}
	// GetWidgets return the map of the available widgets provides by
	// this module
	GetWidgets() []*route.Widget
	// GetTasks return the list of task that need to be added to the
	// task manager
	GetTasks() []ITask
	// GetWPath return the list of wpath to add to the webserver
	GetWPath(*mux.Router) []*route.WPath
}

// LoadModuleDynamicly try to open a shared modules and find
// the function with the signature func GetModule() (string,IModule)
// and return the result of the call to this function is present
func LoadModuleDynamicly(filepath string) (string, IModule, error) {
	p, err := plugin.Open(filepath)
	if err != nil {
		return "", nil, err
	}

	f, err := p.Lookup("GetModule")
	if err != nil {
		return "", nil, err
	}
	if moduleFunc, ok := f.(func() (string, IModule)); ok {
		name, module := moduleFunc()
		AddAvailableModule(name, module)
		return name, module, nil
	}

	return "", nil, errors.New("Invalid signature for function")

}

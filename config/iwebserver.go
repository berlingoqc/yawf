package config

import (
	"log"

	"github.com/berlingoqc/yawf/website/route"
	"github.com/gorilla/mux"
)

// IWebServer is the interface for my webserver
type IWebServer interface {
	// GetConfig return the configuration of the webserver
	GetConfig() *WebSite
	// GetTaskPool grant access to the TaskPool for adding new task
	GetTaskPool() *TaskPool
	// GetLogger return a instance of the logger for the website
	// set with the configuration
	GetLogger() *log.Logger
	// GetMux return the base mux of the website
	GetMux() *mux.Router
	// GetNavigationBar get all the informations from the active module
	// to get all the informations needed for the navigation bar
	GetNavigationBar() *route.NavigationBar
	// AddModule add a new loaded module to the webserver
	AddModule(IModule, Ctx) error
}

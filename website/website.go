// Package website contient le code pour demarrer un site web deja configurer
package website

import (
	"context"
	"fmt"
	"html/template"
	"log"
	"net/http"
	"os"
	"os/signal"
	"syscall"
	"time"

	"github.com/berlingoqc/yawf/module/base"
	"github.com/berlingoqc/yawf/module/user"
	"github.com/berlingoqc/yawf/utility"

	"github.com/berlingoqc/yawf/config"

	"github.com/berlingoqc/yawf/route"

	"github.com/gorilla/mux"
)

var (
	// UserModule instance
	UserModule config.IModule
	// BaseModule instance
	BaseModule config.IModule
)

func init() {
	name, BaseModule := base.GetModule()
	config.AddAvailableModule(name, BaseModule)
	name, UserModule = user.GetModule()
	config.AddAvailableModule(name, UserModule)
}

// WebServer my webserver that work with my modules implements of IWebServer
type WebServer struct {
	// Logger of the web site and its module
	Logger *log.Logger
	// Mux is the base router of the website
	Mux *mux.Router
	// Hs is the http server that run
	Hs *http.Server
	// TaskPool is the manager of my tasks
	TaskPool config.TaskPool
	// ChannelStop is the channel to stop the webserver
	ChannelStop chan os.Signal

	Config *config.WebSite

	NavigationBar *route.NavigationBar
	Routes        map[string]*route.WPath
	ActiveModules map[string]config.IModule
	Widgets       map[string]*route.Widget
}

// GetNavigationBar ...
func (w *WebServer) GetNavigationBar() *route.NavigationBar {
	return w.NavigationBar
}

// GetTaskPool ...
func (w *WebServer) GetTaskPool() *config.TaskPool {
	return &w.TaskPool
}

// GetLogger ...
func (w *WebServer) GetLogger() *log.Logger {
	return w.Logger
}

// GetMux ...
func (w *WebServer) GetMux() *mux.Router {
	return w.Mux
}

// AddModule ...
func (w *WebServer) AddModule(m config.IModule, ctx config.Ctx) error {

	return nil
}

// ImportModuleAsset import the asset from the module inside the root folder
func (w *WebServer) ImportModuleAsset(m config.IModule) error {
	gopath := os.Getenv("GOPATH")
	gopath += "/src" + "/" + m.GetInfo().Package

	files := m.GetNeededAssets()
	for _, f := range files {
		dest := w.Config.File.GetAssetFolderPath() + f
		src := gopath + "/asset" + f
		if err := utility.Copy(src, dest); err != nil {
			return err
		}
	}

	return nil
}

// Setup configure le serveur web doit être appeler avec le reste
func (w *WebServer) Setup(configWs *config.WebSite) error {

	w.NavigationBar = route.GetNavigationBar(configWs.Name)

	w.TaskPool.Tasks = make(map[string]config.ITask)

	w.Routes = make(map[string]*route.WPath)

	w.Config = configWs

	w.Logger = log.New(os.Stdout, "", 0)

	w.Widgets = make(map[string]*route.Widget)

	r := mux.NewRouter()
	w.Mux = r

	// Update le default handler de route pour qu'il ajoute dans la map la navbar
	route.GetHandlerMap = func() map[string]interface{} {
		m := make(map[string]interface{})
		m["Navbar"] = w.NavigationBar
		mf := template.FuncMap{}

		route.AddNavBarTmplFunc(mf)

		m["Func"] = mf
		return m

	}
	// Configure les variables d'env pour les CPath
	route.AssetFolder = configWs.File.GetAssetFolderPath()
	route.FooterTmpl = "/shared/footer.html"
	route.LayoutTmpl = "/shared/layout.html"

	// Creation de mon handler pour servir les fichiers statiques
	r.PathPrefix("/static/").Handler(http.StripPrefix("/static/", http.FileServer(http.Dir(w.Config.File.GetStaticFolderPath()))))

	var err error
	ctx := configWs.ToContext()

	w.Hs, err = config.GetWebServer(ctx)
	if err != nil {
		return err
	}
	w.Hs.Handler = w.Mux

	for k, v := range w.Config.EnableModule {
		if module := config.GetModule(k); module != nil {
			if v != nil {
				vv := v.(map[string]interface{})
				utility.ConcatMap(vv, ctx)
				err = module.Initialize(vv)
			} else {
				err = module.Initialize(ctx)
			}
			if err != nil {
				return err
			}
			// Copy all needed file from the module inside the root
			err = w.ImportModuleAsset(module)
			if err != nil {
				return err
			}
			widgets := module.GetWidgets()
			for _, widg := range widgets {
				w.Widgets[widg.Name] = widg
			}
			tasks := module.GetTasks()
			for _, t := range tasks {
				w.TaskPool.Tasks[t.GetName()] = t
			}
			navItems := module.GetNavigationItems()
			for _, item := range navItems {
				switch v := item.(type) {
				case route.Button:
					w.NavigationBar.Buttons = append(w.NavigationBar.Buttons, v)
				default:
					w.NavigationBar.Items = append(w.NavigationBar.Items, v)
				}
			}
			paths := module.GetWPath(w.Mux)
			for _, path := range paths {
				w.Routes[path.Name] = path
			}
		} else {
			return fmt.Errorf("Can't find the module %v in the AvailableModules", k)
		}
	}

	for _, r := range w.Routes {
		route.InitializeWPath(r, nil)
	}

	return nil
}

// Start demarre le serveur web
func (w *WebServer) Start() {
	// Crée mon channel pour le signal d'arret
	w.ChannelStop = make(chan os.Signal, 1)

	signal.Notify(w.ChannelStop, os.Interrupt, syscall.SIGTERM)

	go func() {
		w.Logger.Printf("Starting yawf.ca at %v\n", w.Hs.Addr)
		if err := w.Hs.ListenAndServe(); err != http.ErrServerClosed {
			w.Logger.Fatal(err)
		}
	}()
}

// Stop arrête le serveur web
func (w *WebServer) Stop() {
	w.Logger.Println("Fermeture du serveur")
	ctx, c := context.WithTimeout(context.Background(), 5*time.Second)
	defer c() // release les ressources du context

	w.Hs.Shutdown(ctx)

	w.Logger.Println("Serveur eteint ...")
}

// StartWebServer start a webserver that as already been configurate
func StartWebServer(websiteRoot string) (*WebServer, error) {

	ws := &WebServer{}

	c, err := config.LoadWebSiteConfig(websiteRoot)
	if err != nil {
		return nil, err
	}
	if err = c.Validate(); err != nil {
		return nil, err
	}

	err = ws.Setup(c)
	if err != nil {
		return nil, err
	}

	ws.Start()

	return ws, nil
}

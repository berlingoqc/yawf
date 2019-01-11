// website contient le code pour demarrer un site web deja configurer
package website

import (
	"context"
	"log"
	"net/http"
	"os"
	"os/signal"
	"syscall"
	"time"

	"github.com/berlingoqc/yawf/module/base"

	"github.com/berlingoqc/yawf/config"
	"github.com/berlingoqc/yawf/module/blog"
	"github.com/berlingoqc/yawf/module/project"
	"github.com/berlingoqc/yawf/module/user"

	"github.com/berlingoqc/yawf/website/route"

	"github.com/gorilla/mux"
)

const (
	KeyAddr    = "addr"
	KeyLogFile = "logfile"
	KeyLogOut  = "logout"
)

var (
	DefaultModule = make(map[string]config.IModule)
)

func init() {
	s, m := blog.GetModule()
	DefaultModule[s] = m
	s, m = base.GetModule()
	DefaultModule[s] = m
	s, m = project.GetModule()
	DefaultModule[s] = m
	s, m = user.GetModule()
	DefaultModule[s] = m
}

// WebServer base de mon web server avec mux , run async et peux
// etre canceller avec un channel
type WebServer struct {
	Logger *log.Logger
	Mux    *mux.Router
	Hs     *http.Server

	TaskPool config.TaskPool

	ChannelStop chan os.Signal

	AssetPath  string
	StaticRoot string

	MainRoutes map[string]*route.WPath
}

func (w *WebServer) GetTaskPool() *config.TaskPool {
	return &w.TaskPool
}

func (w *WebServer) GetLogger() *log.Logger {
	return w.Logger
}

func (w *WebServer) GetMux() *mux.Router {
	return w.Mux
}

func (w *WebServer) AddRoute(r *route.WPath) {
	w.MainRoutes[r.Name] = r
}

// Setup configure le serveur web doit être appeler avec le reste
func (w *WebServer) Setup(assetPath *config.WebSite) error {

	w.TaskPool.Tasks = make(map[string]config.ITask)

	w.MainRoutes = make(map[string]*route.WPath)

	w.AssetPath = assetPath.RootFolder + "/template"
	w.StaticRoot = assetPath.RootFolder + "/static"

	w.Logger = log.New(os.Stdout, "", 0)

	r := mux.NewRouter()
	w.Mux = r

	// Creation de mon handler pour servir les fichiers statiques
	r.PathPrefix("/static/").Handler(http.StripPrefix("/static/", http.FileServer(http.Dir(w.StaticRoot))))

	w.Hs = config.GetWebServer(assetPath.NetOptions)
	w.Hs.Handler = w.Mux

	for k, v := range assetPath.EnableModule {
		if module, ok := DefaultModule[k]; ok {
			module.Initialize(v)
			if err := module.AddToWebServer(w); err != nil {
				return err
			}
		}
	}

	for _, r := range w.MainRoutes {
		if err := r.Initialize(w.AssetPath); err != nil {
			return err
		}
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
func StartWebServer(assetPath string) (*WebServer, error) {

	ws := &WebServer{}

	c, err := config.Load(assetPath)
	if err != nil {
		return nil, err
	}

	err = ws.Setup(c)
	if err != nil {
		return nil, err
	}

	ws.Start()

	return ws, nil
}

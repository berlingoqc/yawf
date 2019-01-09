// website contient le code pour demarrer un site web deja configurer
package website

import (
	"context"
	"fmt"
	"io/ioutil"
	"log"
	"net/http"
	"os"
	"os/signal"
	"syscall"
	"time"

	"github.com/berlingoqc/yawf/website/project"

	"github.com/berlingoqc/yawf/db"

	"github.com/berlingoqc/yawf/config"

	"github.com/berlingoqc/yawf/auth"
	"github.com/berlingoqc/yawf/website/route"

	"github.com/berlingoqc/yawf/api"

	"github.com/gorilla/mux"
)

const (
	KeyAddr    = "addr"
	KeyLogFile = "logfile"
	KeyLogOut  = "logout"
)

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

// Setup configure le serveur web doit être appeler avec le reste
func (w *WebServer) Setup(assetPath config.WebSite, options map[string]interface{}) error {

	w.TaskPool.Tasks = make(map[string]config.ITask)

	w.MainRoutes = make(map[string]*route.WPath)

	w.AssetPath = assetPath.TemplateRoot
	w.StaticRoot = assetPath.StaticRoot
	addr := ":8081"
	w.Logger = log.New(os.Stdout, "", 0)
	if options != nil {
		// regarde pour les clées optionnels
	}
	r := mux.NewRouter()

	// Creation de mon handler pour servir les fichiers statiques
	r.PathPrefix("/static/").Handler(http.StripPrefix("/static/", http.FileServer(http.Dir(w.StaticRoot))))

	w.MainRoutes["yawf"] = route.GetWPath("yawf", r,
		&route.RoutePath{ContentTmplPath: "/index.html", Path: "/"},
		&route.RoutePath{ContentTmplPath: "/about.html", Path: "/about"},
		&route.RoutePath{ContentTmplPath: "/project/project.html", Path: "/project"},
		&route.RoutePath{ContentTmplPath: "/cv/cv.html", Path: "/professional"},
		&route.RoutePath{ContentTmplPath: "/contact.html", Path: "/contact"},
		&route.RoutePath{ContentTmplPath: "/me.html", Path: "/me"},
		&route.RoutePath{ContentTmplPath: "/denied.html", Path: "/denied"},
	)

	w.MainRoutes["yawf"].Route["/personal"] = route.GetMarkdownRoutePath("/personal", func() ([]byte, error) {
		return ioutil.ReadFile("./markdown/personal.md")
	})

	w.MainRoutes["yawf"].Route["/about"] = route.GetMarkdownRoutePath("/about", func() ([]byte, error) {
		return ioutil.ReadFile("./markdown/about.md")
	})

	w.MainRoutes["yawf"].Route["/contact"] = route.GetMarkdownRoutePath("/contact", func() ([]byte, error) {
		return ioutil.ReadFile("./markdown/contact.md")
	})

	w.MainRoutes["yawf"].Route["/project"].Handler = func(r *http.Request) map[string]interface{} {
		idb, err := GetProjectDBInstance("project.db")
		if err != nil {
			log.Println(err)
			return nil
		}
		defer db.CloseDatabse(idb)
		m := make(map[string]interface{})
		m["Projects"], err = idb.GetPProjects()
		if err != nil {
			log.Println(err)
			return nil
		}
		m["ghaccount"], err = idb.GetGHAccount()
		if err != nil {
			log.Println(err)
		}

		return m
	}

	w.MainRoutes["yawf"].Route["/professional"].Handler = func(r *http.Request) map[string]interface{} {
		idb, err := GetProjectDBInstance("project.db")
		if err != nil {
			log.Println(err)
			return nil
		}
		defer db.CloseDatabse(idb)
		m := make(map[string]interface{})
		m["experiences"], err = idb.GetExperience()
		if err != nil {
			log.Println(err)
		}
		m["formations"], err = idb.GetFormation()
		if err != nil {
			log.Println(err)
		}
		m["languages"], err = idb.GetLanguageExperience()
		if err != nil {
			log.Println(err)
		}
		return m
	}

	blogRouter := r.PathPrefix("/blog").Subrouter()
	w.MainRoutes["blog"] = route.GetWPath("blog", blogRouter,
		&route.RoutePath{ContentTmplPath: "/blog/index.html", Path: "/"},
		&route.RoutePath{ContentTmplPath: "/blog/post.html", Path: "/post"},
		&route.RoutePath{ContentTmplPath: "/blog/serie.html", Path: "/serie"},
	)

	// Creation des handlers pour route /auth
	authRouter := r.PathPrefix("/auth").Subrouter()

	w.MainRoutes["auth"] = route.GetWPath("auth", authRouter,
		&route.RoutePath{ContentTmplPath: "/login/login.html", Path: "/login"},
		&route.RoutePath{ContentTmplPath: "/login/new.html", Path: "/new"},
		&route.RoutePath{ContentTmplPath: "/login/confirm.html", Path: "/confirm"},
	)
	w.MainRoutes["auth"].AddMiddleware(api.MiddlewareAuth)

	// Creation des handlers pour route /yase
	yaseRouter := r.PathPrefix("/yase").Subrouter()

	w.MainRoutes["yase"] = route.GetWPath("yase", yaseRouter,
		&route.RoutePath{ContentTmplPath: "/yase/index.html", Path: "/"},
		&route.RoutePath{ContentTmplPath: "/yase/yase_creator.html", Path: "/creator"},
	)
	w.MainRoutes["yase"].StaticRoute["/yase_creator.js"] = "./template/yase/yase_creator.js"
	w.MainRoutes["yase"].StaticRoute["/yase_creator.wasm"] = "./template/yase/yase_creator.wasm"
	w.MainRoutes["yase"].StaticRoute["/yase_creator.data"] = "./template/yase/yase_creator.data"
	w.MainRoutes["yase"].AddMiddleware(api.MiddlewareYase)

	// Creation des handles pour route /accout
	accountRouter := r.PathPrefix("/account").Subrouter()

	w.MainRoutes["account"] = route.GetWPath("account", accountRouter,
		&route.RoutePath{ContentTmplPath: "/account/dashboard_admin.html", Path: "/admin/"},
		&route.RoutePath{ContentTmplPath: "/account/dashboard.html", Path: "/dashboard"},
	)

	w.MainRoutes["account"].Route["/dashboard"].Handler = func(r *http.Request) map[string]interface{} {
		idb, _ := auth.GetAuthDBInstance("auth.db")
		defer idb.CloseDatabase()
		m := make(map[string]interface{})
		u, _ := api.ValidUserCookie(idb, r)
		m["user"] = u
		return m
	}

	w.MainRoutes["account"].AddMiddleware(api.MiddlewareAccount)

	// Fin des handlers de base ma les initialiser
	for _, v := range w.MainRoutes {
		err := v.Initialize(assetPath.TemplateRoot)
		if err != nil {
			return err
		}
	}

	// Initialize l'api d'authentification
	UserAPI := &api.UserAPI{
		Db: &auth.AuthDB{},
	}

	err := UserAPI.Db.OpenDatabase("auth.db")
	if err != nil {
		w.Logger.Panic(err)
	}

	// get le subroute pour initialiser UserAPI
	apiRouter := r.PathPrefix("/api").Subrouter()

	err = UserAPI.Initialize(apiRouter)
	if err != nil {
		w.Logger.Panic(err)
	}

	// Crée ma task pour updater periodiquement mon compte github
	w.TaskPool.AddPeriodicTask("github", 1*time.Hour, func(c chan *config.Signal, args ...interface{}) {
		// Get instance bd
		idb, err := GetProjectDBInstance("project.db")
		if err != nil {
			// log l'erreur
			fmt.Printf("Error task github %v \n", err)
			return
		}
		defer db.CloseDatabse(idb)

		// Get les nouvelles shits
		ua, err := project.UpdateAccountInfo("berlingoqc")
		if err != nil {
			fmt.Printf("Error task github %v \n", err)
			return
		}
		idb.UpdateGHAccount(ua)
		repos, err := project.UpdateRepositoryInfo("berlingoqc", "YASE", "yawf")
		if err != nil {
			fmt.Printf("Error task github %v \n", err)
			return
		}
		// update dans la bd
		for _, r := range repos {
			idb.UpdateGitHubRepo(r)
		}
		w.Logger.Println("Github task over...")
	})
	w.TaskPool.Tasks["github"].Enable()
	w.TaskPool.Tasks["github"].Launch()

	w.Hs = &http.Server{
		Handler:      r,
		Addr:         addr,
		WriteTimeout: 15 * time.Second,
		ReadTimeout:  15 * time.Second,
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

// StartWebServer retourne une instance du serveur web deja
// instancier
func StartWebServer(assetPath string) (*WebServer, error) {

	ws := &WebServer{}

	c := config.WebSite{
		TemplateRoot: assetPath,
		StaticRoot:   "/home/wq/go/src/github.com/berlingoqc/yawf/cmd/yawf_website/static",
	}
	err := ws.Setup(c, nil)
	if err != nil {
		return nil, err
	}

	ws.Start()

	return ws, nil
}

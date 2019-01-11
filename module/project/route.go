package project

import (
	"errors"
	"net/http"

	"github.com/berlingoqc/yawf/config"
	"github.com/berlingoqc/yawf/db"
	"github.com/berlingoqc/yawf/website/route"
)

var (
	ModuleInfo = config.ModuleInfo{
		Name:        "project",
		Description: "With this module you can link your github account and show information about it",
	}
)

func GetModule() (string, config.IModule) {
	return "project", &Module{}
}

type Module struct {
	GitHubAccount string
	WatchedRepo   []string

	root  string
	getDb func() *ProjectDB
}

func (b *Module) Initialize(data map[string]interface{}) error {
	dbPath, err := config.GetDBFullPath(data)
	if err != nil {
		return err
	}
	// Recoit les informations sur la personne dans la map
	if d, ok := data[config.KeyRootFolder]; ok {
		b.root = d.(string)
	} else {
		return errors.New("Missing key " + config.KeyRootFolder)
	}

	b.getDb = func() *ProjectDB {
		idb, _ := b.GetDBInstanceFactory(dbPath)()
		return idb.(*ProjectDB)
	}

	// Doit recevoir les informations du user github
	// et les repo a watcher depuis la db ( ceux des projects )

	return nil
}

func (b *Module) GetInfo() config.ModuleInfo {
	return ModuleInfo
}

func (b *Module) GetNeededAssets() []string {
	return []string{
		"/template/project/index.html",
	}
}

func (b *Module) GetPackageName() string {
	return "github.com/berlingoqc/yawf/module/project"
}

func (b *Module) GetDBInstanceFactory(filepath string) func() (db.IDB, error) {
	return func() (db.IDB, error) {
		idb := &ProjectDB{}
		idb.Initialize(filepath)
		return idb, db.OpenDatabase(idb)
	}
}

type GitHubUser struct {
	Name string
}

func (b *Module) AddToWebServer(ws config.IWebServer) error {
	projectRouter := ws.GetMux().PathPrefix("/project").Subrouter()

	wPath := route.GetWPath("project", projectRouter,
		&route.RoutePath{ContentTmplPath: "/project/index.html", Path: "/"},
	)

	nb := ws.GetNavigationBar()

	li := route.Link{
		Name: "Project", URL: "/project/",
	}

	nb.Items = append(nb.Items, li)

	wPath.Route["/"].Handler = func(m map[string]interface{}, r *http.Request) {
		var err error
		idb := b.getDb()
		defer db.CloseDatabse(idb)
		m["Projects"], err = idb.GetPProjects()
		if err != nil {
			m["Error"] = err
			return
		}
		m["ghaccount"], err = idb.GetGHAccount()
		if err != nil {
			m["Error"] = err
		}
	}

	// Enregistre nos widgets
	mW := make(map[string]route.Widget)
	mW["gh_user_info"] = route.Widget{
		File:   "/shared/github_user_info.html",
		Name:   "gh_user_info",
		Struct: &GitHubUser{},
		Render: func(t interface{}, w http.ResponseWriter, r *http.Request) interface{} {

			// Get ma structure voulu
			gh := t.(*GitHubUser)
			print(gh)

			idb := b.getDb()
			defer db.CloseDatabse(idb)

			u, err := idb.GetGHAccount()
			if err != nil {
				route.RespondWithError(w, http.StatusBadRequest, err.Error())
				return nil
			}
			return u
		},
	}

	routeWidget := projectRouter.Path("/widget")

	route.ModuleWidgetAPI(routeWidget, b.root, mW)

	ws.AddRoute(wPath)

	return nil
}

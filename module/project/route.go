package project

import (
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

	getDb func() *ProjectDB
}

func (b *Module) Initialize(data map[string]interface{}) error {
	dbPath, err := config.GetDBFullPath(data)
	if err != nil {
		return err
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

func (b *Module) AddToWebServer(ws config.IWebServer) error {
	projectRouter := ws.GetMux().PathPrefix("/project").Subrouter()

	wPath := route.GetWPath("project", projectRouter,
		&route.RoutePath{ContentTmplPath: "/project/index.html", Path: "/"},
	)

	wPath.Route["/"].Handler = func(r *http.Request) map[string]interface{} {
		m := make(map[string]interface{})
		var err error
		idb := b.getDb()
		defer db.CloseDatabse(idb)
		m["Projects"], err = idb.GetPProjects()
		if err != nil {
			m["Error"] = err
			return m
		}
		m["ghaccount"], err = idb.GetGHAccount()
		if err != nil {
			m["Error"] = err
		}
		return m
	}

	ws.AddRoute(wPath)

	return nil
}

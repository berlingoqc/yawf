package user

import (
	"github.com/berlingoqc/yawf/api"
	"github.com/berlingoqc/yawf/auth"
	"github.com/berlingoqc/yawf/config"
	"github.com/berlingoqc/yawf/db"
	"github.com/berlingoqc/yawf/website/route"
)

var (
	ModuleInfo = config.ModuleInfo{
		Name:        "user",
		Description: "Add authentification for the website similar to Unix Right",
	}
)

func GetModule() (string, config.IModule) {
	return "user", &Module{}
}

type Module struct {
	getDb func() *auth.AuthDB
}

func (b *Module) Initialize(data map[string]interface{}) error {
	dbPath, err := config.GetDBFullPath(data)
	if err != nil {
		return err
	}
	b.getDb = func() *auth.AuthDB {
		idb, _ := b.GetDBInstanceFactory(dbPath)()
		return idb.(*auth.AuthDB)
	}

	return nil
}

func (b *Module) GetInfo() config.ModuleInfo {
	return ModuleInfo
}
func (b *Module) GetNeededAssets() []string {
	return []string{
		"/template/account/dashboard.html", "/template/account/dashboard_admin.html",
		"/template/login/confirm.html", "/template/login/login.html", "/template/login/new.html",
	}
}

func (b *Module) GetPackageName() string {
	return "github.com/berlingoqc/yawf/module/user"
}

func (b *Module) GetDBInstanceFactory(filepath string) func() (db.IDB, error) {
	return func() (db.IDB, error) {
		idb := &auth.AuthDB{}
		idb.Initialize(filepath)
		return idb, db.OpenDatabase(idb)
	}
}

func (b *Module) AddToWebServer(ws config.IWebServer) error {
	authRouter := ws.GetMux().PathPrefix("/auth").Subrouter()

	wPath := route.GetWPath("auth", authRouter,
		&route.RoutePath{ContentTmplPath: "/login/login.html", Path: "/login"},
		&route.RoutePath{ContentTmplPath: "/login/new.html", Path: "/new"},
		&route.RoutePath{ContentTmplPath: "/login/confirm.html", Path: "/confirm"},
	)
	// Ajoute notre info a la navbar
	nb := ws.GetNavigationBar()
	nb.Buttons = append(nb.Buttons, route.Button{
		Name:  "Account",
		Style: "btn-success",
		URL:   "/auth/login",
	})

	wPath.AddMiddleware(api.MiddlewareAuth)

	ws.AddRoute(wPath)

	accountRouter := ws.GetMux().PathPrefix("/account").Subrouter()

	aPath := route.GetWPath("account", accountRouter,
		&route.RoutePath{ContentTmplPath: "/account/dashboard_admin.html", Path: "/admin/"},
		&route.RoutePath{ContentTmplPath: "/account/dashboard.html", Path: "/dashboard"},
	)

	aPath.AddMiddleware(api.MiddlewareAccount)

	ws.AddRoute(aPath)

	apiRouter := ws.GetMux().PathPrefix("/api").Subrouter()

	// Ajout l'api pour authenfier les users
	userApi := &api.UserAPI{
		Db: b.getDb(),
	}

	err := userApi.Initialize(apiRouter)

	return err
}

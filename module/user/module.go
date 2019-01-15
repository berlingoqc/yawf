package user

import (
	"github.com/berlingoqc/yawf/api"
	"github.com/berlingoqc/yawf/auth"
	"github.com/berlingoqc/yawf/config"
	"github.com/berlingoqc/yawf/db"
	"github.com/berlingoqc/yawf/website/route"
	"github.com/gorilla/mux"
)

var (
	// ModuleInfo ...
	ModuleInfo = config.ModuleInfo{
		Name:        "user",
		Description: "Add authentification for the website similar to Unix Right",
	}
)

// GetModule ...
func GetModule() (string, config.IModule) {
	return "user", &Module{}
}

// Module ...
type Module struct {
	file *config.FileConfig
}

// Initialize ...
func (b *Module) Initialize(data map[string]interface{}) error {
	var err error
	b.file, err = config.GetFileConfig(data)
	if err != nil {
		return err
	}
	return nil
}

// GetInfo ...
func (b *Module) GetInfo() config.ModuleInfo {
	return ModuleInfo
}

// GetNeededAssets ...
func (b *Module) GetNeededAssets() []string {
	return []string{
		"/template/account/dashboard.html", "/template/account/dashboard_admin.html",
		"/template/login/confirm.html", "/template/login/login.html", "/template/login/new.html",
	}
}

// GetDBInstance ...
func (b *Module) GetDBInstance() (db.IDB, error) {
	idb := &auth.AuthDB{}
	idb.Initialize(b.file.GetDBFilePath())
	return idb, db.OpenDatabase(idb)
}

// GetNavigationItems ...
func (b *Module) GetNavigationItems() []interface{} {
	var ll []interface{}
	ll = append(ll, &route.Button{
		Name:  "Account",
		Style: "btn-success",
		URL:   "/auth/login",
	})
	return ll
}

// GetWidgets ...
func (b *Module) GetWidgets() []*route.Widget {
	return nil
}

// GetTasks ...
func (b *Module) GetTasks() []config.ITask {

	return nil
}

// GetWPath ...
func (b *Module) GetWPath(r *mux.Router) []*route.WPath {
	var ll []*route.WPath
	authRouter := r.PathPrefix("/auth").Subrouter()

	wPath := route.GetWPath("auth", authRouter,
		&route.RoutePath{ContentTmplPath: "/login/login.html", Path: "/login"},
		&route.RoutePath{ContentTmplPath: "/login/new.html", Path: "/new"},
		&route.RoutePath{ContentTmplPath: "/login/confirm.html", Path: "/confirm"},
	)
	wPath.AddMiddleware(api.MiddlewareAuth)

	ll = append(ll, wPath)

	accountRouter := r.PathPrefix("/account").Subrouter()

	aPath := route.GetWPath("account", accountRouter,
		&route.RoutePath{ContentTmplPath: "/account/dashboard_admin.html", Path: "/admin/"},
		&route.RoutePath{ContentTmplPath: "/account/dashboard.html", Path: "/dashboard"},
	)

	aPath.AddMiddleware(api.MiddlewareAccount)

	ll = append(ll, aPath)

	apiRouter := r.PathPrefix("/api").Subrouter()

	// Ajout l'api pour authenfier les users
	idb, _ := b.GetDBInstance()
	userAPI := &api.UserAPI{
		Db: idb.(*auth.AuthDB),
	}

	err := userAPI.Initialize(apiRouter)
	if err != nil {
		panic(err)
	}

	return ll
}

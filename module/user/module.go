package user

import (
	"github.com/berlingoqc/yawf/config"
	"github.com/berlingoqc/yawf/db"
	"github.com/berlingoqc/yawf/route"
	"github.com/berlingoqc/yawf/route/security"
	"github.com/gorilla/mux"
)

var (
	// ModuleInfo ...
	ModuleInfo = config.ModuleInfo{
		Name:        "user",
		Package:     "github.com/berlingoqc/yawf/module/user",
		RootPkg:     "github.com/berlingoqc/yawf",
		SubPkg:      "module/user",
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
		"/account/dashboard.html", "/account/dashboard_admin.html",
		"/login/confirm.html", "/login/login.html", "/login/new.html",
	}
}

// GetDBInstance ...
func (b *Module) GetDBInstance() (db.IDB, error) {
	idb := &AuthDB{}
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

	wPath := route.GetWPath("auth", authRouter)

	route.AddWPathItem(wPath,
		route.GetCPath("/login", "/login/login.html"),
		route.GetCPath("/new", "/login/new.html"),
		route.GetCPath("/confirm", "/login/confirm.html"),
	)

	ll = append(ll, wPath)

	accountRouter := r.PathPrefix("/account").Subrouter()

	aPath := route.GetWPath("/account", accountRouter)

	route.AddWPathItem(aPath,
		route.GetCPath("/admin", "/account/dashboard_admin.html"),
		route.GetCPath("/dashboard", "/account/dashboard.html"),
	)
	// Ajoute de la securiter globale pour le path account avec seulement
	// droit rw au user.
	aPath.Security = &security.PathSecurity{
		Owner:        "admin",
		Group:        "user",
		RoleRequired: security.RoleNormal,
		Right:        [3]security.Right{security.RightWrite & security.RightRead, security.RightRead, security.RightNone},
	}

	// Ajoute droit seulement au role admin au /admin

	ll = append(ll, aPath)

	return ll
}

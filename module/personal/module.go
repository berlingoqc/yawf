package personal

import (
	"net/http"

	"github.com/berlingoqc/yawf/config"
	"github.com/berlingoqc/yawf/db"
	"github.com/berlingoqc/yawf/website/route"
)

var (
	ModuleInfo config.ModuleInfo = config.ModuleInfo{
		Name:        "personal",
		Description: "This module allow you to show professional and personal information about yourself",
	}
)

func GetModule() (string, config.IModule) {
	return "personal", &Module{}
}

// PersonalModule is the module for the personal pages
type Module struct {
	ownerInfo *config.OwnerInformation
	root      string

	getDb func() *DB
}

func (b *Module) Initialize(data map[string]interface{}) error {
	// Recoit les informations sur la personne dans la map
	dbPath, err := config.GetDBFullPath(data)
	if err != nil {
		return err
	}
	b.getDb = func() *DB {
		idb, err := b.GetDBInstanceFactory(dbPath)()
		if err != nil {
			return nil
		}
		return idb.(*DB)
	}
	return nil
}

func (b *Module) GetPackageName() string {
	return "github.com/berlingoqc/yawf/module/personal"
}

func (b *Module) GetInfo() config.ModuleInfo {
	return ModuleInfo
}

func (b *Module) GetDBInstanceFactory(filepath string) func() (db.IDB, error) {
	return func() (db.IDB, error) {
		idb := &DB{}
		idb.Initialize(filepath)
		return idb, db.OpenDatabase(idb)
	}
}

func (b *Module) GetNeededAssets() []string {
	return []string{
		"/template/personal/me.html", "/template/personal/cv.html",
	}
}

func (b *Module) AddToWebServer(ws config.IWebServer) error {
	var err error
	r := ws.GetMux().PathPrefix("/me").Subrouter()
	wPath := route.GetWPath("personal", r,
		&route.RoutePath{ContentTmplPath: "/personal/me.html", Path: "/"},
		&route.RoutePath{ContentTmplPath: "/personal/cv.html", Path: "/cv"},
	)

	nb := ws.GetNavigationBar()

	dl := route.DropLink{
		Title: "Me",
		Links: make([][]route.Link, 2),
	}
	dl.Links[0] = make([]route.Link, 3)

	dl.Links[0][0] = route.Link{
		Name: "Index", URL: "/me/",
	}
	dl.Links[0][1] = route.Link{
		Name: "CV", URL: "/me/cv",
	}

	dl.Links[0][2] = route.Link{
		Name: "Myself", URL: "/me/personal",
	}

	nb.Items = append(nb.Items, dl)

	wPath.Route["/personal"] = route.GetMarkdownRouteFile(b.root+"/markdown/personal.md", "/personal")

	wPath.Route["/cv"].Handler = func(m map[string]interface{}, r *http.Request) {
		idb := b.getDb()
		defer db.CloseDatabse(idb)

		m["experiences"], err = idb.GetExperience()
		if err != nil {
			m["Error"] = err
			return
		}
		m["formations"], err = idb.GetFormation()
		if err != nil {
			m["Error"] = err
			return
		}
		m["languages"], err = idb.GetLanguageExperience()
		if err != nil {
			m["Error"] = err
			return
		}
	}

	ws.AddRoute(wPath)

	return nil
}

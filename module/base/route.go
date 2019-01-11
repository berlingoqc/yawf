package base

import (
	"log"
	"net/http"

	"github.com/berlingoqc/yawf/config"
	"github.com/berlingoqc/yawf/db"
	"github.com/berlingoqc/yawf/website/route"
)

var (
	ModuleInfo config.ModuleInfo = config.ModuleInfo{
		Name:        "base",
		Description: "Enable by default with this module you can display personal and professional informations about yourself",
	}
)

func GetModule() (string, config.IModule) {
	return "base", &Module{}
}

// PersonalModule is the module for the personal pages
type Module struct {
	ownerInfo *config.OwnerInformation
	root      string

	getDb func() *DB
}

func (b *Module) Initialize(data map[string]interface{}) error {
	// Recoit les informations sur la personne dans la map
	/*if d, ok := data[config.KeyOwnerInformation]; ok {
		b.ownerInfo = d.(*config.OwnerInformation)
	} else {
		return errors.New("Missing key " + config.KeyOwnerInformation)
	}
	if d, ok := data[config.KeyRootFolder]; ok {
		b.root = d.(string)
	} else {
		return errors.New("Missing key " + config.KeyRootFolder)
	}
	*/
	b.ownerInfo = &config.OwnerInformation{}
	b.getDb = func() *DB {
		idb, err := b.GetDBInstanceFactory(b.root + "/" + config.DBName)()
		if err != nil {
			return nil
		}
		return idb.(*DB)
	}
	return nil
}

func (b *Module) GetPackageName() string {
	return "github.com/berlingoqc/yawf/module/base"
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
		"/template/index.html", "/template/denied.html", "/template/error.html", "/template/me.html",
		"/template/cv.html", "/template/shared/layout.html", "/template/shared/footer.html",
		"/static/css/setupwizard.css", "/static/css/style_global.css",
	}
}

func (b *Module) AddToWebServer(ws config.IWebServer) error {
	var err error
	wPath := route.GetWPath("base", ws.GetMux(),
		&route.RoutePath{ContentTmplPath: "/index.html", Path: "/"},
		&route.RoutePath{ContentTmplPath: "/denied.html", Path: "/denied"},
		&route.RoutePath{ContentTmplPath: "/error.html", Path: "/error"},
		&route.RoutePath{ContentTmplPath: "/me.html", Path: "/me"},

		&route.RoutePath{ContentTmplPath: "/cv.html", Path: "/cv"},
	)

	if b.ownerInfo.ContactPage {
		wPath.Route["/contact"] = route.GetMarkdownRouteFile(b.root+"/markdown/contact.md", "/contact")
	}

	if b.ownerInfo.PersonalPage {
		wPath.Route["/personal"] = route.GetMarkdownRouteFile(b.root+"/markdown/personal.md", "/personal")
	}

	if b.ownerInfo.FAQ {
		wPath.Route["/faq"] = route.GetMarkdownRouteFile(b.root+"/markdown/faq.md", "/faq")
	}

	if b.ownerInfo.About {
		wPath.Route["/about"] = route.GetMarkdownRouteFile(b.root+"/markdown/about.md", "/about")
	}

	wPath.Route["/cv"].Handler = func(r *http.Request) map[string]interface{} {
		idb := b.getDb()
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

	ws.AddRoute(wPath)

	return nil
}

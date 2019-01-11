package base

import (
	"errors"
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
	if d, ok := data[config.KeyRootFolder]; ok {
		b.root = d.(string)
	} else {
		return errors.New("Missing key " + config.KeyRootFolder)
	}

	b.ownerInfo = &config.OwnerInformation{}
	err := config.MapToStruct(data, b.ownerInfo)
	if err != nil {
		return err
	}
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
		"/template/index.html", "/template/denied.html", "/template/error.html",
		"/template/shared/layout.html", "/template/shared/footer.html",
		"/static/css/setupwizard.css", "/static/css/style_global.css", "/template/shared/markdown_page.html",
	}
}

func (b *Module) AddToWebServer(ws config.IWebServer) error {
	wPath := route.GetWPath("base", ws.GetMux(),
		&route.RoutePath{ContentTmplPath: "/index.html", Path: "/"},
		&route.RoutePath{ContentTmplPath: "/denied.html", Path: "/denied"},
		&route.RoutePath{ContentTmplPath: "/error.html", Path: "/error"},
	)

	nb := ws.GetNavigationBar()

	dl := route.DropLink{
		Title: "Basic",
		Links: make([][]route.Link, 1),
	}
	dl.Links[0] = make([]route.Link, 3)

	dl.Links[0][0] = route.Link{
		Name: "Contact", URL: "/contact",
	}
	dl.Links[0][1] = route.Link{
		Name: "FAQ", URL: "/faq",
	}

	dl.Links[0][2] = route.Link{
		Name: "About", URL: "/about",
	}

	nb.Items = append(nb.Items, dl)

	if b.ownerInfo.ContactPage {
		wPath.Route["/contact"] = &route.RoutePath{ContentTmplPath: "/shared/markdown_page.html", Path: "/contact"}
		wPath.Route["/contact"].Handler = func(m map[string]interface{}, r *http.Request) {
			m["File"] = "contact.md"
		}
	}

	if b.ownerInfo.FAQ {
		wPath.Route["/faq"] = &route.RoutePath{ContentTmplPath: "/shared/markdown_page.html", Path: "/faq"}
		wPath.Route["/faq"].Handler = func(m map[string]interface{}, r *http.Request) {
			m["File"] = "faq.md"
		}
	}

	if b.ownerInfo.About {
		wPath.Route["/about"] = &route.RoutePath{ContentTmplPath: "/shared/markdown_page.html", Path: "/about"}
		wPath.Route["/about"].Handler = func(m map[string]interface{}, r *http.Request) {
			m["File"] = "about.md"
		}
	}

	pathMd := ws.GetMux().Path("/md")
	route.AddMarkdownFolderHandler(pathMd, b.root+"/markdown")

	ws.AddRoute(wPath)

	return nil
}

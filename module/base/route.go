package base

import (
	"github.com/berlingoqc/yawf/config"
	"github.com/berlingoqc/yawf/db"
	"github.com/berlingoqc/yawf/route"
	"github.com/gorilla/mux"
)

var (
	// ModuleInfo return the information about this module
	ModuleInfo config.ModuleInfo = config.ModuleInfo{
		Name:        "base",
		Package:     "github.com/berlingoqc/yawf/module/base",
		Description: "Enable by default with this module you can display personal and professional informations about yourself",
	}
)

// GetModule return the module
func GetModule() (string, config.IModule) {
	return "base", &Module{}
}

// Module is the module for the personal pages
type Module struct {
	ownerInfo *config.OwnerInformation
	config    *config.BaseConfig
	file      *config.FileConfig
}

// Initialize ...
func (b *Module) Initialize(data map[string]interface{}) error {
	// Recoit les informations sur la personne dans la map
	var err error
	b.file, err = config.GetFileConfig(data)
	if err != nil {
		return err
	}
	b.ownerInfo, err = config.GetOwnerInformation(data)
	if err != nil {
		return err
	}
	b.config, err = config.GetBaseConfig(data)
	if err != nil {
		return err
	}

	return nil
}

// GetInfo ...
func (b *Module) GetInfo() config.ModuleInfo {
	return ModuleInfo
}

// GetDBInstance ...
func (b *Module) GetDBInstance() (db.IDB, error) {
	idb := &DB{}
	idb.Initialize(b.file.GetDBFilePath())
	return idb, db.OpenDatabase(idb)
}

// GetNeededAssets ...
func (b *Module) GetNeededAssets() []string {
	return []string{
		"/index.html", "/denied.html", "/error.html",
		"/shared/layout.html", "/shared/footer.html",
		"/static/css/setupwizard.css", "/static/css/style_global.css", "/shared/markdown_page.html",
	}
}

// GetNavigationItems ...
func (b *Module) GetNavigationItems() []interface{} {
	var listItems []interface{}
	dl := route.DropLink{
		Title: "Home",
		Links: make([][]route.Link, 1),
	}
	dl.Links[0] = make([]route.Link, 1)

	dl.Links[0][0] = route.Link{
		Name: "Home", URL: "/",
	}

	if b.config.Contact {
		dl.Links[0] = append(dl.Links[0], route.Link{
			Name: "Contact", URL: "/contact",
		})
	}
	if b.config.FAQ {
		dl.Links[0] = append(dl.Links[0], route.Link{
			Name: "FAQ", URL: "/faq",
		})
	}
	if b.config.About {
		dl.Links[0] = append(dl.Links[0], route.Link{
			Name: "About", URL: "/about",
		})
	}
	listItems = append(listItems, dl)
	return listItems
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
	wPath := route.GetWPath("/", r)
	route.AddWPathItem(wPath,
		route.GetCPath("/", "/index.html"),
		route.GetMPathFile("/markdown", b.file.GetMarkdownPath()),
	)

	if b.config.Contact {
		c := route.GetCPathMarkdown("/contact", "contact.md")
		route.AddWPathItem(wPath, c)
	}

	if b.config.FAQ {
		c := route.GetCPathMarkdown("/faq", "faq.md")
		route.AddWPathItem(wPath, c)
	}

	if b.config.About {
		c := route.GetCPathMarkdown("/about", "about.md")
		route.AddWPathItem(wPath, c)
	}

	var ll []*route.WPath
	ll = append(ll, wPath)

	return ll
}

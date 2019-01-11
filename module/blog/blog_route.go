package blog

import (
	"net/http"
	"strconv"

	"github.com/berlingoqc/yawf/config"

	"github.com/berlingoqc/yawf/db"
	"github.com/russross/blackfriday/v2"

	"github.com/berlingoqc/yawf/website/route"
)

var (
	ModuleInfo config.ModuleInfo = config.ModuleInfo{
		Name:        "blog",
		Description: "With this module you can add a simple bloging engine to your website. The post are written in markdown",
	}
)

func GetModule() (string, config.IModule) {
	return "blog", &BlogModule{}
}

// BlogModule is the module to write blog on this website
type BlogModule struct {
	getDB func() *DB
}

func (b *BlogModule) Initialize(data map[string]interface{}) error {
	dbPath, err := config.GetDBFullPath(data)
	if err != nil {
		return err
	}
	b.getDB = func() *DB {
		idb, _ := b.GetDBInstanceFactory(dbPath)()
		return idb.(*DB)
	}
	return nil
}

func (b *BlogModule) GetInfo() config.ModuleInfo {
	return ModuleInfo
}
func (b *BlogModule) GetNeededAssets() []string {
	return []string{
		"/template/blog/index.html", "/template/blog/serie.html", "/template/blog/post.html",
	}
}
func (b *BlogModule) GetPackageName() string {
	return "github.com/berlingoqc/yawf/module/blog"
}

func (b *BlogModule) GetDBInstanceFactory(filepath string) func() (db.IDB, error) {
	return func() (db.IDB, error) {
		blogDb := &DB{}
		blogDb.Initialize(filepath)
		return blogDb, db.OpenDatabase(blogDb)
	}
}

func (b *BlogModule) AddToWebServer(ws config.IWebServer) error {
	blogRoute := ws.GetMux().PathPrefix("/blog").Subrouter()

	wp := route.GetWPath("blog", blogRoute,
		&route.RoutePath{ContentTmplPath: "/blog/index.html", Path: "/"},
		&route.RoutePath{ContentTmplPath: "/blog/serie.html", Path: "/serie"},
		&route.RoutePath{ContentTmplPath: "/blog/post.html", Path: "/post"},
	)

	wp.Route["/"].Handler = func(r *http.Request) map[string]interface{} {
		m := make(map[string]interface{})
		var err error
		idb := b.getDB()
		defer db.CloseDatabse(idb)
		m["Series"], err = idb.GetSerieList(true)
		if err != nil {
			//return nil
		}
		m["posts"], err = idb.GetBlogDescriptionList()
		if err != nil {
			return nil
		}
		return m
	}
	wp.Route["/post"].Handler = func(r *http.Request) map[string]interface{} {
		m := make(map[string]interface{})
		idb := b.getDB()
		defer db.CloseDatabse(idb)

		var err error
		// Si pas de query on retourne la listes des postes seulement
		idPost, ok := r.URL.Query()["id"]
		if !ok || len(idPost) == 0 {
			m["Posts"], _ = idb.GetBlogDescriptionList()
			return m
		}
		id, err := strconv.Atoi(idPost[0])
		if err != nil {
			m["Error"] = err
			return m
		}

		post, err := idb.GetBlogContent(id)
		if err != nil {
			m["Error"] = err
			return m
		}
		m["Post"] = post
		// Convertie le content byte markdown vers string(html)
		m["Content"] = string(blackfriday.Run(post.PostMarkdown))
		return m
	}

	wp.Route["/serie"].Handler = func(r *http.Request) map[string]interface{} {
		m := make(map[string]interface{})
		return m
	}

	ws.AddRoute(wp)

	return nil
}

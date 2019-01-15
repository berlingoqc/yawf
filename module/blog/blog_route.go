package blog

import (
	"io/ioutil"
	"net/http"
	"strconv"

	"github.com/berlingoqc/yawf/config"
	"github.com/gorilla/mux"

	"github.com/berlingoqc/yawf/db"
	"github.com/russross/blackfriday/v2"

	"github.com/berlingoqc/yawf/website/route"
)

var (
	// ModuleInfo ...
	ModuleInfo config.ModuleInfo = config.ModuleInfo{
		Name:        "blog",
		Package:     "github.com/berlingoqc/yawf/module/blog",
		Description: "With this module you can add a simple bloging engine to your website. The post are written in markdown",
	}
)

// Module is the module to write blog on this website
type Module struct {
	file *config.FileConfig
}

// Initialize ...
func (b *Module) Initialize(data map[string]interface{}) error {
	// Recoit les informations sur la personne dans la map
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
		"/template/blog/index.html", "/template/blog/serie.html", "/template/blog/post.html",
	}
}

// GetDBInstance ...
func (b *Module) GetDBInstance() (db.IDB, error) {
	blogDb := &DB{}
	blogDb.Initialize(b.file.GetDBFilePath())
	return blogDb, db.OpenDatabase(blogDb)
}

// GetNavigationItems ...
func (b *Module) GetNavigationItems() []interface{} {
	var ll []interface{}

	dl := route.DropLink{
		Title: "Blog",
		Links: make([][]route.Link, 1),
	}
	dl.Links[0] = make([]route.Link, 3)

	dl.Links[0][0] = route.Link{
		Name: "All", URL: "/blog/",
	}
	dl.Links[0][1] = route.Link{
		Name: "Series", URL: "/blog/serie",
	}

	dl.Links[0][2] = route.Link{
		Name: "Post", URL: "/blog/post",
	}

	ll = append(ll, dl)
	return ll
}

// GetTasks ...
func (b *Module) GetTasks() []config.ITask {
	return nil
}

// GetWPath ...
func (b *Module) GetWPath(r *mux.Router) []*route.WPath {
	var p []*route.WPath
	blogRouter := r.PathPrefix("/blog").Subrouter()
	wp := route.GetWPath("blog", blogRouter,
		&route.RoutePath{ContentTmplPath: "/blog/index.html", Path: "/"},
		&route.RoutePath{ContentTmplPath: "/blog/serie.html", Path: "/serie"},
		&route.RoutePath{ContentTmplPath: "/blog/post.html", Path: "/post"},
	)

	wp.Route["/"].Handler = func(m map[string]interface{}, r *http.Request) {
		var err error
		iidb, _ := b.GetDBInstance()
		idb := iidb.(*DB)
		defer db.CloseDatabse(idb)
		m["Series"], err = idb.GetSerieList(true)
		if err != nil {
			m["Error"] = err
			return
		}
		m["posts"], err = idb.GetBlogDescriptionList()
		if err != nil {
			m["Error"] = err
		}
	}
	wp.Route["/post"].Handler = func(m map[string]interface{}, r *http.Request) {
		iidb, _ := b.GetDBInstance()
		idb := iidb.(*DB)
		defer db.CloseDatabse(idb)

		var err error
		// Si pas de query on retourne la listes des postes seulement
		idPost, ok := r.URL.Query()["id"]
		if !ok || len(idPost) == 0 {
			m["Posts"], _ = idb.GetBlogDescriptionList()
			return
		}
		id, err := strconv.Atoi(idPost[0])
		if err != nil {
			m["Error"] = err
			return
		}

		post, err := idb.GetBlogContent(id)
		if err != nil {
			m["Error"] = err
			return
		}
		m["Post"] = post
		// Convertie le content byte markdown vers string(html)
		data, _ := ioutil.ReadFile("/home/wq/test.md")
		m["md"] = string(blackfriday.Run(data))
		return
	}

	wp.Route["/serie"].Handler = func(m map[string]interface{}, r *http.Request) {
	}

	// Ajoutes mes routs d'api
	apiBlog := r.PathPrefix("/api/blog").Subrouter()

	apiBlog.Path("/post").HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		// get seulement le content depuis la bd
		n, ok := r.URL.Query()["id"]
		if !ok || len(n) != 1 {
			w.Write([]byte("Error"))
			return
		}
		id, err := strconv.Atoi(n[0])
		if err != nil {
			w.Write([]byte(err.Error()))
			return
		}
		iidb, _ := b.GetDBInstance()
		idb := iidb.(*DB)
		defer db.CloseDatabse(idb)
		post, err := idb.GetBlogContent(id)
		if err != nil {
			w.Write([]byte(err.Error()))
			return
		}
		// Convertie le content byte markdown vers string(html)
		data := blackfriday.Run(post.PostMarkdown)
		w.Write(data)
	})
	p = append(p, wp)
	return p
}

package website

import (
	"log"
	"os"
	"path/filepath"
	"testing"

	"github.com/berlingoqc/yawf/auth"
	"github.com/berlingoqc/yawf/config"
	"github.com/berlingoqc/yawf/db"
	"github.com/berlingoqc/yawf/module/base"
	"github.com/berlingoqc/yawf/module/blog"
	"github.com/berlingoqc/yawf/module/project"
)

var (
	goPath         = "/home/wq/go"
	directoryAsset = ""
	wantedModule   = make([]config.IModule, 5)
)

func init() {
	directoryAsset, _ = os.Getwd()
	i := 0
	for _, v := range DefaultModule {
		wantedModule[i] = v
		i++
	}

}

func getAuthDB() *auth.AuthDB {
	a := &auth.AuthDB{}
	a.Initialize(config.GetDBFullPathRoot(directoryAsset + "/root"))
	if err := db.OpenDatabase(a); err != nil {
		log.Fatal(err)
	}

	return a
}

func getBlogDB() *blog.DB {
	a := &blog.DB{}
	a.Initialize(config.GetDBFullPathRoot(directoryAsset + "/root"))
	if err := db.OpenDatabase(a); err != nil {
		panic(err)
	}
	return a
}

func getProjectDB() *project.ProjectDB {
	a := &project.ProjectDB{}
	a.Initialize(config.GetDBFullPathRoot(directoryAsset + "/root"))
	if err := db.OpenDatabase(a); err != nil {
		panic(err)
	}
	return a
}

func getBaseDB() *base.DB {
	a := &base.DB{}
	a.Initialize(config.GetDBFullPathRoot(directoryAsset + "/root"))
	if err := db.OpenDatabase(a); err != nil {
		panic(err)
	}
	return a
}

func RemoveContents(dir string) error {
	files, err := filepath.Glob(filepath.Join(dir, "*"))
	if err != nil {
		return err
	}
	for _, file := range files {
		err = os.RemoveAll(file)
		if err != nil {
			return err
		}
	}
	return nil
}

func TestInitializeWebSite(t *testing.T) {
	RemoveContents(directoryAsset + "/root")
	err := config.Copy(directoryAsset+"/asset", directoryAsset+"/root")
	if err != nil {
		t.Fatal(err)
	}
	// Crée mon config
	wsConfig := &config.WebSite{
		Name:         "wquintal.ca",
		EnableModule: make(map[string]map[string]interface{}),
		NetOptions:   make(map[string]interface{}),
		RootFolder:   directoryAsset + "/root",
		AppUsers:     make(map[string]interface{}),
		Owner: config.OwnerInformation{
			About:        true,
			FAQ:          true,
			ContactPage:  true,
			Birth:        "1995/07/05",
			FullName:     "William Quintal",
			Location:     "Quebec",
			ThumbnailURL: "/static/img/me.png",
		},
	}

	data := config.GetMapFromConfig(wsConfig)

	for k, _ := range DefaultModule {
		wsConfig.EnableModule[k] = data
	}

	wsConfig.NetOptions[config.KeyPort] = "9090"

	err = config.Save(directoryAsset+"/root", wsConfig)
	if err != nil {
		t.Fatal(err)
	}

	wsConfig, err = config.Load(directoryAsset + "/root")
	if err != nil {
		t.Fatal(err)
	}
	t.Log(wsConfig)

	err = config.ImportModuleAsset(goPath, directoryAsset+"/root", wantedModule)
	if err != nil {
		t.Fatal(err)
	}

	// Get toutes les bd
	baseDb := getBaseDB()
	db.CloseDatabse(baseDb)
	authDb := getAuthDB()
	db.CloseDatabse(authDb)
	pdb := getProjectDB()
	db.CloseDatabse(pdb)
	bdb := getBlogDB()
	db.CloseDatabse(bdb)

	ws, err := StartWebServer(directoryAsset + "/root")
	if err != nil {
		t.Fatal(err)
	}

	ws.Stop()

}

func TestWebSiteDefault(t *testing.T) {
	// Valide que le nombre de module disponible est egale a 4
	if len(DefaultModule) != 5 {
		t.Fatal("There should be 5 default module")
	}

	wantedModule := make([]config.IModule, 5)
	i := 0
	for _, v := range DefaultModule {
		wantedModule[i] = v
		i++
	}

	directoryAsset, _ := os.Getwd()
	directoryAsset += "/root"
	goPath := "/home/wq/go"

	err := config.ImportModuleAsset(goPath, directoryAsset, wantedModule)
	if err != nil {
		t.Fatal(err)
	}

}

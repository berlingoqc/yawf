package config

import (
	"encoding/json"
	"errors"
	"io/ioutil"
	"log"
	"net/http"
	"time"

	"github.com/berlingoqc/yawf/db"
	"github.com/berlingoqc/yawf/website/route"
	"github.com/gorilla/mux"
)

const (
	KeyOwnerInformation = "OwnerInformation"
	KeyRootFolder       = "key_root_folder"
	KeyNetOption        = "key_net"
	KeyUserOption       = "key_user"

	KeyPort = "key_net_port"
	KeyAddr = "key_net_addr"

	ConfigFileName = "yawf.conf"
	DBName         = "yawf.db"
)

type IModule interface {
	Initialize(data map[string]interface{}) error
	GetInfo() ModuleInfo
	GetDBInstanceFactory(filepath string) func() (db.IDB, error)

	AddToWebServer(ws IWebServer) error
	GetNeededAssets() []string
	GetPackageName() string
}

var (
	AvailableModule = make(map[string]IModule)
)

type IWebServer interface {
	GetTaskPool() *TaskPool
	GetLogger() *log.Logger
	GetMux() *mux.Router

	AddRoute(r *route.WPath)
}

// ModuleInfo contient l'informations sur les modules
// disponible dans le framework
type ModuleInfo struct {
	Name        string
	Description string
}

// OwnerInformation contient les informations a afficher de base
// sur le owner du site
type OwnerInformation struct {
	FullName     string
	Birth        string
	Location     string
	ThumbnailURL string

	ContactPage  bool
	PersonalPage bool
	About        bool
	FAQ          bool
}

type WebSite struct {
	// RootFolder where all the files are store
	RootFolder string

	// Name is the website
	Name string

	// Owner contains informations about the owner of the website
	Owner OwnerInformation

	// NetOptions contains all the networking information about the webserver
	NetOptions map[string]interface{}

	// AppUsers contains informations about the connection informations
	// for all the service of the website like github or linked
	AppUsers map[string]interface{}

	// EnableModule contains the module enable key = name and ther options
	EnableModule map[string]map[string]interface{}
}

func Save(root string, w *WebSite) error {
	b, err := json.Marshal(w)
	if err != nil {
		return err
	}
	return ioutil.WriteFile(root+"/"+ConfigFileName, b, 0640)
}

func Load(root string) (*WebSite, error) {
	b, err := ioutil.ReadFile(root + "/" + ConfigFileName)
	if err != nil {
		return nil, err
	}
	ws := &WebSite{}
	return ws, json.Unmarshal(b, ws)
}

func ValidFolder() error {

	return nil
}

func GetWebServer(opt map[string]interface{}) *http.Server {
	ws := &http.Server{
		WriteTimeout: 15 * time.Second,
		ReadTimeout:  15 * time.Second,
		Addr:         ":8081",
	}
	if d, ok := opt[KeyAddr]; ok {
		ws.Addr = d.(string)
	}
	if d, ok := opt[KeyPort]; ok {
		ws.Addr = ":" + d.(string)
	}

	return ws
}

func GetMapFromConfig(w *WebSite) map[string]interface{} {
	m := make(map[string]interface{})
	m[KeyOwnerInformation] = w.Owner
	m[KeyRootFolder] = w.RootFolder
	m[KeyNetOption] = w.NetOptions
	m[KeyUserOption] = w.AppUsers
	return m
}

func GetDBFullPath(m map[string]interface{}) (string, error) {
	if d, ok := m[KeyRootFolder]; ok {
		return d.(string) + "/" + DBName, nil
	}
	return "", errors.New("Key dont exists for RootFolder")
}

func GetDBFullPathRoot(root string) string {
	return root + "/" + DBName
}

// ImportModuleAsset import the assets of all module to the rootPath
func ImportModuleAsset(gopath string, rootPath string, wantedModules []IModule) error {
	assetPath := rootPath
	gopath += "/src"

	for _, m := range wantedModules {
		goPath := gopath + "/" + m.GetPackageName()

		files := m.GetNeededAssets()
		for _, f := range files {
			// Copy tout les fichiers vers la destination
			dest := assetPath + f
			f = goPath + "/asset" + f
			// Copy vers le asset path
			err := Copy(f, dest)
			if err != nil {
				return err
			}
		}

	}

	return nil
}

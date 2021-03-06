package config

import (
	"errors"
	"fmt"
	"log"
	"net/http"
	"time"

	"github.com/berlingoqc/yawf/conv"

	"github.com/berlingoqc/yawf/utility"
)

// Ctx represent the context that is pass across the application
// to configure correctly the different part with the informations
// that should be provides to them
type Ctx map[string]interface{}

// GetWebServer try to find is the webConfig struct is present
// in the context
func GetWebServer(c Ctx) (*http.Server, error) {
	no := &NetworkOptions{}
	err := conv.FindStructMap(c, no)
	if err != nil {
		return nil, err
	}
	ws := &http.Server{
		Addr:         fmt.Sprintf("%v:%v", no.Addr, no.Port),
		WriteTimeout: time.Duration(no.TimeoutLength) * time.Second,
		ReadTimeout:  time.Duration(no.TimeoutLength) * time.Second,
	}

	return ws, nil
}

// GetOwnerInformation try to find if the OwnerInformations is in the map
func GetOwnerInformation(c Ctx) (*OwnerInformation, error) {
	t := &OwnerInformation{}
	return t, conv.FindStructMap(c, t)
}

// GetFileConfig try to find if the FileConfig is in the map ( should always be )
func GetFileConfig(c Ctx) (*FileConfig, error) {
	f := &FileConfig{}
	return f, conv.FindStructMap(c, f)
}

// GetBaseConfig try to find the BaseConfig
func GetBaseConfig(c Ctx) (*BaseConfig, error) {
	b := &BaseConfig{}
	return b, conv.FindStructMap(c, b)
}

// GetEnableModule try to find the map of informations to provides to the
// modules enabled
func GetEnableModule(c Ctx) (map[string]map[string]interface{}, error) {
	if d, ok := c["EnableModule"]; ok {
		return d.(map[string]map[string]interface{}), nil
	}

	return nil, &conv.KeyError{
		Name:   "EnableModule",
		Status: conv.NotFound,
	}
}

// OwnerInformation contient les informations a afficher de base
// sur le owner du site
type OwnerInformation struct {
	// FullName of the owner of the website
	FullName string `map:"required"`
	// Birth is the day of birth of the owner
	Birth string
	// Location is the approximal location of the owner like (state, country or town)
	Location string
	// ThumbnailURL is a link to the user thumbnail for the home page
	ThumbnailURL string
}

// NetworkOptions are the configuration for the http.Server networking configuration
type NetworkOptions struct {
	Addr          string
	Port          int
	TimeoutLength int
}

// FileConfig contains the information about the location
// of the different composant of the website
type FileConfig struct {
	RootFolder     string
	DBFile         string
	AssetFolder    string
	StaticFolder   string
	ConfigFile     string
	MarkdownFolder string
	ModuleFolder   string
}

// GetModulePath return the full path of the module folder
func (f *FileConfig) GetModulePath() string {
	return f.RootFolder + "/" + f.ModuleFolder
}

// GetMarkdownPath return the full path of the configuration file
func (f *FileConfig) GetMarkdownPath() string {
	return f.RootFolder + "/" + f.MarkdownFolder
}

// GetConfigPath return the full path of the configuration file
func (f *FileConfig) GetConfigPath() string {
	return f.RootFolder + "/" + f.ConfigFile
}

// GetRootFolder retourne the root folder where the files are stored
func (f *FileConfig) GetRootFolder() string {
	return f.RootFolder
}

// GetDBFilePath return the full path of the database file
func (f *FileConfig) GetDBFilePath() string {
	return f.RootFolder + "/" + f.DBFile
}

// GetAssetFolderPath return the full path of the asset folder
func (f *FileConfig) GetAssetFolderPath() string {
	return f.RootFolder + "/" + f.AssetFolder
}

// GetStaticFolderPath return the fulll path of the static folder
func (f *FileConfig) GetStaticFolderPath() string {
	return f.GetAssetFolderPath() + "/" + f.StaticFolder
}

// GetFolderPathList return the fullpath of all required folder
func (f *FileConfig) GetFolderPathList() []string {
	return []string{
		f.GetRootFolder(),
		f.GetAssetFolderPath(),
		f.GetStaticFolderPath(),
		f.GetMarkdownPath(),
		f.GetMarkdownPath(),
	}
}

// BaseConfig represent the configuration of the base module
type BaseConfig struct {
	Contact bool
	About   bool
	FAQ     bool
}

// BuildModule represent the information for a build module that can be loaded
type BuildModule struct {
	Name    string
	Version string
	Path    string

	Info *ModuleInfo
}

// WebSite is the struct that hold all the information that is provides by the main configuration file
type WebSite struct {
	// File is the struct that contains the informations where the config
	// files are located
	File *FileConfig
	// Name is the website
	Name string

	// Owner contains informations about the owner of the website
	Owner *OwnerInformation

	// NetOptions contains all the networking information about the webserver
	NetOptions *NetworkOptions

	// AppUsers contains informations about the connection informations
	// for all the service of the website like github or linked
	AppUsers map[string]interface{}

	// EnableModule contains the module enable key = name and ther options
	EnableModule map[string]interface{}

	// AvailableModule contains the informations of the shared module install in this computer
	AvailableModule map[string]*BuildModule
}

// ToContext convert the struct to a my context map
func (w *WebSite) ToContext() Ctx {
	m := make(map[string]interface{})
	conv.AddStructToMap(m, w.File)
	conv.AddStructToMap(m, w.Owner)
	conv.AddStructToMap(m, w.NetOptions)
	m["Name"] = w.Name
	m["AppUsers"] = w.AppUsers
	m["EnableModule"] = w.EnableModule
	m["AvailableModule"] = w.AvailableModule
	return m
}

// Validate ensure that all required struct are well configurate and that the root folder exists
func (w *WebSite) Validate() error {

	// Validate that the root folder have all the required things

	// Ensure that all the module are loaded
	for k := range w.EnableModule {
		// Regarde s'il est deja loader ( comme pour les modules de base)
		if m := GetModule(k); m == nil {
			// if not load watch if it's in the available shared module
			// list of the config
			if bm, ok := w.AvailableModule[k]; ok {
				fmt.Printf("Loading module %v at %v\n", bm.Name, bm.Path)
				_, mod, err := LoadModuleDynamicly(bm.Path)
				if err != nil {
					// Si l'erreur est de type Version
					// on peut le rebuilder
					if _, ok := err.(*VersionError); ok {
						log.Printf("Versioning error with existing module rebuilding it\n")
						info := bm.Info
						bm, err = GoBuildModule(info.RootPkg, info.SubPkg, w.File.GetModulePath())
						if err != nil {
							return err
						}
						log.Printf("Build successful, loading again\n")
						if _, _, err = LoadModuleDynamicly(bm.Path); err != nil {
							return err
						}

					}
					return err
				}
				modInfo := mod.GetInfo()
				bm.Info = &modInfo
			} else {
				return errors.New("Can't find module " + k)
			}
		}
	}
	log.Printf("Module is loaded\n")

	return nil
}

// Save the configuration file to path inside File
func (w *WebSite) Save() error {
	ctx := w.ToContext()
	return conv.Save(w.File.GetConfigPath(), ctx)
}

// LoadWebSiteConfig load the configuration of the website
func LoadWebSiteConfig(filePath string) (*WebSite, error) {
	m, err := conv.Load(filePath)
	if err != nil {
		return nil, err
	}

	ws := &WebSite{
		Owner:      &OwnerInformation{},
		NetOptions: &NetworkOptions{},
		File:       &FileConfig{},
	}
	if n, ok := m["Name"]; ok {
		ws.Name = n.(string)
	} else {
		return nil, errors.New("No name")
	}

	err = conv.FindStructMap(m, ws.Owner)
	if err != nil {
		return nil, err
	}
	err = conv.FindStructMap(m, ws.File)
	if err != nil {
		return nil, err
	}
	err = conv.FindStructMap(m, ws.NetOptions)
	if err != nil {
		return nil, err
	}

	ws.AppUsers = m["AppUsers"].(map[string]interface{})
	ws.EnableModule = m["EnableModule"].(map[string]interface{})
	ws.AvailableModule = make(map[string]*BuildModule)
	for k, v := range m["AvailableModule"].(map[string]interface{}) {
		t := &BuildModule{
			Info: &ModuleInfo{},
		}
		err := conv.MapToStruct(v.(map[string]interface{}), t)
		if err != nil {
			return nil, err
		}
		mapInfo := v.(map[string]interface{})["Info"]
		err = conv.MapToStruct(mapInfo.(map[string]interface{}), t.Info)
		if err != nil {
			return nil, err
		}
		log.Print(*t)
		ws.AvailableModule[k] = t

	}

	return ws, nil
}

// ImportModuleAsset import the assets of all module to the rootPath
func ImportModuleAsset(gopath string, rootPath string, wantedModules []IModule) error {
	assetPath := rootPath
	gopath += "/src"

	for _, m := range wantedModules {
		goPath := gopath + "/" + m.GetInfo().Package

		files := m.GetNeededAssets()
		for _, f := range files {
			// Copy tout les fichiers vers la destination
			dest := assetPath + f
			f = goPath + "/asset" + f
			// Copy vers le asset path
			err := utility.Copy(f, dest)
			if err != nil {
				return err
			}
		}

	}

	return nil
}

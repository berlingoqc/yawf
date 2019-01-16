package website

import (
	"os"

	"github.com/berlingoqc/yawf/auth"
	"github.com/berlingoqc/yawf/config"
	"github.com/berlingoqc/yawf/conv"
	"github.com/berlingoqc/yawf/db"
	"github.com/berlingoqc/yawf/utility"
)

func init() {
	// A l'importation du module ajoute les modules de bases a la listes
	// des disponibles
}

// NewWebSiteConfig initialize a new configuration for the website
func NewWebSiteConfig(name string, rootdirectory string, password string, baseConfig *config.BaseConfig, owner *config.OwnerInformation) (*config.WebSite, error) {
	wsConfig := &config.WebSite{
		Name:            name,
		EnableModule:    make(map[string]interface{}),
		AvailableModule: make(map[string]*config.BuildModule),
		NetOptions: &config.NetworkOptions{
			Addr:          "",
			Port:          8081,
			TimeoutLength: 15,
		},
		File: &config.FileConfig{
			RootFolder:     rootdirectory,
			AssetFolder:    "asset",
			StaticFolder:   "static",
			MarkdownFolder: "markdown",
			ModuleFolder:   "module",
			DBFile:         name + ".db",
			ConfigFile:     name + ".json",
		},
		AppUsers: make(map[string]interface{}),
		Owner:    owner,
	}
	// Ajout le base module avec ca configuration
	optBase := make(map[string]interface{})
	conv.AddStructToMap(optBase, baseConfig)
	wsConfig.EnableModule["base"] = optBase

	// Valide le dossier ( doit etre vide )
	if err := utility.IsDirectoryEmpty(rootdirectory); err != nil {
		return nil, err
	}
	folder := []string{wsConfig.File.AssetFolder, "asset/template", "markdown", "module"}
	for _, f := range folder {
		f = rootdirectory + "/" + f
		err := os.MkdirAll(f, 0744)
		if err != nil {
			return nil, err
		}
	}

	ctx := wsConfig.ToContext()
	BaseModule.Initialize(ctx)
	UserModule.Initialize(ctx)

	// Create the database file with the table from our module base and auth
	bdb, err := BaseModule.GetDBInstance()
	if err != nil {
		return nil, err
	}
	db.CloseDatabse(bdb)
	bdb, err = UserModule.GetDBInstance()
	if err != nil {
		return nil, err
	}
	defer db.CloseDatabse(bdb)
	userdb := bdb.(*auth.AuthDB)
	_, err = userdb.CreateAdminAccount(password)

	return wsConfig, nil
}

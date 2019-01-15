package website

import (
	"os"

	"github.com/berlingoqc/yawf/config"
	"github.com/berlingoqc/yawf/conv"
	"github.com/berlingoqc/yawf/utility"
)

func init() {
	// A l'importation du module ajoute les modules de bases a la listes
	// des disponibles
}

// NewWebSiteConfig initialize a new configuration for the website
func NewWebSiteConfig(name string, rootdirectory string, baseConfig *config.BaseConfig, owner *config.OwnerInformation) (*config.WebSite, error) {
	wsConfig := &config.WebSite{
		Name:         name,
		EnableModule: make(map[string]interface{}),
		NetOptions: &config.NetworkOptions{
			Addr:          "",
			Port:          8081,
			TimeoutLength: 15,
		},
		File: &config.FileConfig{
			RootFolder:   rootdirectory,
			AssetFolder:  "asset",
			StaticFolder: "static",
			DBFile:       name + ".db",
			ConfigFile:   name + ".json",
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
	folder := []string{wsConfig.File.AssetFolder, "asset/template", "markdown"}
	for _, f := range folder {
		f = rootdirectory + "/" + f
		err := os.MkdirAll(f, 0744)
		if err != nil {
			return nil, err
		}
	}

	return wsConfig, nil
}

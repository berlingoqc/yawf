package website

import (
	"os"
	"path/filepath"
	"testing"

	"github.com/berlingoqc/yawf/config"
	"github.com/berlingoqc/yawf/module/base"
	"github.com/berlingoqc/yawf/module/user"
)

var (
	directoryAsset = ""
	wantedModule   = make([]config.IModule, 5)
)

func init() {
	directoryAsset, _ = os.Getwd()
	var name string
	name, BaseModule = base.GetModule()
	config.AddAvailableModule(name, BaseModule)
	name, UserModule = user.GetModule()
	config.AddAvailableModule(name, UserModule)
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

func TestLoadWebSite(t *testing.T) {
	ws, err := StartWebServer("/home/wq/go/src/github.com/berlingoqc/yawf/website/root/wquintal.json")

	if err != nil {
		t.Fatal(err)
	}

	ws.Stop()
}
func TestInitializeWebSite(t *testing.T) {
	directoryAsset += "/root"
	RemoveContents(directoryAsset)

	/*
		buildBlog := &config.BuildModule{
			Name: "blog",
			Path: "/home/wq/go/src/github.com/berlingoqc/yawf/website/root/module/blog.so",

			Info: &config.ModuleInfo{
				RootPkg: "github.com/berlingoqc/yawf",
				SubPkg:  "module/blog",
				Package: "github.com/berlingoqc/yawf/module/blog",
			},
		}

	*/
	wsConfig, err := NewWebSiteConfig("wquintal", directoryAsset, "office", &config.BaseConfig{
		About: true,
		FAQ:   true,
	}, &config.OwnerInformation{})
	if err != nil {
		t.Fatal(err)
	}

	/*
		wsConfig.AvailableModule["blog"] = buildBlog
		wsConfig.EnableModule["blog"] = make(map[string]interface{})
	*/
	err = wsConfig.Save()
	if err != nil {
		t.Fatal(err)
	}
	ws, err := StartWebServer(wsConfig.File.GetConfigPath())
	if err != nil {
		t.Fatal(err)
	}

	ws.Stop()

}

package website

import (
	"os"
	"path/filepath"
	"testing"

	"github.com/berlingoqc/yawf/config"
)

var (
	directoryAsset = ""
	wantedModule   = make([]config.IModule, 5)
)

func init() {
	directoryAsset, _ = os.Getwd()
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
	directoryAsset += "/root"
	RemoveContents(directoryAsset)

	n, _, err := config.LoadModuleDynamicly("/home/wq/go/src/github.com/berlingoqc/yawf/module/blog/module/module.so")
	if err != nil {
		t.Fatal(err)
	}
	t.Log(n)

	wsConfig, err := NewWebSiteConfig("wquintal", directoryAsset, &config.BaseConfig{
		About: true,
		FAQ:   true,
	}, &config.OwnerInformation{})
	if err != nil {
		t.Fatal(err)
	}
	wsConfig.EnableModule[n] = make(map[string]interface{})
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

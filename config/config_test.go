package config

import (
	"os"
	"testing"
)

const (
	testDir = "/test"
)

func TestWebSiteConfig(t *testing.T) {
	// Crée un repertoire pour mes affaires

	// Defer le delete du repertoire
	workingDir, _ := os.Getwd()
	workingDir += testDir
	err := os.Mkdir(workingDir, 0644)
	if err != nil {
		t.Fatal(err)
	}
	defer func() {
		os.RemoveAll(workingDir)
	}()
	/*
		baseConfig := &BaseConfig{}
		ownerInfo := &OwnerInformation{
			Birth:        "1995/07/05",
			FullName:     "William Quintal",
			Location:     "Quebec",
			ThumbnailURL: "/static/img/me.png",
		}
		wsConfig, err := NewWebSiteConfig("wquintal.ca", workingDir, baseConfig, ownerInfo)
		if err != nil {
			t.Fatal(err)
		}
		err = wsConfig.Save()
		if err != nil {
			t.Fatal(err)
		}

		wsConfig, err = LoadWebSiteConfig("./")
	*/
}

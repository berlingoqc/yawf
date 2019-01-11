package cli

import (
	"errors"
	"io"
	"log"
	"net/http"
	"os"

	"github.com/berlingoqc/yawf/db"
	"github.com/berlingoqc/yawf/website"

	"github.com/berlingoqc/yawf/api"
	"github.com/berlingoqc/yawf/auth"
	"github.com/berlingoqc/yawf/config"
	"github.com/berlingoqc/yawf/module/base"
	"github.com/berlingoqc/yawf/module/blog"
	"github.com/berlingoqc/yawf/module/project"
	ishell "gopkg.in/abiosoft/ishell.v2"
)

var (
	// Dossier ou son sauvegarder les trucs
	AppFileFolder  = "./web_root"
	ExecutableRoot = "/home/wq/go/src/github.com/berlingoqc/yawf/cmd/wquintal_website"

	// AuthCookie est mon cookis pour m'authentifier lors des requetes
	AuthCookie *http.Cookie
)

func getAuthDB() *auth.AuthDB {
	a := &auth.AuthDB{}
	a.Initialize(config.GetDBFullPathRoot(AppFileFolder))
	if err := db.OpenDatabase(a); err != nil {
		log.Fatal(err)
	}

	return a
}

func getBlogDB() *blog.DB {
	a := &blog.DB{}
	a.Initialize(config.GetDBFullPathRoot(AppFileFolder))
	db.OpenDatabase(a)
	return a
}

func getProjectDB() *project.ProjectDB {
	a := &project.ProjectDB{}
	a.Initialize(config.GetDBFullPathRoot(AppFileFolder))
	db.OpenDatabase(a)
	return a
}

func getBaseDB() *base.DB {
	a := &base.DB{}
	a.Initialize(config.GetDBFullPathRoot(AppFileFolder))
	db.OpenDatabase(a)
	return a
}

func isDirectoryEmpty(direc string) error {
	f, err := os.Open(direc)
	if err != nil {
		return err
	}
	defer f.Close()
	_, err = f.Readdir(1)
	if err == io.EOF {
		return nil
	}
	return errors.New(direc + " is not empty")
}

func InitNewWebsite(rootDirectory string, pw string, c *ishell.Context) error {
	// Valide que le repertoire existes et qu'il est vide
	if err := isDirectoryEmpty(rootDirectory); err != nil {
		return err
	}

	wsConfig := &config.WebSite{
		EnableModule: make(map[string]map[string]interface{}),
	}

	// Demande pour le nom du site web
	c.Print("Domain name of the website (ex : wafg.ca) : ")
	wsConfig.Name = c.ReadLine()

	c.Print("Enter your full name : ")
	wsConfig.Owner.FullName = c.ReadLine()

	// Demande qu'elle page de base je veux d'enable

	// Demande pour le nom de l'usager

	// Demade les modules qu'on veut activer
	var name_list []string
	for k, _ := range website.DefaultModule {
		name_list = append(name_list, k)
	}

	choicesModule := c.Checklist(name_list, "Choice the module you wan't", nil)
	for _, i := range choicesModule {
		wsConfig.EnableModule[name_list[i]] = nil
	}

	bd := getAuthDB()
	defer bd.CloseDatabase()

	// copie les dossiers vers le repertoire
	err := Copy("../wquintal_website/"+"markdown", rootDirectory+"/markdown")
	if err != nil {
		return err
	}
	err = Copy("../wquintal_website/"+"template", rootDirectory+"/template")
	if err != nil {
		return err
	}

	err = Copy("../wquintal_website/"+"static", rootDirectory+"/static")
	if err != nil {
		return err
	}

	_, err = bd.CreateAdminAccount(pw)
	if err != nil {
		return err
	}

	c.Print("Admin account create with success")

	return nil
}

func InitWebsiteCMD(c *ishell.Context) {
	c.ShowPrompt(false)
	defer c.ShowPrompt(true)

	c.Println("Welcome to the YAWF Shell. You can initialize a new website or load a existing one")
	choice := c.MultiChoice([]string{
		"Load context",
		"Initialize a new",
	}, "You can initialize a new website or load a existing one")
	c.Print("Enter the root directory\n")
	rootDirectory := c.ReadLine()
	AppFileFolder = rootDirectory
	c.Print("Enter admin password: ")
	pw := c.ReadPassword()

	var bd *auth.AuthDB

	if choice == 1 {
		if err := InitNewWebsite(rootDirectory, pw, c); err != nil {
			c.Printf("Error %v\n", err)
			return
		}
	}
	bd = getAuthDB()
	defer bd.CloseDatabase()
	user, err := bd.LoginUser("admin", pw)
	if err != nil {
		c.Println("Error login admin : %v\n", err.Error())
		return
	}

	AuthCookie, err = api.SetCookieForUser(user)
	if err != nil {
		c.Println("Error optaining cookie : %v\n", err.Error())
	}

	c.Println("Login as admin")

}

// LoginAdminCMD doit etre executer pour faire les operations sur le site
func LoginAdminCMD(c *ishell.Context) {
	c.ShowPrompt(false)
	defer c.ShowPrompt(true)

	bd := getAuthDB()

	defer bd.CloseDatabase()

	c.Print("Enter admin password : ")
	pw := c.ReadPassword()

	user, err := bd.LoginUser("admin", pw)
	if err != nil {
		c.Printf("Error login : %v\n", err.Error())
		return
	}
	AuthCookie, _ = api.SetCookieForUser(user)
	c.Println("Admin account login without error")
}

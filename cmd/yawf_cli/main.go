package main

import (
	"net/http"

	"github.com/berlingoqc/yawf/cli"
	ishell "gopkg.in/abiosoft/ishell.v2"
)

var (
	// AuthCookie est mon cookie pour m'authentifier lors des requetes
	AuthCookie *http.Cookie
)

func main() {

	shell := ishell.New()

	shell.AddCmd(&ishell.Cmd{
		Name: "init",
		Help: "initialize the context",
		Func: cli.InitWebsiteCMD,
	})

	shell.AddCmd(&ishell.Cmd{
		Name: "login",
		Help: "Login with the user admin to performe task on website",
		Func: cli.LoginAdminCMD,
	})

	shell.Run()
}

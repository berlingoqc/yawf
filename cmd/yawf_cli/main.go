package main

import (
	ishell "gopkg.in/abiosoft/ishell.v2"
)

func main() {

	shell := ishell.New()

	shell.AddCmd(&ishell.Cmd{
		Name: "init",
		Help: "initialize the context",
		Func: InitNewWebsite,
	})

	shell.AddCmd(&ishell.Cmd{
		Name: "load",
		Help: "Load a website context information to make operation on it",
		Func: nil,
	})

	shell.Run()
}

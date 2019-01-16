package main

import (
	"net/http"

	ishell "gopkg.in/abiosoft/ishell.v2"
)

var (
	// AuthCookie est mon cookis pour m'authentifier lors des requetes
	AuthCookie *http.Cookie
)

// InitNewWebsite ask question for the basic information to start a website
func InitNewWebsite(c *ishell.Context) {

	// Demande pour le nom du site web
	c.Print("Domain name of the website (ex : wafg.ca) : ")

	c.Print("Enter your full name : ")

	// Demande qu'elle page de base je veux d'enable

	c.Print("Admin account create with success")

}

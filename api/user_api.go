package api

import (
	"errors"
	"log"
	"net/http"
	"net/url"

	"github.com/berlingoqc/yawf/auth"
	"github.com/gorilla/mux"
)

type ReqUser struct {
	Username string `schema:"username"`
	Password string `schema:"password"`
}

// GetUserURLQuery get les infos du user depuis la query url
func GetUserURLQuery(url *url.URL) (*ReqUser, error) {
	u, ok := url.Query()["username"]
	if !ok || len(u) == 0 {
		return nil, errors.New("No username")
	}
	p, ok := url.Query()["password"]
	if !ok || len(p) == 0 {
		log.Println("Erreur pas de password dans la query")
		return nil, errors.New("No password")
	}
	log.Printf("%v %v", u, p)
	return &ReqUser{Username: u[0], Password: p[0]}, nil
}

// UserAPI est l'api pour les informations des users
type UserAPI struct {
	Db *auth.AuthDB
}

// Initialize ajoute les path dans notre route et initialize la db
func (a *UserAPI) Initialize(route *mux.Router) error {
	route.Path("/auth").HandlerFunc(a.HandlerGetAuth)

	route.Use(MiddlewareAPI)

	return nil
}

// HandlerGetAuth operation GET sur /api/auth, pour login le user. Ajout le cookie securiser
func (a *UserAPI) HandlerGetAuth(w http.ResponseWriter, r *http.Request) {
	user, err := GetUserURLQuery(r.URL)
	if err != nil {
		respondWithError(w, http.StatusBadRequest, err.Error())
		return
	}

	// Valide avec la bd que le user est valide
	dbUser, err := a.Db.LoginUser(user.Username, user.Password)
	if err != nil {
		respondWithError(w, http.StatusNotFound, err.Error())
		return
	}

	// Bon user je crée un token pour qu'il s'identifie
	cookie, err := SetCookieForUser(dbUser)
	if err != nil {
		respondWithError(w, http.StatusNotFound, err.Error())
		return
	}
	http.SetCookie(w, cookie)
}

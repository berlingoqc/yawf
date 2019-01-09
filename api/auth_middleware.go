package api

import (
	"net/http"
	"strings"
	"time"

	"github.com/berlingoqc/yawf/auth"
)

// ValidUserCookie valide que le cookie de la request contient
// l'information d'un user valide
func ValidUserCookie(idb *auth.AuthDB, r *http.Request) (*auth.User, error) {
	user, err := DecodeCookieForUser(r)
	if err != nil {
		return nil, err
	}

	if err = idb.IsValidUser(user.ID, user.SaltedPW); err != nil {
		return nil, err
	}

	return user, nil
}

// GetMiddlewareCustom crée un middleware custom qui restraint l'acces a un certain groupe et un role minimal
func GetMiddlewareCustom(minimalRole auth.Role, valid_groups ...auth.Group) func(http.Handler) http.Handler {
	idb, err := auth.GetAuthDBInstance("auth.db")
	if err != nil {
		// log l'erreur
		return nil
	}
	return func(next http.Handler) http.Handler {
		return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
			u, err := ValidUserCookie(idb, r)
			if err != nil {
				// Si le cookie n'est pas valide redirige vers la page d'authentification
				http.Redirect(w, r, "/login", http.StatusMovedPermanently)
				return
			}
			if u.Role < minimalRole {
				// role invalide acces denied mother fucker
				http.Redirect(w, r, "/denied", http.StatusMovedPermanently)
				return
			}

			next.ServeHTTP(w, r)
		})
	}
}

// MiddlewareYase est le midleware pour la route /yase
// Pour accedera YASE sauf / il faut soit etre admin ou etre dans le groupe yase
func MiddlewareYase(next http.Handler) http.Handler {
	return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		next.ServeHTTP(w, r)
	})
}

// MiddlewareAuth est mon middleware pour /auth
func MiddlewareAuth(next http.Handler) http.Handler {
	idb, err := auth.GetAuthDBInstance("auth.db")
	if err != nil {
		// log l'erreur
		return nil
	}
	return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		// regarde si on demande pour logout
		if strings.Contains(r.RequestURI, "/auth/logout") {
			// delete le cookie et redirige vers home
			c := &http.Cookie{Name: authCookieName, Value: "", Path: "/", Expires: time.Unix(0, 0), HttpOnly: true}
			http.SetCookie(w, c)
			http.Redirect(w, r, "/", http.StatusFound)
			return
		}
		// Valide si la personne est deja logger si oui redirgie vers son dashboard
		_, err := ValidUserCookie(idb, r)
		if err == nil {
			http.Redirect(w, r, "/account/dashboard", http.StatusFound)
			return
		}
		next.ServeHTTP(w, r)
	})
}

// MiddlewareAccount est le middleware pour les pages de settings
// du compte et du compte admin
func MiddlewareAccount(next http.Handler) http.Handler {
	idb, err := auth.GetAuthDBInstance("auth.db")
	if err != nil {
		// log l'erreur
		return nil
	}
	return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		u, err := ValidUserCookie(idb, r)
		if err != nil {
			http.Redirect(w, r, "/denied", http.StatusMovedPermanently)
			return
		}
		// si on veux allez au dashboard admin faut etre admin
		if strings.Contains(r.RequestURI, "/account/admin/") {
			if u.Role != auth.RoleAdmin {
				http.Redirect(w, r, "/denied", http.StatusMovedPermanently)
				return
			}
		}

		next.ServeHTTP(w, r)
	})
}

// MiddlewareAPI est le handler pour ma requete GET d'une user
func MiddlewareAPI(next http.Handler) http.Handler {
	idb, err := auth.GetAuthDBInstance("auth.db")
	if err != nil {
		// log l'erreur
		return nil
	}
	return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		// permet l'acces a /api/auth
		if strings.Contains(r.RequestURI, "/api/auth") && r.Method == "GET" {
			next.ServeHTTP(w, r)
			return
		}
		_, err := ValidUserCookie(idb, r)
		if err != nil {
			http.Redirect(w, r, "/auth/login", http.StatusMovedPermanently)
			return
		}
		next.ServeHTTP(w, r)
	})
}

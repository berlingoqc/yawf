package api

import (
	"net/http"

	"github.com/berlingoqc/yawf/auth"
	"github.com/gorilla/schema"
	"github.com/gorilla/securecookie"
)

const (
	authCookieName = "idc"
)

var hashKey = securecookie.GenerateRandomKey(64)
var blockKey = securecookie.GenerateRandomKey(32)

var secureCookie = securecookie.New(hashKey, blockKey)

var decoder = schema.NewDecoder()
var encoder = schema.NewEncoder()

// SetCookieForUser crée une cookie de securiter avec les informations d'un account
func SetCookieForUser(user *auth.User) (*http.Cookie, error) {
	if encoded, err := secureCookie.Encode(authCookieName, user); err == nil {
		cookie := &http.Cookie{
			Name:  authCookieName,
			Value: encoded,
			Path:  "/",
		}
		return cookie, nil
	} else {
		return nil, err
	}
}

// DecodeCookieForUser decode le cookie
func DecodeCookieForUser(r *http.Request) (*auth.User, error) {
	if cookie, e := r.Cookie(authCookieName); e == nil {
		u := &auth.User{}
		return u, secureCookie.Decode(authCookieName, cookie.Value, &u)
	} else {
		return nil, e
	}
}

package api

import (
	"testing"

	"github.com/berlingoqc/yawf/auth"
)

func TestCookie(t *testing.T) {
	user := &auth.User{SaltedPW: "dsdaasdasfa", Role: auth.RoleNormal, Username: "yawf"}
	cookie, err := SetCookieForUser(user)
	if err != nil {
		t.Fatal(err)
	}
	u2 := &auth.User{}
	err = secureCookie.Decode(authCookieName, cookie.Value, u2)
	if err != nil {
		t.Fatal(err)
	}
}

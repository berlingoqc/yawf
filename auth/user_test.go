package auth

import (
	"testing"
)

const password = "pomme"

func TestPassword(t *testing.T) {
	hash, err := getSaltedHash(password)
	if err != nil {
		panic(err)
	}
	t.Logf("Hash de pomme : %s", hash)
	if err = validPassword(password, hash); err != nil {
		panic(err)
	}
	t.Logf("Mot de passe et le hash match")
}

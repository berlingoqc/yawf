package user

import (
	"os"
	"testing"

	"github.com/berlingoqc/yawf/route/security"
)

const dbName = "authdb.sql"

func TestAuthDB(t *testing.T) {
	os.Remove(dbName)
	defer func() {
		err := os.Remove(dbName)
		if err != nil {
			t.Fatal(err)
		}
	}()
	// TEST OUVERTURE
	authdb := &AuthDB{}
	err := authdb.OpenDatabase(dbName)
	if err != nil {
		t.Fatal(err)
	}
	// TEST CREATION ACCOUNT ADMIN
	user, err := authdb.CreateAdminAccount("pomme1234")
	if err != nil {
		t.Fatal(err)
	}
	t.Logf("ID %v Username %v role %v saltedpw %v", user.ID, user.Username, user.Role, user.SaltedPW)

	// TEST CREATION COMPTE ADMIN QUI EXISTE DEJA
	_, err = authdb.CreateAdminAccount("pomme1234")
	if err == nil {
		t.Fatal("Erreur pas d'exection seconde creation account")
	}

	// TEST LOGIN
	user, err = authdb.LoginUser("admin", "pomme1234")
	if err != nil {
		t.Fatal(err)
	}

	// TEST LOGIN AVEC ID
	err = authdb.IsValidUser(user.ID, user.SaltedPW)
	if err != nil {
		t.Fatal(err)
	}

	// TEST CHANGEMENT MOT DE PASSE

	err = authdb.UpdateAccountPassword(user.ID, user.SaltedPW, "poire4321")
	if err != nil {
		t.Fatal(err)
	}

	err = authdb.UpdateAccountPassword(user.ID, user.SaltedPW, "poire4321")
	if err == nil {
		t.Fatal("Devrait throw erreur le mot de passe doit avec ete changer")
	}

	u, err := authdb.CreateAccount("yawf", "office", security.RoleNormal)
	if err != nil {
		t.Fatal(err)
	}
	t.Logf("Create user %v  pw  %v", u.ID, u.SaltedPW)

	listUser, err := authdb.GetListAccount()
	if err != nil {
		t.Fatal(err)
	}
	if len(listUser) != 2 {
		t.Fatalf("Devrait avoir deux elements dans la liste mais %v", len(listUser))
	}

	err = authdb.DeleteAccount(user.ID)
	if err != nil {
		t.Fatal(err)
	}
}

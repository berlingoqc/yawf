package auth

import (
	"golang.org/x/crypto/bcrypt"
)

// Role que l'usager possede sur le site
type Role int

// Group pour permettre l'acces au sous-application
type Group string

const (
	RoleAdmin  Role = 0
	RoleNormal Role = 1
)

// User model pour les usagers de mon site
type User struct {
	ID        int
	Username  string
	SaltedPW  string
	Role      Role
	Group     []Group
	CreateAt  int64
	LastLogin int64
}

// getSaltedHash retourne le mot de passe hasher avec le salt par bcrypt
func getSaltedHash(password string) (string, error) {
	data, err := bcrypt.GenerateFromPassword([]byte(password), bcrypt.DefaultCost)
	if err != nil {
		return "", err
	}
	return string(data), nil
}

// validPassword valide le mot de passe avec le hash entreposer avec bcrypt
func validPassword(password string, stored string) error {
	byte_pw := []byte(password)
	byte_store := []byte(stored)
	return bcrypt.CompareHashAndPassword(byte_store, byte_pw)
}

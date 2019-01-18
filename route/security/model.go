package security

import "golang.org/x/crypto/bcrypt"

// Role is the type of account
type Role int

// Group pour permettre l'acces au sous-application
type Group string

// Right like the unix right
type Right int

const (
	// RightRead is the right the read (GET)
	RightRead Right = 4
	// RightWrite is the right to write (post/put)
	RightWrite Right = 2
	// RightExecute is the right to execute (?delete?)
	RightExecute Right = 1
	// RightNone is a no right
	RightNone Right = 0
)

// Have tell if this right contains this one
func (r Right) Have(right Right) bool {
	return (r & right) > 0
}

const (
	// RoleAdmin is the admin role for all the admin tasks dah
	RoleAdmin Role = 1
	// RoleNormal is for all the normal users
	RoleNormal Role = 2
	// RoleApplication is for application that want to make calls to the paths
	RoleApplication Role = 4
	// RoleSingleLogin is a role for the single login
	RoleSingleLogin Role = 8
)

// IsOneSet tell if one of the bit are set after the mask
func IsOneSet(requiredRole, role Role) bool {
	r := requiredRole & role
	return (r > 0)
}

// SingleLogin is the structure that hold the information for a single login
// access token for certain path
type SingleLogin struct {
	ExpiredAt int64
	Context   map[string]interface{}
}

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

// InGroup tell if user is in a given group
func InGroup(u *User, g Group) bool {
	for _, ug := range u.Group {
		if ug == g {
			return true
		}
	}
	return false
}

// GetSaltedHash retourne le mot de passe hasher avec le salt par bcrypt
func GetSaltedHash(password string) (string, error) {
	data, err := bcrypt.GenerateFromPassword([]byte(password), bcrypt.DefaultCost)
	if err != nil {
		return "", err
	}
	return string(data), nil
}

// ValidPassword valide le mot de passe avec le hash entreposer avec bcrypt
func ValidPassword(password string, stored string) error {
	bytePw := []byte(password)
	byteStore := []byte(stored)
	return bcrypt.CompareHashAndPassword(byteStore, bytePw)
}

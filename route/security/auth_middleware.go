package security

import (
	"fmt"
	"net/http"

	"github.com/berlingoqc/yawf/conv"
)

var (
	// GetUserRepo is a function that is call to obtain a instance of the
	// IUserRepo use by this application
	GetUserRepo = func() IUserRepo {
		return nil
	}
)

// RedirectLogin is the struct that is pass to the login page when trying to login
type RedirectLogin struct {
	URL     string
	Message string
}

// DeniedAccess is the struct pass to the denied page when access to a unauthorized page
type DeniedAccess struct {
	IPAddr   string
	At       string
	URL      string
	Reported string

	BadField string
	Message  string
}

func (d *DeniedAccess) Error() string {
	return fmt.Sprintf("Denied acces at %v by %v on %v", d.At, d.IPAddr, d.URL)
}

// IUserRepo is an interface to access a user repository ( like a database or a text file )
type IUserRepo interface {
	IsValidUser(*User) error
	Close()
}

// PathSecurity is the security information use by the middleware
// There are assign to the WPath and child IPath. The WPath security
// if not null is apply to the IPath that have a null security
// is wpath security is nil the security only apply to the path
// that are not nil
type PathSecurity struct {
	RoleRequired Role
	Owner        string
	Group        Group
	Right        [3]Right
}

// ValidUserCookie valide que le cookie de la request contient
// l'information d'un user valide
func ValidUserCookie(r *http.Request) (*User, error) {
	user, err := DecodeCookieForUser(r)
	if err != nil {
		return nil, err
	}

	irepo := GetUserRepo()
	defer irepo.Close()
	if err = irepo.IsValidUser(user); err != nil {
		return nil, err
	}
	return user, nil
}

// RedirectTo redirect to an url with the structure given as query
func RedirectTo(w http.ResponseWriter, r *http.Request, dest string, data interface{}) {
	if data != nil {
		queryParam, err := conv.StructToQuery(data)
		if err != nil {
			print("Error redirect " + err.Error())
			http.Redirect(w, r, "/", http.StatusSeeOther)
		}
		dest += queryParam
	}
	http.Redirect(w, r, dest, http.StatusSeeOther)
}

// PathSecurityValidation valid the path access
func PathSecurityValidation(ps *PathSecurity, w http.ResponseWriter, r *http.Request) bool {
	// Regarde si authentifier
	user, err := ValidUserCookie(r)
	if err != nil {
		// redirige vers login avec message d'erreur et l'url qu'on voulait allez
		rl := &RedirectLogin{
			Message: err.Error(),
			URL:     r.URL.String(),
		}
		RedirectTo(w, r, "/auth/login", rl)
		return false
	}

	var errAccess *DeniedAccess
	// valide que le token est encore bon

	// Valide que le role est correcte
	if !IsOneSet(ps.RoleRequired, user.Role) {
		errAccess = &DeniedAccess{
			BadField: "role",
			Message:  "Role incorrect",
		}
	}
	// Regarde s'il s'agit d'un usager ou group ou autre
	var right Right
	var needed Right
	if user.Username == ps.Owner {
		right = ps.Right[0]
	} else if InGroup(user, ps.Group) {
		right = ps.Right[1]
	} else {
		right = ps.Right[2]
	}
	// Regarde le type de request et regarde si on peut la faire avec les droits
	switch r.Method {
	case "GET":
		needed = RightRead
	case "POST", "PUT", "DELETE":
		needed = RightWrite
	default:
		needed = RightExecute
	}
	if !right.Have(needed) {
		errAccess = &DeniedAccess{
			BadField: "right",
			Message:  fmt.Sprintf("Your right are %v but needed %v", right, needed),
		}
	}

	if errAccess != nil {
		RedirectTo(w, r, "/denied", errAccess)
		return false
	}

	return true
}

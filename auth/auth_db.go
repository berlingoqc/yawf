package auth

import (
	"database/sql"
	"fmt"

	"github.com/berlingoqc/yawf/db"
)

type AccountErrType string

const (
	DontExists        AccountErrType = "account don't exists"
	AlreadyExists     AccountErrType = "username already exists"
	DontMatchPassword AccountErrType = "password don't match"
	NotSafePassword   AccountErrType = "password is not safe enough "

	UsernameAdmin string = "admin"
)

// Definition des erreurs

type AccountError struct {
	Username string
	Type     AccountErrType
}

func (a *AccountError) Error() string {
	return fmt.Sprintf("Error %v : %v", a.Username, a.Type)
}

type AuthDB struct {
	DB *db.BaseDB

	FilePath string

	tables []string
}

func (p *AuthDB) Initialize(filePath string) {
	p.FilePath = filePath
	p.DB = &db.BaseDB{}

	p.tables = []string{authSqltable}
}

func (p *AuthDB) GetFilePath() string {
	return p.FilePath
}

func (p *AuthDB) GetTables() []string {
	return p.tables
}

func (p *AuthDB) GetDB() *sql.DB {
	return p.DB.DB
}

func (p *AuthDB) SetDB(s *sql.DB) {
	p.DB.DB = s
}

func (a *AuthDB) OpenDatabase(filePath string) error {
	a.DB = &db.BaseDB{}
	a.DB.AddTables(authSqltable)
	return a.DB.OpenDatabase(filePath)
}

func (a *AuthDB) CloseDatabase() {
	a.DB.CloseDatabase()
}

func (a *AuthDB) DoesAccountExists(username string) (bool, error) {
	var i int
	err := a.DB.DB.QueryRow(QUERY_isusernameexists, username).Scan(&i)
	if err != nil {
		return false, err
	}

	return (i == 1), nil
}

func (a *AuthDB) CreateAdminAccount(password string) (*User, error) {
	return a.CreateAccount(UsernameAdmin, password, RoleAdmin)
}

func (a *AuthDB) CreateAccount(username string, password string, role Role) (*User, error) {
	if b, e := a.DoesAccountExists(username); e != nil || b {
		if e == nil {
			e = &AccountError{
				Username: username,
				Type:     AlreadyExists,
			}
		}
		return nil, e
	}

	hashPw, e := getSaltedHash(password)
	if e != nil {
		return nil, e
	}

	e = a.DB.StmtQuery(QUERY_createaccount, username, role, hashPw)
	if e != nil {
		return nil, e
	}

	return a.LoginUser(username, password)
}

func (a *AuthDB) LoginUser(username string, password string) (*User, error) {
	user := &User{Username: username}
	err := a.DB.DB.QueryRow(QUERY_accountlogin, username).Scan(&user.ID, &user.SaltedPW, &user.Role)
	if err != nil {
		return nil, err
	}
	return user, validPassword(password, user.SaltedPW)
}

func (a *AuthDB) IsValidUser(id int, hash string) error {
	var i int
	if err := a.DB.DB.QueryRow(QUERY_validaccount, id, hash).Scan(&i); err != nil {
		return err
	}
	if i == 0 {
		return &AccountError{Username: string(id), Type: DontExists}
	}
	return nil
}

func (a *AuthDB) UpdateAccountPassword(id int, oldhash string, newpw string) error {
	if e := a.IsValidUser(id, oldhash); e != nil {
		return e
	}
	hash, e := getSaltedHash(newpw)
	if e != nil {
		return e
	}
	return a.DB.StmtQuery(QUERY_update_pw, hash, id)
}

func (a *AuthDB) DeleteAccount(id int) error {
	return a.DB.StmtQuery(QUERY_delete_account, id)
}

func (a *AuthDB) GetListAccount() ([]*User, error) {
	rows, err := a.DB.DB.Query(QUERY_listaccout)
	if err != nil {
		return nil, err
	}
	defer rows.Close()
	var users []*User
	for rows.Next() {
		u := &User{}
		err = rows.Scan(&u.ID, &u.Username, &u.Role)
		if err != nil {
			return users, err
		}
		users = append(users, u)
	}
	return users, nil
}

func GetAuthDBInstance(filePath string) (*AuthDB, error) {
	idb := &AuthDB{}
	return idb, idb.OpenDatabase(filePath)
}

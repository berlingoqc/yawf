package db

import (
	"database/sql"
	"log"

	// mes drivers pour databse/sql
	_ "github.com/mattn/go-sqlite3"
)

type IDB interface {
	GetFilePath() string
	GetTables() []string
	GetDB() *sql.DB
	SetDB(*sql.DB)
}

func OpenDatabase(db IDB) error {
	var err error
	sdb, err := sql.Open("sqlite3", db.GetFilePath())
	if err != nil {
		return err
	}
	db.SetDB(sdb)
	sdb.SetMaxOpenConns(1)
	return CreateTables(db)
}

func CloseDatabse(db IDB) error {
	return db.GetDB().Close()
}

// CreateTables crée les tables de la db
func CreateTables(db IDB) error {
	for _, t := range db.GetTables() {
		_, err := db.GetDB().Exec(t)
		if err != nil {
			return err
		}
	}
	return nil
}

// ChainStmtQuery ...
func ChainStmtQuery(db IDB, m map[string][]interface{}) error {
	var err error
	for k, v := range m {
		err = StmtQuery(db, k, v)
		if err != nil {
			return err
		}
	}
	return nil
}

// StmtQuery ...
func StmtQuery(db IDB, query string, args ...interface{}) error {
	stmt, err := db.GetDB().Prepare(query)
	if err != nil {
		return err
	}
	defer stmt.Close()
	_, err = stmt.Exec(args...)
	return err
}

func QueryStringArray(db IDB, query string, args ...interface{}) ([]string, error) {
	rows, err := db.GetDB().Query(query, args...)
	if err != nil {
		return nil, err
	}
	defer rows.Close()
	var s []string
	if err != nil {
		return nil, err
	}
	for rows.Next() {
		var l string
		err = rows.Scan(&l)
		if err != nil {
			return nil, err
		}
		s = append(s, l)
	}
	return s, nil
}

// BaseDB classe de base pour travailler sur la bd
type BaseDB struct {
	DB       *sql.DB
	FilePath string
	Tables   []string
}

// OpenDatabase ouvre la base de donnée et assure que les tables soit crées
func (b *BaseDB) OpenDatabase(filepath string) error {
	var err error
	b.DB, err = sql.Open("sqlite3", filepath)
	if err != nil {
		return err
	}
	b.DB.SetMaxOpenConns(1)
	return b.CreateTables()
}

// CloseDatabase doit etre appeler a la fin de l'utilisation pour la fermer
func (b *BaseDB) CloseDatabase() {
	err := b.DB.Close()
	if err != nil {
		log.Printf(err.Error())
	}
}

// AddTables ajout des tables dans la bd qui seront crée a l'ouverture de la db
func (b *BaseDB) AddTables(tables ...string) {
	if b.Tables == nil {
		b.Tables = make([]string, 0)
	}
	b.Tables = append(b.Tables, tables...)
}

// CreateTables crée les tables de la db
func (b *BaseDB) CreateTables() error {
	for _, t := range b.Tables {
		_, err := b.DB.Exec(t)
		if err != nil {
			return err
		}
	}
	return nil
}

// ChainStmtQuery ...
func (b *BaseDB) ChainStmtQuery(m map[string][]interface{}) error {
	var err error
	for k, v := range m {
		err = b.StmtQuery(k, v)
		if err != nil {
			return err
		}
	}
	return nil
}

// StmtQuery ...
func (b *BaseDB) StmtQuery(query string, args ...interface{}) error {
	stmt, err := b.DB.Prepare(query)
	if err != nil {
		return err
	}
	defer stmt.Close()
	_, err = stmt.Exec(args...)
	return err
}

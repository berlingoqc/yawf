package base

import (
	"database/sql"

	"github.com/berlingoqc/yawf/db"
)

const (
	SQLTable = `
	CREATE TABLE IF NOT EXISTS subject (
		Name VARCHAR(50) PRIMARY KEY
	);

	CREATE TABLE IF NOT EXISTS language (
		Name VARCHAR(50) PRIMARY KEY
	);


	`

	QueryAddSubject  = "INSERT INTO subject (Name) VALUES (?)"
	QueryAddLanguage = "INSERT INTO language (Name) VALUES (?)"

	QuerySelectLanguages = "SELECT Name FROM language"
	QuerySelectSubjects  = "SELECT Name FROM subject"
)

type DB struct {
	FilePath string
	tables   []string

	Db *sql.DB
}

func (p *DB) Initialize(filePath string) {
	p.FilePath = filePath

	p.tables = []string{SQLTable}
}

func (p *DB) GetFilePath() string {
	return p.FilePath
}

func (p *DB) GetTables() []string {
	return p.tables
}

func (p *DB) GetDB() *sql.DB {
	return p.Db
}

func (p *DB) SetDB(s *sql.DB) {
	p.Db = s
}

func (p *DB) AddLanguage(name string) error {
	return db.StmtQuery(p, QueryAddLanguage, name)
}

func (p *DB) AddSubject(name string) error {
	return db.StmtQuery(p, QueryAddSubject, name)
}

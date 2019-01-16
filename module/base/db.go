package base

import (
	"database/sql"

	"github.com/berlingoqc/yawf/db"
)

const (
	// SQLTable contains all the table of the main module
	SQLTable = `
	CREATE TABLE IF NOT EXISTS subject (
		Name VARCHAR(50) PRIMARY KEY
	);

	CREATE TABLE IF NOT EXISTS language (
		Name VARCHAR(50) PRIMARY KEY
	);

	CREATE TABLE IF NOT EXISTS module (
		Name VARCHAR(50) NOT NULL,
		Version VARCHAR(10) NOT NULL,

		Path VARCHAR(255) NOT NULL,

		PRIMARY KEY (Name)
	);

	`

	// QueryAddSubject add a new subject
	QueryAddSubject = "INSERT INTO subject (Name) VALUES (?)"
	// QueryAddLanguage add a new language
	QueryAddLanguage = "INSERT INTO language (Name) VALUES (?)"

	// QuerySelectLanguages select all the languages
	QuerySelectLanguages = "SELECT Name FROM language"
	// QuerySelectSubjects select all the subjects
	QuerySelectSubjects = "SELECT Name FROM subject"

	// QueryAddModule add a new release of a module
	QueryAddModule = "INSERT INTO module (Name, Version, Path) VALUES (?,?,?)"

	// QuerySelectModules select all the modules
	QuerySelectModules = "SELECT Name, Version, Path FROM module"
	// QuerySelectModule select a module with the name
	QuerySelectModule = QuerySelectModules + " WHERE Name = ?"
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

package base

import (
	"database/sql"
	"strings"

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

	CREATE TABLE IF NOT EXISTS formation (
		ID INTEGER PRIMARY KEY AUTOINCREMENT,
		name VARCHAR(100),
		diploma VARCHAR(100),
		school VARCHAR(100),
		startdate VARCHAR(100),
		enddate VARCHAR(100),
		length INTEGER,
		mention VARCHAR(500),
		description VARCHAR(1000)
	);

	CREATE TABLE IF NOT EXISTS experience (
		ID INTEGER PRIMARY KEY AUTOINCREMENT,
		job VARCHAR(255),
		corporation VARCHAR(255),
		location VARCHAR(255),
		startdate VARCHAR(100),
		enddate VARCHAR(100),
		length INTEGER,
		description VARCHAR(1000)
	);

	CREATE TABLE IF NOT EXISTS language_experience (
		language_id VARCHAR(50) PRIMARY KEY,
		level VARCHAR(50),
		description VARCHAR(1000),
		length INTEGER
	);
	`

	QueryAddSubject            = "INSERT INTO subject (Name) VALUES (?)"
	QueryAddLanguage           = "INSERT INTO language (Name) VALUES (?)"
	QueryAddFormation          = "INSERT INTO formation (name,diploma,school,startdate,enddate,length,mention,description) VALUES (?,?,?,?,?,?,?,?)"
	QueryAddExperience         = "INSERT INTO experience (job,corporation,location,startdate,enddate,length,description) VALUES (?,?,?,?,?,?,?)"
	QueryAddLanguageExperience = "INSERT INTO language_experience (language_id,level,description,length) VALUES (?,?,?,?)"

	QueryFormation          = "SELECT name,diploma,school,startdate,enddate,length,mention,description FROM formation"
	QueryExperience         = "SELECT job,corporation,location,startdate,enddate,length,description FROM experience"
	QueryLanguageExperience = "SELECT language_id, level, description, length FROM language_experience"

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

func (p *DB) AddFormation(r *Formation) error {
	m := strings.Join(r.Mention, ";")
	return db.StmtQuery(p, QueryAddFormation, r.Name, r.NameDiploma, r.School, r.StartDate, r.EndDate, r.LengthYear, m, r.Description)
}

func (p *DB) AddExperience(r *ProfessionalExperience) error {
	return db.StmtQuery(p, QueryAddExperience, r.Job, r.Corporation, r.Location, r.StartDate, r.EndDate, r.LengthYear, r.Description)
}

func (p *DB) AddLanguageExperience(r *LanguageExperience) error {
	return db.StmtQuery(p, QueryAddLanguageExperience, r.Name, r.Level, r.Description, r.Year)
}

func (p *DB) GetFormation() ([]*Formation, error) {
	rows, err := p.Db.Query(QueryFormation)
	if err != nil {
		return nil, err
	}
	defer rows.Close()
	var ff []*Formation
	for rows.Next() {
		var m string
		f := &Formation{}
		err = rows.Scan(&f.Name, &f.NameDiploma, &f.School, &f.StartDate, &f.EndDate, &f.LengthYear, &m, &f.Description)
		if err != nil {
			return nil, err
		}
		f.Mention = strings.Split(m, ";")
		ff = append(ff, f)
	}
	return ff, nil
}

func (p *DB) GetExperience() ([]*ProfessionalExperience, error) {
	rows, err := p.Db.Query(QueryExperience)
	if err != nil {
		return nil, err
	}
	var ee []*ProfessionalExperience
	for rows.Next() {
		e := &ProfessionalExperience{}
		err = rows.Scan(&e.Job, &e.Corporation, &e.Location, &e.StartDate, &e.EndDate, &e.LengthYear, &e.Description)
		if err != nil {
			return nil, err
		}
		ee = append(ee, e)
	}
	return ee, nil
}

func (p *DB) GetLanguageExperience() ([]*LanguageExperience, error) {
	rows, err := p.Db.Query(QueryLanguageExperience)
	if err != nil {
		return nil, err
	}
	var ll []*LanguageExperience
	for rows.Next() {
		l := &LanguageExperience{}
		err = rows.Scan(&l.Name, &l.Level, &l.Description, &l.Year)
		if err != nil {
			return nil, err
		}
		ll = append(ll, l)

	}
	return ll, nil
}

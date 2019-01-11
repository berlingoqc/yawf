package personal

import (
	"database/sql"
)

type Formation struct {
	Name        string
	NameDiploma string
	School      string
	StartDate   string
	EndDate     string
	LengthYear  int
	Mention     []string
	Description string
	Active      bool
}

type ProfessionalExperience struct {
	Job         string
	Corporation string
	Location    string
	StartDate   string
	EndDate     string
	LengthYear  int
	Active      bool

	Description string
}

type LanguageExperience struct {
	Name        string
	Level       string
	Description string
	Year        int
}

func FromRowProfessionalExperience(r *sql.Rows) (*ProfessionalExperience, error) {
	p := &ProfessionalExperience{}
	return nil, r.Scan(&p.Job, &p.Corporation, &p.Location, &p.StartDate, &p.EndDate, &p.LengthYear, &p.Active, &p.Description)
}

type CV struct {
	Name               string
	Emplacement        string
	ImgURL             string
	LettrePresentation []byte
}

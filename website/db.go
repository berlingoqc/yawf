package website

import (
	"database/sql"
	"strings"
	"time"

	"github.com/berlingoqc/yawf/db"
	"github.com/berlingoqc/yawf/website/cv"
	"github.com/berlingoqc/yawf/website/project"
)

type ProjectDB struct {
	FilePath string
	tables   []string
	Db       *sql.DB
}

func (p *ProjectDB) Initialize(filePath string) {
	p.FilePath = filePath

	p.tables = []string{projectSqlTable}
}

func (p *ProjectDB) GetFilePath() string {
	return p.FilePath
}

func (p *ProjectDB) GetTables() []string {
	return p.tables
}

func (p *ProjectDB) GetDB() *sql.DB {
	return p.Db
}

func (p *ProjectDB) SetDB(s *sql.DB) {
	p.Db = s
}

func (p *ProjectDB) AddGHAccount(a *project.GitHubAccount) error {
	return db.StmtQuery(p, QueryAddProjectInfo, a.Name, a.URL, a.ImgURL, a.Location, a.Email, a.Bio, a.NbrPublicRepo, a.NbrPublicGists, a.NbrFollorwers, a.NbrFollowing)
}

func (p *ProjectDB) AddGHOrganization(a *project.GitHubOrganistaion) error {
	return nil
}

func (p *ProjectDB) AddGHRepo(r *project.GitHubRepo) error {
	update := r.LastUpdateOn.Format(timeFormat)
	created := r.CreatedOn.Format(timeFormat)
	return db.StmtQuery(p, QueryAddGitHubRepo, r.Name, r.URl, r.Description, r.StarCount, r.ForksCount, r.CommitNumber, update, created, r.ReadMe)
}

func (p *ProjectDB) AddPProject(pr *project.ProgrammingProject, repoName string) error {
	e := db.StmtQuery(p, QueryAddProject, pr.Name, pr.ThumbnailURL, pr.SiteURL, pr.DocURL, repoName)
	if e != nil {
		return e
	}

	for _, s := range pr.Subjects {
		e = db.StmtQuery(p, QueryAddSubjectProject, pr.Name, s)
		if e != nil {
			return e
		}
	}

	for _, l := range pr.Language {
		e = db.StmtQuery(p, QueryAddLanguageProject, pr.Name, l)
		if e != nil {
			return e
		}
	}

	return nil
}

func (p *ProjectDB) AddLanguage(name string) error {
	return db.StmtQuery(p, QueryAddLanguage, name)
}

func (p *ProjectDB) AddSubject(name string) error {
	return db.StmtQuery(p, QueryAddSubject, name)
}

func (p *ProjectDB) AddFormation(r *cv.Formation) error {
	m := strings.Join(r.Mention, ";")
	return db.StmtQuery(p, QueryAddFormation, r.Name, r.NameDiploma, r.School, r.StartDate, r.EndDate, r.LengthYear, m, r.Description)
}

func (p *ProjectDB) AddExperience(r *cv.ProfessionalExperience) error {
	return db.StmtQuery(p, QueryAddExperience, r.Job, r.Corporation, r.Location, r.StartDate, r.EndDate, r.LengthYear, r.Description)
}

func (p *ProjectDB) AddLanguageExperience(r *cv.LanguageExperience) error {
	return db.StmtQuery(p, QueryAddLanguageExperience, r.Name, r.Level, r.Description, r.Year)
}

func (p *ProjectDB) UpdateGitHubRepo(r *project.GitHubRepo) error {
	return db.StmtQuery(p, QueryUpdateGitHubRepo, r.Description, r.StarCount, r.ForksCount, r.CommitNumber, r.LastUpdateOn, r.ReadMe, r.Name)
}

func (p *ProjectDB) UpdateGHAccount(a *project.GitHubAccount) error {
	return db.StmtQuery(p, QueryUpdateGitHubAccount, a.ImgURL, a.Location, a.Bio, a.Email, a.NbrPublicRepo, a.NbrPublicGists, a.NbrFollorwers, a.NbrFollowing, a.Name)
}

func (p *ProjectDB) GetGHAccount() (*project.GitHubAccount, error) {
	i := &project.GitHubAccount{}
	e := p.Db.QueryRow(QueryProjectInfo).Scan(&i.Name, &i.URL, &i.ImgURL, &i.Location, &i.Email, &i.Bio, &i.NbrPublicRepo, &i.NbrPublicGists, &i.NbrFollorwers, &i.NbrFollowing)
	return i, e
}

func (p *ProjectDB) GetFormation() ([]*cv.Formation, error) {
	rows, err := p.Db.Query(QueryFormation)
	if err != nil {
		return nil, err
	}
	defer rows.Close()
	var ff []*cv.Formation
	for rows.Next() {
		var m string
		f := &cv.Formation{}
		err = rows.Scan(&f.Name, &f.NameDiploma, &f.School, &f.StartDate, &f.EndDate, &f.LengthYear, &m, &f.Description)
		if err != nil {
			return nil, err
		}
		f.Mention = strings.Split(m, ";")
		ff = append(ff, f)
	}
	return ff, nil
}

func (p *ProjectDB) GetExperience() ([]*cv.ProfessionalExperience, error) {
	rows, err := p.Db.Query(QueryExperience)
	if err != nil {
		return nil, err
	}
	var ee []*cv.ProfessionalExperience
	for rows.Next() {
		e := &cv.ProfessionalExperience{}
		err = rows.Scan(&e.Job, &e.Corporation, &e.Location, &e.StartDate, &e.EndDate, &e.LengthYear, &e.Description)
		if err != nil {
			return nil, err
		}
		ee = append(ee, e)
	}
	return ee, nil
}

func (p *ProjectDB) GetLanguageExperience() ([]*cv.LanguageExperience, error) {
	rows, err := p.Db.Query(QueryLanguageExperience)
	if err != nil {
		return nil, err
	}
	var ll []*cv.LanguageExperience
	for rows.Next() {
		l := &cv.LanguageExperience{}
		err = rows.Scan(&l.Name, &l.Level, &l.Description, &l.Year)
		if err != nil {
			return nil, err
		}
		ll = append(ll, l)

	}
	return ll, nil
}

func (p *ProjectDB) GetLanguage() ([]string, error) {
	return db.QueryStringArray(p, QuerySelectLanguages)
}

func (p *ProjectDB) GetSubject() ([]string, error) {
	return db.QueryStringArray(p, QuerySelectSubjects)
}

func (p *ProjectDB) GetPProjects() ([]*project.ProgrammingProject, error) {
	rows, err := p.Db.Query(QuerySelectProjects)
	if err != nil {
		return nil, err
	}
	defer rows.Close()
	if err != nil {
		return nil, err
	}
	var ppl []*project.ProgrammingProject
	for rows.Next() {

		pp := &project.ProgrammingProject{
			GitHub: &project.GitHubRepo{},
		}
		var update, created string
		err = rows.Scan(&pp.Name, &pp.ThumbnailURL, &pp.SiteURL, &pp.DocURL, &pp.GitHub.Name, &pp.GitHub.URl, &pp.GitHub.Description, &pp.GitHub.StarCount, &pp.GitHub.ForksCount, &pp.GitHub.CommitNumber, &update, &created, &pp.GitHub.ReadMe)
		if err != nil {
			return nil, err
		}
		t1, _ := time.Parse(timeFormat, update)
		pp.GitHub.LastUpdateOn = t1
		t2, _ := time.Parse(timeFormat, created)
		pp.GitHub.CreatedOn = t2

		ppl = append(ppl, pp)
	}

	for _, i := range ppl {
		i.Language, err = db.QueryStringArray(p, QueryLanguageProject, i.Name)
		if err != nil {
			return nil, err
		}
		i.Subjects, err = db.QueryStringArray(p, QuerySubjectProject, i.Name)
		if err != nil {
			return nil, err
		}
	}
	return ppl, nil
}

func (p *ProjectDB) GetProjectDetail(key string) (*project.ProgrammingProject, error) {

	return nil, nil
}

func (p *ProjectDB) LinkBlogPost(projectName, blogName string) error {

	return nil
}

func GetProjectDBInstance(filePath string) (*ProjectDB, error) {
	idb := &ProjectDB{}
	idb.Initialize(filePath)

	return idb, db.OpenDatabase(idb)
}

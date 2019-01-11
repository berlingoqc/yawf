package project

import (
	"database/sql"
	"time"

	"github.com/berlingoqc/yawf/db"
)

const (
	timeFormat      = "2006/01/02 03:04"
	projectSqlTable = `
	CREATE TABLE IF NOT EXISTS user_project_info (
		ID 			VARCHAR(50) PRIMARY KEY,
		user_url	VARCHAR(255),
		user_img_url VARCHAR(255),

		location VARCHAR(50),
		email VARCHAR(50),
		bio VARCHAR(1000),

		nbr_public_repo INTEGER,
		nbr_public_gists INTEGER,

		nbr_followers INTEGER,
		nbr_following INTEGER
	);


	CREATE TABLE IF NOT EXISTS organization (
		ID VARCHAR(10) PRIMARY KEY,
		INFO_ID VARCHAR(10) NOT NULL,
		url VARCHAR(255),
		img_url VARCHAR(255),
		name VARCHAR(255),

		bio VARCHAR(1000),

		myrole VARCHAR(255),

		FOREIGN KEY(INFO_ID) REFERENCES user_project_info(ID)
	);

	CREATE TABLE IF NOT EXISTS git_hub_repo (
		ID VARCHAR(50) PRIMARY KEY,
		url VARCHAR(255),
		description VARCHAR(255),

		star INTEGER,
		forks INTEGER,
		commit_nbr INTEGER,

		update_on VARCHAR(255),
		create_on VARCHAR(255),

		readme BLOB
	);

	CREATE TABLE IF NOT EXISTS subject_pproject(
		pp_id VARCHAR(50) NOT NULL,
		subject_id VARCHAR(50) NOT NULL,

		FOREIGN KEY(pp_id) REFERENCES pproject(ID)
		FOREIGN KEY(subject_id) REFERENCES subject(ID)
	);

	CREATE TABLE IF NOT EXISTS language_pproject(
		pp_id VARCHAR(50) NOT NULL,
		language_id VARCHAR(50) NOT NULL,

		FOREIGN KEY(pp_id) REFERENCES pproject(ID)
		FOREIGN KEY(language_id) REFERENCES subject(ID)

	);

	CREATE TABLE IF NOT EXISTS blog_pproject(
		pp_id VARCHAR(50) NOT NULL,
		blog_id VARCHAR(50) NOT NULL,

		FOREIGN KEY(pp_id) REFERENCES pproject(ID)
	);

	CREATE TABLE IF NOT EXISTS pproject (
		ID VARCHAR(50) PRIMARY KEY,
		img_url VARCHAR(255),
		site_url VARCHAR(255),
		doc_url VARCHAR(255),
		ID_GIT VARCHAR(50),
		FOREIGN KEY(ID_GIT) REFERENCES git_hub_repo(ID)
	);

	`

	QueryAddProjectInfo     = "INSERT INTO user_project_info (ID,user_url,user_img_url,location,email,bio,nbr_public_repo,nbr_public_gists,nbr_followers,nbr_following) VALUES (?,?,?,?,?,?,?,?,?,?)"
	QueryAddSubject         = "INSERT INTO subject (Name) VALUES (?)"
	QueryAddSubjectProject  = "INSERT INTO subject_pproject (pp_id,subject_id) VALUES (?,?)"
	QueryAddLanguage        = "INSERT INTO language (Name) VALUES (?)"
	QueryAddLanguageProject = "INSERT INTO language_pproject(pp_id,language_id) VALUES (?,?)"
	QuertAddBlogProject     = "INSERT INTO blog_pproject (pp_id,blog_id) VALUES (?,?)"
	QueryAddGitHubRepo      = "INSERT INTO git_hub_repo (ID,url,description,star,forks,commit_nbr,update_on,create_on,readme) VALUES (?,?,?,?,?,?,?,?,?)"
	QueryAddProject         = "INSERT INTO pproject (ID,img_url,site_url,doc_url,ID_GIT) VALUES (?,?,?,?,?)"

	QueryProjectInfo = "SELECT ID,user_url,user_img_url,location,email,bio,nbr_public_repo,nbr_public_gists,nbr_followers,nbr_following FROM user_project_info"

	QueryUpdateGitHubRepo    = "UPDATE git_hub_repo SET description = ?, star = ?, forks = ?, commit_nbr = ?, update_on = ?, readme = ? WHERE ID = ?"
	QueryUpdateGitHubAccount = "UPDATE user_project_info SET user_img_url = ?, location = ?, bio = ?, email = ?, nbr_public_repo = ?, nbr_public_gists = ?, nbr_followers = ?, nbr_following = ? WHERE ID = ?"

	QuerySelectLanguages = "SELECT Name FROM language"
	QuerySelectSubjects  = "SELECT Name FROM subject"

	QuerySelectProjects  = "SELECT pp.ID, pp.img_url, pp.site_url, pp.doc_url, gh.ID, gh.url, gh.description, gh.star, gh.forks, gh.commit_nbr, gh.update_on, gh.create_on, readme FROM pproject AS pp JOIN git_hub_repo AS gh ON pp.ID_GIT = gh.ID"
	QueryLanguageProject = "SELECT language.Name FROM language_pproject JOIN language ON language.Name = language_pproject.language_id WHERE language_pproject.pp_id = ?"
	QuerySubjectProject  = "SELECT subject.Name FROM subject_pproject JOIN subject ON subject.Name = subject_pproject.subject_id WHERE subject_pproject.pp_id = ?"
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

func (p *ProjectDB) AddGHAccount(a *GitHubAccount) error {
	return db.StmtQuery(p, QueryAddProjectInfo, a.Name, a.URL, a.ImgURL, a.Location, a.Email, a.Bio, a.NbrPublicRepo, a.NbrPublicGists, a.NbrFollorwers, a.NbrFollowing)
}

func (p *ProjectDB) AddGHOrganization(a *GitHubOrganistaion) error {
	return nil
}

func (p *ProjectDB) AddGHRepo(r *GitHubRepo) error {
	update := r.LastUpdateOn.Format(timeFormat)
	created := r.CreatedOn.Format(timeFormat)
	return db.StmtQuery(p, QueryAddGitHubRepo, r.Name, r.URl, r.Description, r.StarCount, r.ForksCount, r.CommitNumber, update, created, r.ReadMe)
}

func (p *ProjectDB) AddPProject(pr *ProgrammingProject, repoName string) error {
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

func (p *ProjectDB) UpdateGitHubRepo(r *GitHubRepo) error {
	return db.StmtQuery(p, QueryUpdateGitHubRepo, r.Description, r.StarCount, r.ForksCount, r.CommitNumber, r.LastUpdateOn, r.ReadMe, r.Name)
}

func (p *ProjectDB) UpdateGHAccount(a *GitHubAccount) error {
	return db.StmtQuery(p, QueryUpdateGitHubAccount, a.ImgURL, a.Location, a.Bio, a.Email, a.NbrPublicRepo, a.NbrPublicGists, a.NbrFollorwers, a.NbrFollowing, a.Name)
}

func (p *ProjectDB) GetGHAccount() (*GitHubAccount, error) {
	i := &GitHubAccount{}
	e := p.Db.QueryRow(QueryProjectInfo).Scan(&i.Name, &i.URL, &i.ImgURL, &i.Location, &i.Email, &i.Bio, &i.NbrPublicRepo, &i.NbrPublicGists, &i.NbrFollorwers, &i.NbrFollowing)
	return i, e
}

func (p *ProjectDB) GetLanguage() ([]string, error) {
	return db.QueryStringArray(p, QuerySelectLanguages)
}

func (p *ProjectDB) GetSubject() ([]string, error) {
	return db.QueryStringArray(p, QuerySelectSubjects)
}

func (p *ProjectDB) GetPProjects() ([]*ProgrammingProject, error) {
	rows, err := p.Db.Query(QuerySelectProjects)
	if err != nil {
		return nil, err
	}
	defer rows.Close()
	if err != nil {
		return nil, err
	}
	var ppl []*ProgrammingProject
	for rows.Next() {

		pp := &ProgrammingProject{
			GitHub: &GitHubRepo{},
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

func (p *ProjectDB) GetProjectDetail(key string) (*ProgrammingProject, error) {

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

package website

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

	CREATE TABLE IF NOT EXISTS subject (
		Name VARCHAR(50) PRIMARY KEY
	);

	CREATE TABLE IF NOT EXISTS language (
		Name VARCHAR(50) PRIMARY KEY
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

	QueryAddFormation          = "INSERT INTO formation (name,diploma,school,startdate,enddate,length,mention,description) VALUES (?,?,?,?,?,?,?,?)"
	QueryAddExperience         = "INSERT INTO experience (job,corporation,location,startdate,enddate,length,description) VALUES (?,?,?,?,?,?,?)"
	QueryAddLanguageExperience = "INSERT INTO language_experience (language_id,level,description,length) VALUES (?,?,?,?)"

	QueryFormation          = "SELECT name,diploma,school,startdate,enddate,length,mention,description FROM formation"
	QueryExperience         = "SELECT job,corporation,location,startdate,enddate,length,description FROM experience"
	QueryLanguageExperience = "SELECT language_id, level, description, length FROM language_experience"

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

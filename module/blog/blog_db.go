package blog

import (
	"database/sql"

	"github.com/berlingoqc/yawf/db"
)

const sqlTableBD = `
	CREATE TABLE IF NOT EXISTS blog_language(
		blog_id INTEGER NOT NULL,
		language_id VARCHAR(50) NOT NULL,

		FOREIGN KEY(blog_id) REFERENCES blog_post(Title)
		FOREIGN KEY(language_id) REFERENCES language(Name)
	);

	CREATE TABLE IF NOT EXISTS blog_subject(
		blog_id INTEGER NOT NULL,
		subject_id VARCHAR(50) NOT NULL,

		FOREIGN KEY(blog_id) REFERENCES blog_post(Title)
		FOREIGN KEY(subject_id) REFERENCES subject(Name) 
	);

	CREATE TABLE IF NOT EXISTS blog_serie(
		ID INTEGER PRIMARY KEY AUTOINCREMENT,
		Title VARCHAR(255) NOT NULL,
		Description VARCHAR(1000),
		ThumbnailURL VARCHAR(255),
		Over boolean,
		CurrentPublication int,
		TotalPublication int
	);

	CREATE TABLE IF NOT EXISTS blog_post(
		ID INTEGER PRIMARY KEY AUTOINCREMENT,
		Title VARCHAR(255) NOT NULL,
		PostOn VARCHAR(100) NOT NULL,
		UpdateOn VARCHAR(100) NOT NULL,
		ThumbnailURL VARCHAR(100) NOT NULL,
		Description VARCHAR(1000) NOT NULL,

		Author VARCHAR(100) NOT NULL,
		VideoURL VARCHAR(255),

		SerieID INTEGER,

		Content BLOB,

		FOREIGN KEY(SerieID) REFERENCES blog_serie(ID)
	);

`

const (
	QueryAddBlogLanguage = "INSERT INTO blog_language(blog_id,language_id) VALUES (?,?)"
	QueryAddBlogSubject  = "INSERT INTO blog_subject(blog_id,subject_id) VALUES (?,?)"
	QueryAddBlogPost     = "INSERT INTO blog_post(Title,PostOn,UpdateOn,ThumbnailURL,Description,Author,VideoURL,Content) VALUES (?,?,?,?,?,?,?,?)"
	QueryAddBlogSerie    = "INSERT INTO blog_serie(Title,Description,ThumbnailURL,Over) VALUES (?,?,?,false)"

	QueryUpdateBlogPost = "UPDATE blog_post SET Title = ?, PostOn = ?,UpdateOn = ?,ThumbnailURL = ?,Description = ?.Author = ?,VideoURL = ?,Content = ? WHERE ID = ?"
	QueryAddBlogToSerie = "UPDATE blog_post SET SerieID = ? WHERE ID = ?"

	QueryBlogLanguage = "SELECT language_id FROM blog_language WHERE blog_id = ?"
	QueryBlogSubject  = "SELECT subject_id FROM blog_subject WHERE blog_id = ?"

	QueryBlogPostList        = "SELECT ID, Title, PostOn FROM blog_post"
	QueryBlogPostDescription = "SELECT ID, Title, PostOn, UpdateOn, ThumbnailURL, Description, Author, VideoURL FROM blog_post"
	QueryBlogContent         = "SELECT ID, Title, PostOn, UpdateOn, ThumbnailURL, Author, VideoURL, Content FROM blog_post WHERE ID = ?"
	QueryBlogSerieList       = "SELECT ID, Title, Description, ThumbnailURL, Over FROM blog_serie"
	QueryBlogSeriePostList   = QueryBlogPostDescription + " WHERE SerieID = ?"
	QueryBlogSinglePostList  = QueryBlogPostDescription + " WHERE SerieID IS NULL"
)

type DB struct {
	FilePath string
	tables   []string
	Db       *sql.DB
}

func (p *DB) Initialize(filePath string) {
	p.FilePath = filePath

	p.tables = []string{sqlTableBD}
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

func (p *DB) AddBlogSerie(serie *Serie) error {
	return db.StmtQuery(p, QueryAddBlogSerie, serie.Title, serie.Description, serie.ThumbnailURL)
}

func (p *DB) AddBlogToSerie(serie int, blog int) error {
	return db.StmtQuery(p, QueryAddBlogToSerie, serie, blog)
}

func (p *DB) AddBlogPost(post *Post) error {
	e := db.StmtQuery(p, QueryAddBlogPost, post.Name, post.PostOn, post.UpdateOn, post.ThumbnailURL, post.Description, post.Author, post.LinkVideoURL, post.PostMarkdown)
	if e != nil {
		return e
	}

	for _, s := range post.Subjects {
		e = db.StmtQuery(p, QueryAddBlogSubject, post.Name, s)
		if e != nil {
			return e
		}
	}

	for _, s := range post.Language {
		e = db.StmtQuery(p, QueryAddBlogLanguage, post.Name, s)
		if e != nil {
			return e
		}
	}

	return nil
}

func (p *DB) UpdateBlogPost(post *Post) error {
	// VA MANQUERA UPDATE LES LANGUAGES OU SUJETS
	return db.StmtQuery(p, QueryUpdateBlogPost, post.Name, post.PostOn, post.UpdateOn, post.ThumbnailURL, post.Description, post.Author, post.LinkVideoURL, post.PostMarkdown, post.ID)
}

func (p *DB) GetSerieList(loadPostDescription bool) ([]*Serie, error) {
	rows, err := p.Db.Query(QueryBlogSerieList)
	if err != nil {
		return nil, err
	}
	var ss []*Serie
	for rows.Next() {
		s := &Serie{}
		err = rows.Scan(&s.ID, &s.Title, &s.Description, &s.ThumbnailURL, &s.Over)
		if err != nil {
			return nil, err
		}
		ss = append(ss, s)
	}
	rows.Close()
	if loadPostDescription {
		for _, s := range ss {
			rows, err = p.Db.Query(QueryBlogSeriePostList, s.ID)
			if err != nil {
				return nil, err
			}
			for rows.Next() {
				pd, err := parseRowBlogDescription(rows)
				if err != nil {
					rows.Close()
					return nil, err
				}
				s.Posts = append(s.Posts, pd)
			}
			rows.Close()
			err = p.LoadLanguageSubjectPosts(s.Posts)
			if err != nil {
				return nil, err
			}
		}

	}
	return ss, nil
}

func (p *DB) GetBlogList() ([]*Post, error) {

	return nil, nil
}

func (p *DB) GetBlogDescriptionList() ([]*Post, error) {
	rows, err := p.Db.Query(QueryBlogSinglePostList)
	if err != nil {
		return nil, err
	}
	defer rows.Close()
	var pp []*Post
	for rows.Next() {
		p, err := parseRowBlogDescription(rows)
		if err != nil {
			return nil, err
		}
		pp = append(pp, p)
	}
	err = p.LoadLanguageSubjectPosts(pp)
	return pp, err
}

func (p *DB) LoadLanguageSubjectPosts(pl []*Post) error {
	var err error
	for _, post := range pl {
		post.Language, err = db.QueryStringArray(p, QueryBlogLanguage, post.Name)
		if err != nil {
			return err
		}
		post.Subjects, err = db.QueryStringArray(p, QueryBlogSubject, post.Name)
		if err != nil {
			return err
		}
	}
	return nil

}

func (p *DB) GetBlogContent(id int) (*Post, error) {
	post := &Post{}
	err := p.Db.QueryRow(QueryBlogContent, id).Scan(&post.ID, &post.Name, &post.PostOn, &post.UpdateOn, &post.ThumbnailURL, &post.Author, &post.LinkVideoURL, &post.PostMarkdown)
	if err != nil {
		return nil, err
	}
	return post, p.LoadLanguageSubjectPosts([]*Post{post})
}

func parseRowBlogDescription(r *sql.Rows) (*Post, error) {
	p := &Post{}
	return p, r.Scan(&p.ID, &p.Name, &p.PostOn, &p.UpdateOn, &p.ThumbnailURL, &p.Description, &p.Author, &p.LinkVideoURL)

}

func GetBlogDBInstance(file string) *DB {
	idb := &DB{}
	idb.Initialize(file)
	if err := db.OpenDatabase(idb); err != nil {
		panic(err)
	}
	return idb
}

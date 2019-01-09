package blog

const sqlTableBD = `
	CREATE TABLE IF NOT EXISTS blog_language(
		blog_id INTEGER NOT NULL,
		language_id VARCHAR(50) NOT NULL,

		FOREIGN KEY(blog_id) REFERENCES blog_post(Title)
		FOREIGN KEY(language_id) REFERENCES language(Name)
	);

	CREATE TABLE IF NOT EXISTS blog_subject(
		blog_id INTEGER NOT NULL,
		language_id VARCHAR(50) NOT NULL,

		FOREIGN KEY(blog_id) REFERENCES blog_post(Title)
		FOREIGN KEY(language_id) REFERENCES subject(Name) 
	);

	CREATE TABLE IF NOT EXISTS blog_post(
		ID INTEGER PRIMARY KEY AUTOINCREMENT,
		Title VARCHAR(255) NOT NULL
		PostOn VARCHAR(100) NOT NULL,
		UpdateOn VARCHAR(100) NOT NULL,
		ThumbnailURL VARCHAR(100) NOT NULL,
		Description VARCHAR(1000) NOT NULL,

		Author VARCHAR(100) NOT NULL,
		VideoURL VARCHAR(255),

		Content BLOB
	);
`

const (
	QueryAddBlogLanguage = "INSERT INTO blog_language(blog_id,language_id) VALUES (?,?)"
	QueryAddBlogSubject  = "INSERT INTO blog_subject(blog_id,subject_id) VALUES (?,?)"
	QueryAddBlogPost     = "INSERT INTO blog_post(Title,PostOn,UpdateOn,ThumbnailURL,Description,Author,VideoURL,Content) VALUES (?,?,?,?,?,?,?,?)"

	QueryBlogPostList        = "SELECT ID, Title, PostOn FROM blog_post"
	QueryBlogPostDescription = "SELECT ID, Title, PostOn, UpdateOn, ThumbnailURL, Description, Author, VideoURL FROM blog_post"
	QueryBlogContent         = "SELECT ID, Title, PostOn, UpdateOn, Author, VideoURL, Content FROM blog_post WHERE Title = ?"
)

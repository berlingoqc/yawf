package blog

// Post ..
type Post struct {
	ID           int
	Name         string
	PostOn       string
	UpdateOn     string
	ThumbnailURL string
	Description  string

	Author   string
	Subjects []string
	Language []string

	LinkVideoURL string

	PostMarkdown []byte
}

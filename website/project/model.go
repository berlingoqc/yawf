package project

import "time"

const (
	GitHubTimeFormat = ""
)

type GitHubOrganistaion struct {
	URL         string
	Name        string
	Description string

	MyRole string
}

// GitHubAccount contient l'informations de mon account
type GitHubAccount struct {
	Name   string
	URL    string
	ImgURL string

	Location       string
	Email          string
	Bio            string
	NbrPublicRepo  int
	NbrPublicGists int
	NbrFollorwers  int
	NbrFollowing   int
}

// GitHubRepo contient l'informations du des mes repo sur github
type GitHubRepo struct {
	URl         string
	Name        string
	Description string

	StarCount  int
	ForksCount int

	LastUpdateOn time.Time
	CreatedOn    time.Time
	CommitNumber int

	ReadMe []byte
}

type Projects struct {
	Account      *GitHubAccount
	Organization []*GitHubOrganistaion
	Programming  []*ProgrammingProject
}

type ProgrammingProject struct {
	Name string
	//DateStarted  time.Time
	ThumbnailURL string
	Subjects     []string
	Language     []string
	GitHub       *GitHubRepo
	SiteURL      string
	DocURL       string
	BlogsRelated []string
}

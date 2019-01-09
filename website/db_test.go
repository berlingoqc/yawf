package website

import (
	"os"
	"testing"
	"time"

	"github.com/berlingoqc/yawf/website/cv"

	"github.com/berlingoqc/yawf/db"
	"github.com/berlingoqc/yawf/website/project"
)

var dbName = "project.db"

func Terror(t *testing.T, err error) {
	if err != nil {
		t.Fatal(err)
	}
}

func TestProjectDB(t *testing.T) {
	os.Remove(dbName)
	defer func() {
		//err := os.Remove(dbName)
		//Terror(t, err)
	}()

	pdb, e := GetProjectDBInstance(dbName)
	Terror(t, e)
	defer db.CloseDatabse(pdb)

	lang := []string{"go", "jquery", "c++", "wasm", "glsl"}
	for _, l := range lang {
		Terror(t, pdb.AddLanguage(l))
	}

	cat := []string{"website", "auth", "blog", "sql", "3d", "sound", "emscripten"}
	for _, c := range cat {
		Terror(t, pdb.AddSubject(c))
	}

	user_info := &project.GitHubAccount{
		Name:           "berlingoqc",
		Bio:            "Berlingoqc aka William Quintal a hobyist from Quebec",
		Email:          "berlingoqc@gmail.com",
		ImgURL:         "https://avatars2.githubusercontent.com/u/11835662?s=460&v=4",
		URL:            "https://github.com/berlingoqc",
		Location:       "Quebec",
		NbrFollorwers:  6,
		NbrFollowing:   4,
		NbrPublicGists: 0,
		NbrPublicRepo:  9,
	}

	Terror(t, pdb.AddGHAccount(user_info))

	git_repo := &project.GitHubRepo{
		Name:         "yawf",
		CommitNumber: 5,
		CreatedOn:    time.Now(),
		LastUpdateOn: time.Now(),
		Description:  `yawf is a website a build to expose myself to the internet. By this i mean blogging, professional and personal information`,
		ForksCount:   0,
		StarCount:    5,
		URl:          "https://github.com/berlingoqc/yawf",
		ReadMe:       []byte("# Header 1"),
	}

	git_repo_2 := &project.GitHubRepo{
		Name:         "YASE",
		CommitNumber: 55,
		CreatedOn:    time.Now(),
		LastUpdateOn: time.Now(),
		Description:  `YASE stand for Yet Another Small Engine, it's a collection of software and librairie to create and reader 3D scene for game , video or animation`,
		ForksCount:   0,
		StarCount:    5,
		URl:          "https://github.com/berlingoqc/YASE",
		ReadMe:       []byte("# Header 1"),
	}

	Terror(t, pdb.AddGHRepo(git_repo))
	Terror(t, pdb.AddGHRepo(git_repo_2))

	p := &project.ProgrammingProject{
		Name:         "yawf website",
		DocURL:       "https://yawf.ca/doc/golang/pkg/github.com/berlingoqc/yawf",
		Language:     []string{"golang", "jquery"},
		Subjects:     []string{"website", "auth", "blog", "sql"},
		ThumbnailURL: "/static/image/website.png",
		SiteURL:      "https://yawf.ca/about",
	}

	p1 := &project.ProgrammingProject{
		Name:         "YASE",
		Language:     []string{"c++", "wasm", "glsl"},
		Subjects:     []string{"3d", "sound", "file"},
		ThumbnailURL: "/static/image/yase_tn.jpg",
		SiteURL:      "https://yawf.ca/yase/",
	}

	Terror(t, pdb.AddPProject(p, "yawf"))
	Terror(t, pdb.AddPProject(p1, "YASE"))

	formation := &cv.Formation{
		Name:        "Diplome d'etude secondaire en Art-Sport Etude",
		NameDiploma: "DES Diplome d'etude secondaire",
		School:      "Ecole Paul-Hubert",
		StartDate:   "2009/08/30",
		EndDate:     "2012/06/20",
		LengthYear:  3,
		Description: "High school diploma with the music program and science courses",
		Mention:     []string{"Sport and Art Program", ""},
	}

	formation_1 := &cv.Formation{
		Name:        "Technique d'informatique industrielle",
		NameDiploma: "DEC Diplome d'etude collegial",
		School:      "Cegep Levis-Lauzon",
		StartDate:   "2014/08/30",
		EndDate:     "2019/06/20",
		LengthYear:  5,
		Description: "Collegial Diploma in industrial computer science : vision , robotics , 3d , ia and game with Calculus Courses )",
		Mention:     []string{"D1 VolleyBall Team"},
	}

	Terror(t, pdb.AddFormation(formation))
	Terror(t, pdb.AddFormation(formation_1))

	experience := &cv.ProfessionalExperience{
		Job:         "Line Cook",
		Corporation: "La Cage, Brasserie Sportive",
		Location:    "Lévis and Boucherville",
		StartDate:   "2017/04/03",
		EndDate:     "2018/07/05",
		Description: `Cook the meal with a focus on speed , quality and high accuracy`,
	}

	Terror(t, pdb.AddExperience(experience))

	le_go := &cv.LanguageExperience{
		Name:  "golang",
		Level: "novice++",
		Year:  2,
		Description: `Go is one of my favorite programming language for my personal project at the moment,
		i love the simplicity and the tools that this language provide.
		`,
	}
	le_cpp := &cv.LanguageExperience{
		Name:        "c++",
		Level:       "beginer++",
		Year:        0,
		Description: `Language use nat school most often, favorite for bigger, modular and cross-plateform library`,
	}
	le_jquery := &cv.LanguageExperience{
		Name:        "jquery",
		Level:       "toadler",
		Year:        0,
		Description: "Not very confident but i create this website to improve myself in front-end development so maybe in the near future",
	}

	Terror(t, pdb.AddLanguageExperience(le_go))
	Terror(t, pdb.AddLanguageExperience(le_cpp))
	Terror(t, pdb.AddLanguageExperience(le_jquery))

	languages, e := pdb.GetLanguage()
	Terror(t, e)
	for _, l := range languages {
		t.Logf("Language : %v\n", l)
	}
	subject, e := pdb.GetSubject()
	Terror(t, e)
	for _, s := range subject {
		t.Logf("Subject : %v\n", s)
	}

	ghuser, e := pdb.GetGHAccount()
	Terror(t, e)
	t.Logf("GHUser : %v\n", ghuser)

	items, e := pdb.GetPProjects()
	Terror(t, e)
	for _, i := range items {
		t.Logf("Project : %v \n", i)
		t.Logf("GitHub : %v\n", i.GitHub)

	}

	fl, e := pdb.GetFormation()
	Terror(t, e)
	for _, i := range fl {
		t.Logf("Formation : %v\n", i)
	}

	el, e := pdb.GetExperience()
	Terror(t, e)
	for _, i := range el {
		t.Logf("Experience : %v\n", i)
	}

	lel, e := pdb.GetLanguageExperience()
	Terror(t, e)
	for _, i := range lel {
		t.Logf("Language Experience : %v\n", i)
	}

	// Get les nouvelles shits
	ua, err := project.UpdateAccountInfo("berlingoqc")
	if err != nil {
		t.Fatalf("Error task github %v \n", err)
		return
	}
	err = pdb.UpdateGHAccount(ua)
	if err != nil {
		t.Fatal(err)
	}
	repos, err := project.UpdateRepositoryInfo("berlingoqc", "YASE", "yawf")
	if err != nil {
		t.Fatalf("Error task github %v \n", err)
	}
	// update dans la bd
	for _, r := range repos {
		err = pdb.UpdateGitHubRepo(r)
		if err != nil {
			t.Fatal(err)
		}
	}

}

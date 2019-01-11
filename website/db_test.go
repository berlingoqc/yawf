package website

import (
	"os"
	"testing"
	"time"

	"github.com/berlingoqc/yawf/module/base"
	"github.com/berlingoqc/yawf/module/blog"

	"github.com/berlingoqc/yawf/db"
	"github.com/berlingoqc/yawf/module/project"
)

var dbName = "root/yawf.db"

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
	basedb := &base.DB{}
	basedb.Initialize(dbName)
	err := db.OpenDatabase(basedb)
	if err != nil {
		t.Fatal(err)
	}

	pdb, e := project.GetProjectDBInstance(dbName)
	Terror(t, e)
	defer db.CloseDatabse(pdb)

	bdb := blog.GetBlogDBInstance(dbName)

	lang := []string{"go", "jquery", "c++", "wasm", "glsl", "shell", "cli", "bash"}
	for _, l := range lang {
		Terror(t, basedb.AddLanguage(l))
	}

	cat := []string{"website", "auth", "blog", "sql", "3d", "sound", "emscripten", "linux", "archlinux", "wine"}
	for _, c := range cat {
		Terror(t, basedb.AddSubject(c))
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

	formation := &base.Formation{
		Name:        "Diplome d'etude secondaire en Art-Sport Etude",
		NameDiploma: "DES Diplome d'etude secondaire",
		School:      "Ecole Paul-Hubert",
		StartDate:   "2009/08/30",
		EndDate:     "2012/06/20",
		LengthYear:  3,
		Description: "High school diploma with the music program and science courses",
		Mention:     []string{"Sport and Art Program", ""},
	}

	formation_1 := &base.Formation{
		Name:        "Technique d'informatique industrielle",
		NameDiploma: "DEC Diplome d'etude collegial",
		School:      "Cegep Levis-Lauzon",
		StartDate:   "2014/08/30",
		EndDate:     "2019/06/20",
		LengthYear:  5,
		Description: "Collegial Diploma in industrial computer science : vision , robotics , 3d , ia and game with Calculus Courses )",
		Mention:     []string{"D1 VolleyBall Team"},
	}

	Terror(t, basedb.AddFormation(formation))
	Terror(t, basedb.AddFormation(formation_1))

	experience := &base.ProfessionalExperience{
		Job:         "Line Cook",
		Corporation: "La Cage, Brasserie Sportive",
		Location:    "Lévis and Boucherville",
		StartDate:   "2017/04/03",
		EndDate:     "2018/07/05",
		Description: `Cook the meal with a focus on speed , quality and high accuracy`,
	}

	Terror(t, basedb.AddExperience(experience))

	le_go := &base.LanguageExperience{
		Name:  "golang",
		Level: "novice++",
		Year:  2,
		Description: `Go is one of my favorite programming language for my personal project at the moment,
		i love the simplicity and the tools that this language provide.
		`,
	}
	le_cpp := &base.LanguageExperience{
		Name:        "c++",
		Level:       "beginer++",
		Year:        0,
		Description: `Language use nat school most often, favorite for bigger, modular and cross-plateform library`,
	}
	le_jquery := &base.LanguageExperience{
		Name:        "jquery",
		Level:       "toadler",
		Year:        0,
		Description: "Not very confident but i create this website to improve myself in front-end development so maybe in the near future",
	}

	Terror(t, basedb.AddLanguageExperience(le_go))
	Terror(t, basedb.AddLanguageExperience(le_cpp))
	Terror(t, basedb.AddLanguageExperience(le_jquery))

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

	fl, e := basedb.GetFormation()
	Terror(t, e)
	for _, i := range fl {
		t.Logf("Formation : %v\n", i)
	}

	el, e := basedb.GetExperience()
	Terror(t, e)
	for _, i := range el {
		t.Logf("Experience : %v\n", i)
	}

	lel, e := basedb.GetLanguageExperience()
	Terror(t, e)
	for _, i := range lel {
		t.Logf("Language Experience : %v\n", i)
	}

	blog_p1 := &blog.Post{
		Author:       "William Quintal",
		Description:  "In this post i talk about my favorite application on Linux. From gamming to graphics to programming",
		Subjects:     []string{"linux", "wine"},
		Language:     []string{"bash"},
		ThumbnailURL: "/static/blog/1/thumbnail.png",
		PostOn:       time.Now().String(),
		UpdateOn:     time.Now().String(),
		Name:         "X Favorite App on Linux",
		PostMarkdown: []byte(`
		# Hello i love you

		won't you tell me your name
		
		`),
	}

	blog_serie := &blog.Serie{
		Title: "Arch Linux From noob to master",
		Description: `Want to ditch that old windows and get in control of your computer ?
		 You love to customize ? You are scared to jump in the world of linux ?
		 
		 Fear no more ! I'm here to teach you the way
		 `,
		ThumbnailURL: "/static/serie/1/thumbnail.png",
	}

	blog_p_s1 := &blog.Post{
		Name:         "The journey from whatever to Linux",
		Author:       "William Quintal",
		Description:  "In this first post , i will explain the history of linux and the basic things to know",
		Subjects:     []string{"linux"},
		Language:     []string{"bash"},
		ThumbnailURL: "/static/blog/2/thumbnail.png",
		PostOn:       time.Now().String(),
		UpdateOn:     time.Now().String(),
		PostMarkdown: []byte(`
			# The journey from whatever to Linux
		`),
	}

	err = bdb.AddBlogPost(blog_p1)
	if err != nil {
		t.Fatal(err)
	}

	err = bdb.AddBlogPost(blog_p_s1)
	if err != nil {
		t.Fatal(err)
	}

	err = bdb.AddBlogSerie(blog_serie)
	if err != nil {
		t.Fatal(err)
	}

	// Get la liste de serie pour avoir l'id
	sl, err := bdb.GetSerieList(false)
	if err != nil {
		t.Fatal(err)
	}
	if len(sl) != 1 {
		t.Fatal("Length of serie list should be one")
	}
	idBs := sl[0].ID
	// get la liste des blog pour avoir l'id
	pl, err := bdb.GetBlogDescriptionList()
	if err != nil {
		t.Fatal(err)
	}
	idBp := pl[1].ID

	err = bdb.AddBlogToSerie(idBs, idBp)
	if err != nil {
		t.Fatal(err)
	}

	// le content du blog post
	pc, err := bdb.GetBlogContent(idBp)
	if err != nil {
		t.Fatal(err)
	}
	t.Log(pc)

	// get serie recursive
	sl, err = bdb.GetSerieList(true)
	if err != nil {
		t.Fatal(err)
	}

}

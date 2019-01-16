package route

import (
	"html/template"
)

type Link struct {
	Name string
	URL  string
}

type DropLink struct {
	Title string
	Links [][]Link
}

type Button struct {
	Style string
	Name  string
	URL   string
}

type NavigationBar struct {
	Title string

	Items   []interface{}
	Buttons []Button
}

func AddNavBarTmplFunc(mf template.FuncMap) {
	mf["islink"] = func(i interface{}) bool {
		_, ok := i.(Link)
		return ok
	}
	mf["isdrop"] = func(i interface{}) bool {
		_, ok := i.(DropLink)
		return ok
	}
	mf["getdrop"] = func(i interface{}) DropLink {
		return i.(DropLink)
	}
	mf["getlink"] = func(i interface{}) Link {
		return i.(Link)
	}

}

func GetNavigationBar(title string) *NavigationBar {
	nb := &NavigationBar{
		Title: title,
	}

	return nb

}

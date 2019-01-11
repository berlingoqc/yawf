package route

type Link struct {
}

type DropLink struct {
	Links [][]Link
}

type NavigationBar struct {
	Title string

	Items []interface{}
}

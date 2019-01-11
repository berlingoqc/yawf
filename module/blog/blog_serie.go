package blog

// Serie represente une series qui regroupe plusieurs blog ensemble
type Serie struct {
	ID           int
	Title        string
	Description  string
	ThumbnailURL string
	Over         bool
	// Represente les postes de ma serie, si le pointeur
	// est vide est qu'il n'est pas encore ecrit mais le titre peux
	// etre deja la
	Posts              []*Post
	LastUpdate         string
	CurrentPublication int
	TotalPublication   int
}

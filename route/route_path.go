package route

/*

// GetRoutePath crée une instance d'une route de base
func GetRoutePath(tmplFile string, path string) *RoutePath {
	// Get le nom du fichier dans le tmplFile
	b := &RoutePath{
		ContentTmplPath: tmplFile,
		Path:            path,
	}
	return b
}

func GetMarkdownRouteFile(fileName string, path string) *RoutePath {
	return GetMarkdownRoutePath(path, func() ([]byte, error) {
		return ioutil.ReadFile(fileName)
	})
}

func GetMarkdownPage(path string, name string) *RoutePath {
	r := &RoutePath{ContentTmplPath: "/shared/markdown_page.html", Path: path}
	r.Handler = func(m map[string]interface{}, r *http.Request) {
		m["File"] = name
	}
	return r
}

// GetMarkdownRoutePath retourne une route qui retour un fichier markdown
func GetMarkdownRoutePath(path string, getContent func() ([]byte, error)) *RoutePath {
	b := &RoutePath{
		ContentTmplPath: "/shared/markdown_page.html",
		Path:            path,
		Handler: func(m map[string]interface{}, r *http.Request) {
			d, e := getContent()
			if e != nil {
				m["Error"] = e
				return
			}
			m["md"] = string(blackfriday.Run(d))
		},
	}
	return b
}
*/

package blog

import (
	"net/http"

	"github.com/berlingoqc/yawf/db"

	"github.com/berlingoqc/yawf/website/route"
)

// ListQuery is the struct to query blog list
type ListQuery struct {
	SerieID int
}

// GetWidgets ...
func (b *Module) GetWidgets() []*route.Widget {
	var ll []*route.Widget
	ll = append(ll, &route.Widget{
		File:   "/blog/blog_list.html",
		Name:   "blog_list",
		Struct: &ListQuery{},
		Render: func(t interface{}, w http.ResponseWriter, r *http.Request) interface{} {
			l := t.(ListQuery)
			iidb, _ := b.GetDBInstance()
			idb := iidb.(*DB)
			defer db.CloseDatabse(idb)
			if l.SerieID > 0 {

			} else {
				ll, err := idb.GetBlogList()
				if err != nil {
					return err
				}
				return ll
			}
			return nil
		},
	})
	return ll

}

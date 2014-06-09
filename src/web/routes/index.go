package routes

import (
	"github.com/martini-contrib/render"
)

//
// Show the index page
//
func Home(r render.Render) {
	r.HTML(200, "index", nil)
}

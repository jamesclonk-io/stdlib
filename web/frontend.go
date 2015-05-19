package web

import (
	"net/http"

	"github.com/unrolled/render"
)

// Frontend has control over routing, rendering and page mastership.
type Frontend struct {
	Router     *Router
	Render     *render.Render
	PageMaster *PageMaster
}

func NewFrontend() *Frontend {
	r := render.New(render.Options{
		IndentJSON: true,
		Layout:     "layout",
		Extensions: []string{".html"},
	})
	router := NewRouter()
	pm := &PageMaster{"jamesclonk.io", "index", http.StatusOK, nil}

	f := &Frontend{router, r, pm}

	// default 404
	router.NotFoundHandler = f.NotFoundHandler("jamesclonk.io")

	return f
}

func (f *Frontend) SetNavigation(navbar NavBar) {
	f.PageMaster.Navbar = navbar
}

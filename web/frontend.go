package web

import (
	"net/http"

	"github.com/unrolled/render"
)

// Frontend has control over routing, rendering and page mastership.
type Frontend struct {
	Title      string
	Router     *Router
	Render     *render.Render
	PageMaster *PageMaster
}

func NewFrontend(title string) *Frontend {
	r := render.New(render.Options{
		IndentJSON: true,
		Layout:     "layout",
		Extensions: []string{".html"},
	})
	router := NewRouter()
	pm := &PageMaster{title, "index", http.StatusOK, nil}

	f := &Frontend{title, router, r, pm}

	// default 404
	router.NotFoundHandler = f.NotFoundHandler(title)

	return f
}

func (f *Frontend) SetNavigation(nav Navigation) {
	f.PageMaster.Navigation = nav
}

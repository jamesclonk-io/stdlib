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

// PageMaster contains default values for Title, Template and StatusCode.
// Also holds the central NavBar data.
type PageMaster struct {
	Title      string
	Template   string
	StatusCode int
	Navbar     NavBar
}

// Page represents a page to be rendered on the browser.
type Page struct {
	Title            string
	Navbar           NavBar
	ActiveNavElement string
	Content          interface{}
	Template         string
	StatusCode       int
	Error            error
}

// NavBar represents the navigation bar for the web page.
type NavBar []NavElement

// NavElement is an element of a NavBar or nested inside another NavElement.
type NavElement struct {
	Name     string
	Link     string
	Dropdown []NavElement
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

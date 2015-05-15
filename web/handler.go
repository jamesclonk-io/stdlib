package web

import (
	"net/http"

	"github.com/unrolled/render"
)

type Page struct {
	Title      string
	Navbar     string
	Content    interface{}
	Template   string
	StatusCode int
	Error      error
}

type Handler func(http.ResponseWriter, *http.Request) *Page

func NewHandler(title, navbar string, r *render.Render, fn Handler) http.HandlerFunc {
	return func(w http.ResponseWriter, req *http.Request) {
		page := fn(w, req)
		page.Title = title
		page.Navbar = navbar

		if len(page.Title) == 0 {
			page.Title = "jamesclonk.io"
		}
		if len(page.Template) == 0 {
			page.Template = "index"
		}
		if page.StatusCode == 0 {
			page.StatusCode = http.StatusOK
		}

		if page.Error != nil {
			r.HTML(w, page.StatusCode, "error", page)
			return
		}
		r.HTML(w, page.StatusCode, page.Template, page)
	}
}

func NotFoundHandler(title, navbar string, r *render.Render) http.HandlerFunc {
	return NewHandler(title, navbar, r, func(http.ResponseWriter, *http.Request) *Page {
		return &Page{
			StatusCode: http.StatusNotFound,
			Template:   "404",
		}
	})
}

func ErrorHandler(title, navbar string, r *render.Render, err error) http.HandlerFunc {
	return NewHandler(title, navbar, r, func(http.ResponseWriter, *http.Request) *Page {
		return &Page{
			StatusCode: http.StatusInternalServerError,
			Error:      err,
		}
	})
}

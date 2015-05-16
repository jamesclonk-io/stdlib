package web

import "net/http"

type Handler func(http.ResponseWriter, *http.Request) *Page

func (f *Frontend) NewHandler(fn Handler) http.HandlerFunc {
	return func(w http.ResponseWriter, req *http.Request) {
		page := fn(w, req)

		if len(page.Title) == 0 {
			page.Title = f.PageMaster.Title
		}
		if len(page.Template) == 0 {
			page.Template = f.PageMaster.Template
		}
		if page.StatusCode == 0 {
			page.StatusCode = f.PageMaster.StatusCode
		}
		if page.Navbar == nil {
			page.Navbar = f.PageMaster.Navbar
		}

		if page.Error != nil {
			f.Render.HTML(w, page.StatusCode, "error", page)
			return
		}
		f.Render.HTML(w, page.StatusCode, page.Template, page)
	}
}

func (f *Frontend) NotFoundHandler(title string) http.HandlerFunc {
	return f.NewHandler(func(http.ResponseWriter, *http.Request) *Page {
		return &Page{
			Title:      title,
			StatusCode: http.StatusNotFound,
			Template:   "404",
		}
	})
}

func Error(title string, code int, err error) *Page {
	return &Page{
		Title:      title,
		StatusCode: code,
		Error:      err,
		Template:   "error",
	}
}

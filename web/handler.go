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

func (b *Backend) NewHandler(fn Handler) http.HandlerFunc {
	return func(w http.ResponseWriter, req *http.Request) {
		page := fn(w, req)

		if page.StatusCode == 0 {
			page.StatusCode = http.StatusOK
		}

		if page.Headers != nil {
			for key, values := range page.Headers {
				w.Header().Del(key)
				for _, value := range values {
					w.Header().Add(key, value)
				}
			}
		}
		b.Render.JSON(w, page.StatusCode, page.Content)
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

func (b *Backend) NotFoundHandler(headers http.Header) http.HandlerFunc {
	return b.NewHandler(func(http.ResponseWriter, *http.Request) *Page {
		return &Page{
			Headers:    headers,
			StatusCode: http.StatusNotFound,
			Content:    nil,
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

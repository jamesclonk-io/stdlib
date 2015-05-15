package web

import (
	"github.com/gorilla/mux"
	"github.com/unrolled/render"
)

type Router struct {
	*mux.Router
	Render *render.Render
}

func NewRouter() *Router {
	r := render.New(render.Options{
		IndentJSON: true,
		Layout:     "layout",
		Extensions: []string{".html"},
	})

	router := mux.NewRouter()
	// default 404
	router.NotFoundHandler = NotFoundHandler("jamesclonk.io", "", r)

	return &Router{router, r}
}

func (r *Router) NewRoute(path, title, navbar string, handler Handler) *mux.Route {
	return r.Handle(path, NewHandler(title, navbar, r.Render, handler))
}

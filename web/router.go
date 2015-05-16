package web

import "github.com/gorilla/mux"

type Router struct {
	*mux.Router
}

func NewRouter() *Router {
	router := mux.NewRouter()
	return &Router{router}
}

func (f *Frontend) NewRoute(path string, handler Handler) *mux.Route {
	return f.Router.Handle(path, f.NewHandler(handler))
}

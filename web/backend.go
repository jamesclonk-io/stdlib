package web

import (
	"github.com/jamesclonk-io/stdlib/env"
	"github.com/unrolled/render"
)

// Backend has control over routing and rendering
type Backend struct {
	Router *Router
	Render *render.Render
}

func NewBackend() *Backend {
	r := render.New(render.Options{
		IndentJSON: true,
	})
	router := NewRouter()

	b := &Backend{router, r}

	// default 404
	router.NotFoundHandler = b.NotFoundHandler(nil)

	// enforce TLS certs for backend
	if len(env.Get("HTTP_CERT_FILE", "")) == 0 ||
		len(env.Get("HTTP_KEY_FILE", "")) == 0 {
		panic("Using web.Backend without TLS certificates is not allowed!")
	}

	return b
}

package web

import (
	"github.com/jamesclonk-io/stdlib/env"
	"github.com/unrolled/render"
)

// Backend has control over routing and rendering
type Backend struct {
	Router         *Router
	Render         *render.Render
	user, password string
}

func NewBackend() *Backend {
	// enforce TLS certs / HTTPS listener for backend
	env.MustGet("JCIO_HTTP_CERT_FILE")
	env.MustGet("JCIO_HTTP_KEY_FILE")

	// backend needs user & password for basic auth
	user := env.MustGet("JCIO_HTTP_AUTH_USER")
	password := env.MustGet("JCIO_HTTP_AUTH_PASSWORD")

	r := render.New(render.Options{
		IndentJSON: true,
	})
	router := NewRouter()

	b := &Backend{router, r, user, password}

	// default 404
	router.NotFoundHandler = b.NotFoundHandler(nil)

	return b
}

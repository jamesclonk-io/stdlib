package web

import (
	"net"
	"net/http"
	"time"

	"github.com/Sirupsen/logrus"
	"github.com/jamesclonk-io/stdlib/env"
	"github.com/jamesclonk-io/stdlib/logger"
	"github.com/tylerb/graceful"
)

type Server struct {
	Logger   *logrus.Logger
	Port     string
	CertFile string
	KeyFile  string
}

func NewServer() *Server {
	return &Server{
		Logger: logger.GetLogger(),
	}
}

func (s *Server) Start(handler http.Handler) {
	if len(s.Port) == 0 {
		s.Port = env.Get("PORT", "3000")
	}
	if len(s.CertFile) == 0 || len(s.KeyFile) == 0 {
		s.CertFile = env.Get("HTTP_CERT_FILE", "")
		s.KeyFile = env.Get("HTTP_KEY_FILE", "")
	}

	address := ":" + s.Port
	s.Logger.WithField("address", address).Info("Start HTTP Server")

	server := &graceful.Server{
		Server:  &http.Server{Addr: address, Handler: handler},
		Timeout: 15 * time.Second,
	}

	var err error
	if len(s.CertFile) > 0 && len(s.KeyFile) > 0 {
		err = server.ListenAndServeTLS(s.CertFile, s.KeyFile)
	} else {
		err = server.ListenAndServe()
	}
	if err != nil {
		if opErr, ok := err.(*net.OpError); !ok || (ok && opErr.Op != "accept") {
			s.Logger.Fatal(err)
		}
	}
}

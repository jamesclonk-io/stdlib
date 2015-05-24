package cms

import (
	"html/template"
	"sync"
	"time"

	"github.com/Sirupsen/logrus"
	"github.com/jamesclonk-io/stdlib/web"
)

type CMS struct {
	data  *CMSData
	input string
	mutex *sync.Mutex
	log   *logrus.Logger
}

type CMSData struct {
	Configuration *CMSConfiguration
	Navigation    *CMSNavigation
	Content       []*CMSContent
	Timestamp     time.Time
}

type CMSConfiguration struct {
}

type CMSNavigation struct {
	NavBar web.NavBar `json:"navbar"`
}

type CMSContent struct {
	Name     string
	Basename string
	Path     string
	Content  template.HTML
}

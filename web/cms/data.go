package cms

import (
	"html/template"
	"time"

	"github.com/jamesclonk-io/stdlib/web"
)

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

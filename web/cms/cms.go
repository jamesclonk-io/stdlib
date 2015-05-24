package cms

import (
	"html/template"
	"sync"
	"time"

	"github.com/Sirupsen/logrus"
	"github.com/jamesclonk-io/stdlib/env"
	"github.com/jamesclonk-io/stdlib/logger"
	"github.com/jamesclonk-io/stdlib/web"
)

type CMS struct {
	title    string
	frontend *web.Frontend
	data     *CMSData
	input    string
	mutex    *sync.Mutex
	log      *logrus.Logger
}

type CMSData struct {
	Configuration *CMSConfiguration
	Navigation    *CMSNavigation
	Content       []*CMSContent
	Timestamp     time.Time
}

type CMSConfiguration struct {
	Configuration map[string]interface{} `json:"configuration"`
}

type CMSNavigation struct {
	Navigation web.Navigation `json:"navigation"`
}

type CMSContent struct {
	Name     string
	Basename string
	Path     string
	Content  template.HTML
}

func NewCMS(frontend *web.Frontend) (*CMS, error) {
	data := &CMSData{}
	input := env.Get("JCIO_CMS_DATA", "https://github.com/jamesclonk-io/content/archive/master.zip")
	mutex := &sync.Mutex{}
	log := logger.GetLogger()

	cms := &CMS{frontend.Title, frontend, data, input, mutex, log}
	if err := cms.checkData(true); err != nil {
		return nil, err
	}
	return cms, nil
}

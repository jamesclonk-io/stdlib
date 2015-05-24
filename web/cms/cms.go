package cms

import (
	"fmt"
	"html/template"
	"net/http"
	"path"
	"sync"
	"time"

	"github.com/Sirupsen/logrus"
	"github.com/jamesclonk-io/stdlib/env"
	"github.com/jamesclonk-io/stdlib/logger"
	"github.com/jamesclonk-io/stdlib/web"
)

func NewCMS() (*CMS, error) {
	data := &CMSData{}
	input := env.Get("JCIO_CMS_DATA", "https://github.com/jamesclonk-io/content/archive/master.zip")
	mutex := &sync.Mutex{}
	log := logger.GetLogger()

	cms := &CMS{data, input, mutex, log}
	return cms, cms.refreshData()
}

func (c *CMS) ViewHandler(title string, w http.ResponseWriter, req *http.Request) *web.Page {
	filename := path.Base(req.RequestURI)

	// find file
	var html template.HTML
	for _, c := range c.data.Content {
		if path.Join("/", c.Path, c.Name) == req.RequestURI {
			html = c.Content
		}
	}

	// wrap into struct
	content := struct {
		Title string
		HTML  template.HTML
	}{
		Title: filename,
		HTML:  html,
	}

	return &web.Page{
		Title:      fmt.Sprintf("%s - %s", title, filename),
		ActiveLink: path.Dir(req.RequestURI),
		Content:    content,
		Template:   "things",
	}
}

func (c *CMS) RefreshHandler(title string, navbar *web.NavBar, thingsIndex int) web.Handler {
	return func(w http.ResponseWriter, req *http.Request) *web.Page {
		c.mutex.Lock()
		defer c.mutex.Unlock()

		if err := c.checkData(navbar, thingsIndex, true); err != nil {
			return web.Error(title, http.StatusInternalServerError, err)
		}

		return &web.Page{
			Title:      title,
			ActiveLink: "/",
			Content:    nil,
			Template:   "index",
		}
	}
}

func (c *CMS) checkData(navbar *web.NavBar, navIndex int, refresh bool) error {
	// reset data if not set
	if c.data == nil {
		c.data = &CMSData{}
	}

	// refresh either every 12 hours, or if refresh parameter set to true
	if time.Since(c.data.Timestamp).Hours() > 12 || refresh {
		if err := c.refreshData(); err != nil {
			return err
		}

		// create new navbar elements
		navElements := make([]web.NavElement, 0)
		for _, content := range c.data.Content {
			navElements = append(navElements, web.NavElement{
				Name:     path.Join("/", content.Path, content.Basename),
				Link:     path.Join("/101", content.Path, content.Name),
				Dropdown: nil,
			})
		}

		// reset navigation bar element for "101"
		(*navbar)[navIndex].Dropdown = navElements
	}
	return nil
}

func (c *CMS) refreshData() (err error) {
	c.data, err = getDataFromZip(c.input)
	if err != nil {
		c.log.WithFields(logrus.Fields{
			"error": err,
			"file":  c.input,
		}).Error("Could not refresh data")
		return err
	}
	return nil
}

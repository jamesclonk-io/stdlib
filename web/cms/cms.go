package cms

import (
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

var (
	cmsData  *CMSData
	cmsInput string
	cmsMutex *sync.Mutex
	log      *logrus.Logger
)

func init() {
	cmsData = &CMSData{}
	cmsInput = env.Get("JCIO_CMS_DATA", "https://github.com/jamesclonk-io/content/archive/master.zip")
	cmsMutex = &sync.Mutex{}
	log = logger.GetLogger()
}

func ViewHandler(w http.ResponseWriter, req *http.Request) *web.Page {
	filename := path.Base(req.RequestURI)

	// find file
	var html template.HTML
	for _, c := range cmsData.Content {
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
		Title:      "jamesclonk.io - 101 Things - " + filename,
		ActiveLink: path.Dir(req.RequestURI),
		Content:    content,
		Template:   "things",
	}
}

func ThingsRefreshHandler(navbar *web.NavBar, thingsIndex int) web.Handler {
	return func(w http.ResponseWriter, req *http.Request) *web.Page {
		cmsMutex.Lock()
		defer cmsMutex.Unlock()

		if err := checkData(navbar, thingsIndex, true); err != nil {
			return web.Error("jamesclonk.io", http.StatusInternalServerError, err)
		}

		return &web.Page{
			Title:      "jamesclonk.io - Refresh",
			ActiveLink: "/",
			Content:    nil,
			Template:   "index",
		}
	}
}

func checkData(navbar *web.NavBar, navIndex int, refresh bool) error {
	// reset data if not set
	if cmsData == nil {
		cmsData = &CMSData{}
	}

	// refresh either every 12 hours, or if refresh parameter set to true
	if time.Since(cmsData.Timestamp).Hours() > 12 || refresh {
		if err := refreshData(cmsInput); err != nil {
			return err
		}

		// create new navbar elements
		navElements := make([]web.NavElement, 0)
		for _, content := range cmsData.Content {
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

func refreshData(input string) (err error) {
	cmsData, err = getDataFromZip(input)
	if err != nil {
		log.WithFields(logrus.Fields{
			"error": err,
			"file":  cmsInput,
		}).Error("Could not refresh data")
		return err
	}
	return nil
}

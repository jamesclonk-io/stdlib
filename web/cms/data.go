package cms

import (
	"encoding/json"
	"html/template"
	"path"
	"path/filepath"
	"strings"
	"time"

	"github.com/Sirupsen/logrus"
	"github.com/jamesclonk-io/stdlib/web"
	"github.com/russross/blackfriday"
)

func (c *CMS) checkData(refresh bool) error {
	// reset data if not set
	if c.data == nil {
		c.data = &CMSData{}
	}

	// refresh either every 24 hours, or if refresh parameter set to true
	if time.Since(c.data.Timestamp).Hours() >= 24 || refresh {
		// refresh cms content data
		if err := c.refreshData(); err != nil {
			return err
		}

		// update navigation
		c.frontend.SetNavigation(c.GetNavigation())
	}
	return nil
}

func (c *CMS) refreshData() (err error) {
	if err := c.getData(); err != nil {
		c.log.WithFields(logrus.Fields{
			"error": err,
			"file":  c.input,
		}).Error("Could not refresh data")
		return err
	}
	return nil
}

func (c *CMS) getData() error {
	var data map[string][]byte
	var err error

	if strings.HasSuffix(c.input, ".zip") {
		if strings.HasPrefix(c.input, "http") {
			// read zip content from url
			data, err = c.readZipFromURL()
		} else {
			// read zip content from local file
			data, err = c.readZipFromFile()
		}
	} else {
		// read zip content from local folder
		data, err = c.readFromFolder()
	}
	if err != nil {
		return err
	}

	c.data = &CMSData{
		Content:   make([]*CMSContent, 0),
		Timestamp: time.Now(),
	}

	// go through all files
	for file, bytes := range data {
		basename := filepath.Base(file)

		if strings.HasSuffix(basename, ".md") {
			html := blackfriday.MarkdownCommon(bytes) // generate HTML from markdown
			content := &CMSContent{
				Name:     path.Base(file),
				Basename: strings.TrimSuffix(basename, filepath.Ext(basename)),
				Path:     path.Dir(file),
				Content:  template.HTML(html),
			}
			c.data.Content = append(c.data.Content, content)

		} else if basename == "navigation.json" {
			var nav CMSNavigation
			if err := json.Unmarshal(bytes, &nav); err != nil {
				return err
			}
			c.data.Navigation = &nav

		} else if basename == "configuration.json" {
			var config CMSConfiguration
			if err := json.Unmarshal(bytes, &config); err != nil {
				return err
			}
			c.data.Configuration = &config
		}
	}
	return nil
}

func (c *CMS) GetNavigation() web.Navigation {
	return c.data.Navigation.Navigation
}

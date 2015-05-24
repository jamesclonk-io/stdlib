package cms

import (
	"time"

	"github.com/Sirupsen/logrus"
	"github.com/jamesclonk-io/stdlib/web"
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
	if err := c.getDataFromZip(); err != nil {
		c.log.WithFields(logrus.Fields{
			"error": err,
			"file":  c.input,
		}).Error("Could not refresh data")
		return err
	}
	return nil
}

func (c *CMS) GetNavigation() web.Navigation {
	return c.data.Navigation.Navigation
}

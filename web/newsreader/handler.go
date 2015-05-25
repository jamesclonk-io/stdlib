package newsreader

import (
	"fmt"
	"net/http"

	"github.com/jamesclonk-io/stdlib/web"
)

func (n *NewsReader) ViewHandler(w http.ResponseWriter, req *http.Request) *web.Page {
	feeds, err := n.getFeeds()
	if err != nil {
		return web.Error("jamesclonk.io - News Error", http.StatusInternalServerError, err)
	}

	return &web.Page{
		Title:      fmt.Sprintf("jamesclonk.io - News"),
		ActiveLink: req.RequestURI,
		Content:    feeds,
		Template:   "news",
	}
}

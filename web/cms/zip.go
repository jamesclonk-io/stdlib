package cms

import (
	"archive/zip"
	"bytes"
	"encoding/json"
	"fmt"
	"html/template"
	"io/ioutil"
	"net/http"
	"path"
	"path/filepath"
	"sort"
	"strings"
	"time"

	"github.com/russross/blackfriday"
)

func (c *CMS) getDataFromZip() error {
	var zipData map[string]string
	var err error

	if strings.HasPrefix(c.input, "http") {
		// read zip content from url
		zipData, err = c.readZipFromURL()
	} else {
		// read zip content from local file
		zipData, err = c.readZipFromFile()
	}
	if err != nil {
		return err
	}

	c.data = &CMSData{
		Content:   make([]*CMSContent, 0),
		Timestamp: time.Now(),
	}

	// collect and sort all files
	var files []string
	for file, _ := range zipData {
		files = append(files, file)
	}
	sort.Sort(sort.StringSlice(files))

	// store root folder
	root := path.Dir(files[0])
	root = root[:strings.Index(root, "/")]

	// strip away root folder from all paths
	for i := range files {
		files[i] = strings.TrimPrefix(files[i], root)
	}

	// go through all files
	for _, file := range files {
		basename := filepath.Base(file)
		data := []byte(zipData[root+file])

		if strings.HasSuffix(basename, ".md") {
			html := blackfriday.MarkdownCommon(data) // generate HTML from markdown
			content := &CMSContent{
				Name:     path.Base(file),
				Basename: strings.TrimSuffix(basename, filepath.Ext(basename)),
				Path:     path.Dir(file),
				Content:  template.HTML(html),
			}
			c.data.Content = append(c.data.Content, content)

		} else if basename == "navigation.json" {
			var nav CMSNavigation
			if err := json.Unmarshal(data, &nav); err != nil {
				return err
			}
			c.data.Navigation = &nav

		} else if basename == "configuration.json" {
			var config CMSConfiguration
			if err := json.Unmarshal(data, &config); err != nil {
				return err
			}
			c.data.Configuration = &config
		}
	}
	return nil
}

func (c *CMS) readZipFromURL() (map[string]string, error) {
	resp, err := http.Get(c.input)
	if err != nil {
		return nil, err
	}
	defer resp.Body.Close()

	data, err := ioutil.ReadAll(resp.Body)
	if err != nil {
		return nil, err
	}
	return readZip(bytes.NewReader(data), int64(len(data)))
}

func (c *CMS) readZipFromFile() (map[string]string, error) {
	data, err := ioutil.ReadFile(c.input)
	if err != nil {
		return nil, err
	}
	return readZip(bytes.NewReader(data), int64(len(data)))
}

func readZip(data *bytes.Reader, size int64) (map[string]string, error) {
	if size < 42 {
		return nil, fmt.Errorf("CMS data invalid, size too small: %d", size)
	}

	r, err := zip.NewReader(data, size)
	if err != nil {
		return nil, err
	}

	// store data in a map of {filename:content}
	contents := make(map[string]string, 0)
	for _, f := range r.File {
		if f.FileInfo().IsDir() {
			continue
		}
		if !strings.HasSuffix(f.Name, ".md") &&
			!strings.HasSuffix(f.Name, ".json") {
			continue
		}

		rc, err := f.Open()
		if err != nil {
			return nil, err
		}
		defer rc.Close()

		data, err := ioutil.ReadAll(rc)
		if err != nil {
			return nil, err
		}
		contents[f.Name] = string(data)
	}
	return contents, nil
}

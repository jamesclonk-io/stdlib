package cms

import (
	"encoding/json"
	"io/ioutil"

	"github.com/jamesclonk-io/stdlib/web"
)

func GetNavBar() (web.NavBar, error) {
	nav, err := readNavigation() // TODO: read navigation data from CMSData
	if err != nil {
		return nil, err
	}
	return nav.NavBar, nil
}

func readNavigation() (*CMSNavigation, error) {
	in, err := ioutil.ReadFile("../content/navigation.json") // TODO: read navigation data from CMSData
	if err != nil {
		return nil, err
	}

	var data CMSNavigation
	if err := json.Unmarshal(in, &data); err != nil {
		panic(err)
		return nil, err
	}
	return &data, nil
}

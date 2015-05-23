package web

import "net/http"

// PageMaster contains default values for Title, Template and StatusCode.
// Also holds the central NavBar data.
type PageMaster struct {
	Title      string
	Template   string
	StatusCode int
	Navbar     NavBar
}

// Page represents a page to be rendered on the browser.
type Page struct {
	Title      string
	Navbar     NavBar
	ActiveLink string
	Headers    http.Header
	Content    interface{}
	Template   string
	StatusCode int
	Error      error
}

// NavBar represents the navigation bar for the web page.
type NavBar []NavElement

// NavElement is an element of a NavBar or nested inside another NavElement.
type NavElement struct {
	Name     string       `json:"name"`
	Link     string       `json:"link,omitempty"`
	Dropdown []NavElement `json:"dropdown,omitempty"`
}

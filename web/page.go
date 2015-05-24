package web

import "net/http"

// PageMaster contains default values for Title, Template and StatusCode.
// Also holds the central Navigation data.
type PageMaster struct {
	Title      string
	Template   string
	StatusCode int
	Navigation Navigation
}

// Page represents a page to be rendered on the browser.
type Page struct {
	Title      string
	Navigation Navigation
	ActiveLink string
	Headers    http.Header
	Content    interface{}
	Template   string
	StatusCode int
	Error      error
}

// Navigation represents the navigation for the web page.
type Navigation []NavigationElement

// NavigationElement is a navigation element of the web page or nested inside another NavigationElement.
type NavigationElement struct {
	Name     string              `json:"name"`
	Link     string              `json:"link,omitempty"`
	Dropdown []NavigationElement `json:"dropdown,omitempty"`
}

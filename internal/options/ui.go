package options

import (
	"net/http"
)

type UIObj struct {
	hndlr     http.Handler
	tmplPath  string
	asstsPath string
	title     string
}

// Get UI handler
func (obj *UIObj) GetUIHandler() http.Handler {
	return obj.hndlr
}

// Set UI handler to obj
func (obj *UIObj) SetUIHandler(h http.Handler) {
	obj.hndlr = h
}

// Get UI title
func (obj *UIObj) GetUITitle() string {
	return obj.title
}

// Get UI path to html templates
func (obj *UIObj) GetUITemplPath() string {
	return obj.tmplPath
}

// Get UI path to assets
func (obj *UIObj) GetUIAssetsPath() string {
	return obj.asstsPath
}

package options

import (
	"github.com/go-chi/chi/v5"
	"html/template"
	"net/http"
)

// Main user-interface object
type UIObj struct {
	server   *http.Server
	router   chi.Router
	htmlNavi []NaviMenu
	tmpl     *template.Template
	//pData	*PagesData
}

// Get UI http server
func (obj *UIObj) GetUIServer() *http.Server {
	return obj.server
}

// Set UI http server to obj
func (obj *UIObj) SetUIServer(s *http.Server) {
	obj.server = s
}

// Get UI http router
func (obj *UIObj) GetUIRouter() chi.Router {
	return obj.router
}

// Set UI server to obj
func (obj *UIObj) SetUIRouter(r chi.Router) {
	obj.router = r
}

// Get UI Navigation menu
func (obj *UIObj) GetUINaviMenu() []NaviMenu {
	return obj.htmlNavi
}

// Set UI Navigation menu value
func (obj *UIObj) SetUINaviMenu(m []NaviMenu) {
	obj.htmlNavi = m
}

// Get UI compilled HTML templates
func (obj *UIObj) GetUIHTMLTemplates() *template.Template {
	return obj.tmpl
}

// Set UI compilled HTML templates
func (obj *UIObj) SetUIHTMLTemplates(t *template.Template) {
	obj.tmpl = t
}


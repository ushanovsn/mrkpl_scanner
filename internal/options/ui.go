package options

import (
	"github.com/go-chi/chi/v5"
	"net/http"
)

type UIObj struct {
	server *http.Server
	router chi.Router
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

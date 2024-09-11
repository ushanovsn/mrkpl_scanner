package ui

import (
	"mrkpl_scanner/internal/handlers"
	opt "mrkpl_scanner/internal/options"
	"net/http"

	"github.com/go-chi/chi/v5"
)

// Init new UI object
func NewUI() *opt.UIObj {
	ui := &opt.UIObj{}
	ui.SetDefaultUI()
	ui.SetUIServer(&http.Server{})
	ui.SetUIRouter(chi.NewRouter())
	return ui
}

// Starting UI Server process with configuration
func StartUIServer(scnr *opt.ScannerObj) error {
	log := scnr.GetLogger()
	log.Out("Starting Web UI...")

	setRouter(scnr)

	return http.ListenAndServe(scnr.GetAddr(), scnr.GetUIObj().GetUIRouter())
}

// Custom mux with router interface
func setRouter(scnr *opt.ScannerObj) {
	route := scnr.GetUIObj().GetUIRouter()

	// fileserver for static files and css
	fs := http.FileServer(http.Dir("./static/assets"))
	// use file server for all addresses with "assets"
	route.Handle("/assets/*", http.StripPrefix("/assets/", fs))

	// set handlers for routes
	route.Route("/", func(r chi.Router) {
		r.Get("/", handlers.IndexPage(scnr))
		r.Route("/gparams", func(r chi.Router) {
			r.Get("/", handlers.GParamsPage(scnr))
			r.Post("/", handlers.GParamsPageUpload(scnr))
		})
		r.Route("/pparams", func(r chi.Router) {
			r.Get("/", handlers.PParamsPage(scnr))
		})
		r.Route("/status", func(r chi.Router) {
			r.Get("/", handlers.StatusPage(scnr))
		})
	})
}

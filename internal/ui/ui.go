package ui

import (
	"context"
	"fmt"
	"mrkpl_scanner/internal/handlers"
	opt "mrkpl_scanner/internal/options"
	"net/http"
	"html/template"

	"github.com/go-chi/chi/v5"
)

// Init new UI object
func NewUI() *opt.UIObj {
	ui := &opt.UIObj{}
	ui.SetDefaultUI()
	ui.SetUIServer(&http.Server{})
	ui.SetUIRouter(chi.NewRouter())
	ui.SetUINaviMenu(opt.GetNavigationMenu())
	ui.SetUIHTMLTemplates(template.Must(template.ParseGlob("./static/htmltemplates/*")))
	return ui
}

// Starting UI Server process with configuration
func StartUIServer(scnr *opt.ScannerObj) {
	log := scnr.GetLogger()
	log.Out("Starting Web UI...")

	// init routers
	setRouter(scnr)

	srv := scnr.GetUIObj().GetUIServer()
	// set server parameters
	srv.Addr = scnr.GetAddr()
	srv.Handler = scnr.GetUIObj().GetUIRouter()

	scnr.GetWG().Add(1)
	// start server
	go func() {
		err := srv.ListenAndServe()
		if err != http.ErrServerClosed {
			log.Out(fmt.Sprintf("Error while WEB UI stopping: %s", err.Error()))
		}
		scnr.GetWG().Done()
	}()
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
		r.Route("/task_param_scan", func(r chi.Router) {
			r.Get("/", handlers.ParamsScanPage(scnr))
			r.Post("/", handlers.ParamsScanPageUpload(scnr))
		})
		r.Route("/pparams", func(r chi.Router) {
			r.Get("/", handlers.PParamsPage(scnr))
		})
		r.Route("/status", func(r chi.Router) {
			r.Get("/", handlers.StatusPage(scnr))
		})
	})
}

// Stopping UI Server
func StopUIServer(scnr *opt.ScannerObj) error {
	log := scnr.GetLogger()
	log.Info("Stopping Web UI...")

	ctx := context.Background()
	// set server parameters
	err := scnr.GetUIObj().GetUIServer().Shutdown(ctx)

	log.Info("Web UI Stopped...")

	return err
}

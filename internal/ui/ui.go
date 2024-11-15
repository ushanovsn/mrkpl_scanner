package ui

import (
	"context"
	"fmt"
	"html/template"
	"mrkpl_scanner/internal/options"
	"net/http"
	"github.com/go-chi/chi/v5"
)

// Init new UI object
func NewUI() *options.UIObj {
	ui := &options.UIObj{}
	ui.SetDefaultUI()
	ui.SetUIServer(&http.Server{})
	ui.SetUIRouter(chi.NewRouter())
	ui.SetUINaviMenu(GetNavigationMenu())
	ui.SetUIHTMLTemplates(template.Must(template.ParseGlob("./static/htmltemplates/*")))
	return ui
}

// Starting UI Server process with configuration
func StartUIServer(scnr *options.ScannerObj) {
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
func setRouter(scnr *options.ScannerObj) {
	route := scnr.GetUIObj().GetUIRouter()

	// fileserver for static files and css
	fs := http.FileServer(http.Dir("./static/assets"))
	// use file server for all addresses with "assets"
	route.Handle("/assets/*", http.StripPrefix("/assets/", fs))

	// set handlers for routes
	route.Route("/", func(r chi.Router) {
		r.Route("/",  func(r chi.Router) {
			r.Handle("/", UIHndlr{scnr, IndexPageHndlr})
			r.Get("/cmd/{cmd}", ProcessCMD(scnr))
		})
		r.Route("/task_param_scan", func(r chi.Router) {
			r.Handle("/", UIHndlr{scnr, ParamsScanPage})
		})
		r.Route("/notify", func(r chi.Router) {
			r.Get("/", TemplatePage(scnr))
		})
		r.Route("/parser_config_wb", func(r chi.Router) {
			r.Handle("/", UIHndlr{scnr, ConfWB})
		})
		r.Route("/log", func(r chi.Router) {
			r.Handle("/", UIHndlr{scnr, LogPageHndlr})
		})

		// templates in design
		r.Route("/parser_config_ozon", func(r chi.Router) {
			r.Get("/", TemplatePage(scnr))
		})
		r.Route("/parser_config_yandex_market", func(r chi.Router) {
			r.Get("/", TemplatePage(scnr))
		})
		r.Route("/parser_config_avito", func(r chi.Router) {
			r.Get("/", TemplatePage(scnr))
		})
		r.Route("/task_param_monitor", func(r chi.Router) {
			r.Get("/", TemplatePage(scnr))
		})
		r.Route("/task_param_search", func(r chi.Router) {
			r.Get("/", TemplatePage(scnr))
		})
	})
}

// Stopping UI Server
func StopUIServer(scnr *options.ScannerObj) error {
	log := scnr.GetLogger()
	log.Info("Stopping Web UI...")

	ctx := context.Background()
	// set server parameters
	err := scnr.GetUIObj().GetUIServer().Shutdown(ctx)

	log.Info("Web UI Stopped...")

	return err
}

package ui

import (
	"fmt"
	"mrkpl_scanner/internal/options"
	"net/http"
)

// Processor for Index page "/"
func IndexPageHndlr(scnr *options.ScannerObj, w http.ResponseWriter, r *http.Request) (int, error) {
	
	if r.Method != "GET" {
		return http.StatusMethodNotAllowed, fmt.Errorf("Method %s not allowed", r.Method)
	}

	tmpl := scnr.GetUIObj().GetUIHTMLTemplates()

	// struct for index page
	data := struct {
		Title      string
		NaviMenu   []options.NaviMenu
		ActiveMenu options.NaviActiveMenu
	}{
		Title:    options.DefUIPageTitle,
		NaviMenu: scnr.GetUIObj().GetUINaviMenu(),
		ActiveMenu: options.NaviActiveMenu{
			ActiveTabVal:    "/",
			ActiveDMenuVal:  "",
			PageDescription: "Главная",
		},
	}

	// base headers
	header := http.StatusOK
	w.Header().Add("Content-Type", "text/html; charset=utf-8")

	// OK. Processing template - replace\substitute values in template and send to front
	w.WriteHeader(header)

	return http.StatusOK, tmpl.ExecuteTemplate(w, "index", data)
}



package ui

import (
	"fmt"
	"mrkpl_scanner/internal/options"
	"net/http"
)

// Processor for page with configuration WB Parser
func ConfWB(scnr *options.ScannerObj, w http.ResponseWriter, r *http.Request) (int, error) {
	if r.Method != "POST" && r.Method != "GET" {
		return http.StatusMethodNotAllowed, fmt.Errorf("Method %s not allowed", r.Method)
	}

	tmpl := scnr.GetUIObj().GetUIHTMLTemplates()

	// struct for WB parser config
	data := struct {
		Title		string
		NaviMenu	[]options.NaviMenu
		ActiveMenu	options.NaviActiveMenu
		ConfWBData	*options.ConfWBPageData
	}{
		Title:    options.DefUIPageTitle,
		NaviMenu: scnr.GetUIObj().GetUINaviMenu(),
		ActiveMenu: options.NaviActiveMenu{
			ActiveTabVal:    "/config",
			ActiveDMenuVal:  "/parser_config_wb",
			PageDescription: "Конфигурация парсера \"Wildberries\"",
		},
		ConfWBData: getPageConfWBData(scnr),
	}

	// process form from page
	if r.Method == "POST" {
		if err := procConfWB(data.ConfWBData, r, scnr); err != nil {
			return http.StatusInternalServerError, fmt.Errorf("Error processing uploaded data from page \"parser_config_wb\", error: %s", err.Error())
		}
	}

	w.Header().Add("Content-Type", "text/html; charset=utf-8")
	// execute template and apply data
	return http.StatusOK, tmpl.ExecuteTemplate(w, "parser_config_wb", data)
}


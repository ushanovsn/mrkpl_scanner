package ui

import (
	"fmt"
	"mrkpl_scanner/internal/options"
	"net/http"
)

// Processor for page with Scan parameters
func ParamsScanPage(scnr *options.ScannerObj, w http.ResponseWriter, r *http.Request) (int, error) {
	if r.Method != "POST" && r.Method != "GET" {
		return http.StatusMethodNotAllowed, fmt.Errorf("Method %s not allowed", r.Method)
	}

	tmpl := scnr.GetUIObj().GetUIHTMLTemplates()

	// struct for Scanner Parameters page
	data := struct {
		Title          string
		NaviMenu       []options.NaviMenu
		ActiveMenu     options.NaviActiveMenu
		ParamsScanData *options.ParamsScanPageData
	}{
		Title:    options.DefUIPageTitle,
		NaviMenu: scnr.GetUIObj().GetUINaviMenu(),
		ActiveMenu: options.NaviActiveMenu{
			ActiveTabVal:    "/task_param",
			ActiveDMenuVal:  "/task_param_scan",
			PageDescription: "Параметры задачи \"Сканирование\"",
		},
		ParamsScanData: getPageTaskScanData(scnr),
	}

	// process form from page
	if r.Method == "POST" {
		if err := procTaskParamScan(data.ParamsScanData, r, scnr); err != nil {
			return http.StatusInternalServerError, fmt.Errorf("Error processing uploaded data from page \"task_param_scan\", error: %s", err.Error())
		}
	}

	w.Header().Add("Content-Type", "text/html; charset=utf-8")
	// execute template and apply data
	return http.StatusOK, tmpl.ExecuteTemplate(w, "task_param_scan", data)
}


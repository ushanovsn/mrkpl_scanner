package ui

import (
	"encoding/json"
	"fmt"
	"mrkpl_scanner/internal/options"
	"net/http"

	"github.com/go-chi/chi/v5"
)

const (
	defIndexName = "IndexPageData"
)

// Processor for Index page "/"
func IndexPageHndlr(scnr *options.ScannerObj, w http.ResponseWriter, r *http.Request) (int, error) {
	
	if r.Method != "GET" {
		return http.StatusMethodNotAllowed, fmt.Errorf("Method %s not allowed", r.Method)
	}

	tmpl := scnr.GetUIObj().GetUIHTMLTemplates()

	// struct for index page
	data := struct {
		Title		string
		NaviMenu	[]options.NaviMenu
		ActiveMenu	options.NaviActiveMenu
		Scanner		*options.IndexPageData
	}{
		Title:    options.DefUIPageTitle,
		NaviMenu: scnr.GetUIObj().GetUINaviMenu(),
		ActiveMenu: options.NaviActiveMenu{
			ActiveTabVal:    "/",
			ActiveDMenuVal:  "",
			PageDescription: "Главная",
		},
		Scanner: getPageIndexData(scnr),
	}

	// base headers
	header := http.StatusOK
	w.Header().Add("Content-Type", "text/html; charset=utf-8")

	// OK. Processing template - replace\substitute values in template and send to front
	w.WriteHeader(header)

	return http.StatusOK, tmpl.ExecuteTemplate(w, "index", data)
}



// Aggregate and return actual values for page "page-task_param_scan"
func getPageIndexData(scnr *options.ScannerObj) *options.IndexPageData {
	log := scnr.GetLogger()

	// init empty data
	d := options.IndexPageData{}
	p , err := scnr.GetParam().GetValue(defIndexName)

	if err == nil {
		// Parse json
		if err = json.Unmarshal([]byte(p), &d); err != nil {
			log.Error("Can't unmarshal json parameter " + defIndexName + ", error: " + err.Error())
		}
	} else {
		log.Info("Can't read parameter " + defIndexName + ", error: " + err.Error())
	}

	return &d
}




func ProcessCMD(scnr *options.ScannerObj) http.HandlerFunc {
	return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		log := scnr.GetLogger()
		
		p := chi.URLParam(r, "cmd")
		log.Debug(fmt.Sprintf("Uploaded value of %s : %v", "cmd",  p))

		respCmd := struct {
			CurrentCmd	string
		}{
			CurrentCmd: "Получена команда: " + p,
		}

		bAr, err := json.Marshal(respCmd)
		if err != nil {
			log.Error("Error when marshalling response CMD: "+ err.Error())
		}

		w.Header().Add("Content-Type", "text/html; charset=utf-8")
		w.WriteHeader(http.StatusOK)

		w.Write(bAr)
	})
}

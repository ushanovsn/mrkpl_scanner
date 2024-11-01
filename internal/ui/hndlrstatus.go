package ui

import (
	opt "mrkpl_scanner/internal/options"
	"net/http"
)


func StatusPage(scnr *opt.ScannerObj) http.HandlerFunc {
	return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		log := scnr.GetLogger()
		tmpl := scnr.GetUIObj().GetUIHTMLTemplates()

		// struct for current page
		data := struct {
			Title      string
			NaviMenu   []opt.NaviMenu
			ActiveMenu opt.NaviActiveMenu
		}{
			Title:    opt.DefUIPageTitle,
			NaviMenu: scnr.GetUIObj().GetUINaviMenu(),
			ActiveMenu: opt.NaviActiveMenu{
				ActiveTabVal:    "/status",
				ActiveDMenuVal:  "",
				PageDescription: "Статус процессов сервера",
			},
		}

		w.Header().Add("Content-Type", "text/html; charset=utf-8")
		w.WriteHeader(http.StatusOK)

		err := tmpl.ExecuteTemplate(w, "status", data)

		if err != nil {
			log.Error("Error executing status template: " + err.Error())
		}

	})
}

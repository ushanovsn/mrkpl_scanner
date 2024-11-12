package ui

import (
	opt "mrkpl_scanner/internal/options"
	"net/http"
)


func TemplatePage(scnr *opt.ScannerObj) http.HandlerFunc {
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
				ActiveTabVal:    "",
				ActiveDMenuVal:  "",
				PageDescription: "",
			},
		}

		w.Header().Add("Content-Type", "text/html; charset=utf-8")
		w.WriteHeader(http.StatusOK)

		err := tmpl.ExecuteTemplate(w, "template", data)

		if err != nil {
			log.Error("Error executing status template: " + err.Error())
		}

	})
}

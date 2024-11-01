package ui

import (
	"fmt"
	"mrkpl_scanner/internal/options"
	"net/http"
)


// Universal handler for web-interface
type UIHndlr struct {
	scnr *options.ScannerObj
	hnd func (*options.ScannerObj, http.ResponseWriter, *http.Request) (int, error)
}

// Interface http.Handler implementation
func (h UIHndlr) ServeHTTP(w http.ResponseWriter, r *http.Request) {
	log := h.scnr.GetLogger()

	// execute handler func
	status, err := h.hnd(h.scnr, w, r)

	log.Debug(fmt.Sprintf("UI processed: %s \"%s%s\" from %s. Result: %s", r.Method, r.Host, r.RequestURI, r.RemoteAddr, http.StatusText(status)))

	// check handler error
	if err != nil {
		log.Error(fmt.Sprintf("Error when processed web handler %s:%s. Error: %s", r.Method, r.RequestURI, err.Error()))
		switch status {
		case http.StatusNotFound:
			http.NotFound(w, r)
            // И если нам необходимо отобразить пользователю
            // вменяемую страницу ошибки, то мы можем использовать
            // наш контекст, например:
            // err := ah.renderTemplate(w, "http_404.tmpl", nil)
		case http.StatusInternalServerError:
			http.Error(w, http.StatusText(status), status)
		default:
			http.Error(w, http.StatusText(status), status)
		}
	}
}
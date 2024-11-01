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


// Parse file data from HTTP request
func getFileFromRequst(r *http.Request, searchName string, fType string) (content []byte, err error, msg string) {
	// Get handler for file, size and headers
	file, handler, err := r.FormFile(searchName)
	if err != nil {
		// if no file - nothing to do
		if err != http.ErrMissingFile {
			err = fmt.Errorf("Error when retrieving the file \"%s\": %s", searchName, err.Error())
			msg = "Ошибка чтения загруженного файла"
		}
		return
	}

	defer file.Close()

	// check content  type
	if v := handler.Header.Get("Content-Type"); fType != "" && v != fType {
		err = fmt.Errorf("Wrong content type (%s) for uploaded file, need \"%s\"", v, fType)
		msg = fmt.Sprintf("Некорректный тип файла, требуется \"%s\"", fType)
		return
	}

	// check size
	if handler.Size == 0 {
		err = fmt.Errorf("Uploaded file is empty")
		msg = "Загружен пустой файл"
		return
	}

	// Create a byte slice with the same size as the file
	content = make([]byte, handler.Size)
	// Read the file content into the byte slice
	if _, err = file.Read(content); err != nil {
		err = fmt.Errorf("Error while read file content: %s", err.Error())
		msg = "Ошибка чтения содержимого загруженного файла"
		return
	}

	return
}




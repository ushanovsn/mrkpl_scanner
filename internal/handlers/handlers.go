package handlers

import (
	"fmt"
	"html/template"
	opt "mrkpl_scanner/internal/options"
	"net/http"
)

// Processor for Index page "/"
func IndexPage(scnr *opt.ScannerObj) http.HandlerFunc {
	return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		log := scnr.GetLogger()
		tmpl := scnr.GetUIObj().GetUIHTMLTemplates()

		// struct for index page
		data := struct {
			Title                   string
			NaviMenu				[]opt.NaviMenu
			ActiveMenu				opt.NaviActiveMenu
			GoogleParamPage         string
			ParserParamPage         string
			ScannerParserStatusPage string
		}{
			Title:                   opt.DefUIPageTitle,
			NaviMenu: scnr.GetUIObj().GetUINaviMenu(),
			ActiveMenu: opt.NaviActiveMenu{
				ActiveTabVal: "/",
				ActiveDMenuVal: "",
				PageDescription: "Главная",
			},
			GoogleParamPage:         "gparams",
			ParserParamPage:         "pparams",
			ScannerParserStatusPage: "status",
		}

		// base headers
		header := http.StatusOK
		w.Header().Add("Content-Type", "text/html; charset=utf-8")


		// OK. Processing template - replace\substitute values in template and send to front
		w.WriteHeader(header)
		
		if err := tmpl.ExecuteTemplate(w, "index", data); err != nil {
			log.Error(fmt.Sprintf("Error executing template: %v\n", err.Error()))
		}

	})
}

// Processor for page with Google parameters
func GParamsPage(scnr *opt.ScannerObj) http.HandlerFunc {
	return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		log := scnr.GetLogger()

		// execute template and apply data
		if err := writeGParams(w, scnr, nil); err != nil {
			log.Error(err.Error())
		}
	})
}

// Processor for update Google parameters
func GParamsPageUpload(scnr *opt.ScannerObj) http.HandlerFunc {
	return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		log := scnr.GetLogger()
		// errors for show at page
		var errList []string

		// Maximum upload file size in bytes
		var maxSize int64 = 10240
		// Parse posted data
		if err := r.ParseMultipartForm(maxSize); err != nil {
			log.Error("Error when parsing form getting the file. Error: " + err.Error())
			errList = append(errList, "Ошибка обработки загруженных данных")
			err := writeGParams(w, scnr, errList)
			if err != nil {
				log.Error(err.Error())
			}
			return
		}

		// Get handler for filename, size and headers
		if file, handler, err := r.FormFile("file"); err != nil {
			// if no file - nothing to do
			if err != http.ErrMissingFile {
				log.Error("Error when retrieving the file. Error: " + err.Error())
				errList = append(errList, "Ошибка чтения загруженного файла аутентификации")
			} else {
				log.Warn("Uploaded file is empty")
			}
		} else {
			defer file.Close()

			log.Info("Page \"GParams\" uploaded the file: " + handler.Filename)
			log.Debug(fmt.Sprintf("Uploaded File Size: %v", handler.Size))
			log.Debug(fmt.Sprintf("MIME Header: %v", handler.Header))

			// check content fo JSON type
			for i, v := range handler.Header {
				if i == "Content-Type" && v[0] != "application/json" {
					log.Error("Wrong content type for uploaded file. Content: " + v[0])
					errList = append(errList, "Ошибка чтения загруженного файла аутентификации")
				} else if i == "Content-Type" && v[0] == "application/json" {
					// check size
					if handler.Size >= maxSize {
						log.Error(fmt.Sprintf("Can't upload the file more than %v Kb.", maxSize/1024))
						errList = append(errList, "Загружен некорректный формат файла аутентификации, файл должен быть JSON типа")
					} else {
						if handler.Size == 0 {
							log.Error("Uploaded file is empty")
							errList = append(errList, "Загружен пустой файл аутентификации")
						} else {
							// Create a byte slice with the same size as the file
							fileContent := make([]byte, handler.Size)
							// Read the file content into the byte slice
							if _, err = file.Read(fileContent); err != nil {
								log.Error(fmt.Sprintf("Error while read file content. Error: %s", err.Error()))
								errList = append(errList, "Ошибка чтения содержимого загруженного файла аутентификации")
							} else {
								if err := scnr.GetGDocSvc().SetAuthKeyFile(string(fileContent)); err != nil {
									log.Error(fmt.Sprintf("Error when set json content to GDoc object. Error: %s", err.Error()))
									errList = append(errList, "Ошибка применения аутентификационных данных загруженного файла")
								} else {
									// save data to param file
									err := scnr.GetParam().SetValue(opt.DefParamGAuthNameP, string(fileContent))
									if err != nil {
										log.Error(fmt.Sprintf("Error when save json content to parameters file. Error: %s", err.Error()))
										errList = append(errList, "Ошибка сохранения загруженного файла аутентификации Google-аккаунта в стстеме")
									}
								}
							}
						}
					}
				}
			}
		}

		// Get text value of URL
		val := r.FormValue("sheeturl")

		log.Debug("Page \"GParams\" uploaded the URL: " + val)

		if val == "" {
			log.Warn("Uploaded URL is empty")
		} else {
			// Set URL to GDoc object
			if err := scnr.GetGDocSvc().SetGSheetURL(val); err != nil {
				log.Error(fmt.Sprintf("Error when parse URL content to GDoc object. Error: %s", err.Error()))
				errList = append(errList, "В загруженных данных указан некорректный адрес Google-документа")
			} else {
				// save data to param file
				err := scnr.GetParam().SetValue(opt.DefParamGURLNameP, string(val))
				if err != nil {
					log.Error(fmt.Sprintf("Error when save URL content to parameters file. Error: %s", err.Error()))
					errList = append(errList, "Ошибка сохранения загруженного адреса Google-документа в стстеме")
				}
			}
		}

		// execute template and apply data
		if err := writeGParams(w, scnr, errList); err != nil {
			log.Error(err.Error())
		}
	})
}

// Main writer for http page with Google Params. Processing data and execute base template
func writeGParams(w http.ResponseWriter, scnr *opt.ScannerObj, errList []string) error {

	// prepare data values for page
	data := opt.GParamsPageData{
		Title:      opt.DefUIPageTitle,
		GPageURL:   scnr.GetGDocSvc().GetGSheetURL(),
		AuthClient: scnr.GetGDocSvc().GetCurClien(),
	}

	// check current values and set indicators
	if data.GPageURL != "" {
		data.GPageURLOkPref = "✔️"
	} else {
		data.GPageURLOkPref = "❌"
		data.ErrLog = append(data.ErrLog, "В системе не задан URL адрес Google-документа")
	}
	if data.AuthClient != "" {
		data.AuthClientOkPref = "✔️"
	} else {
		data.AuthClientOkPref = "❌"
		data.ErrLog = append(data.ErrLog, "В системе отсутствуют аутентификационные данные Google-аккаунта")
	}

	data.ErrLog = append(data.ErrLog, errList...)

	header := http.StatusOK
	w.Header().Add("Content-Type", "text/html; charset=utf-8")

	tmpl, err := template.ParseFiles(opt.DefUITemplatesPath + "/google_params_page.html")

	if err != nil {
		err = fmt.Errorf("Error while loading HTML template file: %v\n", err.Error())
		header = http.StatusInternalServerError
		w.WriteHeader(header)
	} else {
		w.WriteHeader(header)
		err = tmpl.Execute(w, data)
		if err != nil {
			err = fmt.Errorf("Error executing template: %v\n", err.Error())
		}
	}

	return err
}

func PParamsPage(scnr *opt.ScannerObj) http.HandlerFunc {
	return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		log := scnr.GetLogger()
		data := opt.PParamsPageData{
			Title: opt.DefUIPageTitle,
		}

		header := http.StatusOK
		w.Header().Add("Content-Type", "text/html; charset=utf-8")

		tmpl, err := template.ParseFiles(opt.DefUITemplatesPath + "/pparams_page.html")

		if err != nil {
			log.Error(fmt.Sprintf("Error while loading HTML template file: %v\n", err.Error()))
			header = http.StatusInternalServerError
			w.WriteHeader(header)
		} else {
			w.WriteHeader(header)
			err = tmpl.Execute(w, data)
			if err != nil {
				log.Error(fmt.Sprintf("Error executing template: %v\n", err.Error()))
			}
		}

	})
}

func StatusPage(scnr *opt.ScannerObj) http.HandlerFunc {
	return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		log := scnr.GetLogger()
		data := opt.PParamsPageData{
			Title: opt.DefUIPageTitle,
		}

		header := http.StatusOK
		w.Header().Add("Content-Type", "text/html; charset=utf-8")

		tmpl, err := template.ParseFiles(opt.DefUITemplatesPath + "/status.html")

		if err != nil {
			log.Error(fmt.Sprintf("Error while loading HTML template file: %v\n", err.Error()))
			header = http.StatusInternalServerError
			w.WriteHeader(header)
		} else {
			w.WriteHeader(header)
			err = tmpl.Execute(w, data)
			if err != nil {
				log.Error(fmt.Sprintf("Error executing template: %v\n", err.Error()))
			}
		}

	})
}

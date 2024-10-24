package handlers

import (
	"fmt"
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
			Title      string
			NaviMenu   []opt.NaviMenu
			ActiveMenu opt.NaviActiveMenu
		}{
			Title:    opt.DefUIPageTitle,
			NaviMenu: scnr.GetUIObj().GetUINaviMenu(),
			ActiveMenu: opt.NaviActiveMenu{
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

		if err := tmpl.ExecuteTemplate(w, "index", data); err != nil {
			log.Error(fmt.Sprintf("Error executing template: %v\n", err.Error()))
		}

	})
}

// Processor for page with Scan parameters
func ParamsScanPage(scnr *opt.ScannerObj) http.HandlerFunc {
	return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		log := scnr.GetLogger()

		if r.Method != "POST" && r.Method != "GET" {
			w.WriteHeader(http.StatusMethodNotAllowed)
			log.Warn(fmt.Sprintf("Wrong %s method for page \"task_param_scan\"", r.Method ))
			return
		}

		tmpl := scnr.GetUIObj().GetUIHTMLTemplates()

		// struct for Scanner Parameters page
		data := struct {
			Title          string
			NaviMenu       []opt.NaviMenu
			ActiveMenu     opt.NaviActiveMenu
			ParamsScanData opt.ParamsScanPageData
		}{
			Title:    opt.DefUIPageTitle,
			NaviMenu: scnr.GetUIObj().GetUINaviMenu(),
			ActiveMenu: opt.NaviActiveMenu{
				ActiveTabVal:    "/task_param",
				ActiveDMenuVal:  "/task_param_scan",
				PageDescription: "Параметры задачи \"Сканирование\"",
			},
			ParamsScanData: *scnr.GetUIObj().GetPageTaskScanData(scnr),
		}

		// process form from page
		if r.Method == "POST" {
			if err := procTaskParamScan(&data.ParamsScanData, r, scnr); err != nil {
				w.WriteHeader(http.StatusInternalServerError)
				log.Error(fmt.Sprintf("Error processing uploaded data from page \"task_param_scan\", error: %s", err.Error()))
				return
			}
		}

		w.Header().Add("Content-Type", "text/html; charset=utf-8")
	
		// execute template and apply data
		if err := tmpl.ExecuteTemplate(w, "task_param_scan", data); err != nil {
			err = fmt.Errorf("Error executing template: %v\n", err.Error())
		}
	
	})
}

// Processor for update Scan parameters
func ParamsScanPageUpload(scnr *opt.ScannerObj) http.HandlerFunc {
	return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		/*
		log := scnr.GetLogger()
		// errors for show at page
		var errList []string

		// Maximum upload file size in bytes
		var maxSize int64 = 10240
		// Parse posted data
		if err := r.ParseMultipartForm(maxSize); err != nil {
			log.Error("Error when parsing form getting the file. Error: " + err.Error())
			errList = append(errList, "Ошибка обработки загруженных данных")
			err := writeParamsScan(w, scnr, errList)
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

			log.Info("Page \"ParamsScan\" uploaded the file: " + handler.Filename)
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

		log.Debug("Page \"ParamsScan\" uploaded the URL: " + val)

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
		if err := writeParamsScan(w, scnr, errList); err != nil {
			log.Error(err.Error())
		}
*/
	})
}

/*
// Main writer for http page with Scan Params. Processing data and execute base template
func writeParamsScan(w http.ResponseWriter, scnr *opt.ScannerObj, errList []string) error {
	tmpl := scnr.GetUIObj().GetUIHTMLTemplates()

	// struct for Scanner Parameters page
	data := struct {
		Title          string
		NaviMenu       []opt.NaviMenu
		ActiveMenu     opt.NaviActiveMenu
		ParamsScanData opt.ParamsScanPageData
	}{
		Title:    opt.DefUIPageTitle,
		NaviMenu: scnr.GetUIObj().GetUINaviMenu(),
		ActiveMenu: opt.NaviActiveMenu{
			ActiveTabVal:    "/task_param",
			ActiveDMenuVal:  "/task_param_scan",
			PageDescription: "Параметры задачи \"Сканирование\"",
		},
		ParamsScanData: *scnr.GetUIObj().GetPageTaskScanData(scnr),
	}

	// check current values and set indicators
	if data.ParamsScanData.GPageURL != "" {
		data.ParamsScanData.GPageURLOkPref = "✔️"
	} else {
		data.ParamsScanData.GPageURLOkPref = "❌"
		data.ParamsScanData.ErrLog = append(data.ParamsScanData.ErrLog, "В системе не задан URL адрес Google-документа")
	}
	if data.ParamsScanData.AuthClient != "" {
		data.ParamsScanData.AuthClientOkPref = "✔️"
	} else {
		data.ParamsScanData.AuthClientOkPref = "❌"
		data.ParamsScanData.ErrLog = append(data.ParamsScanData.ErrLog, "В системе отсутствуют аутентификационные данные Google-аккаунта")
	}

	data.ParamsScanData.ErrLog = append(data.ParamsScanData.ErrLog, errList...)

	w.Header().Add("Content-Type", "text/html; charset=utf-8")
	w.WriteHeader(http.StatusOK)

	err := tmpl.ExecuteTemplate(w, "task_param_scan", data)

	if err != nil {
		err = fmt.Errorf("Error executing template: %v\n", err.Error())
	}

	return err
}
	*/

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

package handlers

import (
	"fmt"
	"html/template"
	opt "mrkpl_scanner/internal/options"
	"net/http"
	"regexp"
	"strconv"
)

// Processor for Index page "/"
func IndexPage(scnr *opt.ScannerObj) http.HandlerFunc {
	return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		log := scnr.GetLogger()

		// struct for index page
		data := struct {
			Title                   string
			GoogleParamPage         string
			ParserParamPage         string
			ScannerParserStatusPage string
		}{
			Title:                   scnr.GetUIObj().GetUITitle(),
			GoogleParamPage:         "gparams",
			ParserParamPage:         "pparams",
			ScannerParserStatusPage: "status",
		}

		// base headers
		header := http.StatusOK
		w.Header().Add("Content-Type", "text/html; charset=utf-8")

		// load template file
		tmpl, err := template.ParseFiles(scnr.GetUIObj().GetUITemplPath() + "/index.html")

		if err != nil {
			log.Error(fmt.Sprintf("Error while loading HTML template file: %v\n", err.Error()))
			header = http.StatusInternalServerError
			w.WriteHeader(header)
		} else {
			// OK. Processing template - replace\substitute values in template and send to front
			w.WriteHeader(header)
			err = tmpl.Execute(w, data)
			if err != nil {
				log.Error(fmt.Sprintf("Error executing template: %v\n", err.Error()))
			}
		}
	})
}

// Processor for page with Google parameters
func GParamsPage(scnr *opt.ScannerObj) http.HandlerFunc {
	return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		log := scnr.GetLogger()

		// prepare data values for page
		data := opt.GParamsPageData{
			Title:      scnr.GetUIObj().GetUITitle(),
			GPageURL:   scnr.GetGDocSvc().GetGSheetURL(),
			AuthClient: scnr.GetGDocSvc().GetCurClien(),
		}

		// check incoming URL parameters
		updUrl, _ := strconv.ParseBool(r.URL.Query().Get("u"))
		updFile, _ := strconv.ParseBool(r.URL.Query().Get("f"))
		referAdr := r.Referer()
		isUploading, _ := regexp.Match("^*/gparams(\\?|\\z)+", []byte(referAdr))

		// When load parameters and redirect back to here
		if isUploading {
			// mark upload new params result
			if updUrl {
				data.GPageURLOkPref = "✔️"
			} else {
				data.GPageURLOkPref = "❌"
			}

			if updFile {
				data.AuthClientOkPref = "✔️"
			} else {
				data.AuthClientOkPref = "❌"
			}

		} else {
			// check current values
			if data.GPageURL != "" {
				data.GPageURLOkPref = "✔️"
			} else {
				data.GPageURLOkPref = "❌"
			}

			if data.AuthClient != "" {
				data.AuthClientOkPref = "✔️"
			} else {
				data.AuthClientOkPref = "❌"
			}
		}

		header := http.StatusOK
		w.Header().Add("Content-Type", "text/html; charset=utf-8")

		tmpl, err := template.ParseFiles(scnr.GetUIObj().GetUITemplPath() + "/google_params_page.html")

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

// Processor for update Google parameters
func GParamsPageUpload(scnr *opt.ScannerObj) http.HandlerFunc {
	return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		log := scnr.GetLogger()
		// flag for file processing
		var fFileRes bool = false
		// flag for URL processing
		var fUrlRes bool = false

		// Maximum upload file size in bytes
		var maxSize int64 = 10240
		if err:= r.ParseMultipartForm(maxSize); err != nil {
			log.Error("Error when parsing form getting the file. Error: " + err.Error())
		}

		// Get handler for filename, size and headers
		file, handler, err := r.FormFile("file")
		if err != nil {
			log.Error("Error when retrieving the file. Error: " + err.Error())
		} else {
			defer file.Close()

			log.Info("Page \"GParams\" uploaded the file: " + handler.Filename)
			log.Debug(fmt.Sprintf("Uploaded File Size: %v", handler.Size))
			log.Debug(fmt.Sprintf("MIME Header: %v", handler.Header))

			// check content fo JSON type
			for i, v := range handler.Header {
				if i == "Content-Type" && v[0] != "application/json" {
					log.Error("Wrong content type for uploaded file. Content: " + v[0])
				} else if i == "Content-Type" && v[0] == "application/json" {
					// check size
					if handler.Size >= maxSize {
						log.Error(fmt.Sprintf("Can't upload the file more than %v Kb.", maxSize/1024))
					} else {
						if handler.Size == 0 {
							log.Error("Uploaded file is empty")
						} else {
							// Create a byte slice with the same size as the file
							fileContent := make([]byte, handler.Size)
							// Read the file content into the byte slice
							if _, err = file.Read(fileContent); err != nil {
								log.Error(fmt.Sprintf("Error while read file content. Error: %s", err.Error()))
							} else {
								if err := scnr.GetGDocSvc().SetAuthKeyFile(string(fileContent)); err != nil {
									log.Error(fmt.Sprintf("Error when set json content to GDoc object. Error: %s", err.Error()))
								} else {
									// save data to param file
									err := scnr.GetParam().SetValue(opt.DefParamGAuthNameP, string(fileContent))
									if err != nil {
										log.Error(fmt.Sprintf("Error when save json content to parameters file. Error: %s", err.Error()))
									} else {
										// all is ok
										fFileRes = true
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
		if val == "" {
			log.Error("Uploaded URL is empty")
		}

		log.Debug("Page \"GParams\" uploaded the URL: " + val)

		// Set URL to GDoc object
		if err := scnr.GetGDocSvc().SetGSheetURL(val); err != nil {
			log.Error(fmt.Sprintf("Error when parse URL content to GDoc object. Error: %s", err.Error()))
		} else {
			// save data to param file
			err := scnr.GetParam().SetValue(opt.DefParamGURLNameP, string(val))
			if err != nil {
				log.Error(fmt.Sprintf("Error when save URL content to parameters file. Error: %s", err.Error()))
			} else {
				// all is ok
				fUrlRes = true
			}
		}

		http.Redirect(w, r, fmt.Sprintf("/gparams?f=%v&u=%v", fFileRes, fUrlRes), http.StatusMovedPermanently)
	})
}

func PParamsPage(scnr *opt.ScannerObj) http.HandlerFunc {
	return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		log := scnr.GetLogger()
		data := opt.PParamsPageData{
			Title: scnr.GetUIObj().GetUITitle(),
		}

		header := http.StatusOK
		w.Header().Add("Content-Type", "text/html; charset=utf-8")

		tmpl, err := template.ParseFiles(scnr.GetUIObj().GetUITemplPath() + "/pparams_page.html")

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
			Title: scnr.GetUIObj().GetUITitle(),
		}

		header := http.StatusOK
		w.Header().Add("Content-Type", "text/html; charset=utf-8")

		tmpl, err := template.ParseFiles(scnr.GetUIObj().GetUITemplPath() + "/status.html")

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

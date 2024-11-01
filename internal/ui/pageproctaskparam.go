package ui

import (
	"encoding/json"
	"fmt"
	"mrkpl_scanner/internal/options"
	"net/http"
	"regexp"
	"strconv"
	"strings"
)

const (
	defParamScanName = "ParamsScanPageData"
)

// Processor for page-task_param_scan.gohtml uploads
func procTaskParamScan(d *options.ParamsScanPageData, r *http.Request, scnr *options.ScannerObj) (error) {
	log := scnr.GetLogger()
	
	var (
		err			error
		// processing parameter name on Form
		pName 		string
		// parsed value from Form
		val 		string
		// errors text for Source data
		errListSrc 	[]string
	)

	const (
		// Maximum upload file size in bytes
		maxSize 	int64	= 10240
		// Checkbox value when checked
		chbChecked 	string	= "on"
	)

	// Parse posted data
	if err = r.ParseMultipartForm(maxSize); err != nil {
		return fmt.Errorf("Error when parsing page form. Error: %s", err.Error())
	}

	// **************** Checkbox Equal Source and Save ****************

	// Get value of Saving parsed data in source
	pName = "sourceequal"
	val = r.FormValue(pName)
	log.Debug(fmt.Sprintf("Uploaded value of %s : %v", pName,  val))
	// set parsed value
	d.SaveInOne = (val == chbChecked)

	// **************** CLOUD TABLE TAB (GOOGLE) ****************
	
	// Get text value of selected tab
	pName = "source_tabs"
	val = r.FormValue(pName)
	log.Debug(fmt.Sprintf("Uploaded value of %s : %v", pName,  val))

	// current params tab is GOOGLE TAB
	if val == "source_tab_google" {
		//set paкsed value
		d.TabSrcActive = val

		// Get text value of document URL value
		pName = "source_gsheeturl"
		val = r.FormValue(pName)
		log.Debug(fmt.Sprintf("Uploaded value of %s : %v", pName,  val))
		// Set URL to GDoc object
		if err = scnr.GetGDocSvc().SetGSheetURL(val); err != nil {
			log.Error(fmt.Sprintf("Error when parse URL content to GDoc object. Error: %s", err.Error()))
			errListSrc = append(errListSrc, "Указан некорректный адрес Online-документа")
		} else {
			// set parsed value
			d.GParamSrc.GAuth.GPageURL = val
			d.GParamSrc.GAuth.GPageURLOk = true
		}

		// Get handler for file, size and headers
		if file, handler, err := r.FormFile("filename_google_source"); err != nil {
			// if no file - nothing to do
			if err != http.ErrMissingFile {
				log.Error("Error when retrieving the file \"filename_google_source\". Error: " + err.Error())
				errListSrc = append(errListSrc, "Ошибка чтения загруженного файла аутентификации")
			} else {
				// when file parameters already applied before
				if  !(d.GParamSrc.GAuth.AuthClientOk) {
					log.Warn("Uploaded file \"filename_google_source\" is empty")
					errListSrc = append(errListSrc, "Файл аутентификации не был загружен")
				}
			}
		} else {
			defer file.Close()

			log.Info(fmt.Sprintf("File \"filename_google_source\" uploaded: %s, size: %v, MIME Header: %v", handler.Filename, handler.Size, handler.Header))

			// check content for JSON type
			for i, v := range handler.Header {
				if i == "Content-Type" && v[0] != "application/json" {
					log.Error("Wrong content type for uploaded file. Content: " + v[0])
					errListSrc = append(errListSrc, "Некорректный тип файла аутентификации, требуется JSON файл")
				} else if i == "Content-Type" && v[0] == "application/json" {
					// check size
					if handler.Size == 0 {
						log.Error("Uploaded file is empty")
						errListSrc = append(errListSrc, "Загружен пустой файл аутентификации")
					} else {
						// Create a byte slice with the same size as the file
						fileContent := make([]byte, handler.Size)
						// Read the file content into the byte slice
						if _, err = file.Read(fileContent); err != nil {
							log.Error(fmt.Sprintf("Error while read file content. Error: %s", err.Error()))
							errListSrc = append(errListSrc, "Ошибка чтения содержимого загруженного файла аутентификации")
						} else {
							if err = scnr.GetGDocSvc().SetAuthKeyFile(string(fileContent)); err != nil {
								log.Error(fmt.Sprintf("Error when set json content to GDoc object. Error: %s", err.Error()))
								errListSrc = append(errListSrc, "Ошибка применения аутентификационных данных загруженного файла")
							} else {
								d.GParamSrc.GAuth.AuthClient = scnr.GetGDocSvc().GetCurClient()
								// save data to param file
								if err = scnr.GetParam().SetValue(options.DefParamGAuthNameP, string(fileContent)); err != nil {
									log.Error(fmt.Sprintf("Error when save json content to parameters file. Error: %s", err.Error()))
									errListSrc = append(errListSrc, "Ошибка сохранения загруженного файла аутентификации Google-аккаунта в системе")
								} else {
									d.GParamSrc.GAuth.AuthClientOk = true
								}
							}
						}
					}
				}
			}
		}

		// Get text value of start row
		pName = "source_g_startrow"
		val = r.FormValue(pName)
		log.Debug(fmt.Sprintf("Uploaded value of %s : %v", pName,  val))
		// Check value
		if num, err := strconv.Atoi(val); err != nil || num < 1 {
			if err != nil {
				log.Error(fmt.Sprintf("Error when parse start row value. Error: %s", err.Error()))
			} else if num < 1 {
				log.Error(fmt.Sprintf("Incorrect start row value: %v", num))
			}
			errListSrc = append(errListSrc, "Указан некорректный номер строки начала данных для парсинга")
		} else {
			// set parsed value
			d.GParamSrc.TblParam.StartRowNum = num
			d.GParamSrc.TblParam.StartRowNumOk = true
		}

		// Get text value of column with data for parsing
		pName = "source_g_coldata"
		val = r.FormValue(pName)
		log.Debug(fmt.Sprintf("Uploaded value of %s : %v", pName,  val))
		// Check value
		if ok, _ := regexp.MatchString(`\A([A-Z]){1,2}\z`, strings.ToUpper(val)); !ok {
			log.Error(fmt.Sprintf("Wrong source column \"%s\" value: %s", pName, val))
			errListSrc = append(errListSrc, "Указан некорректный столбец с данными для парсинга")
		} else {
			// set parsed value
			d.GParamSrc.TblParam.ColSourceId = strings.ToUpper(val)
			d.GParamSrc.TblParam.ColSourceIdOk = true
		}

		// saving parsing into source
		if d.SaveInOne {
			// Get text value of column with data for parsing
			for i := range d.GParamSrc.TblParam.Params {
				// column name
				pName = fmt.Sprintf("source_g_param%v", i)
				val = r.FormValue(pName)
				log.Debug(fmt.Sprintf("Uploaded value of %s : %v", pName,  val))
				// Check value
				if ok, _ := regexp.MatchString(`\A([A-Z]){0,2}\z`, strings.ToUpper(val)); !ok {
					log.Error(fmt.Sprintf("Wrong parameter \"%s\" column value: %s", pName, val))
					errListSrc = append(errListSrc, fmt.Sprintf("Указан некорректный столбец с параметрами: %v", i))
				} else {
					// set parsed value
					d.GParamSrc.TblParam.Params[i].ColParamValue = strings.ToUpper(val)
				}

				// data type from selector
				if val == "" {
					// reset type if no column
					d.GParamSrc.TblParam.Params[i].ColParamType = val
				} else {
					pName = fmt.Sprintf("source_g_paramtype%v", i)
					val = r.FormValue(pName)
					log.Debug(fmt.Sprintf("Uploaded value of %s : %v", pName,  val))
					d.GParamSrc.TblParam.Params[i].ColParamType = val
				}
			}
		}
	}
	
	// set errors list for source
	d.GParamSrc.ErrLog = errListSrc

	//																										//
	//												 SAVE PARAMETERS										//
	//																										//
	if !d.SaveInOne {
		// **************** CLOUD TABLE TAB (GOOGLE) ****************
		
		// Get text value of selected tab
		pName = "save_data"
		val = r.FormValue(pName)
		log.Debug(fmt.Sprintf("Uploaded value of %s : %v", pName,  val))

		// current params tab is GOOGLE TAB
		if val == "save_tab_google" {
			//set paкsed value
			d.TabSvActive = val

			// Get text value of document URL value
			pName = "save_gsheeturl"
			val = r.FormValue(pName)
			log.Debug(fmt.Sprintf("Uploaded value of %s : %v", pName,  val))
			// Set URL to GDoc object
			
			// need second document service!!!!

			/*
			if err = scnr.GetGDocSvc().SetGSheetURL(val); err != nil {
				log.Error(fmt.Sprintf("Error when parse URL content to GDoc object. Error: %s", err.Error()))
				errListSrc = append(errListSrc, "Указан некорректный адрес Online-документа")
			} else {
				// set parsed value
				d.GParamSv.GAuth.GPageURL = val
				d.GParamSv.GAuth.GPageURLOk = true
			}
			*/

			// Get handler for file, size and headers
			if file, handler, err := r.FormFile("filename_google_save"); err != nil {
				// if no file - nothing to do
				if err != http.ErrMissingFile {
					log.Error("Error when retrieving the file \"filename_google_save\". Error: " + err.Error())
					errListSrc = append(errListSrc, "Ошибка чтения загруженного файла аутентификации")
				} else {
					// when file parameters already applied before
					if  !(d.GParamSv.GAuth.AuthClientOk) {
						log.Warn("Uploaded file \"filename_google_save\" is empty")
						errListSrc = append(errListSrc, "Файл аутентификации не был загружен")
					}
				}
			} else {
				defer file.Close()

				log.Info(fmt.Sprintf("File \"filename_google_save\" uploaded: %s, size: %v, MIME Header: %v", handler.Filename, handler.Size, handler.Header))

				// check content for JSON type
				for i, v := range handler.Header {
					if i == "Content-Type" && v[0] != "application/json" {
						log.Error("Wrong content type for uploaded file. Content: " + v[0])
						errListSrc = append(errListSrc, "Некорректный тип файла аутентификации, требуется JSON файл")
					} else if i == "Content-Type" && v[0] == "application/json" {
						// check size
						if handler.Size == 0 {
							log.Error("Uploaded file is empty")
							errListSrc = append(errListSrc, "Загружен пустой файл аутентификации")
						} else {
							// Create a byte slice with the same size as the file
							fileContent := make([]byte, handler.Size)
							// Read the file content into the byte slice
							if _, err = file.Read(fileContent); err != nil {
								log.Error(fmt.Sprintf("Error while read file content. Error: %s", err.Error()))
								errListSrc = append(errListSrc, "Ошибка чтения содержимого загруженного файла аутентификации")
							} else {
								if err = scnr.GetGDocSvc().SetAuthKeyFile(string(fileContent)); err != nil {
									log.Error(fmt.Sprintf("Error when set json content to GDoc object. Error: %s", err.Error()))
									errListSrc = append(errListSrc, "Ошибка применения аутентификационных данных загруженного файла")
								} else {

									// need second document service!!!!

									/*
									d.GParamSv.GAuth.AuthClient = scnr.GetGDocSvc().GetCurClient()
									// save data to param file
									if err = scnr.GetParam().SetValue(options.DefParamGAuthNameP, string(fileContent)); err != nil {
										log.Error(fmt.Sprintf("Error when save json content to parameters file. Error: %s", err.Error()))
										errListSrc = append(errListSrc, "Ошибка сохранения загруженного файла аутентификации Google-аккаунта в системе")
									} else {
										d.GParamSv.GAuth.AuthClientOk = true
									}
									*/
								}
							}
						}
					}
				}
			}

			// Get text value of start row
			pName = "save_g_startrow"
			val = r.FormValue(pName)
			log.Debug(fmt.Sprintf("Uploaded value of %s : %v", pName,  val))
			// Check value
			if num, err := strconv.Atoi(val); err != nil || num < 1 {
				if err != nil {
					log.Error(fmt.Sprintf("Error when parse start row value. Error: %s", err.Error()))
				} else if num < 1 {
					log.Error(fmt.Sprintf("Incorrect start row value: %v", num))
				}
				errListSrc = append(errListSrc, "Указан некорректный номер строки начала данных для парсинга")
			} else {
				// set parsed value
				d.GParamSv.TblParam.StartRowNum = num
				d.GParamSv.TblParam.StartRowNumOk = true
			}

			// Get text value of column with data for parsing
			for i := range d.GParamSv.TblParam.Params {
				// column name
				pName = fmt.Sprintf("save_g_param%v", i)
				val = r.FormValue(pName)
				log.Debug(fmt.Sprintf("Uploaded value of %s : %v", pName,  val))
				// Check value
				if ok, _ := regexp.MatchString(`\A([A-Z]){0,2}\z`, strings.ToUpper(val)); !ok {
					log.Error(fmt.Sprintf("Wrong parameter \"%s\" column value: %s", pName, val))
					errListSrc = append(errListSrc, fmt.Sprintf("Указан некорректный столбец с параметрами: %v", i))
				} else {
					// set parsed value
					d.GParamSv.TblParam.Params[i].ColParamValue = strings.ToUpper(val)
				}

				// data type from selector
				if val == "" {
					// reset type if no column
					d.GParamSv.TblParam.Params[i].ColParamType = val
				} else {
					pName = fmt.Sprintf("save_g_paramtype%v", i)
					val = r.FormValue(pName)
					log.Debug(fmt.Sprintf("Uploaded value of %s : %v", pName,  val))
					d.GParamSv.TblParam.Params[i].ColParamType = val
				}
			}

			
			// set errors list for save
			d.GParamSv.ErrLog = errListSrc
		}
		






	}








	// convert parameters data to json string
	jsonString, err := json.Marshal(d)
	if err != nil {
		return err
	}

	//save data to PARAMETERS
	return scnr.GetParam().SetValue(defParamScanName, string(jsonString))
}



// Aggregate and return actual values for page "page-task_param_scan"
func getPageTaskScanData(scnr *options.ScannerObj) *options.ParamsScanPageData {
	log := scnr.GetLogger()

	// init empty data
	d := options.ParamsScanPageData{}
	d.GParamSrc.TblParam.Params = make([]options.ColParams, 5)
	d.FileSrc.TblParam.Params = make([]options.ColParams, 5)
	d.DBSrc.TblParam.Params = make([]options.ColParams, 5)
	d.GParamSv.TblParam.Params = make([]options.ColParams, 5)
	d.FileSv.TblParam.Params = make([]options.ColParams, 5)
	d.DBSv.TblParam.Params = make([]options.ColParams, 5)

	p , err := scnr.GetParam().GetValue(defParamScanName)

	if err == nil {
		// Parse json
		if err = json.Unmarshal([]byte(p), &d); err != nil {
			log.Error("Can't unmarshal json parameter " + defParamScanName + ", error: " + err.Error())
		}
	} else {
		log.Info("Can't read parameter " + defParamScanName + ", error: " + err.Error())
	}

	return &d
}
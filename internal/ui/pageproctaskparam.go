package ui

import (
	"encoding/json"
	"fmt"
	"mrkpl_scanner/internal/options"
	"net/http"
	"regexp"
	"strconv"
	"strings"
	"time"
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
		errList 	[]string
		// error string
		errStr 		string
		// file data
		fContent	[]byte
	)

	const (
		// Maximum upload size in bytes
		maxSize 	int64	= 10240
		// Checkbox value when checked
		chbChecked 	string	= "on"
	)

	// Parse posted data
	if err = r.ParseMultipartForm(maxSize); err != nil {
		return fmt.Errorf("Error when parsing page form. Error: %s", err.Error())
	}

	// for column name checking
	colNameReg := regexp.MustCompile(`\A([A-Z]){1,2}\z`)


								//////////////////////////////////////////////////////////////////////////////
								//								SOURCE PART									//
								//////////////////////////////////////////////////////////////////////////////

	// *** Get value of Saving parsed data in source
	pName = "sourceequal"
	val = r.FormValue(pName)
	log.Debug(fmt.Sprintf("Uploaded value of %s : %v", pName,  val))
	// set parsed value
	d.SaveInOne = (val == chbChecked)

	// *** Get text value of selected tab
	pName = "source_tabs"
	val = r.FormValue(pName)
	log.Debug(fmt.Sprintf("Uploaded value of %s : %v", pName,  val))

	////////////////////////////////////////////////////////
	//             CLOUD TABLE TAB (GOOGLE)               //
	////////////////////////////////////////////////////////

	if val == "source_tab_google" {
		errList = make([]string, 0)
		// set parsed value of active tab
		d.TabSrcActive = val

		// *** Get text value of document URL value
		pName = "source_gsheeturl"
		val = r.FormValue(pName)
		log.Debug(fmt.Sprintf("Uploaded value of %s : %v", pName,  val))
		// Set URL to GDoc object
		if val == "" {
			errList = append(errList, "Необходимо указать адрес Online-документа")
		}else if !scnr.GetCloudDocBaseSvc().CheckSheetURL(val) {
			log.Warn("Error when parse URL content to CloudDoc object")
			errList = append(errList, "Указан некорректный адрес Online-документа")
		} else {
			// set parsed value
			d.GParamSrc.GAuth.GPageURL = val
			d.GParamSrc.GAuth.GPageURLOk = true
		}

		// *** get file content from form
		fContent, err, errStr = getFileFromRequst(r, "filename_google_source", "application/json")
		if err != nil {
			if err != http.ErrMissingFile {
				log.Warn(err.Error())
				errList = append(errList, ("Файл аутентификации: " + errStr))
			}
			// error when auth file is not exist and not received now
			if err == http.ErrMissingFile && !(d.GParamSrc.GAuth.AuthClientOk) {
				log.Warn("Uploaded file \"filename_google_source\" is empty")
				errList = append(errList, "Файл аутентификации не был загружен")
			}
		} else {
			// apply and save file content
			if user, err := scnr.GetCloudDocBaseSvc().CheckAuthKeyFile(string(fContent)); err != nil {
				log.Error(fmt.Sprintf("Error when check structure json auth file. Error: %s", err.Error()))
				errList = append(errList, "Ошибка проверки аутентификационных данных загруженного файла")
			} else {
				d.GParamSrc.GAuth.AuthClient = user
				// save data to param file
				if err = scnr.GetParam().SetValue(options.DefParamGAuthNameP, string(fContent)); err != nil {
					log.Error(fmt.Sprintf("Error when save json content to parameters file. Error: %s", err.Error()))
					errList = append(errList, "Ошибка сохранения загруженного файла аутентификации аккаунта в системе")
				} else {
					d.GParamSrc.GAuth.AuthClientOk = true
				}
			}
		}

		// *** Get text value of start row
		pName = "source_g_startrow"
		val = r.FormValue(pName)
		log.Debug(fmt.Sprintf("Uploaded value of %s : %v", pName,  val))
		// Check value
		if num, err := strconv.Atoi(val); err != nil || num < 1 {
			if err != nil {
				log.Warn(fmt.Sprintf("Error when parse start row value. Error: %s", err.Error()))
			} else {
				log.Warn(fmt.Sprintf("Incorrect start row value: %v", num))
			}
			errList = append(errList, "Указан некорректный номер строки начала данных для парсинга")
		} else {
			// set parsed value
			d.GParamSrc.TblParam.StartRowNum = num
			d.GParamSrc.TblParam.StartRowNumOk = true
		}

		// *** Get text value of column with data for parsing
		pName = "source_g_coldata"
		val = r.FormValue(pName)
		log.Debug(fmt.Sprintf("Uploaded value of %s : %v", pName,  val))
		// Check value
		if !colNameReg.MatchString(strings.ToUpper(val)) {
			log.Warn(fmt.Sprintf("Wrong source column \"%s\" value: %s", pName, val))
			errList = append(errList, "Указан некорректный столбец с данными для парсинга")
		} else {
			// set parsed value
			d.GParamSrc.TblParam.ColSourceId = strings.ToUpper(val)
			d.GParamSrc.TblParam.ColSourceIdOk = true
		}

		// saving parsing into source
		if d.SaveInOne {
			// *** Get text value of column with data for parsing
			for i := range d.GParamSrc.TblParam.Params {
				// column name
				pName = fmt.Sprintf("source_g_param%v", i)
				val = r.FormValue(pName)
				log.Debug(fmt.Sprintf("Uploaded value of %s : %v", pName,  val))
				// Check value
				if val != "" && !colNameReg.MatchString(strings.ToUpper(val)) {
					log.Warn(fmt.Sprintf("Wrong parameter \"%s\" column value: %s", pName, val))
					errList = append(errList, fmt.Sprintf("Указан некорректный столбец с параметрами: %v", i))
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
		// set errors list for tab
		d.GParamSrc.ErrLog = errList
	}

								//////////////////////////////////////////////////////////////////////////////
								//								SAVE PART									//
								//////////////////////////////////////////////////////////////////////////////

	if !d.SaveInOne {
		// *** Get text value of selected tab
		pName = "save_tabs"
		val = r.FormValue(pName)
		log.Debug(fmt.Sprintf("Uploaded value of %s : %v", pName,  val))

		////////////////////////////////////////////////////////
		//             CLOUD TABLE TAB (GOOGLE)               //
		////////////////////////////////////////////////////////

		// current params tab is GOOGLE TAB
		if val == "save_tab_google" {
			errList = make([]string, 0)
			//set paкsed value
			d.TabSvActive = val

			// *** Get text value of document URL value
			pName = "save_gsheeturl"
			val = r.FormValue(pName)
			log.Debug(fmt.Sprintf("Uploaded value of %s : %v", pName,  val))
			// Set URL to GDoc object
			if val == "" {
				errList = append(errList, "Необходимо указать адрес Online-документа")
			}else if !scnr.GetCloudDocWriteSvc().CheckSheetURL(val) {
				log.Warn("Error when parse URL content to CloudDoc object")
				errList = append(errList, "Указан некорректный адрес Online-документа")
			} else {
				// set parsed value
				d.GParamSv.GAuth.GPageURL = val
				d.GParamSv.GAuth.GPageURLOk = true
			}


			// *** get file content from form
			fContent, err, errStr = getFileFromRequst(r, "filename_google_save", "application/json")
			if err != nil {
				if err != http.ErrMissingFile {
					log.Warn(err.Error())
					errList = append(errList, ("Файл аутентификации: " + errStr))
				}
				// error when auth file is not exist and not received now
				if err == http.ErrMissingFile && !(d.GParamSv.GAuth.AuthClientOk) {
					log.Warn("Uploaded file \"filename_google_save\" is empty")
					errList = append(errList, "Файл аутентификации не был загружен")
				}
			} else {
				// apply and save file content
				if user, err := scnr.GetCloudDocWriteSvc().CheckAuthKeyFile(string(fContent)); err != nil {
					log.Error(fmt.Sprintf("Error when check structure json auth file. Error: %s", err.Error()))
					errList = append(errList, "Ошибка проверки аутентификационных данных загруженного файла")
				} else {
					d.GParamSv.GAuth.AuthClient = user
					// save data to param file
					if err = scnr.GetParam().SetValue(options.DefParamGAuthNameP, string(fContent)); err != nil {
						log.Error(fmt.Sprintf("Error when save json content to parameters file. Error: %s", err.Error()))
						errList = append(errList, "Ошибка сохранения загруженного файла аутентификации аккаунта в системе")
					} else {
						d.GParamSv.GAuth.AuthClientOk = true
					}
				}
			}

			// *** Get text value of start row
			pName = "save_g_startrow"
			val = r.FormValue(pName)
			log.Debug(fmt.Sprintf("Uploaded value of %s : %v", pName,  val))
			// Check value
			if num, err := strconv.Atoi(val); err != nil || num < 1 {
				if err != nil {
					log.Warn(fmt.Sprintf("Error when parse start row value. Error: %s", err.Error()))
				} else {
					log.Warn(fmt.Sprintf("Incorrect start row value: %v", num))
				}
				errList = append(errList, "Указан некорректный номер строки начала данных")
			} else {
				// set parsed value
				d.GParamSv.TblParam.StartRowNum = num
				d.GParamSv.TblParam.StartRowNumOk = true
			}

			// *** Get text value of column with data for parsing
			for i := range d.GParamSv.TblParam.Params {
				// column name
				pName = fmt.Sprintf("save_g_param%v", i)
				val = r.FormValue(pName)
				log.Debug(fmt.Sprintf("Uploaded value of %s : %v", pName,  val))
				// Check value
				if val != "" && !colNameReg.MatchString(strings.ToUpper(val)) {
					log.Warn(fmt.Sprintf("Wrong parameter \"%s\" column value: %s", pName, val))
					errList = append(errList, fmt.Sprintf("Указан некорректный столбец с параметрами: %v", i))
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
			// set errors list for tab
			d.GParamSv.ErrLog = errList
		}
	}

	//////////////////////////////////////////////////////////////////////////////
	//								COMMON PART									//
	//////////////////////////////////////////////////////////////////////////////

	// *** Get text value of selected radiobutton
	pName = "start_type"
	val = r.FormValue(pName)
	log.Debug(fmt.Sprintf("Uploaded value of %s : %v", pName,  val))
	errList = make([]string, 0)
	d.Starting.SelectName = val

	// current start parameter is Start At Time
	if val == "start_at_time" {
		// *** Get text value of start time
		pName = "start_at_time_value"
		val = r.FormValue(pName)
		log.Debug(fmt.Sprintf("Uploaded value of %s : %v", pName,  val))
		// Set time param
		if val == "" {
			d.Starting.AtTimeOk = false
			errList = append(errList, "Необходимо указать время запуска сканера")
		}else if _, err := time.Parse("15:04", val); err != nil {
			d.Starting.AtTimeOk = false
			log.Warn(fmt.Sprintf("Error when parse start time value (%s). Error: %s", val, err.Error()))
			errList = append(errList, "Указано некорректное время запуска сканера")
		} else {
			// set parsed value
			d.Starting.AtTime = val
			d.Starting.AtTimeOk = true
		}
	} else if val == "start_period" {
		// current start parameter is Start With Period
	
		// *** Get text value of Period Duration
		pName = "start_period_value"
		val = r.FormValue(pName)
		log.Debug(fmt.Sprintf("Uploaded value of %s : %v", pName,  val))
		// Set time param
		if val == "" {
			d.Starting.PeriodicOk = false
			errList = append(errList, "Необходимо указать период запуска сканера")
		}else if v, err := strconv.Atoi(val); err != nil || v < 1 {
			if err != nil {
				log.Warn(fmt.Sprintf("Error when parse period time value (%s). Error: %s", val, err.Error()))
			} else {
				log.Warn(fmt.Sprintf("Incorrect period time value: %v", v))
			}
			d.Starting.PeriodicOk = false
			errList = append(errList, "Указано некорректное значение периода запуска сканера")
		} else {
			// set parsed value
			d.Starting.Periodic = v
			d.Starting.PeriodicOk = true
		}
	}

	d.Starting.ErrLog = errList




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
package ui

import (
	"encoding/json"
	"mrkpl_scanner/internal/options"
	"net/http"
	"fmt"
	"strconv"
)

const (
	defConfWBName = "ConfigWBPageData"
)

// Processor for page-task_param_scan.gohtml uploads
func procConfWB(d *options.ConfWBPageData, r *http.Request, scnr *options.ScannerObj) (error) {
	log := scnr.GetLogger()

	var (
		err			error
		// processing parameter name on Form
		pName 		string
		// parsed value from Form
		val 		string
		// errors text for Source data
		errList 	[]string
	)

	const (
		// Maximum upload size in bytes
		maxSize 	int64	= 10240
	)

	// Parse posted data
	if err = r.ParseMultipartForm(maxSize); err != nil {
		return fmt.Errorf("Error when parsing page form. Error: %s", err.Error())
	}

	// *** Get value of query delay
	pName = "query_delay"
	val = r.FormValue(pName)
	log.Debug(fmt.Sprintf("Uploaded value of %s : %v", pName,  val))
	// Check value
	if num, err := strconv.Atoi(val); err != nil {
		log.Warn(fmt.Sprintf("Error when parse query delay value. Error: %s", err.Error()))
		errList = append(errList, "Указано некорректное время задержки между запросами")
	} else {
		// set parsed value
		d.RequestDelay = num
		d.RequestDelayOk = true
	}

	// *** Get value of address identification type
	pName = "address_ident"
	val = r.FormValue(pName)
	log.Debug(fmt.Sprintf("Uploaded value of %s : %v", pName,  val))
	d.AddressType = val

	// *** Get value of discount type
	pName = "discount_calc"
	val = r.FormValue(pName)
	log.Debug(fmt.Sprintf("Uploaded value of %s : %v", pName,  val))

	if val == "discount_auto" {
		d.DiscountType = val
		d.DiscountValueOk = true
	} else {
		d.DiscountType = val
		// *** Get value of discount value
		pName = "discount_manual_value"
		val = r.FormValue(pName)
		log.Debug(fmt.Sprintf("Uploaded value of %s : %v", pName,  val))
		// Check value
		if num, err := strconv.ParseFloat(val, 64); err != nil {
			log.Warn(fmt.Sprintf("Error when parse discount value. Error: %s", err.Error()))
			errList = append(errList, "Указано некорректное числовое значение скидки")
		} else {
			// set parsed value
			d.DiscountValue = num
			d.RequestDelayOk = true
		}
	}


	d.ErrLog = errList

	// convert parameters data to json string
	jsonString, err := json.Marshal(d)
	if err != nil {
		return err
	}

	//save data to PARAMETERS
	return scnr.GetParam().SetValue(defConfWBName, string(jsonString))
}



// Aggregate and return actual values for page "page-task_param_scan"
func getPageConfWBData(scnr *options.ScannerObj) *options.ConfWBPageData {
	log := scnr.GetLogger()

	// init empty data
	d := options.ConfWBPageData{}

	p , err := scnr.GetParam().GetValue(defConfWBName)

	if err == nil {
		// Parse json
		if err = json.Unmarshal([]byte(p), &d); err != nil {
			log.Error("Can't unmarshal json parameter " + defConfWBName + ", error: " + err.Error())
		}
	} else {
		log.Info("Can't read parameter " + defConfWBName + ", error: " + err.Error())
	}

	return &d
}
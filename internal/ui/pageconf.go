package ui

import (
	"encoding/json"
	"mrkpl_scanner/internal/options"
	"net/http"
)

const (
	defConfWBName = "ConfigWBPageData"
)

// Processor for page-task_param_scan.gohtml uploads
func procConfWB(d *options.ConfWBPageData, r *http.Request, scnr *options.ScannerObj) (error) {
	log := scnr.GetLogger()


	log.Debug("PROC_CONF_WB")

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
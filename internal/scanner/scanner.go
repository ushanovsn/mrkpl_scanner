package scanner

import (
	"fmt"
	"github.com/ushanovsn/golanglogger"
	"mrkpl_scanner/internal/options"
	"mrkpl_scanner/internal/ui"
	"mrkpl_scanner/pkg/gdoc"

	"github.com/ushanovsn/goutils/baseconf"
	"github.com/ushanovsn/goutils/params"
)

// Starting processes
func StartService(scnr *options.ScannerObj) {
	var err error
	log := scnr.GetLogger()
	log.Out("Service starting...")

	// start GUI
	err = ui.StartUIServer(scnr)
	if err != nil {
		log.Error(err.Error())
		return
	}

	doc := scnr.GetGDocSvc()
	err = doc.GoogleAuth()
	if err != nil {
		log.Error(err.Error())
		return
	}

	n, err := gdoc.GoogleSetSheet(doc, 0, "")
	if err != nil {
		log.Error(err.Error())
		return
	}

	fmt.Printf("start with sheet: %s \n", n)

	// scanner START!

	// test
	dErr := gdoc.GoogleTest(doc)
	if dErr != nil {
		log.Error(dErr.Error())
		return
	}

}

// Stop all process
func StopService(scnr *options.ScannerObj) {
	log := scnr.GetLogger()
	log.Out("Service stopping...")

	// stopping logger
	scnr.GetLogger().StopLog()
}

// Init data and configurations.
//
// Load default values when no config found
func InitService() (*options.ScannerObj, error) {
	// create service object (options)
	var scnr options.ScannerObj

	// base initializing
	initCF(&scnr)

	log := scnr.GetLogger()

	// init parameters object. This critical object and app must stopping when error
	param, err := params.New(scnr.GetParamFName(), params.DataEncrypt, scnr.GetAdminPass())
	if err != nil {
		log.Error(fmt.Sprintf("Error when parameters object init. file: %s. Error: %s", scnr.GetParamFName(), err.Error()))
		log.StopLog()
		return nil, err
	}

	// save parameters object to scanner
	scnr.SetParam(param)

	// init and save emty google doc object
	scnr.SetGDocSvc(gdoc.NewGDoc())

	//todo init data when starting GDoc operations????

	// set URL for GDoc
	gURL, err := param.GetValue(options.DefParamGURLNameP)
	if err != nil {
		log.Error(fmt.Sprintf("Error when read parameter : %s. Error: %s", options.DefParamGURLNameP, err.Error()))
	} else {
		if err := scnr.GetGDocSvc().SetGSheetURL(gURL); err != nil {
			log.Error(fmt.Sprintf("Error when Google sheet URL aplying. URL: %s. Error: %s", "", err.Error()))
		}
	}
	// set Auth data for GDoc
	gAuth, err := param.GetValue(options.DefParamGAuthNameP)
	if err != nil {
		log.Error(fmt.Sprintf("Error when read parameter : %s. Error: %s", options.DefParamGAuthNameP, err.Error()))
	} else {
		if err := scnr.GetGDocSvc().SetAuthKeyFile(gAuth); err != nil {
			log.Error(fmt.Sprintf("Error when Google sheet URL aplying. URL: %s. Error: %s", "", err.Error()))
		}
	}

	// init and save Web UI object
	scnr.SetUIObj(ui.NewUI())

	return &scnr, nil
}

// Init flags and configurations and apply them. Return no error - Load default values when wrong data or no config found
func initCF(scnr *options.ScannerObj) {
	// start config with default values
	scnr.SetDefaultConf()
	// receive flags at start and use it
	setCmdFlags(scnr.GetConfigPtr())

	// start logger with init values (flag received or default value)
	lvl, _ := golanglogger.LoggingLevelValue(scnr.GetConfigPtr().LogLevel)
	log := golanglogger.NewSync(lvl, scnr.GetConfigPtr().LogFile)
	// save logger to object
	scnr.SetLogger(log)

	log.Out("The service is being initialized now...")

	// load and process configuration file
	ok := baseconf.ProcConfig(scnr, scnr.GetLogger())
	if !ok {
		log.Error("Error while read configuration from file. Missing parameters was set to default values")
	}

	setCmdFlags(scnr.GetConfigPtr())
	log.Info("Updated config by received flags")

	// apply the configuration
	log.Out("Now applying configuration parameters")

	if lvl, _ := golanglogger.LoggingLevelValue(scnr.GetConfigPtr().LogLevel); log.CurrentLevel() != lvl {
		log.SetLevel(lvl)
	}
	if szm, szd := log.CurrentFileControl(); szm != int(scnr.GetConfigPtr().LogSizeMb) || szd != int(scnr.GetConfigPtr().LogSizeD) {
		log.SetFileParam(int(szm), int(szd))
	}
}

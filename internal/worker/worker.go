package worker

import (
	//"fmt"
	"github.com/ushanovsn/goutils/baseconf"
	"mrkpl_scanner/internal/options"
)

// Start service for management and interaction
func StartWorker(scnr *options.ScannerObj) {

	updateWPConf(scnr)

	scnr.GetWG().Add(1)

	go runWorker(scnr)
}

// Init configurations and apply them. Return no error - Load default values when wrong data or no config found
func updateWPConf(scnr *options.ScannerObj) {
	log := scnr.GetLogger()
	log.Debug("Init configuration for worker-parser")

	// load and process configuration file
	ok := baseconf.ProcConfig(scnr.GetWPObj(), scnr.GetLogger())
	if !ok {
		log.Error("Error while read workwer-parser configuration from file. Missing parameters was set to default values")
	}
}

// Run worker goroutine process
func runWorker(scnr *options.ScannerObj) {
	log := scnr.GetLogger()
	log.Out("Starting Worker-Parser...")

	wrkr := scnr.GetWPObj()
	ch := *wrkr.GetCMDChan()

	// main worker cycle
	for {
		select {
		case cmd := <-ch:
			if cmd == options.CMDExit {
				log.Info("Command Exit recrived. Stopping Worker-Parser...")
				goto exitWP
			}
		default:
			// work
		}
	}

	// close application
exitWP:

	scnr.GetWG().Done()
}

// Stopping worker-parser process
func StopWorker(scnr *options.ScannerObj) {
	log := scnr.GetLogger()
	log.Info("Stopping Worker-Parser...")

	ch := *scnr.GetWPObj().GetCMDChan()

	ch <- options.CMDExit

	log.Info("Worker-Parser Stopped...")
}

package worker

import (
	//"fmt"
	"fmt"
	"mrkpl_scanner/internal/options"
	"strconv"
	"strings"
	"time"

	"github.com/ushanovsn/goutils/baseconf"
)

// Start service for management and interaction
func StartWorker(scnr *options.ScannerObj) {

	// read config
	updateWPConf(scnr)

	// run worker
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
	stat := wrkr.GetStatus()
	doc := scnr.GetGDocSvc()

	var startRow uint
	var curRow uint
	var lastRow uint
	if v1, v2, ok := getRowRange(wrkr.GetConfigPtr()); ok {
		startRow = v1
		lastRow = v2
	}

	// delay for parsing
	waitTime := wrkr.GetConfigPtr().ReqDelay
	// start with first row
	curRow = startRow

	// main worker cycle
	for {
		select {
		case cmd := <-ch:
			if cmd == options.CMDExit {
				log.Info("Command Exit recrived. Stopping Worker-Parser...")
				goto exitWP
			}
		case <-time.After(time.Second * time.Duration(waitTime)):

			// get row number for reading
			if curRow == 0 {
				stat.Set("Stop", "Значение диапазона строк исходных данных отсутствует или некорректно")
				log.Warn("Have no row number to start reading doc")
				continue
			}

			// reading row data from document
			// check status service and init it
			if doc.GetCurSheetName() == "" {
				if err := doc.GoogleInit(); err != nil {
					stat.Set("Stop", "Ошибка использования сервиса Google Docs")
					log.Warn("Have no active gdoc sheet, or error acces it: " + err.Error())
					continue
				}
			}

			cells, err := doc.ReadRow(curRow)
			if err != nil {
				stat.Set("Stop", "Ошибка чтения строки Google Docs")
				log.Warn("Error read row gdoc sheet: " + err.Error())
				continue
			}

			// parse data
			//todo
			fmt.Printf("TEST!!! parsing value: %v\n", cells)

			c := map[int]string{
				wrkr.GetConfigPtr().PriceCol:   "100",
				wrkr.GetConfigPtr().ErrCol:     "No err",
				wrkr.GetConfigPtr().ErrCol + 1: "",
			}

			// write result into document
			if err := doc.UpdateRowVal(curRow, c); err != nil {
				stat.Set("Stop", "Ошибка записи строки Google Docs")
				log.Warn("Error update row gdoc sheet: " + err.Error())
				continue
			}

			// next
			curRow++
			if curRow > lastRow && lastRow > 0 {
				curRow = startRow
				log.Out("NEXT: " + fmt.Sprint(curRow))
			} else {
				// check empty line from document
				log.Out("NEXT: " + fmt.Sprint(curRow))
			}
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

// Get values of row range for parsing in document
func getRowRange(cfg *options.WorkerConfig) (startRowNum uint, endRowNum uint, ok bool) {
	values := strings.Split(cfg.LinesProc, "-")
	// split for one or two values
	if len(values) > 0 && len(values) < 3 {
		// parse first
		v, err := strconv.ParseUint(values[0], 10, 64)
		if err == nil && v > 0 {
			startRowNum = uint(v)
			ok = true

			// have two values
			if len(values) == 2 {
				// parse second
				v, err := strconv.ParseUint(values[1], 10, 64)
				if err != nil {
					// if have second value but it wrong - whore result is false
					ok = false
				} else {
					endRowNum = uint(v)
				}
			}
		}
	}

	return startRowNum, endRowNum, ok
}

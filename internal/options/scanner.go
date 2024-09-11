package options

import (
	"mrkpl_scanner/pkg/gdoc"
	"sync"

	"github.com/ushanovsn/golanglogger"
	"github.com/ushanovsn/goutils/params"
)

// Scanner object
type ScannerObj struct {
	conf     ScannerConfig
	logger   golanglogger.Golanglogger
	gDoc     *gdoc.GDocObj
	uiObj    *UIObj
	paramObj *params.ParamsObj
	wg       *sync.WaitGroup
	stpChan  chan int
}

//
// LOGGER //

// Getting the logger interface object (the interface is actually a pointer)
func (obj *ScannerObj) GetLogger() golanglogger.Golanglogger {
	return obj.logger
}

// Set the logger object (an interface object or pointer to object that imlement interface Golanglogger)
func (obj *ScannerObj) SetLogger(log golanglogger.Golanglogger) {
	obj.logger = log
}

//
// CONFIGURATION //

// Getting the universal pointer to configurations structure (as Interface)
func (obj *ScannerObj) GetConfigUniversalPtr() interface{} {
	return &obj.conf
}

// Getting the pointer to configurations structure
func (obj *ScannerObj) GetConfigPtr() *ScannerConfig {
	return &obj.conf
}

// Getting the config file path (or just name)
func (obj *ScannerObj) GetConfFileName() string {
	return obj.conf.ConfFile
}

// Getting the description for config file
func (obj *ScannerObj) GetConfigDescr() string {
	return DefScnrConfDescr
}

// Getting the admin password
func (obj *ScannerObj) GetAdminPass() string {
	return obj.conf.AdminPassword
}

// Setting the admin password
func (obj *ScannerObj) SetAdminPass(pass string) {
	obj.conf.AdminPassword = pass
}

//
// GOOGLE DOCS //

// Get google service
func (obj *ScannerObj) GetGDocSvc() *gdoc.GDocObj {
	return obj.gDoc
}

// Set google service
func (obj *ScannerObj) SetGDocSvc(doc *gdoc.GDocObj) {
	obj.gDoc = doc
}

//
// USER INTERFACE //

// Get UI obj to scanner obj
func (obj *ScannerObj) GetUIObj() *UIObj {
	return obj.uiObj
}

// Set UI obj to scanner obj
func (obj *ScannerObj) SetUIObj(ui *UIObj) {
	obj.uiObj = ui
}

//
// PARAMETERS //

// Getting the parameters object
func (obj *ScannerObj) GetParam() *params.ParamsObj {
	return obj.paramObj
}

// Setting the parameters object
func (obj *ScannerObj) SetParam(p *params.ParamsObj) {
	obj.paramObj = p
}

// Getting the parameters file name
func (obj *ScannerObj) GetParamFName() string {
	return DefParamFileName
}

// SYNC OBJECTS

// Get wait group
func (obj *ScannerObj) GetWG() *sync.WaitGroup {
	if obj.wg == nil {
		obj.wg = &sync.WaitGroup{}
	}
	return obj.wg
}

// Get stop proc channel
func (obj *ScannerObj) GetStopChan() chan int {
	if obj.stpChan == nil {
		obj.stpChan = make(chan int)
	}
	return obj.stpChan
}

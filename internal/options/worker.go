package options

/*
import (
	"fmt"
)
*/

// Worker-parser object
type WPObj struct {
	// worker configuration
	conf WorkerConfig
	// command channel
	cmdChnl chan WrkrCMD
}

// Create New Worker-Parser object
func NewWP() *WPObj {
	obj := WPObj{}
	obj.SetDefaultConf()
	return &obj
}

// Getting the config file path (or just name)
func (obj *WPObj) GetConfFileName() string {
	return obj.conf.ConfFile
}

// Getting the pointer to configurations structure
func (obj *WPObj) GetConfigUniversalPtr() interface{} {
	return &obj.conf
}

// Getting the description (common comments) for header of config file
func (obj *WPObj) GetConfigDescr() string {
	return DefWrkrConfDescr
}

// Getting the pointer to configurations structure
func (obj *WPObj) GetConfigPtr() *WorkerConfig {
	return &obj.conf
}

func (obj *WPObj) GetCMDChan() *chan WrkrCMD {
	if obj.cmdChnl == nil {
		obj.cmdChnl = make(chan WrkrCMD)
	}
	return &obj.cmdChnl
}

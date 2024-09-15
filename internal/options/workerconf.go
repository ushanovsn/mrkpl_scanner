package options

const (
	// default config file name
	DefWrkrConfFile string = "parser.conf"
	// default description for config file
	DefWrkrConfDescr string = "The configuration file for \"PARSER SERVICE\""
)

// Configuration for worker parser
type WorkerConfig struct {
	// document type with data
	ReadDocType string `cfg:"read_document_type" descr:"Type of document with incoming data (values: msexcel, googledoc)"`
	// saving to original doc flag
	SaveToOrig bool `cfg:"save_in_original" descr:"Save parsed data to original doc"`

	// file with configuration parameters
	ConfFile string

	// temp params
	// address column number
	AddressCol int `cfg:"address_column" descr:"Column number where marketplace goods address for parsing"`
	// price column number
	PriceCol int `cfg:"price_column" descr:"Column number for saving parsed price"`
	// error column number
	ErrCol int `cfg:"error_column" descr:"Column number for saving parse error when occures"`
}

// worker cmd enum
type WrkrCMD int

const (
	CMDStart WrkrCMD = iota
	CMDStop
	CMDRestart
	CMDExit
)

// default values for worker\parser
func (obj *WPObj) SetDefaultConf() {
	obj.conf.ConfFile = DefWrkrConfFile
	obj.conf.ReadDocType = ""
	obj.conf.AddressCol = -1
	obj.conf.PriceCol = -1
	obj.conf.ErrCol = -1
}

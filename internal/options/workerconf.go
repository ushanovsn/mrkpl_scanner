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
	// range row with data in document
	LinesProc string `cfg:"lines_range" descr:"Row range in document with data for processing (values like: 2-300, 3- (full doc))"`
	// request delay
	ReqDelay float64 `cfg:"request_delay" descr:"Delay between requests to marketplace in seconds(values: 0.2, 1.5, 2, ...)"`

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

// Worker current status
type WorkerStatus struct {
	status string
	msg    string
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

// Set values to status struct
func (st *WorkerStatus) Set(status string, msg string) {
	st.status = status
	st.msg = msg
}

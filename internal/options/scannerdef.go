package options

const (
	// default config file name
	DefScnrConfFile string = "scanner.conf"
	// default log file name
	DefScnrLogFile string = "scanner.log"

	// default logger level
	DefScnrLogLvl string = "Error"
	// default logger name
	DefScnrLogName string = "SCANNER"
	// default log file size in megabytes
	DefScnrLogSizeMb uint = 0
	// default log file size in days
	DefScnrLogSizeDay uint = 0

	// default host address
	DefScnrHost string = "localhost"
	// default host port
	DefScnrPort uint = 3003

	// default description for config file
	DefScnrConfDescr string = "The configuration file for the \"SCANNER\""
)

// Set default values to configuration structure
func (obj *ScannerObj) SetDefaultConf() {
	obj.conf.Host = DefScnrHost
	obj.conf.Port = DefScnrPort
	obj.conf.LogLevel = DefScnrLogLvl
	obj.conf.LogSizeMb = DefScnrLogSizeMb
	obj.conf.LogSizeD = DefScnrLogSizeDay
	obj.conf.LogFile = DefScnrLogFile
	obj.conf.ConfFile = DefScnrConfFile
}

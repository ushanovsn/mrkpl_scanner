package options

import (
	"fmt"
)

// Configuration for scanner
type ScannerConfig struct {
	// host address
	Host string `cfg:"host" descr:"Service host address"`
	// host port
	Port uint `cfg:"port" descr:"Service host port"`
	// logger level
	LogLevel string `cfg:"logging_level" descr:"Logger logging level (Debug/Info/Warning/Error)"`
	// log file size in megabytes
	LogSizeMb uint `cfg:"log_file_size_mb" descr:"Log file size in megabytes (0 - one file/no split)"`
	// log file size in megabytes
	LogSizeD uint `cfg:"log_file_size_day" descr:"Log file size in days (0 - one file/no split)"`
	// logging file (no file if empty string)
	LogFile string `cfg:"log_file" descr:"File name or full path for logging file (without spaces or use quotes)"`
	// file with configuration parameters (can be specified\changed only by a flag params)
	ConfFile string
	// password for admin
	AdminPassword string
}

// Getting the net address string
func (obj *ScannerObj) GetAddr() string {
	return fmt.Sprintf("%s:%v", obj.conf.Host, obj.conf.Port)
}

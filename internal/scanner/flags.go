package scanner

import (
	"flag"
	"mrkpl_scanner/internal/options"
	"os"
)

// Set (update) flags values from cmd into config.
//
// Update only received flags, others configuration parameters are not changes.
func setCmdFlags(conf *options.ScannerConfig) {
	var tmpF options.ScannerConfig
	var flags *flag.FlagSet

	// if base flags was set and parsed - changing flagset
	if !flag.Parsed() {
		flags = flag.CommandLine
	} else {
		flags = flag.NewFlagSet("new flag set", flag.ExitOnError)
	}

	// deffine the parameters
	flags.StringVar(&tmpF.Host, "h", conf.Host, "Service host address")
	flags.UintVar(&tmpF.Port, "p", conf.Port, "Service port")
	flags.StringVar(&tmpF.LogLevel, "loglvl", conf.LogLevel, "Service logging level (text)")
	flags.StringVar(&tmpF.LogFile, "logfile", conf.LogFile, "Service log file path\\name")
	flags.StringVar(&tmpF.ConfFile, "conffile", conf.ConfFile, "Service config file path\\name")

	// service will stopping in this place when error occurs
	_ = flags.Parse(os.Args[1:])

	// now set parameters to config
	conf.Host = tmpF.Host
	conf.Port = tmpF.Port
	conf.LogLevel = tmpF.LogLevel
	conf.LogFile = tmpF.LogFile
	conf.ConfFile = tmpF.ConfFile
}

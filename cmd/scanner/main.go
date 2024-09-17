package main

import (
	"fmt"
	"mrkpl_scanner/internal/scanner"
	"mrkpl_scanner/pkg/contitle"
	_ "mrkpl_scanner/static/rsrc"
)

func main() {
	title := "### Scanner ###"
	fmt.Printf("%s\n", title)
	// no control error
	_ = contitle.SetTitle(title)

	scnr, err := scanner.InitService()

	if err == nil {
		scanner.RunService(scnr)
		// when stop or break - gracefull shutdown
		scanner.StopService(scnr)
	} else {
		// the logger may not be initialized
		fmt.Printf("### Init ERROR: %s\n", err.Error())
	}
}

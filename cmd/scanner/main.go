package main

import (
	"fmt"
	"mrkpl_scanner/internal/scanner"
	"mrkpl_scanner/pkg/contitle"
	_ "mrkpl_scanner/static/rsrc"
)

func main() {
	fmt.Printf("### Scanner ###\n")

	// no control error
	_, err := contitle.SetTitle("### Scanner ###")
	if err != nil {
		fmt.Printf("error when set title: %s\n", err)
	}

	scnr, err := scanner.InitService()
	if err == nil {
		scanner.RunService(scnr)

		scanner.StopService(scnr)
	} else {
		fmt.Printf("### Init ERROR: %s\n", err.Error())
	}
}

package main

import (
	"fmt"
	"mrkpl_scanner/internal/scanner"
)

func main() {
	fmt.Printf("### Scanner ###\n")

	scnr, err := scanner.InitService()
	if err == nil {
		scanner.RunService(scnr)

		scanner.StopService(scnr)
	} else {
		fmt.Printf("### Init ERROR: %s\n", err.Error())
	}
}

package wbparser

import (
	"fmt"
	"mrkpl_scanner/pkg/mrktparsers"
)

// Wildberries price parser main object
type Parser struct {
	result   mrktparsers.ParsedData
	itemName string
}

// Get new initialized parser
func New() *Parser {
	p := Parser{}
	return &p
}

// Get data from marketplace of item by link
func (p *Parser) GetItem(lnk string) error {
	fmt.Printf("Getting data by link: %s\n", lnk)

	return nil
}

// Parse data received from marketplace before
func (p *Parser) ParseItem() (d mrktparsers.ParsedData, err error) {
	fmt.Printf("Parsing data...\n")

	return d, err
}

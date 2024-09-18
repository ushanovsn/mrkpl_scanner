package uniparser

import (
	"fmt"
	"mrkpl_scanner/pkg/mrktparsers"
	"mrkpl_scanner/pkg/mrktparsers/wbparser"
)

// Wildberries price parser main object
type Parser struct {
	result  mrktparsers.ParsedData
	lnk     string
	pStatus bool
	prs     map[mrktparsers.MrktplcType]mrktparsers.MrktParser
}

// Get new initialized parser
func New() *Parser {
	p := Parser{
		prs: make(map[mrktparsers.MrktplcType]mrktparsers.MrktParser),
	}
	return &p
}

// Get data from marketplace of item by link
func (p *Parser) GetItem(lnk string) error {
	// check type
	t := mrktparsers.MrktRecognize(lnk)

	// reset old values
	p.result = mrktparsers.ParsedData{}
	p.lnk = ""
	p.pStatus = false

	if t == mrktparsers.MTYPE_WILDBERRIES {
		pr, ok := p.prs[mrktparsers.MTYPE_WILDBERRIES]
		if !ok {
			pr = wbparser.New()
			p.prs[mrktparsers.MTYPE_WILDBERRIES] = pr
		}

		p.lnk = lnk
		err := pr.GetItem(p.lnk)

		return err
	}

	return fmt.Errorf("Marketplace \"%v\" not accepted", t.Name())
}

// Parse data received from marketplace before
func (p *Parser) ParseItem() error {
	// check type
	t := mrktparsers.MrktRecognize(p.lnk)

	if t == mrktparsers.MTYPE_WILDBERRIES {
		pr, ok := p.prs[mrktparsers.MTYPE_WILDBERRIES]
		if !ok {
			return fmt.Errorf("Marketplace \"%v\" not initialized", t.Name())
		}

		err := pr.ParseItem()
		if err == nil {
			p.result, p.pStatus = pr.GetResult()
		} else {
			p.pStatus = false
		}
		return err
	}

	return fmt.Errorf("Marketplace \"%v\" not accepted", t.Name())
}

// Get result of parsing
func (p *Parser) GetResult() (mrktparsers.ParsedData, bool) {
	return p.result, p.pStatus
}

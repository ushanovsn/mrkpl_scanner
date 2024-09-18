package wbparser

import (
	"encoding/json"
	"fmt"
	"io"
	"mrkpl_scanner/pkg/mrktparsers"
	"net/http"
	"regexp"
	"strings"
)

// Wildberries price parser main object
type Parser struct {
	result  mrktparsers.ParsedData
	pStatus bool
	lnk     string
	reId    regexp.Regexp
	itemId  string
	rawData []byte
}

// Get new initialized parser
func New() *Parser {
	p := Parser{
		//re: *regexp.MustCompile(`^((www.)|(https://www.)|(https://))*wildberries.ru/catalog/(\d+)\S*\z`),
		reId: *regexp.MustCompile(`^((www.)|(https://www.)|(https://))*wildberries.ru/catalog/(\d+)\S*\z`),
	}
	return &p
}

// Get result of parsing
func (p *Parser) GetResult() (mrktparsers.ParsedData, bool) {
	return p.result, p.pStatus
}

// Get data from marketplace of item by link
func (p *Parser) GetItem(lnk string) error {
	// reset old values
	p.result = mrktparsers.ParsedData{}
	p.pStatus = false

	f := p.reId.FindStringSubmatch(lnk)
	if len(f) == 0 {
		return fmt.Errorf("Incorrect link type, can't find item id")
	}

	p.lnk = lnk
	// need last value
	p.itemId = f[len(f)-1]

	// url for request data
	url := fmt.Sprintf("https://card.wb.ru/cards/v2/detail?appType=1&curr=rub&dest=%s&ab_testing=false&nm=%s", "-1576352", p.itemId)

	res, err := http.Get(url)
	if err != nil {
		return fmt.Errorf("Error when Request data: %s", err.Error())
	} else if res.StatusCode > 299 {
		return fmt.Errorf("Error. Get status code: %v", res.StatusCode)
	}

	// check data type
	if c := res.Header.Get("Content-Type"); !strings.Contains(c, "application/json") {
		return fmt.Errorf("Received wrong content type: %v", c)
	}

	// read json body
	p.rawData, err = io.ReadAll(res.Body)
	if err != nil {
		return fmt.Errorf("Error reading body of response: %s", err.Error())
	}

	return nil
}

// Parse data received from marketplace before
func (p *Parser) ParseItem() error {
	fmt.Printf("Parsing data...\n")

	if err := json.Unmarshal(p.rawData, &p); err != nil {
		return err
	}

	fmt.Printf("prd: %v \n", p.result)
	return nil
}

func (p *Parser) UnmarshalJSON(data []byte) error {
	var j FullJson = FullJson{}
	if err := json.Unmarshal(data, &j); err != nil {
		return err
	}

	if len(j.Data.Products) == 0 {
		return fmt.Errorf("Parameter product is empty")
	}

	p.result.Name = j.Data.Products[0].Name
	p.result.Brand = j.Data.Products[0].Brand
	p.result.Supplier = j.Data.Products[0].Supplier
	p.result.RegularPrice = float64(j.Data.Products[0].Sizes[0].Price.Basic) / 100
	p.result.DiscountPrice = float64(j.Data.Products[0].Sizes[0].Price.Product) / 100.0
	p.result.IndividualPrice = (float64(j.Data.Products[0].Sizes[0].Price.Product) / 100.0) * 0.97

	p.pStatus = true

	return nil
}

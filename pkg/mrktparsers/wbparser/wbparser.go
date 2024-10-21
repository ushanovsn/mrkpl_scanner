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
	result      mrktparsers.ParsedData
	pStatus     bool
	lnk         string
	urlTemplate string
	sizeId      string
	reId        regexp.Regexp
	reSize      regexp.Regexp
	itemId      string
	location    int
	discount    float64
	rawData     []byte
}

// Get new initialized parser
func New() *Parser {
	p := Parser{
		urlTemplate: "https://card.wb.ru/cards/v2/detail?appType=1&curr=rub&dest=%s&ab_testing=false&nm=%s",
		reId:        *regexp.MustCompile(`^((www.)|(https://www.)|(https://))*wildberries.ru/catalog/(\d+)\S*\z`),
		reSize:      *regexp.MustCompile(`^*wildberries.ru/catalog/\d+/\S*\?((\S*\&)*size=(\d+)*)\S*\z`),
	}
	return &p
}

// Get new initialized parser
func (p *Parser) SetLocation() error {
	// url for request data
	url := "https://user-geo-data.wildberries.ru/get-geo-info?currency=RUB"

	res, err := http.Get(url)
	if err != nil {
		return fmt.Errorf("Error when Request data: %s", err.Error())
	} else if res.StatusCode > 299 {
		res.Body.Close()
		return fmt.Errorf("Error. Get status code: %v", res.StatusCode)
	}
	defer res.Body.Close()

	// check data type
	if c := res.Header.Get("Content-Type"); !strings.Contains(c, "application/json") {
		return fmt.Errorf("Received wrong content type: %v", c)
	}

	// read json body
	rawData, err := io.ReadAll(res.Body)
	if err != nil {
		return fmt.Errorf("Error reading body of response: %s", err.Error())
	}

	var jData struct {
		Address      string `json:"address"`
		Ip           string `json:"ip"`
		Currency     string `json:"currency"`
		Locale       string `json:"locale"`
		Destinations []int  `json:"destinations"`
	}

	err = json.Unmarshal(rawData, &jData)
	if err != nil {
		return fmt.Errorf("Error unmarshalling json: %s", err.Error())
	}

	p.location = jData.Destinations[len(jData.Destinations)-1]
	return nil
}

// Get current discount
func (p *Parser) SetCurDiscount() error {
	// url for request data
	url := "https://static-basket-01.wb.ru/vol1/global-payment/default-payment.json"

	res, err := http.Get(url)
	if err != nil {
		return fmt.Errorf("Error when Request data: %s", err.Error())
	} else if res.StatusCode > 299 {
		res.Body.Close()
		return fmt.Errorf("Error. Get status code: %v", res.StatusCode)
	}
	defer res.Body.Close()

	// check data type
	if c := res.Header.Get("Content-Type"); !strings.Contains(c, "application/json") {
		return fmt.Errorf("Received wrong content type: %v", c)
	}

	// read json body
	rawData, err := io.ReadAll(res.Body)
	if err != nil {
		return fmt.Errorf("Error reading body of response: %s", err.Error())
	}

	var jData struct {
		Data []struct {
			Discount float64 `json:"discount_value"`
		} `json:"data"`
	}

	err = json.Unmarshal(rawData, &jData)
	if err != nil {
		return fmt.Errorf("Error unmarshalling json: %s", err.Error())
	}

	p.discount = jData.Data[0].Discount
	return nil
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
	if len(f) < 2 {
		return fmt.Errorf("Incorrect link type, can't find item id")
	}

	// correct link - save it
	p.lnk = lnk
	// need last value
	p.itemId = f[len(f)-1]

	// check size id in link string
	s := p.reId.FindStringSubmatch(lnk)
	if len(s) < 2 {
		p.sizeId = ""
	} else {
		p.sizeId = s[len(s)-1]
	}

	// url for request data
	url := fmt.Sprintf(p.urlTemplate, p.location, p.itemId)

	res, err := http.Get(url)
	if err != nil {
		return fmt.Errorf("Error when Request data: %s", err.Error())
	} else if res.StatusCode > 299 {
		res.Body.Close()
		return fmt.Errorf("Error. Get status code: %v", res.StatusCode)
	}
	defer res.Body.Close()

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
	// start unmarshal to structure in parser
	return json.Unmarshal(p.rawData, &p)
}

func (p *Parser) UnmarshalJSON(data []byte) error {
	var j FullJson = FullJson{}
	if err := json.Unmarshal(data, &j); err != nil {
		return err
	}

	if len(j.Data.Products) == 0 {
		return fmt.Errorf("Parameter product is empty")
	}

	// find id of product
	pId := 0
	for k, v := range j.Data.Products {
		if fmt.Sprint(v.Id) == p.itemId {
			pId = k
		}
	}

	// find sort id of product
	sId := 0
	if p.sizeId != "" {
		for k, v := range j.Data.Products[pId].Sizes {
			if fmt.Sprint(v.OptionId) == p.sizeId {
				sId = k
			}
		}
	}

	// save values
	p.result.Name = j.Data.Products[pId].Name
	p.result.Brand = j.Data.Products[pId].Brand
	p.result.Supplier = j.Data.Products[pId].Supplier
	p.result.RegularPrice = float64(j.Data.Products[pId].Sizes[sId].Price.Basic) / 100
	p.result.DiscountPrice = float64(j.Data.Products[pId].Sizes[sId].Price.Product) / 100.0
	p.result.IndividualPrice = (float64(j.Data.Products[pId].Sizes[sId].Price.Product) / 100.0) * (100 - (p.discount / 100))

	p.pStatus = true

	return nil
}

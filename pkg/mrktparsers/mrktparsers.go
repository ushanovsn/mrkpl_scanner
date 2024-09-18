package mrktparsers

import (
	"regexp"
)

// Common marketplaces parsers interface
type MrktParser interface {
	// Get full data by address of item (link to the product)
	GetItem(string) error
	// Parse getted item data by link
	ParseItem() error
	// Get parced result
	GetResult() (ParsedData, bool)
}

// Parsed data of item (product)
type ParsedData struct {
	// Name (description) of item (product)
	Name string
	// Brand
	Brand string
	// Supplier
	Supplier string
	// Base (regular) price
	RegularPrice float64
	// Discounted (common) price
	DiscountPrice float64
	// Price with additional promo-discounts
	IndividualPrice float64
}

// Marketplace types
type MrktplcType int

const (
	MTYPE_NO_DEFINED MrktplcType = iota
	MTYPE_WILDBERRIES
	MTYPE_OZON
	MTYPE_YANDEX_MARKET
)

// Defining the type of marketplace
func MrktRecognize(lnk string) MrktplcType {
	var exprWb string
	var match bool
	var err error

	// wilberries
	exprWb = `^(\S*\.)*wildberries.ru\S*\z`
	match, err = regexp.MatchString(exprWb, lnk)
	if match && err == nil {
		return MTYPE_WILDBERRIES
	}

	// ozon
	exprWb = `^(\S*\.)*ozon.ru\S*\z`
	match, err = regexp.MatchString(exprWb, lnk)
	if match && err == nil {
		return MTYPE_OZON
	}

	// yandex
	exprWb = `^(\S*\.)*market.yandex.ru\S*\z`
	match, err = regexp.MatchString(exprWb, lnk)
	if match && err == nil {
		return MTYPE_YANDEX_MARKET
	}

	return MTYPE_NO_DEFINED
}

// Get name of marketplace type
func (mt MrktplcType) Name() string {
	switch mt {
	case MTYPE_WILDBERRIES:
		return "WILDBERRIES"
	case MTYPE_OZON:
		return "OZON"
	case MTYPE_YANDEX_MARKET:
		return "YANDEX MARKET"
	case MTYPE_NO_DEFINED:
		fallthrough
	default:
		return "UNKNOWN"
	}
}

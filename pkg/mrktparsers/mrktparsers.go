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
	MTYPE_AVITO
)

// Defining the type of marketplace
func MrktRecognize(lnk string) MrktplcType {
	var expr string
	var match bool
	var err error

	// wilberries
	expr = `^(\S*\.)*wildberries.ru\S*\z`
	match, err = regexp.MatchString(expr, lnk)
	if match && err == nil {
		return MTYPE_WILDBERRIES
	}

	// ozon
	expr = `^(\S*\.)*ozon.ru\S*\z`
	match, err = regexp.MatchString(expr, lnk)
	if match && err == nil {
		return MTYPE_OZON
	}

	// yandex
	expr = `^(\S*\.)*market.yandex.ru\S*\z`
	match, err = regexp.MatchString(expr, lnk)
	if match && err == nil {
		return MTYPE_YANDEX_MARKET
	}

	// avito
	expr = `^(\S*\.)*avito.ru\S*\z`
	match, err = regexp.MatchString(expr, lnk)
	if match && err == nil {
		return MTYPE_AVITO
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
	case MTYPE_AVITO:
		return "AVITO"
	case MTYPE_NO_DEFINED:
		fallthrough
	default:
		return "UNKNOWN"
	}
}

package mrktparsers

// Common marketplaces parsers interface
type MrktParser interface {
	// Get full data by address of item (link to the product)
	GetItem(string) error
	// Parse getted item data by link
	ParseItem() (ParsedData, error)
}

// Parsed data of item (product)
type ParsedData struct {
	// Name (description) of item (product)
	Name string
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

	// PARSING 	address link
	//
	//
	//

	return MTYPE_NO_DEFINED
}

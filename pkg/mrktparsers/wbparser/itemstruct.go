package wbparser

// Full received json data for item
type FullJson struct {
	Data Data `json:"data"`
}

// Parameter "data" in full json
type Data struct {
	// List of products info (if request one product - one element list)
	Products []Product `json:"products"`
}

// One of product with inner parameters
type Product struct {
	Id            int     `json:"id"`
	Name          string  `json:"name"`
	Brand         string  `json:"brand"`
	Supplier      string  `json:"supplier"`
	SupplierId    int     `json:"supplierId"`
	TotalQuantity int     `json:"totalQuantity"`
	Colors        []Color `json:"colors"`
	Sizes         []Size  `json:"sizes"`
}

// One of colors parameter of product item
type Color struct {
	Id   int    `json:"id"`
	Name string `json:"name"`
}

// One of sizes parameter of product item
type Size struct {
	Name     string `json:"name"`
	OrigName string `json:"origName"`
	OptionId int    `json:"optionId"`
	Price    Price  `json:"price"`
}

// Price parameter of product/size parameter
type Price struct {
	Basic     int `json:"basic"`
	Product   int `json:"product"`
	Total     int `json:"total"`
	Logistics int `json:"logistics"`
	Return    int `json:"return"`
}

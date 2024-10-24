package options

// Navigation Tab element
type NaviMenu struct {
	ItmName   string
	ItmLink   string
	ItmIsMenu bool
	ItmMenu   []NaviDropMenu
}

// Navigation Drop-Down from Tab element
type NaviDropMenu struct {
	ItmName string
	ItmLink string
}

// Navigation Active menu and tab elements names
type NaviActiveMenu struct {
	ActiveTabVal    string
	ActiveDMenuVal  string
	PageDescription string
}

// Data for pages templates
type PagesData struct {
	pParamScan	ParamsScanPageData

}
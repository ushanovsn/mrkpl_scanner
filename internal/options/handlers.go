package options

// Data for Google parameters page template
type GParamsPageData struct {
	Title            string
	GPageURL         string
	GPageURLOkPref   string
	AuthClient       string
	AuthClientOkPref string
	ErrLog           []string
}

type PParamsPageData struct {
	Title string
}

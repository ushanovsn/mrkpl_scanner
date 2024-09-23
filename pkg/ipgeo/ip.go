package ipgeo

import (
	"encoding/json"
	"io"
	"net/http"
	"regexp"
)

// requests to several servers and compare results. when get equal pair - return result
func getIp() (string, error) {
	//text ip//
	//req, err := http.Get("http://ifconfig.io/ip")
	//req, err := http.Get("http://ident.me/")
	//req, err := http.Get("http://whatismyip.akamai.com/")
	//req, err := http.Get("http://api.db-ip.com/v2/free/self/ipAddress")
	req, err := http.Get("https://myexternalip.com/raw")
	if err != nil {
		return "", err
	}
	defer req.Body.Close()

	body, err := io.ReadAll(req.Body)
	if err != nil {
		return "", err
	}

	match, err := regexp.MatchString(`^((\d{1,3})\.){3}(\d{1,3})$`, string(body))
	if match && err == nil {
		return string(body), nil
	}

	return "", err
}

// requests to several servers and compare results. when get equal pair - return result
func getIpLocation() (IpLocation, error) {
	//json//
	req, err := http.Get("http://ip-api.com/json/")
	//req, err := http.Get("https://ipapi.co/json/")
	//req, err := http.Get("https://www.iplocationtools.com/")
	//req, err := http.Get("https://ipwhois.io/")
	if err != nil {
		return IpLocation{}, err
	}
	defer req.Body.Close()

	body, err := io.ReadAll(req.Body)
	if err != nil {
		return IpLocation{}, err
	}

	type ip_api_com struct {
		IpAddress   string  `json:"query"`
		CountryCode string  `json:"countryCode"`
		CountryName string  `json:"country"`
		Region      string  `json:"regionName"`
		City        string  `json:"city"`
		Latitude    float64 `json:"lat"`
		Longitude   float64 `json:"lon"`
	}

	var jData ip_api_com

	if err := json.Unmarshal(body, &jData); err != nil {
		return IpLocation{}, err
	}

	// convert and return
	return IpLocation(jData), nil
}

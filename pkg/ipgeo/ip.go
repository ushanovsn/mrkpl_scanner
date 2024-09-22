package ipgeo

import (
	"encoding/json"
	"io"
	"net/http"
)

type IP struct {
	Query string
}

// requests to several servers and compare results. when get equal pair - return result
func getip() (string, error) {
	//json//
	//req, err := http.Get("http://ip-api.com/json/")
	//req, err := http.Get("https://ipapi.co/json/")
	//req, err := http.Get("https://www.iplocationtools.com/")
	//req, err := http.Get("https://ipwhois.io/")

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

	var ip IP
	json.Unmarshal(body, &ip)

	return ip.Query, nil
}

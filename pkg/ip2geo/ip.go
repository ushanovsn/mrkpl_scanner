package ip2geo

import (
	"encoding/json"
	"io"
	"net/http"
)
type IP struct {
    Query string
}

func getip2() (string, error) {
    req, err := http.Get("http://ip-api.com/json/")
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
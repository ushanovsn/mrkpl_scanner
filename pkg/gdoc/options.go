package gdoc

import (
	"errors"

	"google.golang.org/api/sheets/v4"
)

type GDocObj struct {
	// service set object when Autentication executing
	svc         *sheets.Service
	sheetId     int64
	sprdSheetId string
	sheetName   string

	// Client email for autentication (extract from current json-auth file)
	authClient string

	// JSON file that contains service account keys
	authKeyFile string
	// google sheet address
	gSheetURL string
}

// Getting the file authentication
func (obj *GDocObj) GetAuthKeyFile() string {
	return obj.authKeyFile
}

// Setting the file authentication
func (obj *GDocObj) SetAuthKeyFile(j string) error {

	// need add check and parse received json

	obj.authKeyFile = j

	return nil
}

// Getting the google sheet URL
func (obj *GDocObj) GetGSheetURL() string {
	return obj.gSheetURL
}

// Setting the google sheet URL (URL checking and parsing values from it)
func (obj *GDocObj) SetGSheetURL(url string) (err error) {

	if valURL, valSSId, valSId, ok := parseSheetURL(url); ok {
		obj.gSheetURL = valURL
		obj.sprdSheetId = valSSId
		obj.sheetId = valSId
	} else {
		err = errors.New("Error checking and parsing URL, URL is incorrect")
	}

	return err
}

// Getting the current client email
func (obj *GDocObj) GetCurClien() string {
	return obj.authClient
}

package gdoc

import (
	"encoding/json"
	"errors"

	"google.golang.org/api/sheets/v4"
)

// Object interaction with Google document
type GDocObj struct {
	// service set object when Autentication executing
	svc         *sheets.Service
	auth        credentials
	sheetId     int64
	sprdSheetId string
	// Current connected sheet name
	sheetName string

	// Client email for autentication (extract from current json-auth file)
	authClient string

	// JSON file that contains service account keys
	authKeyFile string
	// google sheet address
	gSheetURL string
}

// Credentials data for google account
type credentials struct {
	Type                string `json:"type"`
	ProjectID           string `json:"project_id"`
	PrivateKeyID        string `json:"private_key_id"`
	PrivateKey          string `json:"private_key"`
	ClientEmail         string `json:"client_email"`
	ClientID            string `json:"client_id"`
	AuthURI             string `json:"auth_uri"`
	TokenURI            string `json:"token_uri"`
	AuthProviderCertURL string `json:"auth_provider_x509_cert_url"`
	ClientCertURL       string `json:"client_x509_cert_url"`
}

// Setting the file authentication
func (obj *GDocObj) SetAuthKeyFile(j string) error {
	var cred credentials
	// Parse service account key
	err := json.Unmarshal([]byte(j), &cred)
	if err != nil {
		return err
	}

	obj.auth = cred
	obj.authKeyFile = j
	obj.authClient = cred.ClientEmail
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
func (obj *GDocObj) GetCurClient() string {
	return obj.authClient
}

// Getting the current sheet Name
func (obj *GDocObj) GetCurSheetName() string {
	return obj.sheetName
}

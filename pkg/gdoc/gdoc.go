package gdoc

import (
	"context"
	"fmt"
	"golang.org/x/oauth2/google"
	"google.golang.org/api/option"
	"google.golang.org/api/sheets/v4"
)

// Init new empty GDoc object
func NewGDoc() *GDocObj {
	gDoc := &GDocObj{}
	return gDoc
}

// Autentication with google account key
func (gDoc *GDocObj) GoogleAuth() error {
	var err error

	json_key := gDoc.authKeyFile

	ctx := context.Background()

	// Authentication
	if len(json_key) == 0 {
		return fmt.Errorf("Have No Authentication Data")
	}

	config, err := google.JWTConfigFromJSON([]byte(json_key), "https://www.googleapis.com/auth/spreadsheets")
	if err != nil {
		return fmt.Errorf("Error when accepting configuration for Google account. Error: %s", err.Error())
	}

	// Create client
	client := config.Client(ctx)
	// Create service
	svc, err := sheets.NewService(ctx, option.WithHTTPClient(client))
	if err != nil {
		return fmt.Errorf("Error when creating Google spreadSheet service. Error: %s", err.Error())
	}

	gDoc.svc = svc
	return nil
}

// Check sheet exist and access, return name of sheet
func GoogleSetSheet(gDoc *GDocObj, sheetId int64, spreadSheetId string) (name string, err error) {
	svc := gDoc.svc

	if spreadSheetId == "" {
		return "", fmt.Errorf("Have No Spreadsheet Data")
	}

	gDoc.sheetId = sheetId
	gDoc.sprdSheetId = spreadSheetId

	// get sheet params
	resp, err := svc.Spreadsheets.Get(spreadSheetId).Fields("sheets(properties(sheetId,title))").Do()
	if err != nil || resp.HTTPStatusCode != 200 {
		return "", fmt.Errorf("Error when getting Google spreadSheet properties. Error: %s", err.Error())
	}

	// get sheet name
	sheetName := ""
	for _, v := range resp.Sheets {
		prop := v.Properties
		if prop.SheetId == int64(sheetId) {
			sheetName = prop.Title
			gDoc.sheetName = sheetName
			return sheetName, nil
		}
	}

	return "", fmt.Errorf("Error when getting Google spreadSheet name. Error: Sheet in spread \"%s\" with id %v not exit", spreadSheetId, sheetId)
}

func GoogleTest(gDoc *GDocObj) error {
	svc := gDoc.svc
	spreadSheetId := gDoc.sprdSheetId
	sheetName := gDoc.sheetName

	ctx := context.Background()

	//Append value to the sheet.
	recordRow := sheets.ValueRange{
		Values: [][]interface{}{{}, {}, {"1", "ABC", "abc@gmail.com"}},
	}

	resp2, err := svc.Spreadsheets.Values.Append(spreadSheetId, sheetName, &recordRow).ValueInputOption("USER_ENTERED").InsertDataOption("INSERT_ROWS").Context(ctx).Do()
	if err != nil || resp2.HTTPStatusCode != 200 {
		return err
	}

	return nil
}

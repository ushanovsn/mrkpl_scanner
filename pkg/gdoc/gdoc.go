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

// Initializing objects. Autentication with google account key and set tpe sheet
func (gDoc *GDocObj) GoogleInit() error {
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

	// apply current sheet
	if n, err := googleApplySheet(gDoc); err != nil {
		return fmt.Errorf("Error when open spreadSheet. Error: %s", err.Error())
	} else {
		gDoc.sheetName = n
	}

	return nil
}

// Check sheet exist and access, return name of sheet
func googleApplySheet(gDoc *GDocObj) (name string, err error) {
	svc := gDoc.svc
	spreadSheetId := gDoc.sprdSheetId
	sheetId := gDoc.sheetId

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

// read full row from doc
func (gDoc *GDocObj) ReadRow(rNum uint) (val map[int]string, err error) {
	spreadSheetId := gDoc.sprdSheetId
	srv := gDoc.svc
	//todo change func - using request columns
	readRange := fmt.Sprintf("%s!A%v:Z%v", gDoc.sheetName, rNum, rNum)
	resp, err := srv.Spreadsheets.Values.Get(spreadSheetId, readRange).Do()
	if err != nil {
		return nil, fmt.Errorf("Unable to retrieve data from sheet: %v", err)
	}

	// Print the values from the response
	if len(resp.Values) == 1 {
		val = make(map[int]string, len(resp.Values[0]))
		for num, cell := range resp.Values[0] {
			val[num + 1] = fmt.Sprint(cell)
		}
	} else {
		return nil, fmt.Errorf("Retrived %v rows, must be 1", len(resp.Values))
	}

	return val, nil
}

// write values (update) in row
func (gDoc *GDocObj) UpdateRowVal(rNum uint, val map[int]string) error {
	spreadSheetId := gDoc.sprdSheetId
	srv := gDoc.svc

	ctx := context.Background()

	rb := &sheets.BatchUpdateValuesRequest{
		ValueInputOption: "USER_ENTERED",
	}

	for k, v := range val {
		rb.Data = append(rb.Data, &sheets.ValueRange{
			Range:  fmt.Sprintf("%s!%s%v", gDoc.sheetName, string(rune('A'-1+k)), rNum),
			Values: [][]interface{}{{v}},
		})
	}

	if _, err := srv.Spreadsheets.Values.BatchUpdate(spreadSheetId, rb).Context(ctx).Do(); err != nil {
		return err
	}

	return nil
}

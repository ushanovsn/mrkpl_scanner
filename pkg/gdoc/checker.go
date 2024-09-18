package gdoc

import (
	"regexp"
	"strconv"
)

// Checking URL for correctness and parsing values
func parseSheetURL(urlIncom string) (url string, sSheetId string, sheetId int64, ok bool) {
	ok = false

	// extract URL
	re := regexp.MustCompile(`docs.google.com/spreadsheets/d/\S+/edit\S*#gid=\d+\S*`)
	url = re.FindString(urlIncom)
	if url == "" {
		return url, sSheetId, sheetId, ok
	}

	// extract SpreadSheetId
	re = regexp.MustCompile(`docs.google.com/spreadsheets/d/(\S+)/edit\S*#gid=\d+\S*`)
	sSIdResult := re.FindStringSubmatch(urlIncom)
	if sSIdResult == nil {
		return url, sSheetId, sheetId, ok
	}
	sSheetId = sSIdResult[1]

	// extract SheetId
	re = regexp.MustCompile(`docs.google.com/spreadsheets/d/\S+/edit\S*#gid=(\d+)\S*`)
	sIdPreResult := re.FindStringSubmatch(urlIncom)
	if sIdPreResult == nil {
		return url, sSheetId, sheetId, ok
	}

	// parse SheetId to value
	if id, err := strconv.ParseInt(sIdPreResult[1], 10, 64); err != nil {
		return url, sSheetId, sheetId, ok
	} else {
		sheetId = id
	}

	// it's OK!
	ok = true

	return url, sSheetId, sheetId, ok
}

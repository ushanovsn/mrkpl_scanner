package ui

import (
	"fmt"
	"io"
	"mrkpl_scanner/internal/options"
	"net/http"
	"os"
)

// Processor for Index page "/"
func LogPageHndlr(scnr *options.ScannerObj, w http.ResponseWriter, r *http.Request) (int, error) {
	
	if r.Method != "GET" {
		return http.StatusMethodNotAllowed, fmt.Errorf("Method %s not allowed", r.Method)
	}

	tmpl := scnr.GetUIObj().GetUIHTMLTemplates()

	// struct for index page
	data := struct {
		Title		string
		NaviMenu	[]options.NaviMenu
		ActiveMenu	options.NaviActiveMenu
		Log			[]string
	}{
		Title:    options.DefUIPageTitle,
		NaviMenu: scnr.GetUIObj().GetUINaviMenu(),
		ActiveMenu: options.NaviActiveMenu{
			ActiveTabVal:    "/log",
			ActiveDMenuVal:  "",
			PageDescription: "Лог-файл",
		},
		Log: make([]string, 100),
	}

	data.Log = append(data.Log, "Тест лог")

	// base headers
	header := http.StatusOK
	w.Header().Add("Content-Type", "text/html; charset=utf-8")

	// OK. Processing template - replace\substitute values in template and send to front
	w.WriteHeader(header)

	return http.StatusOK, tmpl.ExecuteTemplate(w, "log", data)
}



func getLastLineWithSeek(filepath string) string {
    fileHandle, err := os.Open(filepath)

    if err != nil {
        return ""
    }
    defer fileHandle.Close()

    line := ""
    var cursor int64 = 0
    stat, _ := fileHandle.Stat()
    filesize := stat.Size()
    for { 
        cursor -= 1
        fileHandle.Seek(cursor, io.SeekEnd)

        char := make([]byte, 1)
        fileHandle.Read(char)

        if cursor != -1 && (char[0] == 10 || char[0] == 13) { // stop if we find a line
            break
        }

        line = fmt.Sprintf("%s%s", string(char), line) // there is more efficient way

        if cursor == -filesize { // stop if we are at the begining
            break
        }
    }

    return line
}



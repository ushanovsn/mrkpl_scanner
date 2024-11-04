package options

import (
	//"encoding/json"
)
// Data for Scan parameters page template
type ParamsScanPageData struct {
	// Saving parsed data in source
	SaveInOne		bool
	// Name of selected tab with params
	TabSrcActive	string
	// Source part params
	GParamSrc struct {
		GAuth		GoogleAuthParam
		TblParam	TableStruct
		ErrLog		[]string	`json:"-"`
	}
	FileSrc struct {
		File		FileStruct
		TblParam	TableStruct
		ErrLog		[]string	`json:"-"`
	}
	DBSrc struct {
		DB			DBStruct
		TblParam	TableStruct
		ErrLog		[]string	`json:"-"`
	}
	// Save part params
	TabSvActive		string
	GParamSv struct {
		GAuth		GoogleAuthParam
		TblParam	TableStruct
		ErrLog		[]string	`json:"-"`
	}
	FileSv struct {
		File		FileStruct
		TblParam	TableStruct
		ErrLog		[]string	`json:"-"`
	}
	DBSv struct {
		DB			DBStruct
		TblParam	TableStruct
		ErrLog		[]string	`json:"-"`
	}
	Starting struct {
		SelectName	string
		AtTime		string
		AtTimeOk	bool
		Periodic	int
		PeriodicOk	bool
		ErrLog		[]string	`json:"-"`
	}
}


// Parameters for Google Docs
type GoogleAuthParam struct {
	// URL of google document
	GPageURL        string
	// google doc accepted flag
	GPageURLOk		bool
	// authentication file
	AuthClient 		string
	// authentication file accepted flag
	AuthClientOk	bool
}

// Parameters of table structure
type TableStruct struct {
	StartRowNum			int
	StartRowNumOk		bool
	ColSourceId			string
	ColSourceIdOk		bool
	Params				[]ColParams
}

// Parameters for file
type FileStruct struct {
	FileName		string
	FileNameOk		bool
	FileCreateAllow	bool
	FileNewByTempl	bool
}

// Parameters for database
type DBStruct struct {
	DBType			string
	DBTypeOk		bool
	Addr			string
	AddrOk			bool
	Port			int
	PortOk			bool
	User			string
	UserOk			bool
	Pass			string
	PassOk			bool
	Database		string
	DatabaseOk		bool
	Table			string
	TableOk			bool
	UpdTable		bool
}

// One parameter of table param
type ColParams struct {
	ColParamValue		string
	ColParamType		string
}
/*
// Clear error messages when marshalling
func (p *ParamsScanPageData) MarshalJSON() ([]byte, error) {
	clr := *p
	clr.GParamSrc.ErrLog = nil
	clr.FileSrc.ErrLog = nil
	clr.DBSrc.ErrLog = nil
	clr.GParamSv.ErrLog = nil
	clr.FileSv.ErrLog = nil
	clr.DBSv.ErrLog = nil
	return json.Marshal(clr)
}
*/



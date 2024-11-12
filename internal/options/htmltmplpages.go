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


// Data for WB configuration page template
type ConfWBPageData struct {
	// Delay between requests
	RequestDelay	int
	// Delay between requests access flag
	RequestDelayOk	bool
	// Location addres identification
	AddressType		string
	// Special discount type
	DiscountType	string
	// Special discount manual value
	DiscountValue	float64
	// Special discount manual value access flag
	DiscountValueOk	bool
	// Errors list
	ErrLog			[]string	`json:"-"`
}


// Data for Index page template
type IndexPageData struct {
	// Enable status of Scanner service
	ScannerEnable	bool
	// Status and state Scanner service info
	ScannerSvcState	SeviceInfoData
}

type SeviceInfoData struct {
	AutoStartType	string
	CurrentState	string
	CurrentError	string
}
package options

import (
	"github.com/go-chi/chi/v5"
	"html/template"
	"net/http"
)

type UIObj struct {
	server   *http.Server
	router   chi.Router
	htmlNavi []NaviMenu
	tmpl     *template.Template
	//pData	*PagesData
}

// Get UI http server
func (obj *UIObj) GetUIServer() *http.Server {
	return obj.server
}

// Set UI http server to obj
func (obj *UIObj) SetUIServer(s *http.Server) {
	obj.server = s
}

// Get UI http router
func (obj *UIObj) GetUIRouter() chi.Router {
	return obj.router
}

// Set UI server to obj
func (obj *UIObj) SetUIRouter(r chi.Router) {
	obj.router = r
}

// Get UI Navigation menu
func (obj *UIObj) GetUINaviMenu() []NaviMenu {
	return obj.htmlNavi
}

// Set UI Navigation menu value
func (obj *UIObj) SetUINaviMenu(m []NaviMenu) {
	obj.htmlNavi = m
}

// Get UI compilled HTML templates
func (obj *UIObj) GetUIHTMLTemplates() *template.Template {
	return obj.tmpl
}

// Set UI compilled HTML templates
func (obj *UIObj) SetUIHTMLTemplates(t *template.Template) {
	obj.tmpl = t
}


// Aggregate and return actual values for page "page-task_param_scan"
func (obj *UIObj)GetPageTaskScanData(scnr *ScannerObj) *ParamsScanPageData {
	d := ParamsScanPageData{}

	/*
	d.SaveInOne = true
	d.TabSrcActive = ""
	d.TabSvActive = "save_tab_xls"
	d.GParamSrc.GAuth.AuthClient = "USER"
	d.GParamSrc.GAuth.AuthClientOk = true
	d.GParamSrc.TblParam.StartRowNum = 3
	d.GParamSrc.TblParam.StartRowNumOk = true
	d.GParamSrc.TblParam.Params = []ColParams {
		{
		ColParamValue: "A",
		ColParamType: "data_time",
		},
		{
		ColParamValue: "B",
		},
		{},
	}




	d.FileSrc.File.FileName = "test file.txt"
	d.FileSrc.File.FileNameOk = true
	d.FileSrc.TblParam.Params = []ColParams {
		{
		ColParamValue: "A",
		ColParamType: "data_time",
		},
		{
		ColParamValue: "D",
		},
		{},
	}



	d.DBSrc.DB.Addr = "localhost"
	d.DBSrc.DB.AddrOk = true
	d.DBSrc.TblParam.Params = []ColParams {
		{
		ColParamValue: "A",
		ColParamType: "data_time",
		},
		{
		ColParamValue: "D",
		},
		{},
	}
	
	d.DBSrc.ErrLog = []string{"Test","tttttttt", "dfgkjsfbglhs;glsab asdgfhsad igbhagfhfgh fdgh fj;dsfg npds;f j;dfgijdf p;ogjdsf; gj ;fosd;gf"}







	
	d.GParamSv.TblParam.StartRowNum = 3
	d.GParamSv.TblParam.StartRowNumOk = true
	d.GParamSv.TblParam.Params = []ColParams {
		{
		ColParamValue: "B",
		ColParamType: "brand",
		},
		{
		ColParamValue: "D",
		},
		{},
	}

	
	d.FileSv.File.FileCreateAllow = true
	d.FileSv.TblParam.Params = []ColParams {
		{
		ColParamValue: "B",
		ColParamType: "brand",
		},
		{
		ColParamValue: "D",
		},
		{},
	}

	d.FileSv.ErrLog = []string{"Test","tttttttt"}

	
	d.DBSv.DB.DBType = "pg"
	d.DBSv.DB.DBTypeOk = true
	d.DBSv.TblParam.Params = []ColParams {
		{
		ColParamValue: "B",
		ColParamType: "brand",
		},
		{
		ColParamValue: "D",
		},
		{},
	}
*/


	return &d
}
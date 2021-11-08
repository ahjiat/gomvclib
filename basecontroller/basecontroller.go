package basecontroller

import (
	"net/http"
	"html/template"
	"github.com/gorilla/mux"
)

type BaseController struct {
	Response http.ResponseWriter
	Request *http.Request
	ViewBasePath string
	ViewBag interface{}
	ActionName string
	viewTemplate *template.Template
}
func (self *BaseController) Echo(s string) {
	self.Response.Write([]byte(s))
}
func (self *BaseController) GetUrlVar(s string) string {
	return mux.Vars(self.Request)[s]
}
func (self *BaseController) View(fileName string) {
	var viewFile string
	if fileName == "" {
		viewFile = self.ViewBasePath + "/" + self.ActionName + ".html"
	} else {
		viewFile = self.ViewBasePath + "/" + fileName
	}
	if(self.viewTemplate == nil) {
		self.viewTemplate = template.New(viewFile)
	}
	/*
	t, err := template.New(viewFile).Parse(tpl)
	if err != nil {
		o("===== 1st error ======")
		o(err)
		return
	}
	*/
}

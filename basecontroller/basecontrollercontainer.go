package basecontroller

import (
	"net/http"
	"html/template"
	"github.com/gorilla/mux"
	"os"
	"path/filepath"
	"strings"
	"fmt"
	"bytes"
)

type BaseControllerContainerTemplate struct {
	tpl *template.Template
}
func (self *BaseControllerContainerTemplate) Append(header string, body string) *BaseControllerContainerTemplate {
	body = `{{define "`+strings.TrimSpace(header)+`"}}` + body + `{{end}}`
	if _, err := self.tpl.Parse(body); err != nil { panic(err) }
	return self
}

type BaseControllerContainer struct {
	Response http.ResponseWriter
	Request *http.Request
	ViewBasePath string
	ActionName string
	MasterTemplate *template.Template
	Templates *template.Template
	ViewRootPath string
	ViewBag interface{}
	ChainArgs []interface{}
	NeedNext bool
	OutChainArgs []interface{}
}
func (self *BaseControllerContainer) Echo(s string) {
	self.Response.Write([]byte(s))
}
func (self *BaseControllerContainer) GetUrlVar(s string) string {
	return mux.Vars(self.Request)[s]
}
func (self *BaseControllerContainer) GetView(inputData interface{}, fileNames... string) string {
	buff := self.getViewContent(inputData, fileNames...)
	return buff.String()
}
func (self *BaseControllerContainer) View(fileNames... string) {
	buff := self.getViewContent(self.ViewBag, fileNames...)
	self.Response.Write(buff.Bytes())
}
func (self *BaseControllerContainer) SetMasterView(fileNames... string) *BaseControllerContainerTemplate {
	var file, fileName string
	var dat []byte
	var err error
	if len(fileNames) != 0 { fileName = fileNames[0] }

	file, fileName = self.retriveAbsFile(fileName)

	if self.MasterTemplate == nil {
		self.MasterTemplate = template.New("master")
	}

	dat, err = os.ReadFile(file); if err != nil { panic(err) }
	_, err = self.MasterTemplate.Parse(string(dat)); if err != nil { panic(err) }

	return &BaseControllerContainerTemplate{self.MasterTemplate}
}
func (self *BaseControllerContainer) RouteNext(args... interface{}) {
	self.NeedNext = true
	self.OutChainArgs = args
}
func (self *BaseControllerContainer) ParseHtmlTemplate(inputData interface{}, content string) string {
	var output bytes.Buffer
	tpl, err := template.New("").Parse(content); if err != nil { panic(err) }
	err = tpl.Execute(&output, inputData); if err != nil { panic(err) }
	return output.String()
}

func (self *BaseControllerContainer) getViewContent(inputData interface{}, fileNames... string) bytes.Buffer {
	var file, fileName string
	var err error
	if len(fileNames) != 0 { fileName = fileNames[0] }

	file, fileName = self.retriveAbsFile(fileName)

	var tpl *template.Template
	var dat []byte;
	if(self.MasterTemplate != nil) {
		var data string
		tpl = self.MasterTemplate
		dat, err = os.ReadFile(file); if err != nil { panic(err) }
		if len(dat) != 0 { data = `{{define "webcontent"}}` + string(dat) + `{{end}}` }
		tpl, err = tpl.Parse(data); if err != nil { panic(err) }
	} else {
		tpl = self.Templates.Lookup(fileName)
		if tpl == nil {
			dat, err = os.ReadFile(file); if err != nil { panic(err) }
			tpl, err = self.Templates.New(fileName).Parse(string(dat)); if err != nil { panic(err) }
		}
	}
	var output bytes.Buffer
	err = tpl.Execute(&output, inputData)
	if err != nil { panic(err) }
	return output
}
func (self *BaseControllerContainer) retriveAbsFile(fileName string) (string,string) {
	var file string
	if strings.HasPrefix(fileName, "/") {
		file = self.ViewRootPath
	} else {
		file = self.ViewBasePath + "/"
	}
	if fileName == "" {
		fileName = self.ActionName + ".html"
	}
	file += fileName

	if absfile, _ := filepath.Abs(file); !strings.HasPrefix(absfile, self.ViewRootPath) {
		panic(fmt.Sprintf("filename %s must within %s ", fileName, self.ViewRootPath))
	}
	_, err := os.Stat(file); if err != nil { panic(err) }
	return file, fileName
}

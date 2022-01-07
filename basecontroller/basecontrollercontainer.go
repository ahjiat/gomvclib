package basecontroller

import (
	"net/http"
	"text/template"
	"github.com/gorilla/mux"
	"os"
	"path/filepath"
	"strings"
	"fmt"
	"bytes"
)

type BaseControllerContainerTemplate struct {
	masterTemplates	map[string]*template.Template
	tpl *template.Template
	viewRootPath string
	viewBasePath string
	actionName string
}
func (self *BaseControllerContainerTemplate) DefineTemplate(name string, inputData interface{}, fileNames... string) *BaseControllerContainerTemplate {
	var output bytes.Buffer
	var fileName string
	var ok bool
	var mt *template.Template
	var err error
	if len(fileNames) != 0 { fileName = fileNames[0] }
	file, fileName := self.retriveAbsFile(fileName)

	if mt, ok = self.masterTemplates[fileName]; ! ok {
		dat, err := os.ReadFile(file); if err != nil { panic(err) }
		mt, err = template.New(fileName).Delims("@{", "}").Parse(string(dat)); if err != nil { panic(err) }
		self.masterTemplates[fileName] = mt
	}
	err = mt.Execute(&output, inputData); if err != nil { panic(err) }
	return self.DefineTemplateByString(name, output.String())
}
func (self *BaseControllerContainerTemplate) DefineTemplateByString(name string, body string) *BaseControllerContainerTemplate {
	if _, err := self.tpl.New(name).Delims("@{", "}").Parse(body); err != nil { panic(err) }
	return self
}
func (self *BaseControllerContainerTemplate) retriveAbsFile(fileName string) (string,string) {
	var file string
	if strings.HasPrefix(fileName, "/") {
		file = self.viewRootPath
	} else {
		file = self.viewBasePath + "/"
	}
	if fileName == "" {
		fileName = self.actionName + ".html"
	}
	file += fileName

	if absfile, _ := filepath.Abs(file); !strings.HasPrefix(absfile, self.viewRootPath) {
		panic(fmt.Sprintf("filename %s must within %s ", fileName, self.viewRootPath))
	}
	fileName = strings.TrimPrefix(file, self.viewRootPath)
	return file, fileName
}

type BaseControllerContainer struct {
	Response http.ResponseWriter
	Request *http.Request
	ViewBasePath string
	ActionName string
	Templates map[string]*template.Template
	ViewRootPath string
	ViewBag interface{}
	ChainArgs []interface{}
	NeedNext bool
	OutChainArgs []interface{}
	ContainerTemplate *BaseControllerContainerTemplate
	MasterTemplates	map[string]*template.Template
	MasterTemplate **template.Template
	RoutePath string
	IRouteArgs []interface{}
}
func (self *BaseControllerContainer) Echo(value string, args ...interface{}) {
	if len(args) == 0 { self.Response.Write([]byte(value)); return }
	self.Response.Write([]byte(fmt.Sprintf(value, args...)))
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
func (self *BaseControllerContainer) MasterView(tplName string, inputData interface{}, fileNames... string) {
	mst, ok := self.GetMasterView(); if ! ok { return }
	mst.DefineTemplate(tplName, inputData, fileNames...)
	tpl := *self.MasterTemplate
	err := tpl.Execute(self.Response, inputData); if err != nil { panic(err) }
}
func (self *BaseControllerContainer) CreateMasterView(fileNames... string) *BaseControllerContainerTemplate {
	var fileName string
	if len(fileNames) != 0 { fileName = fileNames[0] }
	file, fileName := self.retriveAbsFile(fileName)
	var otpl *template.Template; var err error
	otpl, ok := self.MasterTemplates[fileName]
	if ! ok {
		data, err := os.ReadFile(file); if err != nil { panic(err) }
		otpl, err = template.New(fileName).Delims("@{", "}").Parse(string(data)); if err != nil { panic(err) }
		self.MasterTemplates[fileName] = otpl
	}
	*self.MasterTemplate, err = otpl.Clone(); if err != nil { panic(err) }
	self.ContainerTemplate = &BaseControllerContainerTemplate{self.MasterTemplates, *self.MasterTemplate, self.ViewRootPath, self.ViewBasePath, self.ActionName}
	return self.ContainerTemplate
}
func (self *BaseControllerContainer) RemoveMasterView() {
	*self.MasterTemplate = nil
}
func (self *BaseControllerContainer) GetMasterView() (*BaseControllerContainerTemplate, bool) {
	if *self.MasterTemplate == nil { return nil, false }
	self.ContainerTemplate = &BaseControllerContainerTemplate{self.MasterTemplates, *self.MasterTemplate, self.ViewRootPath, self.ViewBasePath, self.ActionName}
	return self.ContainerTemplate, true
}
func (self *BaseControllerContainer) RouteNext(args... interface{}) {
	self.NeedNext = true
	self.OutChainArgs = args
}
func (self *BaseControllerContainer) ParseTemplate(inputData interface{}, content string) string {
	var output bytes.Buffer
	tpl, err := template.New("").Delims("@{", "}").Parse(content); if err != nil { panic(err) }
	err = tpl.Execute(&output, inputData); if err != nil { panic(err) }
	return output.String()
}

func (self *BaseControllerContainer) getViewContent(inputData interface{}, fileNames... string) bytes.Buffer {
	var file, fileName string
	var err error
	var ok bool
	var output bytes.Buffer
	if len(fileNames) != 0 { fileName = fileNames[0] }

	file, fileName = self.retriveAbsFile(fileName)
	var tpl *template.Template

	if tpl, ok = self.Templates[fileName]; ! ok {
		data, err := os.ReadFile(file); if err != nil { panic(err) }
		tpl, err = template.New(fileName).Delims("@{", "}").Parse(string(data)); if err != nil { panic(err) }
		self.Templates[fileName] = tpl
	}

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
	fileName = strings.TrimPrefix(file, self.ViewRootPath)
	return file, fileName
}

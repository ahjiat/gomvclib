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
	tpl *template.Template
	viewRootPath string
	viewBasePath string
}
func (self *BaseControllerContainerTemplate) DefineTemplate(name string, inputData interface{}, fileName string) *BaseControllerContainerTemplate {
	var output bytes.Buffer
	if fileName == "" { panic("DefineTemplate filename can not empty") }
	file, fileName := self.retriveAbsFile(fileName)

	dat, err := os.ReadFile(file); if err != nil { panic(err) }
	data := string(dat)
	t, err := template.New("").Parse(data); if err != nil { panic(err) }
	err = t.Execute(&output, inputData); if err != nil { panic(err) }

	return self.DefineTemplateByString(name, output.String())
}
func (self *BaseControllerContainerTemplate) DefineTemplateByString(name string, body string) *BaseControllerContainerTemplate {
	if _, err := self.tpl.New(name).Parse(body); err != nil { panic(err) }
	return self
}
func (self *BaseControllerContainerTemplate) retriveAbsFile(fileName string) (string,string) {
	var file string
	if strings.HasPrefix(fileName, "/") {
		file = self.viewRootPath
	} else {
		file = self.viewBasePath + "/"
	}
	file += fileName

	if absfile, _ := filepath.Abs(file); !strings.HasPrefix(absfile, self.viewRootPath) {
		panic(fmt.Sprintf("filename %s must within %s ", fileName, self.viewRootPath))
	}
	_, err := os.Stat(file); if err != nil { panic(err) }
	fileName = strings.TrimPrefix(file, self.viewRootPath)
	return file, fileName
}

type BaseControllerContainer struct {
	Response http.ResponseWriter
	Request *http.Request
	ViewBasePath string
	ActionName string
	Templates *template.Template
	ViewRootPath string
	ViewBag interface{}
	ChainArgs []interface{}
	NeedNext bool
	OutChainArgs []interface{}
	MasterTemplate *template.Template
	ContainerTemplate *BaseControllerContainerTemplate
	MasterTemplateName **string
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
func (self *BaseControllerContainer) CreateMasterView(fileNames... string) *BaseControllerContainerTemplate {
	var fileName string
	if len(fileNames) != 0 { fileName = fileNames[0] }
	file, fileName := self.retriveAbsFile(fileName)
	tpl := self.MasterTemplate.Lookup(fileName)
	if tpl == nil {
		data, err := os.ReadFile(file); if err != nil { panic(err) }
		tpl, err = self.MasterTemplate.New(fileName).Parse(string(data)); if err != nil { panic(err) }
	}
	*self.MasterTemplateName = &fileName
	self.ContainerTemplate = &BaseControllerContainerTemplate{tpl, self.ViewRootPath, self.ViewBasePath}
	return self.ContainerTemplate
}
func (self *BaseControllerContainer) RemoveMasterView() {
	*self.MasterTemplateName = new(string)
}
func (self *BaseControllerContainer) GetMasterView() *BaseControllerContainerTemplate {
	if *self.MasterTemplateName == nil || self.MasterTemplate.Lookup(**self.MasterTemplateName) == nil {
		return nil
	}
	tpl := self.MasterTemplate.Lookup(**self.MasterTemplateName)
	self.ContainerTemplate = &BaseControllerContainerTemplate{tpl, self.ViewRootPath, self.ViewBasePath}
	return self.ContainerTemplate
}
func (self *BaseControllerContainer) RouteNext(args... interface{}) {
	self.NeedNext = true
	self.OutChainArgs = args
}
func (self *BaseControllerContainer) ParseTemplate(inputData interface{}, content string) string {
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
	if *self.MasterTemplateName != nil && self.MasterTemplate.Lookup(**self.MasterTemplateName) != nil {
		tpl = self.MasterTemplate.Lookup(**self.MasterTemplateName)
		if viewTpl := self.MasterTemplate.Lookup(fileName); viewTpl == nil {
			data, err := os.ReadFile(file); if err != nil { panic(err) }
			_, err = self.MasterTemplate.New("view").Parse(string(data)); if err != nil { panic(err) }
		}
	} else {
		tpl = self.Templates.Lookup(fileName)
		if tpl == nil {
			data, err := os.ReadFile(file); if err != nil { panic(err) }
			tpl, err = self.Templates.New(fileName).Parse(string(data)); if err != nil { panic(err) }
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
	fileName = strings.TrimPrefix(file, self.ViewRootPath)
	return file, fileName
}

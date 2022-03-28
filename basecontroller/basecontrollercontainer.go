package basecontroller

import (
	"net/http"
	"text/template"
	"github.com/ahjiat/gomvclib/global"
	"github.com/gorilla/mux"
	"os"
	"path/filepath"
	"strings"
	"fmt"
	"bytes"
	"errors"
	"syscall"
)

type BaseControllerContainerTemplate struct {
	masterTemplates	map[string]*template.Template
	tpl *template.Template
	viewRootPath string
	viewBasePath string
	actionName string
}
func (self *BaseControllerContainerTemplate) DefineTemplate(name string, args... interface{}) *BaseControllerContainerTemplate {
	var output bytes.Buffer
	var ok bool
	var t *template.Template
	var err error
	var fileName string; var dat interface{} = nil; var masterViewBag interface{} = nil
	if len(args) >= 1 { fileName = args[0].(string) }
	if len(args) >= 2 { dat = args[1] }
	if len(args) >= 3 { masterViewBag = args[2] }

	file, fileName := self.retriveAbsFile(fileName)
	if t, ok = self.masterTemplates[fileName]; ! ok {
		rawFile, err := os.ReadFile(file); if err != nil { panic(err) }
		t, err = template.New(fileName).Delims("@[", "]").Parse(string(rawFile)); if err != nil { panic(err) }
		self.masterTemplates[fileName] = t
	}
	err = t.Execute(&output, dat); if err != nil { panic(err) }
	return self.DefineTemplateByString(name, output.String(), masterViewBag)
}
func (self *BaseControllerContainerTemplate) defineTemplateCoreInternal(inputData interface{}, fileName string, loopLimitCount int, mDat interface{}) string {
	var output,output2 bytes.Buffer
	var ok bool
	var mt *template.Template
	var err error
	file, fileName := self.retriveAbsFile(fileName)

	if mt, ok = self.masterTemplates[fileName]; ! ok {
		dat, err := os.ReadFile(file); if err != nil { panic(err) }
		mt, err = template.New(fileName).Delims("@[", "]").Parse(string(dat)); if err != nil { panic(err) }
		self.masterTemplates[fileName] = mt
	}
	err = mt.Execute(&output, inputData); if err != nil { panic(err) }

	funcMap := template.FuncMap {
		"LoadFile": func(file string, datas ...interface{}) (string, error) {
			var data interface{}
			if len(datas) != 0 { data = datas[0] }
			if loopLimitCount >= 100 {
				return "", errors.New(`Error, infinity loop!, reached max recursive call to function "LoadFile"`)
			}
			return self.defineTemplateCoreInternal(data, file, loopLimitCount + 1, mDat), nil
		},
	}
	t, err := template.New("").Delims("@{", "}").Funcs(funcMap).Parse(output.String()); if err != nil { panic(err) }
	err = t.Execute(&output2, mDat); if err != nil { panic(err) }

	return output2.String()
}
func (self *BaseControllerContainerTemplate) DefineTemplateByString(name string, body string, mDat interface{}) *BaseControllerContainerTemplate {
	funcMap := template.FuncMap {
		"LoadFile": func(file string, datas ...interface{}) (string, error) {
			var data interface{}
			if len(datas) != 0 { data = datas[0] }
			return self.defineTemplateCoreInternal(data, file, 0, mDat), nil
		},
	}
	if _, err := self.tpl.New(name).Delims("@{", "}").Funcs(funcMap).Parse(body); err != nil { panic(err) }
	return self
}
func (self *BaseControllerContainerTemplate) retriveAbsFile(fileName string) (string,string) {
	var file string
	if strings.HasPrefix(fileName, global.SysPathSeparator) {
		file = self.viewRootPath
	} else {
		file = self.viewBasePath + global.SysPathSeparator
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
	MasterViewBag interface{}
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
func (self *BaseControllerContainer) GetView(fileNames... string) string {
	var fileName string
	if len(fileNames) >= 1 { fileName =  fileNames[0] }
	buff := self.getViewContent(fileName)
	return buff.String()
}
func (self *BaseControllerContainer) View(fileNames... string) {
	var fileName string
	if len(fileNames) >= 1 { fileName =  fileNames[0] }
	buff := self.getViewContent(fileName)
	self.Response.Write(buff.Bytes())
}
func (self *BaseControllerContainer) MasterView(args... interface{}) {
	mst, ok := self.GetMasterView(); if ! ok { return }
	var fileName string
	var dat interface{} = nil
	if len(args) >= 1 { fileName = args[0].(string) }
	if len(args) >= 2 { dat = args[1] }
	_, fileName = self.retriveAbsFile(fileName)
	mst.DefineTemplate(fileName, fileName, dat, self.MasterViewBag)
	tpl := *self.MasterTemplate
	if err := tpl.Execute(self.Response, self.MasterViewBag); err != nil  && ! errors.Is(err, syscall.EPIPE) { panic(err) }
}
func (self *BaseControllerContainer) CreateMasterTemplate(args... interface{}) *BaseControllerContainerTemplate {
	var fileName string; var err error
	var dat interface{} = nil
	if len(args) >= 1 { fileName = args[0].(string) }
	if len(args) >= 2 { dat = args[1] }

	// only for explicit declaration while Parse, chained by associated templates will overwite its functionality
	// while calling on MasterView, its FuncMap will be overwritten by DefineTemplateByString
	funcMap := template.FuncMap {
		"LoadFile": func() string { return "" },
	}
	_, fileName = self.retriveAbsFile(fileName)
	rawFile := self.defineMasterTemplateCore(dat, fileName)
	*self.MasterTemplate, err = template.New(fileName).Delims("@{", "}").Funcs(funcMap).Parse(string(rawFile)); if err != nil { panic(err) }
	self.ContainerTemplate = &BaseControllerContainerTemplate{self.MasterTemplates, *self.MasterTemplate, self.ViewRootPath, self.ViewBasePath, self.ActionName}
	return self.ContainerTemplate
}
func (self *BaseControllerContainer) RemoveMasterTemplate() {
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

func (self *BaseControllerContainer) getViewContent(fileName string) bytes.Buffer {
	var output bytes.Buffer
	var tpl *template.Template
	var ok bool
	file, fileName := self.retriveAbsFile(fileName)
	if tpl, ok = self.Templates[fileName]; ! ok {
		rawFile, err := os.ReadFile(file); if err != nil { panic(err) }
		funcMap := template.FuncMap {
			"LoadFile": func(file string) (string, error) {
				return self.defineViewTemplateCoreInternal(file, 0), nil
			},
		}
		tpl, err = template.New(fileName).Delims("@{", "}").Funcs(funcMap).Parse(string(rawFile)); if err != nil { panic(err) }
		self.Templates[fileName] = tpl
	}
	err := tpl.Execute(&output, self.ViewBag); if err != nil { panic(err) }
	return output
}
func (self *BaseControllerContainer) retriveAbsFile(fileName string) (string,string) {
	var file string
	if strings.HasPrefix(fileName, global.SysPathSeparator) {
		file = self.ViewRootPath
	} else {
		file = self.ViewBasePath + global.SysPathSeparator
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
func (self *BaseControllerContainer) defineMasterTemplateCore(inputData interface{}, fileName string) string {
	var output bytes.Buffer
	var ok bool
	var mt *template.Template
	var err error
	file, fileName := self.retriveAbsFile(fileName)
	if mt, ok = self.MasterTemplates[fileName]; ! ok {
		dat, err := os.ReadFile(file); if err != nil { panic(err) }
		mt, err = template.New(fileName).Delims("@[", "]").Parse(string(dat)); if err != nil { panic(err) }
		self.MasterTemplates[fileName] = mt
	}
	err = mt.Execute(&output, inputData); if err != nil { panic(err) }
	return output.String()
}
func (self *BaseControllerContainer) defineViewTemplateCoreInternal(fileName string, loopLimitCount int) string {
	var output bytes.Buffer
	var ok bool
	var tpl *template.Template
	var err error

	file, fileName := self.retriveAbsFile(fileName)
	if tpl, ok = self.Templates[fileName]; ! ok {
		rawFile, err := os.ReadFile(file); if err != nil { panic(err) }
		funcMap := template.FuncMap {
			"LoadFile": func(file string) (string, error) {
				if loopLimitCount >= 100 {
					return "", errors.New(`Error, infinity loop!, reached max recursive call to function "LoadFile"`)
				}
				return self.defineViewTemplateCoreInternal(file, loopLimitCount + 1), nil
			},
		}
		tpl, err = template.New(fileName).Delims("@{", "}").Funcs(funcMap).Parse(string(rawFile)); if err != nil { panic(err) }
		self.Templates[fileName] = tpl
	}
	err = tpl.Execute(&output, self.ViewBag); if err != nil { panic(err) }
	return output.String()
}

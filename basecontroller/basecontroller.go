package basecontroller

import (
	"net/http"
	"html/template"
	"github.com/gorilla/mux"
	"os"
	"path/filepath"
	"strings"
	"fmt"
)

type BaseController struct {
	Response http.ResponseWriter
	Request *http.Request
	ViewBasePath string
	ActionName string
	Templates *template.Template
	ViewRootPath string
	ViewBag interface{}
}
func (self *BaseController) Echo(s string) {
	self.Response.Write([]byte(s))
}
func (self *BaseController) GetUrlVar(s string) string {
	return mux.Vars(self.Request)[s]
}
func (self *BaseController) View(fileNames... string) {
	var file, fileName string
	if len(fileNames) != 0 { fileName = fileNames[0] }
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
	if _, err := os.Stat(file); err != nil {
		panic(err)
	}

	template := self.Templates.Lookup(fileName)
	if template == nil {
		dat, err := os.ReadFile(file)
		if err != nil { panic(err) }
		template, err = self.Templates.New(fileName).Parse(string(dat))
		if err != nil { panic(err) }
	}
	err := template.Execute(self.Response, self.ViewBag)
	if err != nil { panic(err) }
}

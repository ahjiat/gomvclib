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

	if file, _ := filepath.Abs(viewFile); !strings.HasPrefix(file, self.ViewRootPath) {
		panic(fmt.Sprintf("filename %s must within %s ", fileName, self.ViewRootPath))
	}
	if _, err := os.Stat(viewFile); err != nil {
		panic(err)
	}

	/*
	template := self.Templates.Lookup(viewFile)
	if template == nil {
		if _, err := os.Stat(viewFile); err != nil { panic(err) }
		tt1, err := self.Templates.New(viewFile).ParseFiles(viewFile)
		if err != nil { panic(err) }
		template = tt1
	}
	*/

	/*
	t1, err := self.Templates.New("Info.html").ParseFiles("/var/www/go/webserver/gomvc/view/Info/Info.html")
	err = t1.Execute(self.Response, nil)
	if err != nil { panic(err) }
	*/

	/*
	t1, err := template.ParseFiles("/var/www/go/webserver/gomvc/view/Info/Info.html")
	err = t1.Execute(self.Response, nil)
	if err != nil { panic(err) }
	*/

	/*
	a, _ := filepath.Abs("/var/www/go/webserver/gomvclib/../../")
	self.Echo(a + "\n")
	*/
}

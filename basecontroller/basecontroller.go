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
	if _, err := os.Stat(file); err != nil {
		panic(err)
	}

	template := self.Templates.Lookup(fileName)
	if template == nil {
		dat, err := os.ReadFile(file)
		if err != nil { panic(err) }
		template, err = self.Templates.New(fileName).Parse(string(dat))
		if err != nil { panic(err) }
		fmt.Println("\n ===== CRETED ==== \n")
	}
	err := template.Execute(self.Response, nil)
	if err != nil { panic(err) }
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

package basecontroller

import (
	"net/http"
	"github.com/gorilla/mux"
)

type BaseController struct {
	Response http.ResponseWriter
	Request *http.Request
}
func (self *BaseController) Echo(s string) {
	self.Response.Write([]byte(s))
}
func (self *BaseController) GetVar(s string) string {
	return mux.Vars(self.Request)[s]

}

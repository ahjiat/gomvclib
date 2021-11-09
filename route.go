package Web
import (
	"github.com/ahjiat/gomvclib/basecontroller"
	"github.com/ahjiat/gomvclib/global"
	"strings"
	"github.com/gorilla/mux"
	"os"
	"html/template"
)

type Route struct {
	muxRouter *mux.Router
	domains []string
	methods []string
	pathPrefix string
	pt paramtype
	viewDirName string
	viewDirPath string
	controllerDirName string
	controllerDirPath string
}
func (self *Route) SetViewDir(path string) *Route {
	name := path
	path = global.SysPath + "/" + path 
	if _, err := os.Stat(path); os.IsNotExist(err) { errorLog("SetViewPath directory [" + path + "] not exist") }
	newRoute := *self
	newRoute.viewDirPath = path
	newRoute.viewDirName = name
	return &newRoute
}
func (self *Route) SetControllerDir(path string) *Route {
	name := path
	path = global.SysPath + "/" + path 
	if _, err := os.Stat(path); os.IsNotExist(err) { errorLog("SetControllerPath directory [" + path + "] not exist") }
	newRoute := *self
	newRoute.controllerDirPath = path
	newRoute.controllerDirName = name
	return &newRoute
}
func (self *Route) AnyDomain() *Route {
	return self.Domains([]string{}...)
}
func (self *Route) Domains(domains...string) *Route {
	newRoute := *self
	newRoute.domains = domains
	return &newRoute
}
func (self *Route) PathPrefix(path string) *Route {
	newRoute := *self
	newRoute.pathPrefix = self.pathPrefix + path
	return &newRoute
}
func (self *Route) Methods(methods...string) *Route {
	newRoute := *self
	newRoute.methods = methods
	return &newRoute
}
func (self *Route) SupportParameters(in ...interface{}) *Route {
	newRoute := *self
	newRoute.pt = paramtype{}
	newRoute.pt.Process(in...)
	return &newRoute
}
func (self *Route) Route(routeConfig interface{}, icontroller interface{}) {
	var rc []RouteConfig
	switch routeConfig.(type) {
	case RouteConfig:
		rc = []RouteConfig{routeConfig.(RouteConfig)}
	case []RouteConfig:
		rc = routeConfig.([]RouteConfig)
	default:
		errorLog("Route: paramters not support %T", routeConfig)
		return
	}
	for _, row := range rc {
		path := self.pathPrefix + row.Path
		action := row.Action
		if !isMethodExist(&icontroller, action) {
			errorLog("Web.RouteConfig, path:%s controller:%T action:%s not found! ", path, icontroller, action)
		}
		if !isFieldExist(&icontroller, "Response") || !isFieldExist(&icontroller, "Request") {
			errorLog("Web.RouteConfig, controller:%T missing 'BaseController' ", icontroller)
		}
		get, post := retrieveMethodParams(&icontroller, action)
		handler := RouteHandler{
			muxRouter:  self.muxRouter,
			pt: self.pt,
			viewDirName: self.viewDirName,
			viewDirPath: self.viewDirPath,
			controllerDirName: self.controllerDirName,
			controllerDirPath: self.controllerDirPath,
			store: direction{
				&icontroller, &post, &get, &action,
				getBaseViewPath(&icontroller, self.controllerDirName, self.viewDirName),
				template.New("")},
		}
		handler.addMuxRoute(path, self.domains, self.methods)
	}
}
func (self *Route) RouteByController(path string, icontroller interface{}) {
	if !isFieldExist(&icontroller, "Response") || !isFieldExist(&icontroller, "Request") {
		errorLog("Web.RouteConfig, controller:%T missing 'BaseController' ", icontroller)
	}
	baseMethods := listAllMethods(new(basecontroller.BaseController))
	skipMethods := map[string]int{}
	for _, v := range baseMethods {
		skipMethods[v] = 1
	}
	methods := listAllMethods(icontroller)
	rc := []RouteConfig{}
	for _, name := range methods {
		if _, has := skipMethods[name]; has {
			continue
		}
		lowcase := path + "/" + strings.ToLower(name)
		rc = append(rc, RouteConfig{Path: lowcase, Action: name})
	}
	self.Route(rc, icontroller)
}

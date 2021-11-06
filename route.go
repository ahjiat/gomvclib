package Web

import (
	"github.com/ahjiat/gomvclib/basecontroller"
	"strings"
	"github.com/gorilla/mux"
)

type RouteConfig struct {
	Path   string
	Action string
}

type Route struct {
	muxRouter *mux.Router
	domains []string
	pathPrefix string
	mrht mainRouteHandlerType
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
func (self *Route) SupportParameters(in ...interface{}) *Route {
	newRoute := *self
	newRoute.mrht.pt.Process(in...)
	//newRoute.mrht.pt.display()
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
		self.mrht.addToRoute(path, icontroller, post, get, action, self.domains)
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
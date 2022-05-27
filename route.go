package Web
import (
	"github.com/ahjiat/gomvclib/basecontroller"
	"github.com/ahjiat/gomvclib/global"
	"strings"
	"github.com/gorilla/mux"
	"os"
	"text/template"
	"path/filepath"
)

var shareTemplates map[string]*template.Template = map[string]*template.Template{}

type Route struct {
	muxRouter *mux.Router
	domains []string
	methods []string
	pathPrefix string
	pt paramtype
	viewDirName string
	viewDirPath string
	routeChainConfig []RouteChainConfig
	ViewFuncMap template.FuncMap
}
func (self *Route) SetViewDir(path string) *Route {
	name := path
	path = filepath.Join(global.SysPath, path)
	if _, err := os.Stat(path); os.IsNotExist(err) { errorLog("SetViewPath directory [" + path + "] not exist") }
	newRoute := *self
	newRoute.viewDirPath = path
	newRoute.viewDirName = name
	return &newRoute
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
func (self *Route) SetViewFunc(funcMap map[string]any) *Route {
	newRoute := *self
	newRoute.ViewFuncMap = template.FuncMap{}
	for k, f := range funcMap { newRoute.ViewFuncMap[k] = f }
	return &newRoute
}
func (self *Route) AddViewFunc(funcMap map[string]any) *Route {
	newRoute := *self
	if newRoute.ViewFuncMap == nil {
		newRoute.ViewFuncMap = template.FuncMap{}
	}
	for k, f := range funcMap { newRoute.ViewFuncMap[k] = f }
	return &newRoute
}
func (self *Route) SupportParameters(in ...interface{}) *Route {
	newRoute := *self
	newRoute.pt = paramtype{}
	newRoute.pt.Process(in...)
	return &newRoute
}
func (self *Route) Use(actions interface{}, icontroller interface{}) *Route {
	newRoute := *self
	var rcc []RouteChainConfig
	switch actions.(type) {
		case string:
			rcc = append(rcc, RouteChainConfig{actions.(string), icontroller})
		case []string:
			for _, action := range actions.([]string) {
				rcc = append(rcc, RouteChainConfig{action, icontroller})
			}
		default:
			errorLog("web.RouteChain: parameters not support %T", actions)
	}
	for _, config := range rcc {
		if ! isMethodExist(&config.Controller, config.Action) {
			errorLog("web.RouteChain, Controller:%T Action:[%s] not found!", config.Controller, config.Action)
		}
		if ! isFieldExist(&config.Controller, "Base") {
			errorLog("web.RouteChain, controller:%T missing [BaseController] ", config.Controller)
		}
		newRoute.routeChainConfig = append(newRoute.routeChainConfig, config)
	}
	return &newRoute
}
func (self *Route) Route(routeConfig interface{}, icontroller interface{}, iRouteArgs ...interface{}) {
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
		mts := shareTemplates
		handler := RouteHandler {
			muxRouter:  self.muxRouter,
			mainHandle: self.createHandle(&action, icontroller, mts, iRouteArgs),
		}
		for i, _ := range self.routeChainConfig {
			config := self.routeChainConfig[i]
			handler.middlewareHandle = append(handler.middlewareHandle, self.createHandle(&config.Action, config.Controller, mts, iRouteArgs))
		}
		handler.addMuxRoute(path, self.domains, self.methods)
	}
}
func (self *Route) RouteByStaticDir(path string, dir string, defaultIndex string) {
	r := self.PathPrefix(path)
	if defaultIndex != "" {
		r.Route(RouteConfig{"/{n:.*}", "Process"}, new(basecontroller.ServeStaticDir), dir, global.SysPathSeparator + defaultIndex)
	} else {
		r.Route(RouteConfig{"/{n:.*}", "Process"}, new(basecontroller.ServeStaticDir), dir)
	}
}
func (self *Route) RouteByController(path string, icontroller interface{}) {
	if ! isFieldExist(&icontroller, "Base") {
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
		lowcase := filepath.Join(path, strings.ToLower(name))
		rc = append(rc, RouteConfig{Path: lowcase, Action: name})
	}
	self.Route(rc, icontroller)
}
func (self *Route) createHandle(action *string, icontroller interface{}, mts map[string]*template.Template, iRouteArgs []interface{}) *RouteHandle {
		if !isMethodExist(&icontroller, *action) {
			errorLog("Web.RouteConfig, controller:%T action:%s not found! ", icontroller, *action)
		}
		if ! isFieldExist(&icontroller, "Base") {
			errorLog("Web.RouteConfig, controller:%T missing 'BaseController' ", icontroller)
		}
		get, post := retrieveMethodParams(&icontroller, *action)
		cloneViewFuncMap := template.FuncMap{}
		if self.ViewFuncMap != nil { for k, f := range self.ViewFuncMap { cloneViewFuncMap[k] = f } }
		return &RouteHandle {
			pt: self.pt,
			viewDirName: self.viewDirName,
			viewDirPath: self.viewDirPath,
			store: direction{
				&icontroller, &post, &get, action,
				getBaseViewPath(&icontroller, self.viewDirName),
				mts, mts},
			routePath: self.pathPrefix,
			iRouteArgs: iRouteArgs,
			viewFuncMap: cloneViewFuncMap,
		}
}

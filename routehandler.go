package Web
import (
	"net/http"
	"github.com/ahjiat/gomvclib/global"
	"github.com/gorilla/mux"
	"reflect"
	"strconv"
	"strings"
	"github.com/ahjiat/gomvclib/basecontroller"
	"text/template"
)

type RouteHandle struct {
	pt paramtype
	viewDirName string
	viewDirPath string
	store direction
	routePath string
	iRouteArgs []interface{}
	viewFuncMap template.FuncMap
	attr global.Attribute
}

type RouteHandler struct {
	muxRouter *mux.Router
	mainHandle *RouteHandle
	middlewareHandle []*RouteHandle
}
func (self *RouteHandler) addMuxRoute(path string, domains []string, methods []string) {
	if(path == "") { return }
	var routes []*mux.Route

	for _, domain := range domains {
		if(domain == "") { continue }
		routes = append(routes, self.muxRouteExactly(path, self.mainRouteHandler).Host(domain))
	}

	if len(methods) > 0 {
		if len(routes) > 0 {
			for _, route := range routes {
				route.Methods(methods...)
			}
		} else {
			routes = append(routes, self.muxRouteExactly(path, self.mainRouteHandler).Methods(methods...))
		}
	}

	if len(routes) == 0 { self.muxRouteExactly(path, self.mainRouteHandler) }
}
func (self *RouteHandler) muxRouteExactly(path string, f func (http.ResponseWriter, *http.Request)) *mux.Route {
	return self.muxRouter.HandleFunc(path, f)
}
func (self *RouteHandler) muxRouteIgnoreSlash(path string, f func (http.ResponseWriter, *http.Request)) *mux.Route {
	return self.muxRouteExactly(path+"{n:\\/?}", f)
}
func (self *RouteHandler) mainRouteHandler(w http.ResponseWriter, r *http.Request) {
	var args []interface{}
	var isNext bool = true
	mtPtr := new(*template.Template)
	for i, _ := range self.middlewareHandle {
		args, isNext = self.callHandle(w, r, self.middlewareHandle[i], args, mtPtr)
		if ! isNext  { return }
	}
	self.callHandle(w, r, self.mainHandle, args, mtPtr)
}

func (self *RouteHandler) callHandle(w http.ResponseWriter, r *http.Request, handle *RouteHandle, chainArgs []interface{}, mtPtr **template.Template) ([]interface{}, bool) {
	store := handle.store
	va := reflect.ValueOf(*store.ptr)
	v := reflect.New(va.Type().Elem())

	instance := &basecontroller.BaseControllerContainer {
		Response: w,
		Request: r,
		ViewBasePath: store.viewBasePath,
		ActionName: *store.action,
		Templates: store.templates,
		ViewRootPath: handle.viewDirPath,
		ChainArgs: chainArgs,
		MasterTemplates: store.masterTemplates,
		MasterTemplate: mtPtr,
		RoutePath: handle.routePath,
		IRouteArgs: handle.iRouteArgs,
		ViewFuncMap: handle.viewFuncMap,
		Attribute: handle.attr,
	}
	v.Elem().FieldByName("Base").Set(reflect.ValueOf(interface{}(instance)))

	//field := v.Elem().FieldByName("ViewRootPath"); _ = field
    //reflect.NewAt(field.Type(), unsafe.Pointer(field.UnsafeAddr())).Elem().Set(reflect.ValueOf(interface{}(self.viewDirPath)))
    //reflect.NewAt(field.Type(), unsafe.Pointer(field.UnsafeAddr())).Elem().SetString(store.viewBasePath)

	method := v.MethodByName(*store.action);
	if method.Type().NumIn() == 0 {
		method.Call([]reflect.Value{})
		return instance.OutChainArgs, instance.NeedNext
	}
	paramt := method.Type().In(0)
	fields := reflect.New(paramt).Elem()
	if len(*store.get) != 0 {
		for name, t := range *store.get {
			var found bool
			var vals []string
			if t.IsArray {
				vals, found = r.URL.Query()[name+"[]"]
			}
			if ! found {
				vals, found = r.URL.Query()[name]
			}
			if ! found { continue }
			self.setmainRouteHandlerField("GET_", &name, &vals, &fields, &t.Value, handle.pt)
		}
	}
	if len(*store.post) != 0 {
		r.ParseForm();
		for name, t := range *store.post {
			var found bool
			var vals []string
			if t.IsArray {
				vals, found = r.PostForm[name+"[]"]
			}
			if ! found {
				vals, found = r.PostForm[name]
			}
			if ! found { continue }
			self.setmainRouteHandlerField("POST_", &name, &vals, &fields, &t.Value, handle.pt)
		}
	}
	method.Call([]reflect.Value{fields})
	return instance.OutChainArgs, instance.NeedNext
}

func (self *RouteHandler) setmainRouteHandlerField(mtd string, name *string, vals *[]string, fields *reflect.Value, t *string, pt paramtype) {
	var val *string = &(*vals)[0]
	switch *t {
		case "int":
			v, _ := strconv.ParseInt(*val, 10, 64)
			fields.FieldByName(mtd+*name).SetInt(v)
		case "*int":
			v1, _ := strconv.ParseInt(*val, 10, 64); v := int(v1)
			fields.FieldByName(mtd+*name).Set(reflect.ValueOf(&v))
		case "[]int":
			var vInts []int
			for _, s := range *vals {
				v, _ := strconv.Atoi(s)
				vInts = append(vInts, v)
			}
			fields.FieldByName(mtd+*name).Set(reflect.ValueOf(vInts))
		case "string":
			fields.FieldByName(mtd+*name).SetString(*val)
		case "*string":
			fields.FieldByName(mtd+*name).Set(reflect.ValueOf(val))
		case "[]string":
			fields.FieldByName(mtd+*name).Set(reflect.ValueOf(*vals))
		case "float64":
			v, _ := strconv.ParseFloat(*val, 64)
			fields.FieldByName(mtd+*name).SetFloat(v)
		case "float32":
			v, _ := strconv.ParseFloat(*val, 32)
			fields.FieldByName(mtd+*name).SetFloat(v)
		case "bool":
			fields.FieldByName(mtd+*name).SetBool(strings.ToLower(*val) == "true")
		default:
			st, ok := pt[*t];
			if !ok { return }
			va := reflect.ValueOf(*st.iparam)
			v := reflect.New(va.Type().Elem());
			v.MethodByName("Set").Call([]reflect.Value{ reflect.ValueOf(val) })
			if(st.isPtr) {
				fields.FieldByName(mtd+*name).Set( reflect.ValueOf(v.Interface()) )
			} else {
				fields.FieldByName(mtd+*name).Set( reflect.ValueOf(v.Elem().Interface()) )
			}
	}
}

package Web
import (
	"net/http"
	"github.com/gorilla/mux"
	"reflect"
	"strconv"
	"strings"
)

type mainRouteHandlerType struct {
	muxRouter *mux.Router
	pt paramtype
	store direction
}
func (self *mainRouteHandlerType) addMuxRoute(path string, domains []string, methods []string) {
	if(path == "") { return }
	if len(domains) == 0  && len(methods) == 0 {
		self.muxRouteExactly(path+"{n:\\/?}", self.mainRouteHandler)
		return
	}
	if len(domains) > 0 {
		for _, domain := range domains {
			if(domain == "") { continue }
			route := self.muxRouteExactly(path+"{n:\\/?}", self.mainRouteHandler).Host(domain)
			if(len(methods) > 0) { route.Methods(methods...)  }
		}
		return
	}
	if len(methods) > 0 { self.muxRouteExactly(path+"{n:\\/?}", self.mainRouteHandler).Methods(methods...)  }
}
func (self *mainRouteHandlerType) muxRouteExactly(path string, f func (http.ResponseWriter, *http.Request)) *mux.Route {
	return self.muxRouter.HandleFunc(path, f)
}
func (self *mainRouteHandlerType) muxRouteIgnoreSlash(path string, f func (http.ResponseWriter, *http.Request)) *mux.Route {
	return self.muxRouteExactly(path+"{n:\\/?}", f)
}
func (self *mainRouteHandlerType) mainRouteHandler(w http.ResponseWriter, r *http.Request) {
	store := self.store
	va := reflect.ValueOf(*store.ptr)
	v := reflect.New(va.Type().Elem())
	v.Elem().FieldByName("Response").Set(reflect.ValueOf(interface{}(w)))
	v.Elem().FieldByName("Request").Set(reflect.ValueOf(interface{}(r)))
	v.Elem().FieldByName("ViewBasePath").SetString(store.viewBasePath)
	v.Elem().FieldByName("ActionName").SetString(*store.action)
	method := v.MethodByName(*store.action);
	if method.Type().NumIn() == 0 {
		method.Call([]reflect.Value{})
		return
	}
	paramt := method.Type().In(0)
	fields := reflect.New(paramt).Elem()
	if len(*store.get) != 0 {
		for name, t := range *store.get {
			val := r.URL.Query().Get(name)
			if val == "" { continue }
			self.setmainRouteHandlerField("GET_", &name, &val, &fields, &t)
		}
	}
	if len(*store.post) != 0 {
		r.ParseForm();
		for name, t := range *store.post {
			val := r.PostFormValue(name)
			if val == "" { continue }
			self.setmainRouteHandlerField("POST_", &name, &val, &fields, &t)
		}
	}
	method.Call([]reflect.Value{fields})
}
func (self *mainRouteHandlerType) setmainRouteHandlerField(mtd string, name *string, val *string, fields *reflect.Value, t *string) {
	switch *t {
		case "int":
			v, _ := strconv.ParseInt(*val, 10, 64)
			fields.FieldByName(mtd+*name).SetInt(v)
		case "string":
			fields.FieldByName(mtd+*name).SetString(*val)
		case "float64":
			v, _ := strconv.ParseFloat(*val, 64)
			fields.FieldByName(mtd+*name).SetFloat(v)
		case "float32":
			v, _ := strconv.ParseFloat(*val, 32)
			fields.FieldByName(mtd+*name).SetFloat(v)
		case "bool":
			fields.FieldByName(mtd+*name).SetBool(strings.ToLower(*val) == "true")
		default:
			st, ok := self.pt[*t];
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

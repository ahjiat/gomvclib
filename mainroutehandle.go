package Web
import (
	"net/http"
	"github.com/gorilla/mux"
	"reflect"
	"strconv"
)

type direction struct { 
	ptr *interface{}
	post *http_method
	get *http_method
	action *string
}
type mainRouteHandlerType struct {
	muxRouter *mux.Router
	pt paramtype
	domains []string
	storage map[string]direction
}
func (self *mainRouteHandlerType) addToRoute(path string, icontroller interface{}, post http_method, get http_method, action string) {
	storage_name := getRouteName()
	self.storage[storage_name] = direction{&icontroller, &post, &get, &action}
	self.addMuxRoute(path, storage_name)
}
func (self *mainRouteHandlerType) addMuxRoute(path string, name string) {
	if(len(self.domains) == 0) {
		self.muxRouteExactly(path+"{n:\\/?}", self.mainRouteHandler).Name(name)
		return
	}
	for _, domain := range self.domains {
		if(domain == "") { continue }
		self.muxRouteExactly(path+"{n:\\/?}", self.mainRouteHandler).Name(name).Host(domain)
	}
}
func (self *mainRouteHandlerType) muxRouteExactly(path string, f func (http.ResponseWriter, *http.Request)) *mux.Route {
	return self.muxRouter.HandleFunc(path, f)
}
func (self *mainRouteHandlerType) muxRouteIgnoreSlash(path string, f func (http.ResponseWriter, *http.Request)) *mux.Route {
	return self.muxRouteExactly(path+"{n:\\/?}", f)
}
func (self *mainRouteHandlerType) mainRouteHandler(w http.ResponseWriter, r *http.Request) {
	if store, ok := self.storage[mux.CurrentRoute(r).GetName()]; ok {
		va := reflect.ValueOf(*store.ptr)
		v := reflect.New(va.Type().Elem())
		v.Elem().FieldByName("Response").Set(reflect.ValueOf(interface{}(w)))
		v.Elem().FieldByName("Request").Set(reflect.ValueOf(interface{}(r)))
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
}
func (self *mainRouteHandlerType) setmainRouteHandlerField(mtd string, name *string, val *string, fields *reflect.Value, t *string) {
	switch *t {
		case "int":
			v, _ := strconv.ParseInt(*val, 10, 64)
			fields.FieldByName(mtd+*name).SetInt(v)
		case "string":
			fields.FieldByName(mtd+*name).SetString(*val)
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

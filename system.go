package Web
import (
	"fmt"
	"reflect"
	"os"
	"regexp"
	"strings"
	"text/template"
	"path/filepath"
	"log"
)

type RouteConfig struct {
	Path   string
	Action string
}
type RouteChainConfig struct {
	Action string
	Controller interface{}
}

type direction struct {
	ptr *interface{}
	post *http_method
	get *http_method
	action *string
	viewBasePath string
	templates map[string]*template.Template
	masterTemplates map[string]*template.Template
}

type http_method map[string] struct {
	Value string
	IsArray bool
}

func errorLog(str string, msg ...interface{}) {
	log.Fatal(fmt.Sprintf("ERROR => " + str + " ", msg...))
	os.Exit(1)
}
func listAllMethods(icontroller interface{}) []string {
	vt := reflect.TypeOf(icontroller)
	rtn := []string{}
	for i := 0; i < vt.NumMethod(); i++ {
		rtn = append(rtn, vt.Method(i).Name)
	}
	return rtn
}
func isMethodExist(icontroller *interface{}, name string) bool {
	vt := reflect.TypeOf(*icontroller)
	_, ok := vt.MethodByName(name)
	return ok
}
func isFieldExist(icontroller *interface{}, name string) bool {
	vt := reflect.TypeOf(*icontroller).Elem();
	_, ok := vt.FieldByName(name);
	return ok
}
func getTypeName(icontroller interface{}) string {
	if t := reflect.TypeOf(icontroller); t.Kind() == reflect.Ptr {
        return t.Elem().Name()
    } else {
        return t.Name()
    }
}
func retrieveMethodParams(icontroller *interface{}, methodName string) (http_method, http_method) {
	method := reflect.ValueOf(*icontroller).MethodByName(methodName)
	getHttps := http_method{}
	postHttps := http_method{}
	for i := 0; i < method.Type().NumIn(); i++ {
		switch i {
			case 0:
				for j := 0; j < method.Type().In(i).NumField(); j++ {
					field := method.Type().In(i).Field(j)
					if ok, _ := regexp.MatchString("^POST_[a-zA-Z]+[a-zA-Z0-9]+$", field.Name); ok {
						key := field.Name[5:];
						value := field.Type.String();
						isArray := strings.Contains(value, "[]")
						postHttps[key] = struct { Value string; IsArray bool } {value, isArray}
						continue
					}
					if ok, _ := regexp.MatchString("^GET_[a-zA-Z]+[a-zA-Z0-9]+$", field.Name); ok {
						key := field.Name[4:];
						value := field.Type.String();
						isArray := strings.Contains(value, "[]")
						getHttps[key] = struct { Value string; IsArray bool } {value, isArray}
						continue
					}
				}
		}
	}
	return getHttps, postHttps
}
func getBaseViewPath(icontroller *interface{}, viewName string, viewDirPath string) string {
	path := fmt.Sprintf("%T", *icontroller)
	path = strings.Replace(path, "*", "", 1)
	path = path[ strings.LastIndex(path, ".")+1: ]
	path = filepath.Join(filepath.Dir(viewDirPath), viewName, path)
	return path
}

package Web
import (
	"fmt"
	"reflect"
	"os"
	"regexp"
	"strconv"
)

var muxRouteCount int
type http_method map[string]string


func getRouteName() string {
	muxRouteCount++
	return strconv.Itoa(muxRouteCount)
}
func errorLog(str string, msg ...interface{}) {
	fmt.Fprintf(os.Stderr, "ERROR => " + str + " ", msg...)
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
						postHttps[key] = value
						continue
					}
					if ok, _ := regexp.MatchString("^GET_[a-zA-Z]+[a-zA-Z0-9]+$", field.Name); ok {
						key := field.Name[4:];
						value := field.Type.String();
						getHttps[key] = value
						continue
					}
				}
		}
	}
	return getHttps, postHttps
}

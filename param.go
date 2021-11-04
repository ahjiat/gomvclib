package Web
import (
	"fmt"
	"strings"
)
type param_container struct {
	iparam *interface{}
	isPtr bool
}
type paramtype map[string]param_container
func (self paramtype) Process(in ...interface{}) {
	for i := 0; i < len(in); i++ {
		name := fmt.Sprintf("%T", in[i])
		if !isMethodExist(&in[i], "Set") {
			errorLog("SupportParamTypes => '%s' missing 'Set' method", name)
		}
		self[name] =  param_container{ &in[i], true }
		name = strings.Replace(name, "*", "", 1)
		self[name] =  param_container{ &in[i], false }
	}
}
func (self paramtype) display() {
	for key, s := range self {
		fmt.Printf("%s => %T\n", key, *s.iparam)
	}
}

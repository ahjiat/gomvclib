package global

import (
	"os"
	"fmt"
)

func init() {
	curPath, err := os.Getwd();
	if err != nil { panic(err) }
	SysPath = curPath
	SysPathSeparator = fmt.Sprintf("%c", os.PathSeparator)
}

type Attribute struct {
	Message string
	IData interface{}
}

var SysPath string
var SysPathSeparator string

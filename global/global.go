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

var SysPath string
var SysPathSeparator string

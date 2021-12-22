package global

import (
	"os"
)

func init() {
	curPath, err := os.Getwd();
	if err != nil { panic(err) }
	SysPath = curPath
}

var SysPath string

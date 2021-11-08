package Web

import (
	"github.com/gorilla/mux"
	"os"
)

var sysCurPath__ string

func Router() (*Route, *mux.Router) {
	muxRouter := mux.NewRouter()
	route := Route { muxRouter: muxRouter, pt: paramtype{} }
	systemInformation := mainRouteHandlerType{ muxRouter:  muxRouter }
	sys := systeminfo{muxRouter}
	systemInformation.muxRouteExactly("/system/information", sys.pageShowRouteInfoHandler)

	if sysCurPath__ == "" {
		curPath, err := os.Getwd(); 
		if err != nil { panic(err) }
		sysCurPath__ = curPath
	}
	return &route, muxRouter
}

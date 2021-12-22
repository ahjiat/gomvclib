package Web

import (
	"github.com/gorilla/mux"
	"github.com/ahjiat/gomvclib/global"
	"fmt"
)

func Router() (*Route, *mux.Router) {
	muxRouter := mux.NewRouter()
	route := Route { muxRouter: muxRouter, pt: paramtype{} }
	systemInformation := RouteHandler{ muxRouter:  muxRouter }
	sys := systeminfo{muxRouter}
	systemInformation.muxRouteExactly("/system/information", sys.pageShowRouteInfoHandler)

	if global.SysPath == "" { panic("global.SysPath Null") }
	return &route, muxRouter
}

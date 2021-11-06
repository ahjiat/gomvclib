package Web

import (
	"github.com/gorilla/mux"
)

func Router() (*Route, *mux.Router) {
	muxRouter := mux.NewRouter()
	route := Route { muxRouter: muxRouter, pt: paramtype{} }
	systemInformation := mainRouteHandlerType{ muxRouter:  muxRouter }
	sys := systeminfo{muxRouter}
	systemInformation.muxRouteExactly("/system/information", sys.pageShowRouteInfoHandler)
	return &route, muxRouter
}



package Web

import (
	"github.com/gorilla/mux"
)

func Router() (*Route, *mux.Router) {
	muxRouter := mux.NewRouter()
	route := Route {
		muxRouter: muxRouter,
		mrht: mainRouteHandlerType{muxRouter: muxRouter, pt: paramtype{}, storage: map[string]direction{}},
	}
	sys := systeminfo{muxRouter}
	route.mrht.muxRouteExactly("/system/information", sys.pageShowRouteInfoHandler)
	return &route, muxRouter

}



package basecontroller

import (
	"path/filepath"
	"net/http"
	"os"
	"strings"
	"fmt"
)

type ServeStaticDir struct{ BaseController }

func (self *ServeStaticDir) Process() {
	r := self.Base.Request
	w := self.Base.Response
	path, err := filepath.Abs(r.URL.Path)
	if err != nil {
		http.Error(w, err.Error(), http.StatusBadRequest)
		return
	}

	dir := self.Base.IRouteArgs[0].(string)
	baseFile := strings.TrimPrefix(path, self.Base.RoutePath)
	file := baseFile

	if file == "" && len(self.Base.IRouteArgs) == 2 {
		file = self.Base.IRouteArgs[1].(string)
	}

	file = dir + file
	_, err = os.Stat(file)
	if os.IsNotExist(err) {
		http.Error(w, fmt.Sprintf("The requested URL %v was not found on this server.", baseFile), http.StatusNotFound)
		return
	} else if err != nil {
		http.Error(w, err.Error(), http.StatusInternalServerError)
		return
	}
	http.ServeFile(w, r, file)
}

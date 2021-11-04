package Web
import (
	"fmt"
	"net/http"
	"github.com/gorilla/mux"
)

type systeminfo struct {
	r *mux.Router
}
func (s *systeminfo) pageShowRouteInfoHandler(w http.ResponseWriter, r *http.Request) {
	user, pass, ok := r.BasicAuth()
	if !ok || !(user == "user" && pass == "1234") {
		w.Header().Set("WWW-Authenticate", `Basic realm="restricted", charset="UTF-8"`)
		http.Error(w, "Unauthorized login", http.StatusUnauthorized)
		return
	}
	msg := s.getRouteInfo()
	w.Write([]byte(`
		<html>
		<body>
	`));
	w.Write([]byte(msg));
	w.Write([]byte(`
		</body>
		</html>
	`));
}
func (s *systeminfo) getRouteInfo() string {
	var msg string
	s.r.Walk(func(route *mux.Route, router *mux.Router, ancestors []*mux.Route) error {
		host, _ := route.GetHostTemplate();
		if  host == "" { host = ".*" }
		path, _ := route.GetPathTemplate()
		name := route.GetName();
		msg += fmt.Sprintf("<div>%s, [%s][%s]</div>", name, host, path)
		return nil
	})
	return msg
}

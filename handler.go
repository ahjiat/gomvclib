package Web
import (
	"net/http"
)

func PageNotFoundHandler(w http.ResponseWriter, r *http.Request) {
    w.WriteHeader(http.StatusNotFound)
	w.Write([]byte("Page Not Found"));
}



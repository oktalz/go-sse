package server

import (
	"net/http"
	"strings"
)

func (s *Server) ServeHTTP(w http.ResponseWriter, r *http.Request) {
	path := strings.Split(r.URL.Path, "/")
	if len(path) < 2 {
		http.Error(w, "", http.StatusNotFound)
		return
	}
	//log.Println(path, r.URL.RawQuery)
	var requestType string
	if len(path) > 2 {
		requestType = path[len(path)-3]
		if requestType != "bind" {
			requestType = path[len(path)-1]
		}
	}

	switch requestType {
	case "sse":
		switch r.Method {
		case http.MethodGet:
			s.sse.HandlerGet(w, r)
		case http.MethodPost: // POST == GET
			s.sse.HandlerGet(w, r)
		case http.MethodDelete:
			s.sse.HandlerDelete(w, r)
		case http.MethodPatch:
			s.sse.HandlerPatch(w, r)
		default:
			http.Error(w, "", http.StatusMethodNotAllowed)
		}
	case "bind":
		functionName := path[len(path)-1]
		structName := path[len(path)-2]
		var args []string
		if r.URL.RawQuery != "" {
			data := strings.Split(r.URL.RawQuery, "&")
			for _, val := range data {
				d := strings.SplitN(val, "=", 2)
				if len(d) < 2 {
					http.Error(w, "", http.StatusNotAcceptable)
					return
				}
				args = append(args, d[1])
			}
		}
		s.bind.Serve(w, r, structName, functionName, args...)
	default:
		http.Error(w, "", http.StatusNotFound)
	}
}

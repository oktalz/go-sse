func (s *Server) ServeHTTP(w http.ResponseWriter, r *http.Request) {

	path := strings.Split(r.URL.Path, "/")

	functionName := path[len(path)-1] // HL
	structName := path[len(path)-2]   // HL

	var args []string
	if r.URL.RawQuery != "" {
		data := strings.Split(r.URL.RawQuery, "&")
		for _, val := range data {
			d := strings.SplitN(val, "=", 2)
			if len(d) < 2 {
				http.Error(w, "", http.StatusNotAcceptable)
				return
			}
			args = append(args, d[1]) // HL
		}
	}
	s.bind.Serve(w, r, structName, functionName, args...)
	// ..
}

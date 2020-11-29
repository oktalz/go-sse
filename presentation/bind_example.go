func main() {
	//github.com/oktalz/go-sse/server
	srv := server.New(server.ServerOptions{ // HL
		Path: "/api", // HL
	}) // HL
	// srv => http.Handler

	err := srv.Bind(&RandomStruct{}) // HL

	srv := &http.Server{
		Addr:    ":3042",
		Handler: srv,
	}
	_ = srv.ListenAndServe() // available at api/bind/RandomStruct/GetRandom
}

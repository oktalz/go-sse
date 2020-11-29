type RandomStruct struct {
	Event string
}

func (r *RandomStruct) Init(server *server.Server) error {
	go func() {
		for {
			server.Emit(r.Event, r.GetRandom()) // HL
			time.Sleep(4 * time.Second)         // simulate work
		}
	}()
	return nil
}

func main() {

	// history of 5 messages available
	server.AddEvent("randomstruct", 5) // HL
	// api/sse?sub=randomstruct

	_ = server.Bind(&RandomStruct{Event: "randomstruct"}) // HL
}

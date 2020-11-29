func main() {

	server.AddEvent("time", 5) // history of 5 messages available
	server.AddEvent("status", 1) // HL

	go func() {
		for {
			server.Emit("time", time.Now())
			time.Sleep(42 * time.Second)
		}
	}()

	server.Emit("status", Status{ // HL
		Accepted: 98, // HL
		Rejected: 2 // HL
	}) // HL

	// ...
}

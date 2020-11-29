for {

	// Read from our message
	msg, open := <-message

	if !open {
		// If our messageChan was closed, this means that the client has
		// disconnected.
		break
	}

	// Write to the ResponseWriter, `w`.
	fmt.Fprintf(w, "data: %s\n\n", msg)

	// Flush the response.  This is only possible if
	// the response supports streaming.
	f.Flush()
}

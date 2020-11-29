package sse

import (
	"fmt"
	"math/rand"
	"net/http"
	"strings"
	"time"

	jsoniter "github.com/json-iterator/go"
	"github.com/oklog/ulid"
	"github.com/oktalz/go-sse/history"
)

func (s *sse) HandlerGet(w http.ResponseWriter, r *http.Request) {
	w.Header().Set("Access-Control-Allow-Origin", "*")
	sub := r.URL.Query().Get("sub")
	subs := strings.Split(sub, ",")

	f, ok := w.(http.Flusher)
	if !ok {
		http.Error(w, "Streaming unsupported!", http.StatusInternalServerError)
		return
	}

	t := time.Now()
	entropy := rand.New(rand.NewSource(t.UnixNano())) //nolint:gosec
	id := ulid.MustNew(ulid.Timestamp(t), entropy).String()

	messageChan := make(chan string)
	client := Client{
		id:            id,
		data:          messageChan,
		subscriptions: subs,
	}

	s.newClients <- client

	ctx := r.Context()
	go func() {
		<-ctx.Done()
		s.removeClients <- client
		// log.Println("HTTP connection just closed.")
	}()

	// Set the headers related to event streaming.
	w.Header().Set("Content-Type", "text/event-stream")
	w.Header().Set("Cache-Control", "no-cache")
	w.Header().Set("Connection", "keep-alive")

	var json = jsoniter.ConfigCompatibleWithStandardLibrary
	idResult, _ := json.Marshal(&history.Message{
		Event: "ID",
		Data:  id,
	})
	fmt.Fprintf(w, "data: %s\n\n", string(idResult))
	f.Flush()

	for {

		// Read from our messageChan.
		msg, open := <-messageChan

		if !open {
			// If our messageChan was closed, this means that the client has
			// disconnected.
			break
		}

		// Write to the ResponseWriter, `w`.
		fmt.Fprintf(w, "data: %s\n\n", msg)
		fmt.Println("SEND", msg)

		// Flush the response.  This is only possible if
		// the response supports streaming.
		f.Flush()
	}

	// Done.
	// log.Printf("Finished HTTP request at %s for %s\n", r.URL.Path, client.id)
}

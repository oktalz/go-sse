package sse

import (
	"log"
	"net/http"
	"strings"

	jsoniter "github.com/json-iterator/go"
	"github.com/oktalz/go-sse/history"
)

func (s *sse) HandlerDelete(w http.ResponseWriter, r *http.Request) {
	w.Header().Set("Access-Control-Allow-Origin", "*")

	sub := r.URL.Query().Get("sub")
	subs := strings.Split(sub, ",")

	id := r.URL.Query().Get("id")

	w.Header().Set("Content-Type", "application/json")
	w.Header().Set("Cache-Control", "no-cache")
	// w.Header().Set("Connection", "keep-alive")
	log.Println(id, subs)

	messageChan := make(chan string)
	client := Client{
		id:            id,
		data:          messageChan,
		subscriptions: subs,
	}

	s.deleteSubscriptions <- client
	response := <-messageChan

	var json = jsoniter.ConfigCompatibleWithStandardLibrary
	idResult, _ := json.Marshal(&history.Message{
		Event: "DELETE_SUB",
		Data:  response,
	})

	_, _ = w.Write(idResult)
	// Done.
	// log.Println("Finished HTTP request at ", r.URL.Path, id)
}

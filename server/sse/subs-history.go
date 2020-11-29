package sse

import (
	"log"
	"net/http"
	"strconv"
)

func (s *sse) EventGet(w http.ResponseWriter, r *http.Request) {
	w.Header().Set("Access-Control-Allow-Origin", "*")

	event := r.URL.Query().Get("event")

	urlFromStr := r.URL.Query().Get("from")
	urlToStr := r.URL.Query().Get("to")
	id := r.URL.Query().Get("id")

	urlFrom, err := strconv.Atoi(urlFromStr)
	if err != nil {
		http.Error(w, "", http.StatusBadRequest)
	}
	urlTo, err := strconv.Atoi(urlToStr)
	if err != nil {
		http.Error(w, "", http.StatusBadRequest)
	}

	w.Header().Set("Cache-Control", "no-cache")
	// w.Header().Set("Connection", "keep-alive")
	log.Println(id, urlFrom, urlTo, event)

	msg := MessageBulk{
		ClientID: id,
		Event:    event,
		From:     urlFrom,
		To:       urlTo,
	}
	s.messageBulk <- msg
	// Done.
	// log.Println("Finished HTTP request at ", r.URL.Path, id)
}

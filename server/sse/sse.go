package sse

import (
	"net/http"

	"github.com/oktalz/go-sse/history"
)

type Subscription struct {
	clients map[chan string]struct{}
}

type MessageBulk struct {
	ClientID string
	Event    string
	From     int
	To       int
}

type MessageWithID struct {
	ID    int         `json:"id"`
	Event string      `json:"event"`
	Data  interface{} `json:"data"`
}

type Client struct {
	id            string
	data          chan string
	subscriptions []string
}

type SSE interface {
	HandleClients()
	AddEvent(event string, queue int)
	Emit(event string, data interface{})
	HandlerGet(w http.ResponseWriter, r *http.Request)
	HandlerDelete(w http.ResponseWriter, r *http.Request)
	HandlerPatch(w http.ResponseWriter, r *http.Request)
}

type sse struct {
	clients             map[string]Client
	newClients          chan Client
	patchClients        chan Client
	deleteSubscriptions chan Client
	removeClients       chan Client
	history             map[string]history.Simple
	messages            chan history.Message
	messageBulk         chan MessageBulk
	eventsHistory       map[string]*history.Simple
	subscriptions       map[string]Subscription
}

func New() SSE {
	return &sse{
		subscriptions:       map[string]Subscription{},
		clients:             map[string]Client{},
		newClients:          make(chan Client),
		patchClients:        make(chan Client),
		deleteSubscriptions: make(chan Client),
		removeClients:       make(chan Client),
		history:             make(map[string]history.Simple),
		messages:            make(chan history.Message),
		messageBulk:         make(chan MessageBulk),
		eventsHistory:       map[string]*history.Simple{},
	}
}

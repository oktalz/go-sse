package history

import "fmt"

var ErrNoMessages = fmt.Errorf("no messages")

type Message struct {
	Event string      `json:"event"`
	Data  interface{} `json:"data"`
}

type MessageWithID struct {
	ID    int64       `json:"id"`
	Event string      `json:"event"`
	Data  interface{} `json:"data"`
}

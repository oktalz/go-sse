package sse

import (
	"fmt"

	jsoniter "github.com/json-iterator/go"
	"github.com/oktalz/go-sse/history"
)

func (s *sse) AddEvent(event string, queue int) {
	s.subscriptions[event] = Subscription{
		clients: make(map[chan string]struct{}),
	}
	s.eventsHistory[event] = &history.Simple{
		MaxSize: queue,
	}
}

func (s *sse) Emit(event string, data interface{}) {
	go func() {
		s.messages <- history.Message{
			Event: event,
			Data:  data,
		}
	}()
}

func (s *sse) HandleClients() {

	go func() {
		var json = jsoniter.ConfigCompatibleWithStandardLibrary
		for {
			select {
			case newClient := <-s.newClients:
				s.clients[newClient.id] = newClient
				for _, subscritpion := range newClient.subscriptions {
					sub, ok := s.subscriptions[subscritpion]
					if ok {
						sub.clients[newClient.data] = struct{}{}
					}
				}
				// log.Printf("Added new client %s\n", newClient.id)

			case patchClient := <-s.patchClients:

				client, ok := s.clients[patchClient.id]
				if !ok {
					patchClient.data <- "ERR - NO CLIENT WITH ID " + patchClient.id
					continue
				}
				for _, subscritpion := range patchClient.subscriptions {
					sub, ok := s.subscriptions[subscritpion]
					if ok {
						_, subscribed := sub.clients[client.data]
						if !subscribed {
							sub.clients[client.data] = struct{}{}
						}
					}
				}
				patchClient.data <- "OK"
				// log.Printf("patched client %s\n", client.id)

			case deleteSubs := <-s.deleteSubscriptions:

				client, ok := s.clients[deleteSubs.id]
				if !ok {
					deleteSubs.data <- "ERR - NO CLIENT WITH ID " + deleteSubs.id
					continue
				}
				for _, subscritpion := range deleteSubs.subscriptions {
					sub, ok := s.subscriptions[subscritpion]
					if ok {
						_, subscribed := sub.clients[client.data]
						if subscribed {
							delete(sub.clients, client.data)
						}
					}
				}
				deleteSubs.data <- "OK"
				// log.Printf("patched (DELETE subs) client %s\n", client.id)

			case client := <-s.removeClients:
				for _, subscritpion := range client.subscriptions {
					sub, ok := s.subscriptions[subscritpion]
					if ok {
						delete(sub.clients, client.data)
					}
				}
				select {
				case <-client.data:
				default:
					close(client.data)
				}
				delete(s.clients, client.id)
				// log.Printf("Removed client %s\n", client.id)

			case msg := <-s.messages:

				history, ok := s.eventsHistory[msg.Event]
				if !ok {
					continue
				}
				id, _ := history.Add(msg)
				messageWithID := MessageWithID{
					ID:    id,
					Event: msg.Event,
					Data:  msg.Data,
				}
				sub, ok := s.subscriptions[msg.Event]
				if !ok {
					break
				}
				result, err := json.Marshal(&messageWithID)
				if err != nil {
					break
				}
				message := string(result)

				for s := range sub.clients {
					s <- message
					res := message
					if len(res) > 90 {
						res = fmt.Sprintf("%s...+[%d]", res[0:59], len(res)-90)
					}
					// log.Println(msg.Event, " => ", res)
				}

			case bulk := <-s.messageBulk:
				client, ok := s.clients[bulk.ClientID]
				if !ok {
					continue
				}
				evHistory, ok := s.eventsHistory[bulk.Event]
				if !ok {
					continue
				}
				data := []MessageWithID{}
				for i := bulk.From; i < bulk.To; i++ {
					message, err := evHistory.Get(i)
					if err != nil {
						continue
					}
					data = append(data, MessageWithID{
						ID:    i,
						Event: bulk.Event,
						Data:  message.Data,
					})
				}
				msg := history.Message{
					Event: "BULK",
					Data:  data,
				}

				result, err := json.Marshal(&msg)
				if err != nil {
					break
				}
				client.data <- string(result)
			}
		}
	}()
}

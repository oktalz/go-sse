msg := <-s.messages: // HL

id, _ := history.Add(msg)
messageWithID := MessageWithID{
	ID:    id,
	Event: msg.Event,
	Data:  msg.Data,
}
result, err := json.Marshal(&messageWithID)
sub, ok := s.subscriptions[msg.Event]

for s := range sub.clients { // HL
	s <- result
}

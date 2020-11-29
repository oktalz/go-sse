func (s *sse) Emit(event string, data interface{}) {
	s.messages <- history.Message{
		Event: event,
		Data:  data,
	}
}

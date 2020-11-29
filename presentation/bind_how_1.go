func (s *Server) Bind(structure interface{}) error {
	initMethod, err := s.bind.Bind(structure) // HL
	if err != nil {
		return err
	}
	if initMethod != nil {
		in := []reflect.Value{reflect.ValueOf(structure), reflect.ValueOf(s)}
		initMethod.Func.Call(in)
	}
	return nil
}

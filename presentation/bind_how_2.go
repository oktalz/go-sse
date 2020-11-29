func (b *bind) Bind(structure interface{}) (*reflect.Method, error) {

	structureName, err := getType(structure) // HL
	structType := reflect.TypeOf(structure)
	// ...
	for i := 0; i < structType.NumMethod(); i++ {
		method := structType.Method(i) // HL
		bind[method.Name] = method
	}
	return initMethod, nil
}

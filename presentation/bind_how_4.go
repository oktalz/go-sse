func (b *bind) Serve(w ResponseWriter, r *http.Request, structName, functionName string, args ...string) {

	method, structure, err := b.getMethod(structName, functionName)

	numIn := method.Type.NumIn() // HL
	argsShift := 0
	for i := 0; i < numIn; i++ { // HL
		inV := method.Type.In(i)
		in_Kind := inV.Kind() // func

		switch in_Kind { // HL
		case reflect.String:
			in = append(in, reflect.ValueOf(args[i-argsShift]))
		case reflect.Int:
			x, err := strconv.Atoi(args[i-argsShift]) // ...
			in = append(in, reflect.ValueOf(x))
		case reflect.Int64: // ...
		case reflect.Float64: // ...
		}
	}
	// calling method #2 // HL
}

	values := method.Func.Call(in) // HL
	errorInterface := reflect.TypeOf((*error)(nil)).Elem()

	// ...
	// if last return value is error type treat it differently
	http.Error(w, lastValue.Error(), http.StatusNotAcceptable) // HL

	// ...
	_, _ = w.Write(jsonResult) // HL

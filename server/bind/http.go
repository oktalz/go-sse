package bind

import (
	"net/http"
	"reflect"
	"strconv"

	jsoniter "github.com/json-iterator/go"
)

func (b *bind) Serve(w http.ResponseWriter, r *http.Request, structName, functionName string, args ...string) {

	var json = jsoniter.ConfigCompatibleWithStandardLibrary

	w.Header().Set("Access-Control-Allow-Origin", "*")

	w.Header().Set("Cache-Control", "no-cache")
	method, structure, err := b.getMethod(structName, functionName)
	if err != nil {
		http.Error(w, err.Error(), http.StatusBadRequest)
		return
	}

	in := []reflect.Value{}

	numIn := method.Type.NumIn()
	argsShift := 0
	for i := 0; i < numIn; i++ {
		inV := method.Type.In(i)
		in_Kind := inV.Kind() // func
		if i == 0 && in_Kind == reflect.Ptr {
			in = append(in, reflect.ValueOf(structure))
			argsShift = 1
			continue
		}
		switch in_Kind {
		case reflect.String:
			in = append(in, reflect.ValueOf(args[i-argsShift]))
		case reflect.Int:
			x, err := strconv.Atoi(args[i-argsShift])
			if err != nil {
				http.Error(w, err.Error(), http.StatusBadRequest)
				return
			}
			in = append(in, reflect.ValueOf(x))
		case reflect.Int64:
			x, err := strconv.ParseInt(args[i-argsShift], 10, 64)
			if err != nil {
				http.Error(w, err.Error(), http.StatusBadRequest)
				return
			}
			in = append(in, reflect.ValueOf(x))
		case reflect.Float64:
			x, err := strconv.ParseFloat(args[i-argsShift], 64)
			if err != nil {
				http.Error(w, err.Error(), http.StatusBadRequest)
				return
			}
			in = append(in, reflect.ValueOf(x))
		default:
			http.Error(w, "unsupported argument type", http.StatusInternalServerError)
			return
		}
		// fmt.Printf("\nParameter IN: "+strconv.Itoa(i)+"\nKind: %v\nName: %v\n-----------", in_Kind, inV.Name())
	}

	values := method.Func.Call(in)
	result := []interface{}{}
	errorInterface := reflect.TypeOf((*error)(nil)).Elem()
	skipLast := false
	if len(values) > 0 {
		typeOut := method.Type.Out(len(values) - 1)
		if typeOut.Implements(errorInterface) {
			skipLast = true
			lastValue, ok := values[len(values)-1].Interface().(error)
			if ok {
				http.Error(w, lastValue.Error(), http.StatusNotAcceptable)
				return
			}
		}
	}
	for index, value := range values {
		if skipLast && index == len(values)-1 {
			continue
		}
		result = append(result, value.Interface())
	}
	var jsonResult []byte
	if len(result) > 1 {
		jsonResult, err = json.Marshal(&result)
	} else {
		jsonResult, err = json.Marshal(&result[0])
	}
	if err != nil {
		http.Error(w, "", http.StatusInternalServerError)
		return
	}
	_, _ = w.Write(jsonResult)

	// Done.
	/*res := string(jsonResult)
	if len(res) > 60 {
		res = fmt.Sprintf("%s...+[%d]", res[0:59], len(res)-60)
	}
	log.Println(r.URL.Path, structName, functionName, res)*/
}

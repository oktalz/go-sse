package bind

import (
	"fmt"
	"net/http"
	"reflect"

	"github.com/oktalz/go-sse/server/router"
)

type Bind interface {
	Bind(structure interface{}) (*reflect.Method, error)
	Serve(w http.ResponseWriter, r *http.Request, structName, functionName string, args ...string)
}

func New() Bind {
	return &bind{
		binds:          map[string]*map[string]reflect.Method{},
		bindStructures: map[string]interface{}{},
		router:         router.New(),
	}
}

type bind struct {
	binds          map[string]*map[string]reflect.Method
	bindStructures map[string]interface{}
	router         *router.Router
}

func (b *bind) getMethod(structure, method string) (reflect.Method, interface{}, error) {
	structureBinds, ok := b.binds[structure]
	if !ok {
		return reflect.Method{}, nil, fmt.Errorf("structure %s not found", structure)
	}
	reflectMethod, ok := (*structureBinds)[method]
	if !ok {
		return reflect.Method{}, nil, fmt.Errorf("method %s not found in structure %s", method, structure)
	}
	structureData := b.bindStructures[structure]
	return reflectMethod, structureData, nil
}

func (b *bind) Bind(structure interface{}) (*reflect.Method, error) {
	bind := map[string]reflect.Method{}
	structureName, err := getType(structure)
	if err != nil {
		return nil, err
	}
	b.binds[structureName] = &bind
	b.bindStructures[structureName] = structure
	structType := reflect.TypeOf(structure)
	// find if Init() exists
	var initMethod *reflect.Method
	for i := 0; i < structType.NumMethod(); i++ {
		method := structType.Method(i)
		numIn := method.Type.NumIn()
		if method.Name == "Init" && numIn == 2 { // maybe allow non pointer func
			initMethod = &method
			// fmt.Println("Found", structureName, method.Name)
			continue
		}
		bind[method.Name] = method
		// fmt.Println("Binding", structureName, method.Name)
	}
	return initMethod, nil
}

func getType(myvar interface{}) (string, error) {
	if t := reflect.TypeOf(myvar); t.Kind() == reflect.Ptr {
		return t.Elem().Name(), nil
	} else {
		return "", fmt.Errorf("must be a pointer")
	}
}

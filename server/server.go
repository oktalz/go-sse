package server

import (
	"math/rand"
	"net/http"
	"reflect"
	"time"

	"github.com/oklog/ulid"

	"github.com/oktalz/go-sse/server/bind"
	"github.com/oktalz/go-sse/server/sse"
)

type ServerOptions struct {
	Path string
}

type Server struct {
	bind    bind.Bind
	sse     sse.SSE
	options ServerOptions
}

func New(options ServerOptions) *Server {
	server := Server{
		options: options,
	}
	server.sse = sse.New()
	server.bind = bind.New()
	return &server
}

func (s *Server) Start() {
	s.sse.HandleClients()

	/*r.Route(s.options.Path+"/sse/event/{event}", func(r chi.Router) {
		r.Get("/", s.EventGet)
	})
	r.Route(s.options.Path+"/bind/{struct}/{method}", func(r chi.Router) {
		r.Get("/", s.BindGet)
	})
	r.Get(s.options.Path+"/sse", s.HandlerGet)
	r.Patch(s.options.Path+"/sse", s.HandlerPatch)
	r.Delete(s.options.Path+"/sse", s.HandlerDelete)*/
}

func (s *Server) Bind(structure interface{}) error {
	initMethod, err := s.bind.Bind(structure)
	if err != nil {
		return err
	}
	// fmt.Println("Running", initMethod.Name)
	if initMethod != nil {
		in := []reflect.Value{reflect.ValueOf(structure), reflect.ValueOf(s)}
		initMethod.Func.Call(in)
	}
	return nil
}

func (s *Server) AddEvent(event string, queue int) {
	s.sse.AddEvent(event, queue)
}

func (s *Server) Emit(event string, data interface{}) {
	s.sse.Emit(event, data)
}

type TimeMessage struct {
	ID   int       `json:"id"`
	Time time.Time `json:"time"`
}

type LoginHandler struct {
	LoginPath string
}

// ServeHTTP handles HTTP requests.
func (s *LoginHandler) ServeHTTP(w http.ResponseWriter, r *http.Request) {
	// TODO remove this !!!!!!
	w.Header().Set("Access-Control-Allow-Origin", "*")
	// TODO check on DB user/pass ....
	expiration := time.Now().Add(365 * 24 * time.Hour)
	t := time.Now()
	entropy := rand.New(rand.NewSource(t.UnixNano())) // nolint:gosec
	token := ulid.MustNew(ulid.Timestamp(t), entropy).String()

	// cookie := http.Cookie{Name: "sh_username", Value: user.Email, Expires: expiration}
	// http.SetCookie(w, &cookie)
	cookieToken := http.Cookie{Name: "sh_token", Value: token, Expires: expiration}
	http.SetCookie(w, &cookieToken)
}

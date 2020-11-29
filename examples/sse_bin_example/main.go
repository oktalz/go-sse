package main

import (
	"context"
	"fmt"
	"io/ioutil"
	"log"
	"math/rand"
	"net/http"
	"time"

	"github.com/oktalz/go-sse/client"
	"github.com/oktalz/go-sse/server"
)

type SomeStruct struct {
	A int
	B string
}

type RandomStruct struct {
	Event string
}

func (r *RandomStruct) Init(server *server.Server) error {
	// init what you need
	// add one event loop for example
	go func() {
		for {
			server.Emit(r.Event, r.GetRandom())
			time.Sleep(4 * time.Second)
		}
	}()
	return nil
}

// you can call it directly from web also via /api/bind/RandomStruct/GetRandom
func (r *RandomStruct) GetRandom() int {
	return rand.Intn(100) // nolint:gosec
}

// you can call it directly from web also via /api/bind/RandomStruct/GetRandomInStruct?args={argument}
func (r *RandomStruct) GetRandomInStruct(argument string) SomeStruct {
	return SomeStruct{
		A: rand.Intn(100), // nolint:gosec
		B: argument,
	}
}

func randomRequest(url string) {
	response, err := http.Get(url) // nolint:gosec
	if err != nil {
		fmt.Println(err)
		return
	}
	defer response.Body.Close()
	contents, err := ioutil.ReadAll(response.Body)
	if err != nil {
		fmt.Println(err)
		return
	}
	fmt.Printf("%s\n", string(contents))
}

func main() {
	log.SetFlags(log.LstdFlags | log.Lshortfile)
	server := server.New(server.ServerOptions{
		Path: "/api",
	})
	// Start processing events
	server.AddEvent("time", 7)
	server.AddEvent("weather", 1)
	server.AddEvent("randomstruct", 5)

	_ = server.Bind(&RandomStruct{Event: "randomstruct"})

	server.Start()
	go emitLoop(server)

	chn := make(chan client.Event)
	var clientDone context.CancelFunc

	go func() {
		time.Sleep(1 * time.Second)
		id, chn, cancel := client.Get("http://127.0.0.1:3020/api/sse?sub=time,randomstruct", client.SSEOptions{
			Reconnect:      true,
			ReconnectAfter: 5 * time.Second,
		})
		fmt.Println("ID is", id)
		clientDone = cancel
		var data client.Event
		for {
			data = <-chn
			if data.Data.Event == "time" && data.Err == nil {
				time, ok := data.Data.Data.(string)
				if !ok {
					fmt.Println("REC UNEXPECTED DATA", data.Data.ID, data.Data.Event, data.Data.Data, data.Err)
				} else {
					fmt.Println("REC", data.Data.ID, data.Data.Event, time)
				}
			} else {
				fmt.Println("REC", data.Data.ID, data.Data.Event, data.Data.Data, data.Err)
			}
		}
	}()
	sseServer := startSSEServer(":3020", server)
	time.Sleep(5 * time.Second)
	clientDone()
	go func() {
		time.Sleep(5 * time.Second)
		client.Get("http://127.0.0.1:3020/api/sse?sub=time", client.SSEOptions{
			Reconnect:      true,
			ReconnectAfter: 5 * time.Second,
		})
		var data client.Event
		for {
			data = <-chn
			log.Println(data.Data.ID, data.Data.Event, data.Err)
		}
	}()
	time.Sleep(10 * time.Second)
	ctx, cancel := context.WithTimeout(context.Background(), 1*time.Second)
	defer cancel()
	err := sseServer.Shutdown(ctx)
	if err != nil {
		log.Panic(err)
	}
	time.Sleep(59 * time.Second)
}

func startSSEServer(adr string, handler http.Handler) *http.Server {
	srv := &http.Server{
		Addr:    adr,
		Handler: handler,
	}
	go func() {
		err := srv.ListenAndServe()
		if err != http.ErrServerClosed {
			log.Fatalf("ListenAndServe(): %s", err)
		}
		srv.Close()
	}()
	return srv
}

func emitLoop(server *server.Server) {
	go func() {
		for i := 0; ; i++ {
			server.Emit("time", time.Now())
			time.Sleep(1000 * time.Millisecond)
			if i > 1e9 {
				i = 0
			}
		}
	}()
}

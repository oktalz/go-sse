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
			time.Sleep(1 * time.Second)
			server.Emit(r.Event, r.GetRandom())
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
	server.AddEvent("randomstructevent", 3)
	_ = server.Bind(&RandomStruct{Event: "randomstructevent"})

	server.Start()

	var clientDone context.CancelFunc

	go func() {
		time.Sleep(1 * time.Second)
		id, chn, cancel := client.Get("http://127.0.0.1:3020/api/sse?sub=randomstructevent", client.SSEOptions{
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

	// do simple http get at one second
	go func() {
		time.Sleep(800 * time.Millisecond)
		randomRequest(fmt.Sprintf("http://127.0.0.1:3020/api/bind/RandomStruct/%s", "GetRandom"))
		randomRequest(fmt.Sprintf("http://127.0.0.1:3020/api/bind/RandomStruct/%s", "GetRandom"))
		randomRequest(fmt.Sprintf("http://127.0.0.1:3020/api/bind/RandomStruct/%s?args=%s", "GetRandomInStruct", "one"))
		randomRequest(fmt.Sprintf("http://127.0.0.1:3020/api/bind/RandomStruct/%s?args=%s", "GetRandomInStruct", "two"))
	}()

	sseServer := startSSEServer(":3020", server)
	time.Sleep(5 * time.Second)
	log.Println("disconnecting client")
	clientDone()

	time.Sleep(15 * time.Second)
	ctx, cancel := context.WithTimeout(context.Background(), 1*time.Second)
	defer cancel()
	err := sseServer.Shutdown(ctx)
	if err != nil {
		log.Panic(err)
	}
	// wait a bit so we receive sse event
	time.Sleep(15 * time.Second)
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

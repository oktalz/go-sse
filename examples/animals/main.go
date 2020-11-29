package main

import (
	"context"
	"fmt"
	"io/ioutil"
	"log"
	"net/http"
	"time"

	"github.com/oktalz/go-sse/server"
)

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
	// response is always in json format
	fmt.Printf("%s\n", string(contents))
}

func main() {
	log.SetFlags(log.LstdFlags | log.Lshortfile)

	server := server.New(server.ServerOptions{
		Path: "/api",
	})

	_ = server.Bind(&Animals{})

	server.Start()

	// do simple http get
	go func() {
		time.Sleep(800 * time.Millisecond)
		randomRequest(fmt.Sprintf("http://127.0.0.1:3020/api/bind/Animals/%s?name=%s", "GetAnimalSound", "cat"))
		// this is example mentioned, args names can't be read in runtime
		randomRequest(fmt.Sprintf("http://127.0.0.1:3020/api/bind/Animals/%s?hahahaha=%s", "GetAnimalSound", "lion"))
	}()

	sseServer := startSSEServer(":3020", server)
	time.Sleep(5 * time.Second)
	ctx, cancel := context.WithTimeout(context.Background(), 1*time.Second)
	defer cancel()
	err := sseServer.Shutdown(ctx)
	if err != nil {
		log.Panic(err)
	}
	time.Sleep(1 * time.Second)
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

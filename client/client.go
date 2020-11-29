package client

import (
	"bufio"
	"context"
	"errors"
	"fmt"
	"io"
	"net/http"
	"syscall"
	"time"

	jsoniter "github.com/json-iterator/go"
	"github.com/oktalz/go-sse/server/sse"
)

type Event struct {
	Data sse.MessageWithID
	Err  error
}

type SSEOptions struct {
	Reconnect      bool
	ReconnectAfter time.Duration
}

func Get(uri string, options ...SSEOptions) (id string, evCh chan Event, cancel context.CancelFunc) {
	evCh = make(chan Event)
	idCh := make(chan string)
	ctx, clientDone := context.WithCancel(context.Background())

	sseOptions := SSEOptions{
		Reconnect: false,
	}
	if len(options) > 0 {
		sseOptions = SSEOptions{
			Reconnect:      options[0].Reconnect,
			ReconnectAfter: options[0].ReconnectAfter,
		}
	}
	go get(ctx, uri, idCh, evCh, sseOptions)
	return <-idCh, evCh, clientDone
}

func get(ctx context.Context, uri string, idCh chan string, evCh chan<- Event, options SSEOptions) {
	resp, err := http.Get(uri) //nolint:gosec
	if err != nil {
		if errors.Is(err, syscall.ECONNREFUSED) && options.Reconnect {
			evCh <- Event{
				Err: err,
			}
			time.Sleep(options.ReconnectAfter)
			// endless loop of connecting
			go get(ctx, uri, idCh, evCh, options)
		} else {
			evCh <- Event{
				Err: err,
			}
		}
		return
	}
	defer resp.Body.Close()
	events := make(chan Event)
	reconnect := make(chan struct{})

	go streamScan(idCh, events, reconnect, resp.Body, options)

	for {
		select {
		case data := <-events:
			evCh <- data
		case <-reconnect:
			time.Sleep(options.ReconnectAfter)
			go get(ctx, uri, idCh, evCh, options)
			return
		case <-ctx.Done():
			// fmt.Println("halted operation2")
			return
		}
	}
}

func streamScan(idCh chan string, events chan Event, reconnect chan struct{}, body io.ReadCloser, options SSEOptions) {
	var json = jsoniter.ConfigCompatibleWithStandardLibrary
	var err error
	scanner := bufio.NewScanner(body)
	for scanner.Scan() {
		result := scanner.Bytes()
		if len(result) > 6 { // `data: ` is removed
			var data sse.MessageWithID
			err = json.Unmarshal(result[6:], &data)
			if err != nil {
				events <- Event{
					Err: scanner.Err(),
				}
			} else {
				if data.Event == "ID" {
					idStr, ok := data.Data.(string)
					if !ok {
						events <- Event{
							Err: fmt.Errorf("ID in wrong format"),
						}
						close(idCh)
					} else {
						idCh <- idStr
						close(idCh)
					}
				} else {
					events <- Event{
						Data: data,
					}
				}
			}
		}
	}
	if errors.Is(scanner.Err(), io.ErrUnexpectedEOF) {
		if options.Reconnect {
			reconnect <- struct{}{}
			return
		}
	}
	events <- Event{
		Err: scanner.Err(),
	}
}

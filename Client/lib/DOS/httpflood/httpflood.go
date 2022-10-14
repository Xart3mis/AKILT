package httpflood

import (
	"bytes"
	"context"
	"fmt"
	"math/rand"
	"net/http"
	"net/url"
	"sync"
	"time"
)

var waitGroup sync.WaitGroup
var data chan string

func fetchURL(count int, ctx context.Context) {
	fmt.Printf("URL fetcher #%d has started ...\n", count)
	for {
		select {
		case <-ctx.Done():
			fmt.Printf("URL fetcher #%d has stopped...\n", count)
			return
		default:
			url := <-data
			var method string
			var req *http.Request
			var err error

			client := &http.Client{}

			// Random method: mix things up a bit!
			switch rand.Intn(7) {
			case 0:
				method = "HEAD"
				req, err = http.NewRequest(method, url, nil)
			case 1:
				method = "POST"
				randomJson := []byte(`{"why":"I do not know.."}`)
				req, err = http.NewRequest(method, url, bytes.NewBuffer(randomJson))
			case 2:
				method = "PUT"
				randomJson := []byte(`{"why":"Ooh! Ooh! I know!"}`)
				req, err = http.NewRequest(method, url, bytes.NewBuffer(randomJson))
			case 3:
				method = "PATCH"
				randomJson := []byte(`{"why":"Ah damn, I lost it again."}`)
				req, err = http.NewRequest(method, url, bytes.NewBuffer(randomJson))
			case 4:
				method = "HELP"
				randomJson := []byte(`{"why":"Really. I lost it. Send help. Kthxbye."}`)
				req, err = http.NewRequest(method, url, bytes.NewBuffer(randomJson))
			case 5:
				method = "MATTIAS"
				randomJson := []byte(`{"why":"Now I remember. The METHOD can be _anything_."}`)
				req, err = http.NewRequest(method, url, bytes.NewBuffer(randomJson))
			default:
				method = "GET"
				req, err = http.NewRequest(method, url, nil)
			}

			if err != nil {
				fmt.Println(err)
			}

			client.Do(req)

			// Don't care about exit codes, just get that HTTP call out the door
			fmt.Printf("#%d: fetched %s \n", count, url)
		}
	}
}

func FloodUrl(target string, timeout time.Duration, workersN int64) error {
	rand.Seed(time.Now().Unix())

	ctx, cancel := context.WithTimeout(context.Background(), timeout)
	defer cancel()

	_, err := url.ParseRequestURI(target)
	if err != nil {
		return err
	}

	// Loop with some random parameters
	data = make(chan string)

	// Start X amount of concurrent fetchers
	go func(ctx context.Context) {
		for i := int64(1); i < workersN+1; i++ {
			select {
			case <-ctx.Done():
				return
			default:
				waitGroup.Add(1)
				go fetchURL(int(i), ctx)
				data <- (target)
			}
		}
	}(ctx)

	<-time.After(timeout)
	return nil
}

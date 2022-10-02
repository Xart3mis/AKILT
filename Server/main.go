//TODO: GUI for this
//TODO: protocol buffers + gRPC

package main

import (
	"bufio"
	"encoding/json"
	"fmt"
	"net/http"
	"os"
)

type ClientData struct {
	ShouldUpdate bool   `json:"ClientShouldUpdate"`
	OnScreenText string `json:"ClientOnScreenText"`
}

type ClientProductId struct {
	ProductId string `json:"ProductId"`
}

var clients map[string]ClientData = make(map[string]ClientData)
var on_screen_text string

func main() {
	go func() {
		scanner := bufio.NewScanner(os.Stdin)
		for scanner.Scan() {
			on_screen_text = scanner.Text()
		}
		if err := scanner.Err(); err != nil {
			os.Exit(1)
		}
	}()

	// handle route using handler function
	http.HandleFunc("/", func(w http.ResponseWriter, r *http.Request) {
		switch r.Method {
		case http.MethodPost:
			SetProductId(w, r)
		default:
			http.Error(w, "Method Not Allowed", http.StatusMethodNotAllowed)
		}

	})

	// listen to port
	http.ListenAndServe(":5050", nil)
}

func SetProductId(w http.ResponseWriter, req *http.Request) {
	var pid ClientProductId

	if err := json.NewDecoder(req.Body).Decode(&pid); err != nil {
		w.WriteHeader(400)
		fmt.Fprint(w, "Invalid Product Id")
		return
	}

	should_update := len(on_screen_text) > 0
	clients[pid.ProductId] = ClientData{ShouldUpdate: should_update, OnScreenText: on_screen_text}
	client_json, err := json.Marshal(clients)
	if err != nil {
		panic(err)
	}
	fmt.Fprint(w, string(client_json))
}

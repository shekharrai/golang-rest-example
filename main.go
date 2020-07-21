package main

import (
	"encoding/json"
	"fmt"
	"log"
	"net/http"
)

type Message struct {
	Key   string `json:"key, omitempty"`
	Value string `json:"value, omitempty"`
}

func (h *Message) ServeHTTP(w http.ResponseWriter, r *http.Request) {
	w.Write([]byte(h.Key))
}

func main() {

	unmarshal()

}

func marshal() {
	message := &Message{
		Key:   "Hello",
		Value: "World",
	}
	messageJson, error := json.Marshal(message)
	if error != nil {
		log.Fatal(error)
	}
	fmt.Println(string(messageJson))
}

func unmarshal() {
	messageJson := `{
		"key":"Hello",
		"value":"World"
	}`

	message := Message{}
	error := json.Unmarshal([]byte(messageJson), &message)
	if error != nil {
		log.Fatal(error)
	}
	fmt.Println(message.Key)
	fmt.Println(message.Value)
}

func handleHttp() {
	http.Handle("/hello", &Message{Key: "Hello Value"})

	http.HandleFunc("/world", func(w http.ResponseWriter, r *http.Request) {
		w.Write([]byte("Simple Value!!"))
	})

	http.ListenAndServe(":4400", nil)
}

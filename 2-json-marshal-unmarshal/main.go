package main

import (
	"encoding/json"
	"fmt"
	"log"
)

type Message struct {
	Key   string `json:"key, omitempty"`
	Value string `json:"value, omitempty"`
}

func main() {

	marshal()
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

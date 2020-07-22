package main

import (
	"net/http"
)

func main() {

	greetingsAPI()

}

func greetingsAPI() {

	http.HandleFunc("/greetings", func(w http.ResponseWriter, r *http.Request) {
		w.Write([]byte("Yay! Hello World!"))
	})

	http.ListenAndServe(":4400", nil)
}

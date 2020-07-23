package main

import (
	"encoding/json"
	"io/ioutil"
	"log"
	"net/http"
)

type Book struct {
	Id     int     `json:"id"`
	Name   string  `json:"name, omitempty"`
	Author string  `json:"author, omitempty"`
	Price  float32 `json:"price, omitempty"`
}

var books []Book

func init() {
	booksJson := `[
		{
			"id": 1,
			"name": "Clean Code: A Handbook of Agile Software Craftsmanship",
			"author": "Robert C. Martin",
			"price": 28.79
		}
	]`

	err := json.Unmarshal([]byte(booksJson), &books)

	if err != nil {
		log.Fatal(err)
	}
}

func booksHandler(w http.ResponseWriter, r *http.Request) {
	switch r.Method {
	case http.MethodGet:
		getAllBooks(w)
		return

	case http.MethodPost:
		saveBook(w, r)
		return
	}
}

func getAllBooks(w http.ResponseWriter) {
	booksJson, err := json.Marshal(books)

	if err != nil {
		w.WriteHeader(http.StatusInternalServerError)
		return
	}

	w.Header().Set("Content-Type", "application/json")
	w.Write([]byte(booksJson))
}

func saveBook(w http.ResponseWriter, r *http.Request) {

	var newBook Book
	bookBytes, err := ioutil.ReadAll(r.Body)

	if err != nil {
		w.WriteHeader(http.StatusBadRequest)
		return
	}

	err = json.Unmarshal(bookBytes, &newBook)

	if err != nil {
		w.WriteHeader(http.StatusBadRequest)
		return
	}

	books = append(books, newBook)

	w.WriteHeader(http.StatusCreated)
}

func main() {
	http.HandleFunc("/books", booksHandler)
	http.ListenAndServe(":4400", nil)
}

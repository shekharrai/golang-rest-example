package main

import (
	"encoding/json"
	"io/ioutil"
	"log"
	"net/http"
	"strconv"
	"strings"
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
			"id": 0,
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

func bookHandler(w http.ResponseWriter, r *http.Request) {
	pathSegments := strings.Split(r.URL.Path, "books/")
	bookId, err := strconv.Atoi(pathSegments[len(pathSegments)-1])

	if err != nil {
		w.WriteHeader(http.StatusNotFound)
		return
	}

	book, bookIndex := retrieveBook(bookId)

	if book == nil {
		w.WriteHeader(http.StatusNotFound)
		return
	}

	switch r.Method {
	case http.MethodGet:
		getByBookId(w, book)
		return

	case http.MethodPut:
		updateBookById(w, r, bookId, bookIndex)
		return

	case http.MethodDelete:
		books = append(books[:bookIndex], books[bookIndex+1:]...)
		w.WriteHeader(http.StatusOK)
		return

	default:
		w.WriteHeader(http.StatusMethodNotAllowed)
		return
	}
}

func updateBookById(w http.ResponseWriter, r *http.Request, bookId int, bookIndex int) {
	bookByte, err := ioutil.ReadAll(r.Body)
	if err != nil {
		w.WriteHeader(http.StatusBadRequest)
		return
	}
	var updateBook Book
	err = json.Unmarshal(bookByte, &updateBook)

	if err != nil {
		w.WriteHeader(http.StatusBadRequest)
		return
	}
	updateBook.Id = bookId
	books[bookIndex] = updateBook
	w.WriteHeader(http.StatusOK)
}

func getByBookId(w http.ResponseWriter, book *Book) {
	bookJson, err := json.Marshal(book)

	if err != nil {
		w.WriteHeader(http.StatusInternalServerError)
		return
	}

	w.Header().Add("Content-Type", "application/json")
	w.Write(bookJson)
	w.WriteHeader(http.StatusOK)
}

func getAllBooks(w http.ResponseWriter) {
	booksJson, err := json.Marshal(books)

	if err != nil {
		w.WriteHeader(http.StatusInternalServerError)
		return
	}

	w.Header().Set("Content-Type", "application/json")
	w.Write(booksJson)
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

	newBook.Id = getNextId()

	books = append(books, newBook)

	w.WriteHeader(http.StatusCreated)
}

func retrieveBook(bookId int) (*Book, int) {
	for i, book := range books {
		if book.Id == bookId {
			return &book, i
		}
	}
	return nil, 0
}

func getNextId() int {
	highestId := -1

	for _, book := range books {
		if highestId < book.Id {
			highestId = book.Id
		}
	}
	return highestId + 1
}

func main() {
	http.HandleFunc("/books", booksHandler)
	http.HandleFunc("/books/", bookHandler)
	http.ListenAndServe(":4400", nil)
}

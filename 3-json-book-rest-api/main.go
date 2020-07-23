package main

import (
	"encoding/json"
	"fmt"
	"io/ioutil"
	"log"
	"net/http"
	"strconv"
	"strings"
	"time"
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

func booksHandler(writer http.ResponseWriter, request *http.Request) {
	switch request.Method {
	case http.MethodGet:
		getAllBooks(writer)
		return

	case http.MethodPost:
		saveBook(writer, request)
		return
	}
}

func bookHandler(writer http.ResponseWriter, request *http.Request) {
	pathSegments := strings.Split(request.URL.Path, "books/")
	bookId, err := strconv.Atoi(pathSegments[len(pathSegments)-1])

	if err != nil {
		writer.WriteHeader(http.StatusNotFound)
		return
	}

	book, bookIndex := retrieveBook(bookId)

	if book == nil {
		writer.WriteHeader(http.StatusNotFound)
		return
	}

	switch request.Method {
	case http.MethodGet:
		getByBookId(writer, book)
		return

	case http.MethodPut:
		updateBookById(writer, request, bookId, bookIndex)
		return

	case http.MethodDelete:
		books = append(books[:bookIndex], books[bookIndex+1:]...)
		writer.WriteHeader(http.StatusOK)
		return

	default:
		writer.WriteHeader(http.StatusMethodNotAllowed)
		return
	}
}

func updateBookById(writer http.ResponseWriter, request *http.Request, bookId int, bookIndex int) {
	bookByte, err := ioutil.ReadAll(request.Body)
	if err != nil {
		writer.WriteHeader(http.StatusBadRequest)
		return
	}
	var updateBook Book
	err = json.Unmarshal(bookByte, &updateBook)

	if err != nil {
		writer.WriteHeader(http.StatusBadRequest)
		return
	}
	updateBook.Id = bookId
	books[bookIndex] = updateBook
	writer.WriteHeader(http.StatusOK)
}

func getByBookId(writer http.ResponseWriter, book *Book) {
	bookJson, err := json.Marshal(book)

	if err != nil {
		writer.WriteHeader(http.StatusInternalServerError)
		return
	}

	writer.Header().Add("Content-Type", "application/json")
	writer.Write(bookJson)
	writer.WriteHeader(http.StatusOK)
}

func getAllBooks(writer http.ResponseWriter) {
	booksJson, err := json.Marshal(books)

	if err != nil {
		writer.WriteHeader(http.StatusInternalServerError)
		return
	}

	writer.Header().Set("Content-Type", "application/json")
	writer.Write(booksJson)
}

func saveBook(writer http.ResponseWriter, request *http.Request) {

	var newBook Book
	bookBytes, err := ioutil.ReadAll(request.Body)

	if err != nil {
		writer.WriteHeader(http.StatusBadRequest)
		return
	}

	err = json.Unmarshal(bookBytes, &newBook)

	if err != nil {
		writer.WriteHeader(http.StatusBadRequest)
		return
	}

	newBook.Id = getNextId()

	books = append(books, newBook)

	writer.WriteHeader(http.StatusCreated)
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

func logMiddlewareHandler(handler http.Handler) http.Handler {
	return http.HandlerFunc(func(writer http.ResponseWriter, request *http.Request) {
		start := time.Now()
		handler.ServeHTTP(writer, request)
		fmt.Printf("Log Middleware timeout: %s", time.Since(start))
	})
}

func main() {

	bookListHandler := http.HandlerFunc(booksHandler)
	bookItemHandler := http.HandlerFunc(bookHandler)

	http.Handle("/books", logMiddlewareHandler(bookListHandler))
	http.Handle("/books/", logMiddlewareHandler(bookItemHandler))
	http.ListenAndServe(":4400", nil)
}

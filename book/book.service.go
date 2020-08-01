package book

import (
	"encoding/json"
	"fmt"
	"io/ioutil"
	"net/http"
	"strconv"
	"strings"
)

func SetupRoutes(basePath string) {
	http.Handle(fmt.Sprintf("%s/books", basePath), http.HandlerFunc(booksHandler))
	http.Handle(fmt.Sprintf("%s/books/", basePath), http.HandlerFunc(bookHandler))
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

	book := getBook(bookId)

	if book == nil {
		writer.WriteHeader(http.StatusNotFound)
		return
	}

	switch request.Method {
	case http.MethodGet:
		getByBookId(writer, book)
		return

	case http.MethodPut:
		updateBookById(writer, request)
		return

	case http.MethodDelete:
		removeBook(bookId)
		writer.WriteHeader(http.StatusOK)
		return

	case http.MethodOptions:
		return

	default:
		writer.WriteHeader(http.StatusMethodNotAllowed)
		return
	}
}

func updateBookById(writer http.ResponseWriter, request *http.Request) {
	bookByte, err := ioutil.ReadAll(request.Body)
	if err != nil {
		writer.WriteHeader(http.StatusBadRequest)
		return
	}
	var updatedBook Book
	err = json.Unmarshal(bookByte, &updatedBook)

	if err != nil {
		writer.WriteHeader(http.StatusBadRequest)
		return
	}
	updateBook(updatedBook)
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
	books := getBookList()
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
	_, err = updateBook(newBook)
	if err != nil {
		writer.WriteHeader(http.StatusInternalServerError)
		return
	}
	writer.WriteHeader(http.StatusCreated)
}

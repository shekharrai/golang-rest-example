package main

import (
	"github.com/shekharrai/golang-rest-example/book"
	"net/http"
)

const ApiBasePath = "/api/v1"

func main() {
	book.SetupRoutes(ApiBasePath)
	http.ListenAndServe(":4400", nil)
}

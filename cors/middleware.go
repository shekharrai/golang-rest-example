package cors

import (
	"net/http"
)

func MiddleWare(handler http.Handler) http.Handler {
	return http.HandlerFunc(func(writer http.ResponseWriter, request *http.Request) {
		writer.Header().Add("Access-Control-Allow-Origin", "*")
		writer.Header().Add("Content-Type", "application/json")
		writer.Header().Add("Access-Control-Allow-Methods", "POST, GET, OPTIONS, PUT, DELETE")
		handler.ServeHTTP(writer, request)
	})
}

package http

import (
	"net/http"

	"github.com/gorilla/mux"
)

func helloFunc(writer http.ResponseWriter, request *http.Request) {
	writer.Write([]byte("Hello"))
}

// NewRouter creates a new mux.Router.
func NewRouter() *mux.Router {
	router := mux.NewRouter()

	router.HandleFunc("/test", helloFunc)
	return router
}

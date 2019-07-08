package main

import (
	"net/http"

	appHttp "github.com/nieltg/quickshoot-party-match-server/pkg/http"
)

func main() {
	http.ListenAndServe(":8080", appHttp.NewRouter())
}

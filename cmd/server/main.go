package main

import (
	"net/http"
	"time"

	appHttp "github.com/nieltg/quickshoot-party-match-server/pkg/http"
)

func main() {
	server := new(appHttp.Server)
	server.Domain.JoinMaxDuration = 5 * time.Minute

	http.ListenAndServe(":8080", server.Handler())
}

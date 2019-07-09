package main

import (
	"log"
	"net/http"
	"time"

	appHttp "github.com/nieltg/quickshoot-party-match-server/pkg/http"
)

func main() {
	server := new(appHttp.Server)
	server.DeferredRequestMaxDuration = 30 * time.Second
	server.Domain.JoinMaxDuration = 5 * time.Minute

	err := http.ListenAndServe(":8080", server.Handler())
	if err != nil {
		log.Fatal("Unable to listen and serve: ", err)
	}
}

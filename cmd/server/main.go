package main

import (
	"fmt"
	"net/http"
	"time"

	appHttp "github.com/nieltg/quickshoot-party-match-server/pkg/http"
)

func main() {
	fmt.Println("Server is starting...")

	server := new(appHttp.Server)
	server.DeferredRequestMaxDuration = 30 * time.Second
	server.Domain.JoinMaxDuration = 5 * time.Minute

	fmt.Println("Server is started!");

	http.ListenAndServe(":8080", server.Handler())
}

package main

import (
	"fmt"
	"log"
	"net/http"
	"time"

	"github.com/nieltg/quickshoot-party-match-server/pkg/iohttp"
	"github.com/nieltg/quickshoot-party-match-server/pkg/modelmemory"
)

func main() {
	handler := iohttp.Handler{
		DeferredRequestMaxDuration: 30 * time.Second,

		Domain: &modelmemory.Domain{
			JoinMaxDuration: 5 * time.Minute,
		},
	}

	fmt.Println("Server is configured to listen on :8080")

	err := http.ListenAndServe(":8080", handler.Handler())
	if err != nil {
		log.Fatal("Unable to listen and serve: ", err)
	}
}

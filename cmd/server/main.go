package main

import (
	"fmt"
	"log"
	"net/http"
	"os"
	"time"

	_ "github.com/joho/godotenv/autoload"

	"github.com/nieltg/quickshoot-party-match-server/pkg/iohttp"
	"github.com/nieltg/quickshoot-party-match-server/pkg/modelmemory"
)

func findListenAddress() string {
	if envPort := os.Getenv("PORT"); envPort != "" {
		return fmt.Sprintf(":%s", envPort)
	}
	return ":8080"
}

func main() {
	handler := iohttp.Handler{
		DeferredRequestMaxDuration: 30 * time.Second,

		Domain: &modelmemory.Domain{
			JoinMaxDuration: 5 * time.Minute,
		},
	}

	listenAddress := findListenAddress()

	fmt.Println("Server is configured to listen on", listenAddress)
	if err := http.ListenAndServe(listenAddress, handler.Handler()); err != nil {
		log.Fatal("Unable to listen and serve: ", err)
	}
}

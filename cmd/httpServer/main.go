package main

import (
	hndl "../httpServer/internal/handlers"
	"fmt"
	"log"
	"net/http"
	"os"
)

var apiPath = "/api/v1/wallet"

func main() {
	http.HandleFunc(apiPath, hndl.CreateWalletHandler)
	http.HandleFunc(apiPath+"/", hndl.WalletInfoHandler)

	port := os.Getenv("PORT")
	if port == "" {
		port = "8080"
		log.Printf("Defaulting to port %s", port)
	}

	log.Printf("Listening on port %s", port)
	log.Printf("Open http://localhost:%s in the browser", port)
	log.Fatal(http.ListenAndServe(fmt.Sprintf(":%s", port), nil))
}

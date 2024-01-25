package main

import (
	"fmt"
	hndl "infotecsSelection/cmd/httpServer/internal/handlers"
	"log"
	"net/http"
	"os"
)

var apiPath = "/api/v1/wallet"

func main() {
	http.HandleFunc(apiPath, hndl.CreateWalletHandler)
	http.HandleFunc(apiPath+"/", hndl.WalletInfoAndOpHandler)

	port := os.Getenv("PORT")
	if port == "" {
		port = "8080"
		log.Printf("Defaulting to port %s", port)
	}

	log.Printf("Listening on port %s", port)
	log.Printf("Open http://localhost:%s in the browser", port)
	log.Fatal(http.ListenAndServe(fmt.Sprintf(":%s", port), nil))
}

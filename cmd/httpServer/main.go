package main

import (
	"infotecsSelection/internal/app/httpServer"
	"log"
)

func main() {
	config := httpServer.NewConfig()

	if err := httpServer.Start(config); err != nil {
		log.Fatal(err)
	}
}

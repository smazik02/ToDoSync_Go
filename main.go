package main

import (
	"log"
	"os"
	"strconv"
)

func main() {
	if len(os.Args) < 2 {
		log.Fatalf("Usage: %s <port>\n", os.Args[0])
	}

	port, err := strconv.Atoi(os.Args[1])
	if err != nil {
		log.Fatalln(err.Error())
	}

	server, err := NewServer(port)
	if err != nil {
		log.Fatalln(err.Error())
	}

	db := NewDBConnection()
	defer db.Close()

	log.Printf("Server listening on port %d...\n", port)

	server.Run()

	log.Println("Server shutting down...")
}

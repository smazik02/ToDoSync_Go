package main

import (
	"log"
	"os"
	"strconv"
	"todosync_go/internal/database"
	"todosync_go/internal/server"
)

func main() {
	if len(os.Args) < 2 {
		log.Fatalf("Usage: %s <port>\n", os.Args[0])
	}

	port, err := strconv.Atoi(os.Args[1])
	if err != nil {
		log.Fatalln(err.Error())
	}

	srv, err := server.NewServer(port)
	if err != nil {
		log.Fatalln(err.Error())
	}

	db := database.NewDBConnection()
	defer db.Close()

	log.Printf("Server listening on port %d...\n", port)

	srv.Run()

	log.Println("Server shutting down...")
}

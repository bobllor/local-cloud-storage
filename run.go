package main

import (
	"log"

	"github.com/joho/godotenv"
)

const (
	dbName   = "MasterStorage"
	fileUser = "file_user"
)

func main() {
	/*addr := ":8080"

	_, err := server.NewServer(addr)
	if err != nil {
		log.Fatal(err)
	}*/

	err := godotenv.Load(".env")
	if err != nil {
		log.Fatal(err)
	}
}

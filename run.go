package main

import (
	"fmt"
	"log"
	"os"

	dbcon "github.com/bobllor/cloud-project/src/db_con"
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

	fileConfig := dbcon.NewConfig(fileUser, os.Getenv("DB_PASSWORD"), "tcp", "127.0.0.1", dbName)

	fileDbCon, err := dbcon.NewDatabase(fileConfig)
	if err != nil {
		log.Fatal(err)
	}

	filesDb := dbcon.NewFilesDatabase(fileDbCon)

	files, err := filesDb.QueryFiles()
	if err != nil {
		log.Fatal(err)
	}

	fmt.Println(files)
}

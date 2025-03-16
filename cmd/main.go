package main

import (
	"log"

	"github.com/gopherVault/db"
	"github.com/gopherVault/filebased"
)

func main() {
	dbConfig := db.NewDBConfig("sample.db")
	fileBasedDB := filebased.NewFileBasedDB(dbConfig)

	err := fileBasedDB.SaveToFile([]byte("Hello, World! testing 2"))

	if err != nil {
		log.Fatal(err)
	}
}

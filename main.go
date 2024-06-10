package main

import (
	"log"

	ace "github.com/H-ADJI/ace/cmd"
	_ "github.com/mattn/go-sqlite3"
)

func main() {
	db, err := ace.GetDB()
	if err != nil {
		log.Fatalln("couldnt get db connection", err)
	}
	ace.LoadData(db)
	log.Println(ace.Search(db, "array"))
}

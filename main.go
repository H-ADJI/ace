package main

import (
	"log"
	"os"

	ace "github.com/H-ADJI/ace/cmd"
	_ "github.com/mattn/go-sqlite3"
)

func main() {
	os.Setenv("DEV_ENV", "true")
	db, err := ace.GetDB()
	if err != nil {
		log.Fatalln("couldnt get db connection", err)
	}
	ace.LoadData(db)

}

package main

import (
	"fmt"
	"log"

	ace "github.com/H-ADJI/ace/cmd"
	_ "github.com/mattn/go-sqlite3"
)

func main() {
	db, err := ace.GetDB()
	if err != nil {
		log.Fatalln("couldnt get db", err)
	}
	err = ace.DropTable(db)
	if err != nil {
		log.Fatalln("couldnt drop table", err)
	}
	err = ace.CreateTable(db)
	if err != nil {
		log.Fatalln("couldnt create table", err)
	}
	challenges := ace.ParseChallenges()
	for i := range 10 {
		chall := challenges[i]
		url, _ := chall.InsertIntoDB(db)
		fmt.Println(url)
	}
}

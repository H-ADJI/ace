package main

import (
	"log"

	ace "github.com/H-ADJI/ace/cmd"
	tea "github.com/charmbracelet/bubbletea"
	_ "github.com/mattn/go-sqlite3"
)

func main() {
	db, err := ace.GetDB()
	if err != nil {
		log.Fatalln("couldnt get db connection", err)
	}
	ace.LoadData(db)
	// log.Println(ace.Search(db, "array"))
	tui := tea.NewProgram(ace.InitialModel())
	if _, err := tui.Run(); err != nil {
		log.Fatalln("running tui failed")
	}
}

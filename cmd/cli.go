package ace

import (
	"database/sql"
	"log"
)

func LoadData(db *sql.DB) {
	err := CreateTable(db)
	if err != nil {
		log.Fatalln("couldnt create table", err)
	}
	err = CreateSearchTable(db)
	if err != nil {
		log.Fatalln("couldnt create search table", err)
	}
	// read data from db
	challenges, err := readDBChallenges(db)
	if err != nil {
		log.Fatalf("coulndt read db challenges, %s", err)
	}
	if len(challenges) == 0 {
		// if not exist scrape it
		crawlChallenges(db)
	}
	err = PopulateSearchTable(db)
	if err != nil {
		log.Fatalln("couldnt populate search table", err)
	}
}

func Search(db *sql.DB, searchQuery string) []Challenge {
	challs, err := searchChallenge(db, searchQuery)
	if err != nil {
		log.Fatalf("Error getting search result due to %s", err)
	}
	return challs
}

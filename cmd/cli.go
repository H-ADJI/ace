package ace

import (
	"database/sql"
	"fmt"
	"log"
	"os"
	"strconv"
)

func LoadData(db *sql.DB) {
	devENV := "false"
	if fromEnv, ok := os.LookupEnv("DEV_ENV"); ok {
		devENV = fromEnv
	}
	if v, _ := strconv.ParseBool(devENV); v {
		fmt.Println(devENV)
		err := DropTable(db)
		if err != nil {
			log.Fatalln("couldnt drop table", err)
		}
	}
	err := CreateTable(db)
	if err != nil {
		log.Fatalln("couldnt create table", err)
	}
	// read data from db
	challenges, _ := readDBChallenges(db)
	if len(challenges) == 0 {
		// if not exist scrape it
		challenges = ScrapeChallenges(-1)
		// load it back to db
		for _, chall := range challenges {
			chall.InsertIntoDB(db)
		}
	}
}

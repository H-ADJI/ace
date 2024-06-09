package ace

import (
	"database/sql"
	"fmt"
	"log"
	"os"
	"path/filepath"

	_ "github.com/mattn/go-sqlite3"
	gap "github.com/muesli/go-app-paths"
)

type Challenge struct {
	Title       string
	Url         string
	Description string
	Tags        string
	Difficulty  string
}

func (chall Challenge) String() string {
	return fmt.Sprintf("Title : %s  |  Difficulty : %s | url : %s", chall.Title, chall.Difficulty, chall.Url)
}

const appName = "ace"

func GetDB() (*sql.DB, error) {
	userScope := gap.NewScope(gap.User, appName)
	appDataPaths, err := userScope.DataDirs()
	if err != nil {
		fmt.Printf("couldn't list system app directories, %s", err)
		return nil, err
	}
	usedPath := appDataPaths[0]
	err = os.MkdirAll(usedPath, os.ModePerm)
	if err != nil {
		fmt.Printf("couldn't create application directory, %s", err)
		return nil, err
	}
	dbPath := filepath.Join(usedPath, fmt.Sprintf("%s_data.db", appName))
	db, err := sql.Open("sqlite3", dbPath)
	if err != nil {
		fmt.Printf("couldn't open DB connection, %s", err)
	}
	log.Printf("db created at %s", dbPath)
	return db, nil
}

func CreateTable(db *sql.DB) error {
	_, err := db.Exec(`
	create table IF NOT EXISTS challenges (id integer not null primary key, title text,url text, description text, difficulty text,tags text);
	`)
	if err != nil {
		return fmt.Errorf("couldn't prepare sql statement, %s", err)
	}
	return nil
}

func CreateSearchTable(db *sql.DB) error {
	_, err := db.Exec(`
	CREATE VIRTUAL TABLE IF NOT EXISTS challenges_search USING FTS5(title, description, tags);
	`)
	if err != nil {
		return fmt.Errorf("couldn't prepare sql statement, %s", err)
	}
	return nil
}

func PopulateSearchTable(db *sql.DB) error {
	_, err := db.Exec(`
	INSERT INTO challenges_search SELECT title, description, tags FROM challenges;
	`)
	if err != nil {
		return fmt.Errorf("couldn't prepare sql statement, %s", err)
	}
	return nil
}

func DropTable(db *sql.DB) error {
	_, err := db.Exec("DROP TABLE IF EXISTS challenges;")
	return err
}

func (chall Challenge) InsertIntoDB(db *sql.DB) error {
	q, err := db.Prepare("INSERT INTO challenges (title, description, url, difficulty, tags) VALUES (?, ?, ?, ?,?);")
	if err != nil {
		return fmt.Errorf("couldn't prepare sql statement, %s", err)
	}
	defer q.Close()
	_, err = q.Exec(chall.Title, chall.Description, chall.Url, chall.Difficulty, chall.Tags)
	if err != nil {
		return fmt.Errorf("couldnt insert %s", chall)
	}

	return nil
}

func readDBChallenges(db *sql.DB) ([]Challenge, error) {
	res, err := db.Query("Select title, description, url, tags, difficulty from challenges")
	challenges := make([]Challenge, 0)
	if err != nil {
		return challenges, fmt.Errorf("wrong query, %s", err)
	}
	defer res.Close()
	for res.Next() {
		chall := Challenge{}
		res.Scan(&chall.Title, &chall.Description, &chall.Url, &chall.Tags, &chall.Difficulty)
		challenges = append(challenges, chall)
	}
	return challenges, nil
}

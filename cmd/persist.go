package ace

import (
	"database/sql"
	"fmt"
	"os"
	"path/filepath"

	_ "github.com/mattn/go-sqlite3"
)

type Challenge struct {
	Title       string
	Url         string
	Description string
	tags        []string
	Difficulty  string
}

func (chall Challenge) String() string {
	return fmt.Sprintf("Title : %s  |  Difficulty : %s | url : %s", chall.Title, chall.Difficulty, chall.Url)
}

const appName = "ace"

var homePath = os.Getenv("HOME")

func GetDB() (*sql.DB, error) {
	appDir := filepath.Join(homePath, fmt.Sprintf(".%s", appName))
	dbName := "data.db"
	if fromEnv, ok := os.LookupEnv("DB_NAME"); ok {
		dbName = fromEnv
	}
	err := os.MkdirAll(appDir, os.ModePerm)
	if err != nil {
		fmt.Printf("couldn't create application directory, %s", err)
		return nil, err
	}
	dbPath := filepath.Join(appDir, dbName)
	db, err := sql.Open("sqlite3", dbPath)
	if err != nil {
		fmt.Printf("couldn't open DB connection, %s", err)
	}
	return db, nil
}

func CreateTable(db *sql.DB) error {
	_, err := db.Exec(`
	create table IF NOT EXISTS challenges (id integer not null primary key, title text,url text, description text, difficulty text);
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

func (chall Challenge) InsertIntoDB(db *sql.DB) (string, error) {
	q, err := db.Prepare("INSERT INTO challenges (title, description, url, difficulty) VALUES (?, ?, ?, ?);")
	if err != nil {
		return "", fmt.Errorf("couldn't prepare sql statement, %s", err)
	}
	defer q.Close()
	res, err := q.Exec(chall.Title, chall.Description, chall.Url, chall.Difficulty)
	if err != nil {
		return "", fmt.Errorf("couldnt insert %s", chall)
	}
	insertedRowId, _ := res.LastInsertId()
	row := db.QueryRow("SELECT url  FROM challenges WHERE id=?;", insertedRowId)
	var url string
	err = row.Scan(&url)
	if err != nil {
		return "", fmt.Errorf("couldn't scan affected row, %s", err)
	}
	return url, nil
}

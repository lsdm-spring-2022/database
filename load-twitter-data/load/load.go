package load

import (
	"database/sql"
	"encoding/json"
	"fmt"
	"github.com/joho/godotenv"
	"io/ioutil"
	"log"
	"os"
	"path/filepath"
	"strings"
	"sync"
)

type TwitterData struct {
	Country     string `json:"country"`
	CreatedAt   string `json:"created_at"`
	Trend       string `json:"trend"`
	TweetVolume int    `json:"tweet_volume"`
	AsOf        string `json:"as_of"`
}

var (
	db      *sql.DB
	DB_USER string
	DB_PASS string
	DB_HOST string
	DB_PORT string
	DB_NAME string
)

func IterateOverTwitterData(wg *sync.WaitGroup, path string, id int) {
	defer wg.Done()

	writeToLogFile(fmt.Sprintf("Goroutine %d: Started", id), id)

	err := godotenv.Load(".env")
	if err != nil {
		log.Fatalf("Some error occured. Err: %s", err)
	}
	jsonFiles, err := ioutil.ReadDir(path)
	if err != nil {
		log.Fatal(err)
	}
	DB_USER = os.Getenv("DB_USER")
	DB_PASS = os.Getenv("DB_PASS")
	DB_HOST = os.Getenv("DB_HOST")
	DB_PORT = os.Getenv("DB_PORT")
	DB_NAME = os.Getenv("DB_NAME")

	createDbConnection()

	for _, jsonFile := range jsonFiles {
		readJSON(filepath.Join(path, jsonFile.Name()), id)
		writeToLogFile(fmt.Sprintf("Goroutine %d: Processed file %s", id, jsonFile.Name()), id)
	}

	writeToLogFile(fmt.Sprintf("Goroutine %d: Finished", id), id)
}

func createDbConnection() {
	connectionString := fmt.Sprintf("%s:%s@tcp(%s:%s)/%s", DB_USER, DB_PASS, DB_HOST, DB_PORT, DB_NAME)
	var err error
	db, err = sql.Open("mysql", connectionString)
	if err != nil {
		log.Fatal(err)
	}

	db.SetMaxIdleConns(100)
}

func readJSON(path string, id int) {
	f, err := os.Open(path)
	if err != nil {
		log.Fatal(err)
	}

	defer f.Close()

	byteResult, _ := ioutil.ReadAll(f)

	var twitterData []*TwitterData

	if err := json.Unmarshal(byteResult, &twitterData); err != nil {
		log.Fatal(err)
	}

	sqlText := createBatchInsertCommand(twitterData)

	insertBatchRows(sqlText, id)
}

func createBatchInsertCommand(twitterData []*TwitterData) string {
	var values []string

	for _, obj := range twitterData {
		values = append(values, fmt.Sprintf("('%s', '%s', '%s', '%d', '%s')", obj.Country, obj.CreatedAt, obj.Trend, obj.TweetVolume, obj.AsOf))
	}
	joinedInsert := fmt.Sprintf("INSERT IGNORE INTO twitter(country, created_at, trend, tweet_volume, as_of) VALUES %s", strings.Join(values, ","))
	return joinedInsert
}

func insertBatchRows(sqlText string, id int) {
	res, err := db.Exec(sqlText)

	if err != nil {
		writeToErrorFile(fmt.Sprintf("Exec error occured with error: %v\n", err))
		return
	}

	rowsAffected, affectedErr := res.RowsAffected()
	if affectedErr != nil {
		writeToErrorFile(fmt.Sprintf("Rows affected error occured with error: %v\n", err))
		log.Fatal(affectedErr)
	}

	logMessage := fmt.Sprintf("Goroutine %d: Rows affected = %d", id, rowsAffected)

	fmt.Println(logMessage)
	writeToLogFile(logMessage, id)
}

func writeToErrorFile(message string) {
	f, err := os.OpenFile("error-log.txt", os.O_APPEND|os.O_CREATE|os.O_WRONLY, 0644)
	if err != nil {
		log.Println(err)
	}
	defer f.Close()
	if _, err := f.WriteString(fmt.Sprintf("%s\n", message)); err != nil {
		log.Println(err)
	}
}

func writeToLogFile(message string, id int) {
	logName := fmt.Sprintf("logs/log-%d.txt", id)
	f, err := os.OpenFile(logName, os.O_APPEND|os.O_CREATE|os.O_WRONLY, 0644)
	if err != nil {
		log.Println(err)
	}
	defer f.Close()
	if _, err := f.WriteString(fmt.Sprintf("%s\n", message)); err != nil {
		log.Println(err)
	}
}

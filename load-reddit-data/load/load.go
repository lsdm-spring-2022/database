package load

import (
	"database/sql"
	"fmt"
	"github.com/gocarina/gocsv"
	"github.com/joho/godotenv"
	"io/ioutil"
	"log"
	"os"
	"path/filepath"
	"strconv"
	"strings"
	"sync"
)

type CsvRow struct {
	CreatedUTC  string `csv:"created_utc"`
	Subreddit   string `csv:"subreddit"`
	Title       string `csv:"title"`
	Score       string `csv:"score"`
	NumComments string `csv:"num_comments"`
	Region      string `csv:"region"`
	DateStored  string `csv:"date_stored"`
}

var (
	db      *sql.DB
	DB_USER string
	DB_PASS string
	DB_HOST string
	DB_PORT string
	DB_NAME string
)

func IterateOverCleanedCsvs(wg *sync.WaitGroup, path string, id int) {
	defer wg.Done()

	writeToLogFile(fmt.Sprintf("Goroutine %d: Started", id), id)

	err := godotenv.Load(".env")
	if err != nil {
		log.Fatalf("Some error occured. Err: %s", err)
	}
	csvs, err := ioutil.ReadDir(path)
	if err != nil {
		log.Fatal(err)
	}
	DB_USER = os.Getenv("DB_USER")
	DB_PASS = os.Getenv("DB_PASS")
	DB_HOST = os.Getenv("DB_HOST")
	DB_PORT = os.Getenv("DB_PORT")
	DB_NAME = os.Getenv("DB_NAME")

	createDbConnection()

	for _, csv := range csvs {
		readCSV(filepath.Join(path, csv.Name()), id)
		writeToLogFile(fmt.Sprintf("Goroutine %d: Processed file %s", id, csv.Name()), id)
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

func readCSV(path string, id int) {
	f, err := os.Open(path)
	if err != nil {
		log.Fatal(err)
	}

	defer f.Close()

	var rows []*CsvRow

	if err := gocsv.UnmarshalFile(f, &rows); err != nil {
		log.Fatal(err)
	}

	sqlText := createBatchInsertCommand(rows)

	insertBatchRows(sqlText, id)
}

func createBatchInsertCommand(rows []*CsvRow) string {
	var values []string

	for _, row := range rows {
		scoreInt, scoreErr := strconv.Atoi(row.Score)
		if scoreErr != nil {
			writeToErrorFile(fmt.Sprintf("Error converting score %s to int", row.Score))
			continue
		}

		numCommentsInt, numCommentsErr := strconv.Atoi(row.NumComments)
		if numCommentsErr != nil {
			writeToErrorFile(fmt.Sprintf("Error converting numComments %s to int", row.NumComments))
			continue
		}

		if len(row.Title) > 255 {
			row.Title = truncateString(row.Title)
		}

		values = append(values, fmt.Sprintf("('%s', '%s', '%s', '%s', '%d', '%s', '%d')", row.CreatedUTC, row.Region, row.Subreddit, row.Title, scoreInt, row.DateStored, numCommentsInt))
	}
	joinedInsert := fmt.Sprintf("INSERT IGNORE INTO reddit(date_posted, region, subreddit, post_title, upvotes, date_stored, comments) VALUES %s", strings.Join(values, ","))
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

func truncateString(s string) string {
	return s[:255]
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

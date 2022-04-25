package main

import (
	"database/sql"
	"encoding/csv"
	"fmt"
	"io/ioutil"
	"log"
	"os"
	"path/filepath"
	"strconv"
	"strings"
	"time"

	_ "github.com/go-sql-driver/mysql"
	"github.com/gocarina/gocsv"
)

var (
	DB_USER = os.Getenv("DB_USER")
	DB_PASS = os.Getenv("DB_PASS")
	DB_HOST = os.Getenv("DB_HOST")
	DB_PORT = os.Getenv("DB_PORT")
	DB_NAME = os.Getenv("DB_NAME")
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

func formatTimestamp(timestamp int64) string {
	return time.Unix(timestamp, 0).Format("2006-01-02 03:04:05")
}

func checkFloatParse(s string) (float64, error) {
	i, err := strconv.ParseFloat(s, 64)
	if err != nil {
		return 0.0, err
	}
	return i, nil
}

func checkIntParse(s string) (int64, error) {
	i, err := strconv.ParseInt(s, 10, 64)
	if err != nil {
		return 0, err
	}
	return i, nil
}

func truncateString(s string) string {
	return s[:255]
}

func cleanTitle(title string) string {
	noComma := strings.Replace(title, "'", "", -1)
	return strings.ToValidUTF8(noComma, "")
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

func loadRowIntoDatabase(created string, subreddit string, title string, score float64, numComments float64, region string, stored string) {
	connectionString := fmt.Sprintf("%s:%s@tcp(%s:%s)/%s", DB_USER, DB_PASS, DB_HOST, DB_PORT, DB_NAME)
	db, err := sql.Open("mysql", connectionString)

	if err != nil {
		writeToErrorFile(fmt.Sprintf("DB open error occured with subreddit: %s, title: %s, region: %s, err: %v", subreddit, title, region, err))
		return
	}

	defer db.Close()

	db.SetMaxIdleConns(64)
	db.SetMaxOpenConns(64)
	db.SetConnMaxLifetime(time.Minute)

	sql := fmt.Sprintf("INSERT INTO reddit(date_posted, region, subreddit, post_title, upvotes, date_stored, comments) VALUES ('%s', '%s', '%s', '%s', '%f', '%s', '%f')", created, region, subreddit, title, score, stored, numComments)
	res, err := db.Exec(sql)

	if err != nil {
		writeToErrorFile(fmt.Sprintf("Exec error occured with subreddit: %s, title: %s, region: %s, err: %v", subreddit, title, region, err))
		return
	}

	lastId, err := res.LastInsertId()

	if err != nil {
		log.Fatal(err)
	}

	fmt.Printf("Loaded region: %s, subreddit: %s, created: %s, lastId: %d\n", region, subreddit, created, lastId)
}

func writeDataToCsv(created string, subreddit string, title string, score float64, numComments float64, region string, stored string) {
	file, err := os.Create("ready-data.csv")
	if err != nil {
		log.Fatal(err)
	}

	defer file.Close()

	writer := csv.NewWriter(file)

	defer writer.Flush()

	writer.Write([]string{"created_utc", "subreddit", "title", "score", "num_comments", "region", "date_stored"})
}

func checkRowForErrors(row *CsvRow) {
	validCreated, createdErr := checkIntParse(row.CreatedUTC)
	if createdErr != nil {
		return
	}

	validScore, scoreErr := checkFloatParse(row.Score)
	if scoreErr != nil {
		return
	}

	validNumComments, numCommentsErr := checkFloatParse(row.NumComments)
	if numCommentsErr != nil {
		return
	}

	validDateStored, dateStoredErr := checkIntParse(row.DateStored)
	if dateStoredErr != nil {
		return
	}

	validTitle := cleanTitle(row.Title)

	if len(validTitle) > 255 {
		validTitle = truncateString(validTitle)
	}

	createdTimestamp := formatTimestamp(validCreated)
	storedTimestamp := formatTimestamp(validDateStored)

	loadRowIntoDatabase(createdTimestamp, row.Subreddit, validTitle, validScore, validNumComments, row.Region, storedTimestamp)
}

func parseCSV(path string) {
	f, err := os.Open(path)
	if err != nil {
		log.Fatal(err)
	}

	defer f.Close()

	rows := []*CsvRow{}

	if err := gocsv.UnmarshalFile(f, &rows); err != nil {
		log.Fatal(err)
	}

	for _, row := range rows {
		checkRowForErrors(row)
	}
}

func iterateOverCleanedCsvs(path string) {
	csvs, err := ioutil.ReadDir(path)
	if err != nil {
		log.Fatal(err)
	}

	for _, csv := range csvs {
		parseCSV(filepath.Join(path, csv.Name()))
	}
}

func main() {
	cleanedDataDirs, err := ioutil.ReadDir("cleaned-reddit-data")
	if err != nil {
		log.Fatal(err)
	}

	for _, dataDir := range cleanedDataDirs {
		fmt.Println(dataDir.Name())
		if strings.Contains(dataDir.Name(), "-cleaned-reddit-data") {
			go iterateOverCleanedCsvs("cleaned-reddit-data/" + dataDir.Name())
		}
	}

	var input string
	fmt.Scanln(&input)
}

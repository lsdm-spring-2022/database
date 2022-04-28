package main

import (
	"errors"
	"fmt"
	"io/ioutil"
	"loadtwitter/load"
	"log"
	"os"
	"path/filepath"
	"strings"
	"sync"
	"time"

	_ "github.com/go-sql-driver/mysql"
)

func main() {
	start := time.Now()

	createLogsDirectory()

	dataDirs, err := ioutil.ReadDir("twitter-data")
	if err != nil {
		log.Fatal(err)
	}

	var wg sync.WaitGroup

	for idx, dataDir := range dataDirs {
		if strings.Contains(dataDir.Name(), "-twitter-data") {
			writeToLogFile(fmt.Sprintf("Starting goroutine %d", idx))
			wg.Add(1)
			dirPath := filepath.FromSlash("twitter-data/" + dataDir.Name())
			go load.IterateOverTwitterData(&wg, dirPath, idx)
		}
	}

	writeToLogFile("Main: Waiting for goroutines to finish")
	wg.Wait()
	writeToLogFile("Main: Finished")

	duration := time.Since(start)
	dt := time.Now()
	formattedTime := dt.Format("01-02-2006 15:04:05")
	writeToLogFile(fmt.Sprintf("Completed at %s. Took %v", formattedTime, duration))
}

func writeToLogFile(message string) {
	f, err := os.OpenFile("log.txt", os.O_APPEND|os.O_CREATE|os.O_WRONLY, 0644)
	if err != nil {
		log.Println(err)
	}
	defer f.Close()
	if _, err := f.WriteString(fmt.Sprintf("%s\n", message)); err != nil {
		log.Println(err)
	}
}

func createLogsDirectory() {
	path := "logs"
	if _, err := os.Stat(path); errors.Is(err, os.ErrNotExist) {
		err := os.Mkdir(path, os.ModePerm)
		if err != nil {
			log.Println(err)
		}
	}
}

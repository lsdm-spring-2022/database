package main

import (
	"errors"
	"fmt"
	"io/ioutil"
	"loaddata/load"
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

	cleanedDataDirs, err := ioutil.ReadDir("cleaned-reddit-data")
	if err != nil {
		log.Fatal(err)
	}

	var wg sync.WaitGroup

	for idx, dataDir := range cleanedDataDirs {
		fmt.Println(dataDir.Name())
		if strings.Contains(dataDir.Name(), "-cleaned-reddit-data") {
			writeToLogFile(fmt.Sprintf("Starting goroutine %d", idx))
			wg.Add(1)
			dirPath := filepath.FromSlash("cleaned-reddit-data/" + dataDir.Name())
			go load.IterateOverCleanedCsvs(&wg, dirPath, idx)
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

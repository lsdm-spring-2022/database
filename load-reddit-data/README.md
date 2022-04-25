# Load Data

## Purpose
This code is used to load the cleaned Reddit data produced by the [clean-data](https://github.com/lsdm-spring-2022/database/tree/main/clean-reddit-data) code into a MySQL database. This code uses Goroutines and batch insert statements to load the data into a MySQL database.

## Loading Process
1. The main process creates a list of all directories in the `cleaned-reddit-data` directory
2. The main process sets up a `WaitGroup` to wait for all the goroutines to finish
3. The main process launches goroutines to process each directory in the list
4. Each goroutine reads all of the rows from each CSV file that it is processing
5. Each goroutine creates a batch insert statement for the rows and executes the batch insert statement

## Usage
1. Create a `.env` file using the `.env.sample` as a template
1. Run the file using the following command:
```go
go run main.go
```
**Note:** The code assumes that the `cleaned-reddit-data` directory is in the same directory as the `main.go` file.

## Results
- The code will produce logs to the console, a log file called `log.txt` in the current directory, a log file called `error-log.txt` in the current directory, and a directory called `logs` that contains the logs of each goroutine.
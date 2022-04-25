# Clean Reddit Data


## Purpose
This code is used to clean the reddit data retrieved by the [data-scraping](https://github.com/lsdm-spring-2022/data-scraping) repository. This code uses Pandas DataFrames to load the data and clean the data.

## Cleaning Process
1. Load the data from each CSV into a DataFrame
2. Removed all rows that contain `over_18` data
3. Extract necessary columns from the DataFrame
4. Clean the `title` column by removing line endings
5. Remove all rows that do not have a valid `created_utc` Unix timestamp
6. Removed all rows that do not have a valid integer value for `score` and `num_comments`
7. Format the `created_utc` date to be in the format `YYYY-MM-DD HH:MM:SS`
8. Write the cleaned DataFrame to a CSV file

## Usage
1. Install the dependencies
2. Run the code with the following command:
```python
python clean.py [REDDIT_DATA_DIRECTORY]
```
**Note:** The [REDDIT_DATA_DIRECTORY] is the directory containing the CSV files and should be in the format:
```bash
├── [country1]-reddit-data
│   ├── [country1]-[year]-[month]-reddit-data.csv
│   ├── ...
├── [country2]-reddit-data
│   ├── [country2]-[year]-[month]-reddit-data.csv
│   ├── [country2]-[year]-[month]-reddit-data.csv
│   ├── ...
```

## Results
- The code will create a new directory called `cleaned-reddit-data` in this directory. This directory will contain the cleaned CSV files and have the same structure as the `REDDIT_DATA_DIRECTORY`.
- The code will produce logs to the console and a log file called `error_log.txt` in the current directory.
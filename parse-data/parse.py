import os
import pandas as pd
import time


def write_to_error_file(message):
    with open('error_log.txt', 'a') as f:
        f.write(message)


def read_csv_data(file_path, region):
    current_epoch = int(time.time())
    try:
        df = pd.read_csv(file_path, sep='\t')
        df_copy = df[['created_utc', 'subreddit', 'title', 'score', 'num_comments']].copy()
        df_copy['region'] = region
        df_copy['date_stored'] = current_epoch
        return df_copy.replace(r'\r+|\n+|\t+|\,+','', regex=True)
    except Exception as e:
        print(f'Error parsing {file_path}')
        write_to_error_file(f'Error parsing {file_path}\n')
        return None


def write_df_to_csv(df, file_path):
    df.to_csv(file_path, sep=',', encoding='utf-8', index=False)


def main():
    REDDIT_DATA_DIR = 'reddit-data'
    CLEANED_REDDIT_DATA_DIR = 'cleaned-reddit-data'
    reddit_data_dirs = os.listdir(REDDIT_DATA_DIR)
    os.makedirs(CLEANED_REDDIT_DATA_DIR, exist_ok=True)
    for data_dir in reddit_data_dirs:
        if '-reddit-data' in data_dir:
            region = data_dir.split('-')[0]
            cleaned_dir_name = f'{region}-cleaned-reddit-data'
            cleaned_dir_path = os.path.join(CLEANED_REDDIT_DATA_DIR, cleaned_dir_name)
            os.makedirs(cleaned_dir_path, exist_ok=True)
            current_dir_path = os.path.join(REDDIT_DATA_DIR, data_dir)
            files = os.listdir(current_dir_path)
            for file in files:
                current_file_path = os.path.join(current_dir_path, file)
                df = read_csv_data(current_file_path, region)
                if df is not None:
                    write_df_to_csv(df, os.path.join(cleaned_dir_path, file))


if __name__ == '__main__':
    main()

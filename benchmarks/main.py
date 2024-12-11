import subprocess
import re, sys
from pathlib import Path
import matplotlib.pyplot as plt
import numpy as np
import datetime
from pathlib import Path


times = []
timestamps = []
tsets = []

postgres_times = []

def extract_payload(text, pattern, func):
    match = re.search(pattern, text)
    if match:
        payload = func(match.group(1))
    else:
        raise Exception

    return payload


def grep_mysql_time(fname):
    process = subprocess.run(["/home/nastya/dw-query-digest_0.9.6_linux_amd64/dw-query-digest", "-nocache", f"{fname}"], capture_output=True, text=True)
    output = process.stdout
    
    query = output.split("Query")[1] + output.split("Query")[2]
    pattern = r"mean time       : (\d+\.\d+)ms"
    time = extract_payload(query, pattern, float)
    times.append(time)
    
    pattern = r"Capture start      : (\d\d\d\d-\d\d-\d\d \d\d:\d\d:\d\d\.\d+) \+\d+ UTC"
    data = extract_payload(output, pattern, str)
    datetime_obj = datetime.datetime.strptime(data, "%Y-%m-%d %H:%M:%S.%f")
    timestamp_seconds = datetime_obj.timestamp()
    timestamps.append(timestamp_seconds)
    
    tsets.append((time, timestamp_seconds))
    

def calculate_percentiles(data, percentiles):
    return np.percentile(data, percentiles)


def plot_data(data, time_data, percentiles, filename):
    fig, axes = plt.subplots(2, 2, figsize=(10, 8))

    axes[0, 0].plot(time_data, data)
    axes[0, 0].set_xlabel('Время')
    axes[0, 0].set_ylabel('Время выполнения запроса (мс)')
    axes[0, 0].set_title('График времени выполнения запросов во времени')

    axes[0, 1].scatter(range(len(data)), data)
    axes[0, 1].set_xlabel('Номер запроса')
    axes[0, 1].set_ylabel('Время выполнения запроса (мс)')
    axes[0, 1].set_title('Распределение по перцентилям')
    for p in percentiles:
        axes[0, 1].axhline(y=calculate_percentiles(data, [p])[0], color='red', linestyle='--')

    axes[1, 0].hist(data, bins=20)
    axes[1, 0].set_xlabel('Время выполнения запроса (мс)')
    axes[1, 0].set_ylabel('Количество запросов')
    axes[1, 0].set_title('Гистограмма распределения')

    axes[1, 1].axis('off')
    axes[1, 1].text(0.05, 0.5, 
         f'Перцентили:\n'
         f'0.5: {calculate_percentiles(data, 50):.2f} мс\n'
         f'0.75: {calculate_percentiles(data, 75):.2f} мс\n'
         f'0.9: {calculate_percentiles(data, 90):.2f} мс\n'
         f'0.95: {calculate_percentiles(data, 95):.2f} мс\n'
         f'0.99: {calculate_percentiles(data, 99):.2f} мс', 
         fontsize=12)

    plt.tight_layout()
    plt.savefig(filename)
    
    
def process_mysql(path):
    p = Path(path)

    for x in p.rglob("*"):
        grep_mysql_time(x)

    sorted_data = sorted(tsets, key=lambda item: item[1])

    sorted_times = [it[0] for it in sorted_data]
    sorted_timestamps = [it[1] for it in sorted_data]
    percentiles = [0.5, 0.75, 0.9, 0.95, 0.99]

    plot_data(sorted_times, sorted_timestamps, percentiles, 'mysql.png')
    
    
def get_mean_time(fname):
    f = open(fname, "r")
    lines = f.readlines()
    f.close()
    times = []
    for line in lines:
        times.append(float(line.split(" ")[2]))   
        
    return sum(times)/len(times) * 0.001

    
def process_postgres(path):
    p = Path(path)

    paths = []
    for x in p.rglob("*"):
        paths.append(x)

    
    for fname in paths:
        postgres_times.append(get_mean_time(fname))
    percentiles = [0.5, 0.75, 0.9, 0.95, 0.99]
    timestamps = np.arange(10)
    
    plot_data(postgres_times, timestamps, percentiles, 'postgres.png')
    

if __name__ == "__main__":  
    args = sys.argv
    
    process_mysql(args[1])  
    process_postgres(args[2])
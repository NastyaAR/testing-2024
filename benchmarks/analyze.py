import matplotlib.pyplot as plt
import numpy as np
import argparse
import os

def analyze_data(filename):
    cpu_usage = []
    mem_usage = []

    with open(filename, 'r') as f:
        for line in f:
            cpu, mem = line.strip().split()
            cpu_usage.append(float(cpu.strip('%')))
            mem_usage.append(float(mem.strip('%')))

    cpu_max = np.max(cpu_usage)
    cpu_min = np.min(cpu_usage)
    cpu_avg = np.mean(cpu_usage)
    mem_max = np.max(mem_usage)
    mem_min = np.min(mem_usage)
    mem_avg = np.mean(mem_usage)

    return (cpu_max, cpu_min, cpu_avg, mem_max, mem_min, mem_avg)

def plot_data(filename, cpu_name, mem_name):
    cpu_usage = []
    mem_usage = []

    with open(filename, 'r') as f:
        for line in f:
            cpu, mem = line.strip().split()
            cpu_usage.append(float(cpu.strip('%')))
            mem_usage.append(float(mem.strip('%')))

    plt.figure(figsize=(10, 6))
    plt.plot(cpu_usage, label='CPU Usage')
    plt.xlabel('Iteration')
    plt.ylabel('Usage (%)')
    plt.title('CPU Usage')
    plt.legend()
    plt.savefig(cpu_name)
    plt.close()

    plt.figure(figsize=(10, 6))
    plt.plot(mem_usage, label='Memory Usage')
    plt.xlabel('Iteration')
    plt.ylabel('Usage (%)')
    plt.title('Memory Usage')
    plt.legend()
    plt.savefig(mem_name)
    plt.close()

if __name__ == "__main__":
    parser = argparse.ArgumentParser(description="Анализ данных CPU и Memory.")
    parser.add_argument("filename", help="Имя файла с данными.")
    parser.add_argument("cpu_name", help="Имя файла для сохранения графика CPU.")
    parser.add_argument("mem_name", help="Имя файла для сохранения графика Memory.")
    args = parser.parse_args()

    cpu_max, cpu_min, cpu_avg, mem_max, mem_min, mem_avg = analyze_data(args.filename)

    print(f"CPU Max: {cpu_max:.2f}%")
    print(f"CPU Min: {cpu_min:.2f}%")
    print(f"CPU Avg: {cpu_avg:.2f}%")
    print(f"Memory Max: {mem_max:.2f}%")
    print(f"Memory Min: {mem_min:.2f}%")
    print(f"Memory Avg: {mem_avg:.2f}%")

    plot_data(args.filename, args.cpu_name, args.mem_name)
from pathlib import Path
import matplotlib.pyplot as plt
import numpy as np
import datetime
from pathlib import Path

from main import calculate_percentiles
from main import plot_data

maxs_cpu = []
mins_cpu = []
means_cpu = []
    
maxs_ram = []
mins_ram = []
means_ram = []

def get_all_resources(path):
    maxs_cpu = []
    mins_cpu = []
    means_cpu = []
        
    maxs_ram = []
    mins_ram = []
    means_ram = []
    
    p = Path(path)
    paths = []
    for x in p.rglob("*"):
        paths.append(x)
        
    res = []
    for pt in paths:
        with open(pt, "r") as f:
            lines = f.readlines()
            
            cur = []
            for l in lines:
                if l[:3] != "CON":
                    info = l.split("   ")[2:4]
                    info[0] = float(info[0].strip("%"))
                    info[1] =''.join(i for i in info[1].split(" / ")[0] if not i.isalpha())
                    info[1] = float(info[1].strip())
                    res.append(info)
                    cur.append(info)
            maxs_cpu.append(max([cur[i][0] for i in range(len(cur))]))
            mins_cpu.append(min([cur[i][0] for i in range(len(cur))]))
            means_cpu.append(sum([cur[i][0] for i in range(len(cur))]) / len(cur))
            maxs_ram.append(max([cur[i][1] for i in range(len(cur))]))
            mins_ram.append(min([cur[i][1] for i in range(len(cur))]))
            means_ram.append(sum([cur[i][1] for i in range(len(cur))]) / len(cur))
                          
    return res

def postgres_info(path):
    info = get_all_resources(path)
    
    cpu = [info[i][0] for i in range(len(info))]
    ram = [info[i][1] for i in range(len(info))]
    
    percentiles = [0.5, 0.75, 0.9, 0.95, 0.99]
    timestamps = np.arange(len(cpu))
    
    plot_data(cpu, timestamps, percentiles, "cpu_graph.png")
    plot_data(ram, timestamps, percentiles, "ram_graph.png")
    
    plt.figure()    
    plt.plot(timestamps, cpu)
    plt.savefig("postgres_cpu.png")
    
    plt.close()
    
    plt.figure() 
    plt.plot(timestamps, )
    plt.savefig("postgres_ram.png")
    
    with open("postgres_cpu_report", "w") as f:
        f.write("num,max,min,mean\n")
        for i in range(len(maxs_cpu)):
            f.write(f"{i+1},{maxs_cpu[i]},{mins_cpu[i]},{means_cpu[i]}\n")
            
    with open("postgres_ram_report", "w") as f:
        f.write("num,max,min,mean\n")
        for i in range(len(maxs_ram)):
            f.write(f"{i+1},{maxs_ram[i]},{mins_ram[i]},{means_ram[i]}\n")
            
            
def mysql_info(path):
    info = get_all_resources(path)
    
    cpu = [info[i][0] for i in range(len(info))]
    ram = [info[i][1] for i in range(len(info))]
    
    percentiles = [0.5, 0.75, 0.9, 0.95, 0.99]
    timestamps = np.arange(len(cpu))
    
    plot_data(cpu, timestamps, percentiles, "mysql_cpu_graph.png")
    plot_data(ram, timestamps, percentiles, "mysql_ram_graph.png")
    
    plt.figure()    
    plt.plot(timestamps, cpu)
    plt.savefig("mysql_cpu.png")
    
    plt.close()
    
    plt.figure() 
    plt.plot(timestamps, ram)
    plt.savefig("mysql_ram.png")
    
    with open("mysql_cpu_report", "w") as f:
        f.write("num,max,min,mean\n")
        for i in range(len(maxs_cpu)):
            f.write(f"{i+1},{maxs_cpu[i]},{mins_cpu[i]},{means_cpu[i]}\n")
            
    with open("mysql_ram_report", "w") as f:
        f.write("num,max,min,mean\n")
        for i in range(len(maxs_ram)):
            f.write(f"{i+1},{maxs_ram[i]},{mins_ram[i]},{means_ram[i]}\n")
            

if __name__  == "__main__":
    postgres_info("postgres_res_usage2")
    mysql_info("mysql_res_usage")
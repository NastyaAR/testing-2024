import subprocess
import time
import os
import glob

class MySQL_Test:
    def __init__(self, password, iterations=1):
        self.password = password
        self.iterations = iterations
        self.command = "while true; do docker stats --no-stream | grep test-2024-mysql_db-1 | awk '{print $3,$7}' >> mysql_time.txt; done"


    def get_current_file(self, i):
        return self.command.replace("time", str(i))


    def test_database(self, i, func):
        command = f"echo {self.password} | sudo -S rm -rf ../dbdata"
        subprocess.run(command, shell=True)
        
        subprocess.run(["docker", "compose", "up", "-d", "mysql_db", "mysql-exporter"])
        time.sleep(40)

        self.grant_privileges()
        self.set_slow_query_log()
        self.create_test_database()

        process = subprocess.Popen(self.get_current_file(i), shell=True)

        func()
        
        process.kill()

        self.copy_slow_logs()

        subprocess.run(["docker", "compose", "stop", "mysql_db", "mysql-exporter"])
        subprocess.run(["docker", "compose", "rm", "-f", "mysql_db", "mysql-exporter"])

    def grant_privileges(self):
        subprocess.run(
            [
                "mysql",
                "-P",
                "3306",
                "-h",
                "0.0.0.0",
                "-u",
                "root",
                "--password=secret",
                "--execute=GRANT ALL PRIVILEGES ON *.* TO example WITH GRANT OPTION;",
            ]
        )

    def set_slow_query_log(self):
        subprocess.run(
            [
                "mysql",
                "-P",
                "3306",
                "-h",
                "0.0.0.0",
                "-u",
                "example",
                "--password=secret2",
                "--execute=SET GLOBAL slow_query_log = 'ON';",
            ]
        )
        subprocess.run(
            [
                "mysql",
                "-P",
                "3306",
                "-h",
                "0.0.0.0",
                "-u",
                "example",
                "--password=secret2",
                "--execute=SET GLOBAL long_query_time = 0;",
            ]
        )

    def create_test_database(self):
        print("CREATE")
        subprocess.run(
            [
                "mysql",
                "-P",
                "3306",
                "-h",
                "0.0.0.0",
                "-u",
                "example",
                "--password=secret2",
                "-e",
                "create table test (id serial primary key, name text, surname text, fl float, num bigint);",
                "test",
            ]
        )
        print("SUCCESS")

    def insert_test(self):
        print("MEASURE")
        subprocess.run(
            [
                "mysqlslap",
                "--host",
                "0.0.0.0",
                "--port",
                "3306",
                "--user",
                "example",
                "--create-schema",
                "test",
                "--delimiter=';'",
                "--query=insert.sql",
                "--iterations",
                str(self.iterations),
                "--password=secret2",
            ]
        )
        time.sleep(15)
        
    def index_test(self):
        subprocess.run(
            [
                "mysqlslap",
                "--host",
                "0.0.0.0",
                "--port",
                "3306",
                "--user",
                "example",
                "--create-schema",
                "test",
                "--delimiter=';'",
                "--query=SELECT sum(fl) FROM test group by num",
                "--iterations",
                str(self.iterations),
                "--password=secret2",
            ]
        )
        time.sleep(15)


    def copy_slow_logs(self):
        source_dir = "/home/nastya/test-2024/dbdata"
        slow_log_files = glob.glob(os.path.join(source_dir, "*-slow.log"))
        
        for file in slow_log_files:
            os.system(f"echo {self.password} | sudo -S cp {file} .")
 
    def change_ownership(self, username):
        os.system(f"echo {self.password} | sudo -S chown {username}:{username} *-slow.log")
        


if __name__ == "__main__":
    password = os.environ.get("PASS")
    mysql_setup = MySQL_Test(password, 11)

    for i in range(11):
        mysql_setup.test_database(i, mysql_setup.insert_test)
    mysql_setup.change_ownership("nastya")
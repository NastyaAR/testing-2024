FROM ubuntu:jammy-20240911.1
RUN useradd tester -m -d /home/tester
RUN apt update && apt install -y git wget make curl tar
RUN wget https://go.dev/dl/go1.23.1.linux-amd64.tar.gz &&  rm -rf /usr/local/go && tar -C /usr/local -xzf go1.23.1.linux-amd64.tar.gz
RUN mkdir migrate && cd migrate && curl -L https://github.com/golang-migrate/migrate/releases/download/v4.18.1/migrate.linux-amd64.tar.gz | tar xvz
ENV PATH=$PATH:/migrate:/usr/local/go/bin
RUN git clone https://github.com/NastyaAR/testing-2024.git -b lab02
ENTRYPOINT ["/bin/bash", "-c"]


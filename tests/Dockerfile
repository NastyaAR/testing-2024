FROM ubuntu:jammy-20240911.1
RUN useradd tester -m -d /home/tester
RUN apt update && apt install -y git wget
RUN wget https://go.dev/dl/go1.23.1.linux-amd64.tar.gz &&  rm -rf /usr/local/go && tar -C /usr/local -xzf go1.23.1.linux-amd64.tar.gz
ENV PATH=$PATH:/usr/local/go/bin
RUN git clone https://github.com/NastyaAR/testing-2024.git -b lab02
WORKDIR /testing-2024
ENTRYPOINT ["make test"]


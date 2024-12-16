FROM ubuntu:jammy-20240911.1
RUN useradd tester -m -d /home/tester
RUN apt update && apt install -y git wget make curl tar openjdk-11-jdk
RUN wget https://go.dev/dl/go1.23.1.linux-amd64.tar.gz &&  rm -rf /usr/local/go && tar -C /usr/local -xzf go1.23.1.linux-amd64.tar.gz
RUN mkdir migrate && cd migrate && curl -L https://github.com/golang-migrate/migrate/releases/download/v4.18.1/migrate.linux-amd64.tar.gz | tar xvz
ENV PATH=$PATH:/migrate:/usr/local/go/bin
RUN git clone https://github.com/NastyaAR/testing-2024.git -b lab02
ARG postgres_test_host
ARG postgres_test_port
ARG token
ARG password
ENV POSTGRES_TEST_HOST=$postgres_test_host
ENV POSTGRES_TEST_PORT=$postgres_test_port
ENV TOKEN=$token
ENV PASSWORD=$password
RUN wget https://github.com/allure-framework/allure2/releases/download/2.30.0/allure-2.30.0.tgz && tar -xzf allure-2.30.0.tgz
ENV PATH=$PATH:/migrate:/usr/local/go/bin:/allure-2.30.0/bin:/usr/bin
ENTRYPOINT ["/bin/bash", "-c"]


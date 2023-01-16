FROM ubuntu:20.04
RUN apt-get update && apt-get install -y apt-utils
RUN apt-get install -y --no-install-recommends --no-install-suggests --fix-missing curl ca-certificates
RUN curl --proto '=https' --tlsv1.2 -sSf https://starship.rs/install.sh | sh -s -- -y
FROM ubuntu:latest

RUN apt-get update && apt-get install -y \
    sudo \
    nano \
    curl \
    wget \
    vim \
    iputils-ping \
    && apt-get clean

RUN useradd -m user && echo "user:user" | chpasswd && usermod -aG sudo user

WORKDIR /home/user
USER user

CMD ["/bin/bash"]
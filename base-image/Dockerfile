FROM alpine:3.19

RUN apk add --no-cache \
    curl \
    wget \
    git \
    vim \
    nano \
    bash \
    ca-certificates \
    build-base \
    python3 \
    py3-pip \
    nodejs \
    npm \
    coreutils \
    findutils \
    grep \
    sed \
    awk \
    netcat-openbsd \
    bind-tools \
    less \
    jq \
    tar \
    unzip \
    sudo

# non ron user
RUN adduser -D -s /bin/bash box && \
    echo "box ALL=(ALL) NOPASSWD:ALL" >> /etc/sudoers

WORKDIR /workspace
USER box

CMD ["/bin/bash"]

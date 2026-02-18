FROM alpine:3.19


RUN apk add --no-cache \
    bash \
    curl \
    wget \
    vim \
    nano \
    coreutils \
    findutils \
    grep \
    sed \
    gawk \
    netcat-openbsd \
    bind-tools \
    less \
    jq \
    tar \
    unzip \
    sudo \
    && rm -rf /var/cache/apk/*


RUN adduser -D -s /bin/bash box \
    && echo "box ALL=(ALL) NOPASSWD:ALL" >> /etc/sudoers


WORKDIR /box


USER box


ENV PS1="box:\w$ "


CMD ["/bin/bash"]

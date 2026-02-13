FROM alpine:3.19

# Install essential Linux tools for practice
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

# Create non-root user 'box' with sudo access
RUN adduser -D -s /bin/bash box \
    && echo "box ALL=(ALL) NOPASSWD:ALL" >> /etc/sudoers

# Set working directory
WORKDIR /box

# Switch to non-root user
USER box

# Set PS1 prompt to show current directory
ENV PS1="box:\w$ "

# Default command
CMD ["/bin/bash"]

FROM alpine:3.20

# 기본 패키지 설치
RUN apk add --no-cache \
    bash \
    curl \
    wget \
    busybox-extras \
    bind-tools \
    redis \
    mysql-client \
    docker-cli \
    nmap \
    tcpdump \
    iproute2 \
    net-tools

# kubectl 설치
RUN curl -LO https://dl.k8s.io/release/$(curl -sL https://dl.k8s.io/release/stable.txt)/bin/linux/amd64/kubectl && \
    install -m 0755 kubectl /usr/bin/kubectl && \
    rm kubectl

# kafkactl 설치
RUN KAFKACTL_VERSION=v4.0.0 && \
    curl -L -o kafkactl "https://github.com/deviceinsight/kafkactl/releases/download/${KAFKACTL_VERSION}/kafkactl_linux_amd64" && \
    install -m 0755 kafkactl /usr/local/bin/kafkactl && \
    rm kafkactl

# 기본 쉘을 bash로 설정
SHELL ["/bin/bash", "-c"]

ENTRYPOINT ["bash"]

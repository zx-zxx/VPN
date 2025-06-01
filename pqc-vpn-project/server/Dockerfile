FROM golang:1.24

WORKDIR /app

# 安装依赖
RUN apt-get update && apt-get install -y \
    gcc \
    libc-dev \
    iptables \
    tcpdump \
    iproute2 \
    iputils-ping \  
    iperf3 \        
    && rm -rf /var/lib/apt/lists/*  
    # 清理缓存，减少镜像大小

# 配置 Go 代理
ENV GO111MODULE=on
ENV GOPROXY=https://goproxy.cn,direct

# 初始化 Go 模块
RUN go mod init pqc-vpn-server

# 安装依赖
RUN go get github.com/songgao/water 

# 复制代码和证书
COPY server.go .
COPY *.crt .
COPY *.key .

# 编译代码
RUN go build -o pqc-vpn-server server.go

# 运行服务器
CMD ["./pqc-vpn-server"]
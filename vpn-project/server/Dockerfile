FROM golang:1.23-alpine

WORKDIR /app

# 安装依赖
RUN apk add --no-cache gcc musl-dev iptables tcpdump

# 配置 Go 代理
ENV GO111MODULE=on
ENV GOPROXY=https://goproxy.cn,direct

# 初始化 Go 模块
RUN go mod init vpn-server

# 安装依赖
RUN go get github.com/songgao/water

# 复制代码和证书
COPY server.go .
COPY *.crt .
COPY *.key .

# 编译代码
RUN go build -o vpn-server server.go

# 运行服务器
CMD ["./vpn-server"]
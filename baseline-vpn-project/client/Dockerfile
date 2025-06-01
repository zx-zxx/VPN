FROM golang:1.23-alpine

WORKDIR /app

# 安装依赖
RUN apk add --no-cache gcc musl-dev tcpdump

# 配置 Go 代理
ENV GO111MODULE=on
ENV GOPROXY=https://goproxy.cn,direct

# 初始化 Go 模块
RUN go mod init baseline-vpn-client

# 安装依赖
RUN go get github.com/songgao/water

# 复制代码和证书
COPY client.go .
COPY *.crt .
COPY *.key .

# 编译代码
RUN go build -o baseline-vpn-client client.go

# 运行客户端
CMD ["./baseline-vpn-client"]
package main

import (
    "crypto/tls"
    "crypto/x509"
    "io/ioutil"
    "log"
    "net"
    "github.com/songgao/water"
    "os/exec"
)

func main() {
    // 创建TUN接口
    ifce, err := water.New(water.Config{DeviceType: water.TUN})
    if err != nil {
        log.Fatalf("Failed to create TUN interface: %v", err)
    }
    log.Printf("TUN interface name: %s", ifce.Name())

    // 自动配置TUN接口
    cmd := exec.Command("ip", "link", "set", "tun0", "up")
    if err := cmd.Run(); err != nil {
        log.Printf("Failed to set tun0 up: %v", err)
    }
    cmd = exec.Command("ip", "addr", "add", "10.2.0.1/24", "dev", "tun0")
    if err := cmd.Run(); err != nil {
        log.Printf("Failed to add IP to tun0: %v", err)
    }

    cert, err := tls.LoadX509KeyPair("/app/server.crt", "/app/server.key")
    if err != nil {
        log.Fatalf("Failed to load server certificate: %v", err)
    }

    clientCert, err := ioutil.ReadFile("/app/client.crt")
    if err != nil {
        log.Fatalf("Failed to read client certificate: %v", err)
    }
    certPool := x509.NewCertPool()
    certPool.AppendCertsFromPEM(clientCert)

    config := &tls.Config{
        Certificates: []tls.Certificate{cert},
        ClientAuth:   tls.RequireAndVerifyClientCert,
        ClientCAs:    certPool,
        MinVersion:   tls.VersionTLS13,
        MaxVersion:   tls.VersionTLS13, // 限制最高版本为TLS 1.3
        SessionTicketsDisabled: true,
        CurvePreferences: []tls.CurveID{
            tls.X25519, // 只使用经典算法
        },
    }

    listener, err := tls.Listen("tcp", ":7443", config)
    if err != nil {
        log.Fatalf("Failed to listen: %v", err)
    }
    defer listener.Close()

    log.Println("Server listening on :7443")
    for {
        conn, err := listener.Accept()
        if err != nil {
            log.Printf("Failed to accept connection: %v", err)
            continue
        }
        log.Println("Client connected")
        go handleConnection(conn, ifce)
    }
}

func handleConnection(conn net.Conn, ifce *water.Interface) {
    defer conn.Close()
    // 从客户端读取流量并写入TUN接口
    go func() {
        buf := make([]byte, 1500)
        for {
            n, err := conn.Read(buf)
            if err != nil {
                log.Printf("Failed to read from client: %v", err)
                return
            }
            log.Printf("Received %d bytes from client", n)
            _, err = ifce.Write(buf[:n])
            if err != nil {
                log.Printf("Failed to write to TUN: %v", err)
                continue
            }
            log.Printf("Wrote %d bytes to TUN interface", n)
        }
    }()

    // 从TUN接口读取流量并发送回客户端
    buf := make([]byte, 1500)
    for {
        n, err := ifce.Read(buf)
        if err != nil {
            log.Printf("Failed to read from TUN: %v", err)
            continue
        }
        log.Printf("Read %d bytes from TUN", n)
        _, err = conn.Write(buf[:n])
        if err != nil {
            log.Printf("Failed to write to client: %v", err)
            continue
        }
    }
}
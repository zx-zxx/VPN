package main

import (
    "crypto/tls"
    "crypto/x509"
    "io/ioutil"
    "log"
    "github.com/songgao/water"
    "os/exec"
)

func main() {
    // 配置TUN接口
    ifce, err := water.New(water.Config{DeviceType: water.TUN})
    if err != nil {
        log.Fatalf("Failed to create TUN interface: %v", err)
    }
    log.Printf("TUN interface name: %s", ifce.Name())

    // 自动配置TUN接口，使用不同的 TUN 子网，例如 10.1.0.0/24
    cmd := exec.Command("ip", "link", "set", "tun0", "up")
    if err := cmd.Run(); err != nil {
        log.Printf("Failed to set tun0 up: %v", err)
    }
    cmd = exec.Command("ip", "addr", "add", "10.1.0.2/24", "dev", "tun0")   
    if err := cmd.Run(); err != nil {
        log.Printf("Failed to add IP to tun0: %v", err)
    }
    cmd = exec.Command("ip", "route", "del", "default", "via", "192.168.2.1", "dev", "eth0")
    if err := cmd.Run(); err != nil {
        log.Printf("Failed to delete previous default route eth0: %v", err)
    }
    cmd = exec.Command("ip", "route", "add", "default", "via", "10.1.0.1", "dev", "tun0")
    if err := cmd.Run(); err != nil {
        log.Printf("Failed to add route: %v", err)
    }

    // 加载客户端证书和私钥
    cert, err := tls.LoadX509KeyPair("/app/client.crt", "/app/client.key")
    if err != nil {
        log.Fatalf("Failed to load client certificate: %v", err)
    }

    // 加载服务器证书
    serverCert, err := ioutil.ReadFile("/app/server.crt")
    if err != nil {
        log.Fatalf("Failed to read server certificate: %v", err)
    }
    certPool := x509.NewCertPool()
    certPool.AppendCertsFromPEM(serverCert)

    log.Println("TLS VersionTLS13:", tls.VersionTLS13)

    // 配置TLS
    config := &tls.Config{
        Certificates: []tls.Certificate{cert},
        RootCAs:      certPool,
        MinVersion:   tls.VersionTLS13,
        MaxVersion:   tls.VersionTLS13,
        SessionTicketsDisabled: true,
        // 优先使用 X25519Kyber768Draft00，支持 X25519 作为回退
        CurvePreferences: []tls.CurveID{
            tls.X25519MLKEM768, // 后量子混合密钥交换
            tls.X25519,                // 经典密钥交换（回退）
        },
    }

    // 连接到服务器，使用不同的端口，例如 9443
    conn, err := tls.Dial("tcp", "pqc-vpn-server:9443", config)
    if err != nil {
        log.Fatalf("Failed to connect: %v", err)
    }
    log.Println("Supported Versions:", config.MinVersion, config.MaxVersion)
    log.Println("Negotiated Cipher Suite:", conn.ConnectionState().CipherSuite)
    log.Println("Negotiated Version:", conn.ConnectionState().Version)
    // log.Println("Negotiated Curve/KEM:", conn.ConnectionState().CurveID) // 添加日志查看使用的密钥交换机制
    if err := conn.Handshake(); err != nil {
        log.Println("Handshake Error:", err)
    }
    defer conn.Close()

    // 读取TUN接口的流量并发送到服务器
    go func() {
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
                log.Printf("Failed to write to server: %v", err)
                continue
            }
        }
    }()

    // 读取服务器的回复
    buf := make([]byte, 1500)
    for {
        n, err := conn.Read(buf)
        if err != nil {
            log.Fatalf("Failed to read from server: %v", err)
        }
        log.Printf("Received %d bytes from server", n)
        _, err = ifce.Write(buf[:n])
        if err != nil {
            log.Printf("Failed to write to TUN: %v", err)
            continue
        }
    }
}
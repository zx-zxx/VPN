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

    // 自动配置TUN接口
    setupTunInterface()

    // TLS配置
    config := setupTLSConfig()

    // 连接到服务器（使用服务名）
    conn, err := tls.Dial("tcp", "baseline-vpn-server:7443", config)
    if err != nil {
        log.Fatalf("Failed to connect: %v", err)
    }
    log.Println("Supported Versions:", config.MinVersion, config.MaxVersion)
    log.Println("Negotiated Cipher Suite:", conn.ConnectionState().CipherSuite)
    log.Println("Negotiated Version:", conn.ConnectionState().Version)
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

func setupTunInterface() {
    commands := [][]string{
        {"ip", "link", "set", "tun0", "up"},
        {"ip", "addr", "add", "10.2.0.2/24", "dev", "tun0"},
        {"ip", "route", "del", "default", "via", "192.168.3.1", "dev", "eth0"},
        {"ip", "route", "add", "default", "via", "10.2.0.1", "dev", "tun0"},
    }
    
    for _, cmdArgs := range commands {
        cmd := exec.Command(cmdArgs[0], cmdArgs[1:]...)
        if err := cmd.Run(); err != nil {
            log.Printf("Failed to execute %v: %v", cmdArgs, err)
        }
    }
} 

func setupTLSConfig() *tls.Config {
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

    return &tls.Config{
        Certificates: []tls.Certificate{cert},
        RootCAs:      certPool,
        MinVersion:   tls.VersionTLS13,
        MaxVersion:   tls.VersionTLS13,
        SessionTicketsDisabled: true, // 禁用会话复用，确保每次都握手
        CurvePreferences: []tls.CurveID{
            tls.X25519, // 只使用经典算法
        },
    }
}
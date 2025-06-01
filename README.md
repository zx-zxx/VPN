# 写在最前面

压缩包中有三个文件夹，分别为 `baseline-vpn-project`、`vpn-project`、`pqc-vpn-project`。

（`image` 文件夹只是存放该 README.md 的图片）

因为选择的密钥交换算法不同，环境有细微差异，与其把选择算法的代码杂糅在一起，不如创建三个项目，从头开始来的更清晰明了且快捷；又因为使用`Docker`，环境配置可以自己通过代码控制，使这种实现方式成为可能。

三个项目选择的密钥交换算法如下：

- `baseline-vpn-project`：**X25519**
- `vpn-project`：**X25519Kyber768Draft00**
- `baseline-vpn-project`：**X25519MLKEM768**



# 1. baseline-vpn-project

**所需条件：**Docker 环境

## 1.1 Build and Run

1. 进入 /baseline-vpn-project

2. 在该文件夹下进入控制台

3. 输入命令

   ```shell
   docker-compose up
   ```

   - 构建镜像并启动容器

## 1.2 Golang version

**Go 1.23-alphine**

（alphine 版本体积较小）

## 1.3 External Libraries

- `github.com/songgao/water`，该库用于创建 `TUN` 接口
  - 安装方法：运行完构建容器的命令后，会自动运行
  
    ```shell
    go get github.com/songgao/water
    ```
  
    ![water库](.\images\water库.png)

## 1.4 Key Share Algorithm

**X25519**                                                <img src=".\images\x25519.png" alt="x25519" style="zoom:80%;" />



# 2. vpn-project

## 2.1 Build and Run

同 1.1 节

## 2.2 Golang version

同 1.2 节

## 2.3 External Libraries

同 1.3 节

## 2.4 Key Share Algorithm

**X25519Kyber768Draft00**          <img src=".\images\X25519Kyber768Draft00.png" alt="X25519Kyber768Draft00" style="zoom:80%;" />



# 3. pqc-vpn-project

## 3.1 Build and Run

同 1.1 节

## 3.2 Golang version

**Go 1.24**

## 3.3 External Libraries

同 1.3 节

## 3.4 Key Share Algorithm

**X25519MLKEM768**                    <img src=".\images\X25519MLKEM768.png" alt="X25519MLKEM768" style="zoom:80%;" />

可能是 Wireshark 的版本问题，它无法识别 X25519MLKEM768，但根据 4588 的编号，查询 Go 1.24 的[官方文档](https://pkg.go.dev/crypto/tls#CurveID)可知，4588 就是指 X25519MLKEM768。

<img src=".\images\官方文档.png" alt="官方文档" style="zoom:50%;" />
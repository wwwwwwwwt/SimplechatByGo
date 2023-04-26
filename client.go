/*
 * @Author: zzzzztw
 * @Date: 2023-04-26 16:56:56
 * @LastEditors: Do not edit
 * @LastEditTime: 2023-04-26 20:47:36
 * @FilePath: /zhang/SimpleChatByGo/client.go
 */
package main

import (
	"flag"
	"fmt"
	"net"
)

type Client struct {
	ServerIp   string
	ServerPort int
	Name       string
	conn       net.Conn
}

func NewClient(Ip string, Port int) (*Client, error) {

	// 创建客户端对象

	client := &Client{ServerIp: Ip, ServerPort: Port}

	// 连接server
	con, err := net.Dial("tcp", fmt.Sprintf("%s:%d", Ip, Port))

	if err != nil {
		fmt.Println("net.dial is error:", err)
	}

	client.conn = con
	// 返回对象
	return client, err

}

var serverIp string
var serverPort int

func init() {
	flag.StringVar(&serverIp, "ip", "127.0.0.1", "默认ip地址")
	flag.IntVar(&serverPort, "port", 8888, "默认端口地址")
}

func main() {

	flag.Parse()

	client, err := NewClient(serverIp, serverPort)

	if err != nil || client == nil {
		fmt.Println(">>>>> 连接客户端失败")
		return
	}

	fmt.Println(">>>>> 连接客户端成功")

	// 启动客户端的业务

	select {}
}

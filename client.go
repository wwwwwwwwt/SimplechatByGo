/*
 * @Author: zzzzztw
 * @Date: 2023-04-26 16:56:56
 * @LastEditors: Do not edit
 * @LastEditTime: 2023-04-26 17:04:16
 * @FilePath: /zhang/SimpleChatByGo/client.go
 */
package main

import (
	"fmt"
	"net"
)

type Client struct {
	ServerIp   string
	ServerPort int
	Name       string
	conn       net.Conn
}

func NewClient(Ip string, Port int) *Client {

	// 创建客户端对象

	client := &Client{ServerIp: Ip, ServerPort: Port}

	// 连接server
	con, err := net.Dial("tcp", fmt.Sprintf("%s:%d", Ip, Port))

	if err != nil {
		fmt.Println("net.dial is error:", err)
	}

	client.conn = con
	// 返回对象
	return client

}

func main() {
	client := NewClient("127.0.0.1", 8888)

	if client == nil {
		fmt.Println(">>>>> 连接客户端失败")
		return
	}

	fmt.Println(">>>>> 连接客户端成功")

	// 启动客户端的业务

	select {}
}

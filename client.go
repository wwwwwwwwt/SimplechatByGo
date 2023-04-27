/*
 * @Author: zzzzztw
 * @Date: 2023-04-26 16:56:56
 * @LastEditors: Do not edit
 * @LastEditTime: 2023-04-27 16:57:14
 * @FilePath: /zhang/SimpleChatByGo/client.go
 */
package main

import (
	"flag"
	"fmt"
	"io"
	"net"
	"os"
)

type Client struct {
	ServerIp   string
	ServerPort int
	Name       string
	conn       net.Conn
	flag       int
}

func NewClient(Ip string, Port int) (*Client, error) {

	// 创建客户端对象

	client := &Client{ServerIp: Ip, ServerPort: Port}
	client.flag = 999

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

func (client *Client) menu() bool {
	var flag int
	fmt.Println("1.公聊模式")
	fmt.Println("2.私聊模式")
	fmt.Println("3.更新用户名")
	fmt.Println("4.查看当前在线用户")
	fmt.Println("0.退出")

	fmt.Scanln(&flag)

	if flag >= 0 && flag <= 4 {
		client.flag = flag
		return true
	} else {
		fmt.Println("请输入合法范围内的数字")
		return false
	}

}

func (client *Client) DealResponese() {
	// 一旦client.conn有数据,就拷贝到std标准输出上，并永久阻塞等待
	io.Copy(os.Stdout, client.conn)
	/* io.Copy(os.Stdout, client.conn)等同于
	for {
		buf := make([]byte, 4096)
		client.conn.Read(buf)
		fmt.Println(buf)
	}
	*/

}

func (client *Client) UpdateName() bool {
	// 提示用户输入用户名

	fmt.Println(">>>>请输入用户名")
	fmt.Scanln(&client.Name)

	sendMsg := "rename|" + client.Name + "\n"
	_, err := client.conn.Write([]byte(sendMsg))
	if err != nil {
		fmt.Println("更改用户名失败", err)
		return false
	}
	fmt.Println("更改用户名成功")
	return true

}

func (Client *Client) CheckOnline() {

	sendMsg := "who" + "\n"

	_, err := Client.conn.Write([]byte(sendMsg))

	if err != nil {
		fmt.Println("查找当前在线用户失败", err)
	}

}

func (client *Client) PublicChat() {
	// 提示用户输入内容
	var chatMsg string
	fmt.Println(">>>>请输入聊天内容， exit退出")
	fmt.Scanln(&chatMsg)
	//发给服务器
	for chatMsg != "exit" {
		if len(chatMsg) > 0 {
			sendMsg := chatMsg + "\n"
			_, err := client.conn.Write([]byte(sendMsg))
			if err != nil {
				fmt.Println("conn Write err:", err)
				break
			}
		}
		chatMsg = ""
		fmt.Println(">>>>请输入聊天内容， exit退出")
		fmt.Scanln(&chatMsg)

	}

	return
}

func (client *Client) PrivateChat() {
	client.CheckOnline()
	var name string
	var content string
	fmt.Println(">>>>请输入聊天对象name,exit退出")
	fmt.Scanln(&name)
	for name != "exit" {

		fmt.Println(">>>>请输入聊天内容, exit退出")
		fmt.Scanln(&content)
		for content != "exit" {
			sendMsg := "to|" + name + "|" + content + "\n\n"
			_, err := client.conn.Write([]byte(sendMsg))

			if err != nil {
				fmt.Println("私聊发送失败", err)
				break
			}
			content = ""
			fmt.Println(">>>>请输入聊天内容, exit退出")
			fmt.Scanln(&content)
		}
		client.CheckOnline()
		name = ""
		fmt.Println(">>>>请输入聊天对象name,exit退出")
		fmt.Scanln(&name)

	}

}

func (client *Client) Run() {
	for client.flag != 0 {
		for client.menu() != true {

		}

		switch client.flag {
		case 1:
			//公聊
			fmt.Println("公聊模式选择")
			client.PublicChat()
		case 2:
			//私聊
			fmt.Println("私聊模式选择")

			client.PrivateChat()
		case 3:
			//更改用户名
			fmt.Println("请输入更改名字")
			client.UpdateName()
		case 4:
			client.CheckOnline()
		}

	}
}

func main() {

	flag.Parse()

	client, err := NewClient(serverIp, serverPort)

	if err != nil || client == nil {
		fmt.Println(">>>>> 连接客户端失败")
		return
	}

	//单独开一个gorountine去处理server的回执消息
	go client.DealResponese()
	fmt.Println(">>>>> 连接客户端成功")

	// 启动客户端的业务

	client.Run()
}

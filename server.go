/*
 * @Author: zzzzztw
 * @Date: 2023-04-25 13:59:08
 * @LastEditors: Do not edit
 * @LastEditTime: 2023-04-26 16:27:14
 * @FilePath: /zhang/SimpleChatByGo/server.go
 */

package main

import (
	"fmt"
	"io"
	"net"
	"sync"
	"time"
)

type Server struct {
	Ip   string
	Port int

	//用于记录当前在线用户
	OnlineMap map[string]*User
	maplock   sync.RWMutex

	//用于消息广播的channel
	Message chan string
}

// 创建Server接口

func NewServer(ip string, port int) *Server {
	server := &Server{ip, port, make(map[string]*User), sync.RWMutex{}, make(chan string)}
	return server
}

//广播消息

func (t *Server) Broadcast(user *User, msg string) {
	sendMsg := "[" + user.Addr + "]" + user.Name + ":" + msg

	t.Message <- sendMsg
}

// 服务器接收消息

func (t *Server) ListenMessager() {
	for {
		msg := <-t.Message
		// 将msg发送给全部的在线User
		t.maplock.Lock()
		for _, cli := range t.OnlineMap {
			cli.C <- msg

		}
		t.maplock.Unlock()
	}
}

// conn后的handler 方法
func (t *Server) Handler(conn net.Conn) {

	// 处理当前连接业务
	fmt.Println("conn is accept")

	//新连接用户加入onlinemap
	user := NewUser(conn, t)

	// 接受客户端发送的消息，目前通过nc 模拟

	user.Online()
	alive := make(chan bool)
	go func() {
		buf := make([]byte, 4096)
		for {
			n, err := conn.Read(buf)
			if n == 0 {
				user.Offline()
				return
			}

			if err != nil && err != io.EOF {
				fmt.Println("Conn Read err:", err)
				return
			}

			// 提取用户的消息,去除"\n"
			msg := string(buf[:n-1])

			//
			user.Domessage(msg)

			alive <- true

		}
	}()

	//阻塞住，防止user掉线

	for {
		select {
		case <-alive:
		case <-time.After(200 * time.Second): // case判断时会将其重置时间
			//已经超时
			//将当前的User强制关闭
			user.SendMsg("长时间不发消息，已被踢出")
			close(user.C)
			user.Offline()
			conn.Close()

		}
	}
}

// 启动server接口

func (t *Server) Start() {

	//socket
	//bind 直接包含到Listen中了
	//listen

	listener, err := net.Listen("tcp", fmt.Sprintf("%s:%d", t.Ip, t.Port))

	if err != nil {
		fmt.Println("net.Listen error", err)
		return
	}
	//close()

	defer listener.Close()

	// 启动监听message 的goroutine
	go t.ListenMessager()

	//accept

	for {
		conn, err := listener.Accept()
		if err != nil {
			fmt.Println("listener accept error", err)
			continue
		}

		go t.Handler(conn) // 开一个协程去操作这个连接

	}

}

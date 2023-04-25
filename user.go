/*
 * @Author: zzzzztw
 * @Date: 2023-04-25 14:56:39
 * @LastEditors: Do not edit
 * @LastEditTime: 2023-04-25 20:49:09
 * @FilePath: /zhang/SimpleChatByGo/user.go
 */
package main

import "net"

type User struct {
	Name   string
	Addr   string
	C      chan string
	conn   net.Conn
	server *Server
}

// 创建监听消息的Handler，基于channel

func (t *User) ListenMessage() {

	for {
		msg := <-t.C

		t.conn.Write([]byte(msg + "\n"))
	}

}

// 创建一个新的客户端，绑定客户端地址
func NewUser(conn net.Conn, server *Server) *User {

	userAddr := conn.RemoteAddr().String()

	user := &User{userAddr, userAddr, make(chan string), conn, server}
	go user.ListenMessage() // 给每个对象都创造一个监听gorountine
	return user
}

func (t *User) Online() {
	t.server.maplock.Lock()
	t.server.OnlineMap[t.Name] = t
	t.server.maplock.Unlock()

	//用户上线广播
	t.server.Broadcast(t, "已上线")
}

func (t *User) Offline() {

	t.server.maplock.Lock()
	delete(t.server.OnlineMap, t.Name)
	t.server.maplock.Unlock()

	//用户上线广播
	t.server.Broadcast(t, "已下线")
}

func (t *User) Domessage(msg string) {
	t.server.Broadcast(t, msg)
}

/*
 * @Author: zzzzztw
 * @Date: 2023-04-25 14:56:39
 * @LastEditors: Do not edit
 * @LastEditTime: 2023-04-26 14:29:47
 * @FilePath: /zhang/SimpleChatByGo/user.go
 */
package main

import (
	"net"
	"strings"
)

type User struct {
	Name   string
	Addr   string
	C      chan string
	conn   net.Conn
	server *Server
}

// 创建监听消息的Handler，基于channel

func (t *User) ListenMessage() {
	//msg := range t.C
	for {
		//msg := <-t.C
		msg, ok := <-t.C
		if !ok {
			return
		}
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

	if msg == "who" {
		// 查询当前在线用户
		t.server.maplock.Lock()

		for _, user := range t.server.OnlineMap {
			onlineMsg := "[" + user.Addr + "]" + user.Name + ":" + "在线...\n"
			t.SendMsg(onlineMsg)
		}

		t.server.maplock.Unlock()

	} else if len(msg) > 7 && msg[:7] == "rename|" {
		//定义消息格式： rename|wgl
		newname := strings.Split(msg, "|")[1]
		if _, ok := t.server.OnlineMap[newname]; ok {
			t.SendMsg("当前用户名已经被占用\n")
		} else {
			t.server.maplock.Lock()
			delete(t.server.OnlineMap, t.Name)
			t.server.OnlineMap[newname] = t // 修改在线key-val
			t.server.maplock.Unlock()

			t.Name = newname
			t.SendMsg("您已经更新用户名" + newname + "\n")
		}

	} else {
		t.server.Broadcast(t, msg)
	}
}

// 给当前用户发送消息
func (t *User) SendMsg(msg string) {
	t.conn.Write([]byte(msg))
}

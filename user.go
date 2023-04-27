/*
 * @Author: zzzzztw
 * @Date: 2023-04-25 14:56:39
 * @LastEditors: Do not edit
 * @LastEditTime: 2023-04-27 17:10:33
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
		if !ok { // 管道关闭
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

	} else if len(msg) >= 3 && msg[:3] == "to|" {
		// 消息格式：to|name|消息内容

		// 1. 获取对方用户名
		n := len(strings.Split(msg, "|"))
		if n < 3 {
			t.SendMsg("消息格式不正确，请输入\"to|name|content\"格式。\n")
			return
		}
		remoteName := strings.Split(msg, "|")[1]

		if remoteName == "" {
			t.SendMsg("消息格式不正确，请输入\"to|name|content\"格式。\n")
			return
		}

		//2.根据用户名获取对方user

		if remoteUser, ok := t.server.OnlineMap[remoteName]; !ok {
			t.SendMsg("用户不存在或不在线\n")
		} else {
			content := strings.Split(msg, "|")[2]
			if content == "" {
				t.SendMsg("无内容请重发\n")
				return
			} else {
				remoteUser.SendMsg(t.Name + "对您说：" + content + "\n")
			}
		}

	} else {
		t.server.Broadcast(t, msg)
	}
}

// 给t用户发送消息
func (t *User) SendMsg(msg string) {
	t.conn.Write([]byte(msg))
}

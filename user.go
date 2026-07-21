package main

import (
	"net"
	"strings"
)

type User struct {
	Name   string
	Addr   string
	Ch     chan string
	conn   net.Conn
	server *Server
}

func (this *User) ListenMessage() {
	for {
		msg := <-this.Ch
		this.conn.Write([]byte(msg + "\n"))
	}
}

func (this *User) Online() {
	this.server.mapLock.Lock()
	this.server.OnlineMap[this.Name] = this
	this.server.mapLock.Unlock()

	this.server.BroadCast(this, "now online!")
}

func (this *User) Offline() {
	this.server.mapLock.Lock()
	delete(this.server.OnlineMap, this.Name)
	this.server.mapLock.Unlock()

	this.server.BroadCast(this, "now offline!")
}

func (this *User) DoMessage(msg string) {
	if msg == "/online" {
		// query for all online users
		this.server.mapLock.RLock()
		for _, user := range this.server.OnlineMap {
			temp := "[" + user.Addr + "]" + user.Name + ":" + "is online...\n"
			this.SendMessage(temp)
		}
		this.server.mapLock.RUnlock()
	} else if len(msg) >= 8 && msg[:8] == "/rename " {
		NewName := msg[8:]
		_, ok := this.server.OnlineMap[NewName]
		if ok {
			this.SendMessage("Error: name " + NewName + " already exist! \n")
		} else {
			this.server.mapLock.Lock()
			delete(this.server.OnlineMap, this.Name)
			this.server.OnlineMap[NewName] = this
			this.server.mapLock.Unlock()
			this.Name = NewName
			this.SendMessage("Your name has been successfully changed, new name : " + NewName + "\n")
		}
	} else if len(msg) >= 4 && msg[:4] == "/to " {
		parts := strings.SplitN(msg, " ", 3)
		if len(parts) < 3 {
			this.SendMessage("Usage: /to <username> <message>\n")
		} else {
			to := parts[1]
			pmsg := parts[2]
			target, ok := this.server.OnlineMap[to]
			if !ok {
				this.SendMessage("Error: user " + to + " is not online!\n")
			} else {
				target.SendMessage("[Private Message] " + this.Name + ": " + pmsg + "\n")
			}
		}
	} else {
		this.server.BroadCast(this, msg)
	}
}

func (this *User) SendMessage(msg string) {
	this.conn.Write([]byte(msg))
}

func NewUser(conn net.Conn, server *Server) *User {
	userAddr := conn.RemoteAddr().String()

	user := &User{
		Name:   userAddr,
		Addr:   userAddr,
		Ch:     make(chan string),
		conn:   conn,
		server: server,
	}

	go user.ListenMessage()
	return user
}

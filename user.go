package main

import "net"

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

func NewUser(conn net.Conn) *User {
	userAddr := conn.RemoteAddr().String()

	user := &User{
		Name: userAddr,
		Addr: userAddr,
		Ch:   make(chan string),
		conn: conn,
	}

	go user.ListenMessage()
	return user
}

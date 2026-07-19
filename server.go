package main

import (
	"fmt"
	"io"
	"net"
	"sync"
)

type Server struct {
	Ip   string
	Port int
	// Online user list
	OnlineMap map[string]*User
	mapLock   sync.RWMutex
	// broadcast channel
	Message chan string
}

func NewServer(ip string, port int) *Server {
	server := &Server{
		Ip:        ip,
		Port:      port,
		OnlineMap: make(map[string]*User),
		Message:   make(chan string),
	}
	return server
}

func (this *Server) Handler(conn net.Conn) {
	// ... Logics here
	// fmt.Printf("New connection accepted!")
	user := NewUser(conn)

	// when user gets online, put user info into online map
	this.mapLock.Lock()
	this.OnlineMap[user.Name] = user
	this.mapLock.Unlock()

	// broadcast this user now online!
	this.BroadCast(user, "now online!")

	// this goroutine is set for recieving the message from the client
	// current function is TCP-Based Multi-User Chatroom
	// Private chatting function will be added in the future development
	go func() {
		buf := make([]byte, 4096)
		for {
			n, err := conn.Read(buf)
			if n == 0 {
				this.BroadCast(user, "now offline!")
				return
			}

			if err != nil && err != io.EOF {
				fmt.Println("Conn Read err:", err)
				return
			}

			// n-1 is set for remove the tail '\n'
			msg := string(buf[:n-1])

			this.BroadCast(user, msg)
		}
	}()

	// current handler blocked
	select {}
}

func (this *Server) BroadCast(user *User, msg string) {
	sendMsg := "[" + user.Addr + "] " + user.Name + ":" + msg
	this.Message <- sendMsg
}

func (this *Server) ListenMessage() {
	for {
		msg := <-this.Message

		// spread msg to every online user
		this.mapLock.Lock()
		for _, cli := range this.OnlineMap {
			cli.Ch <- msg
		}
		this.mapLock.Unlock()

	}
}

func (this *Server) Start() {
	listener, err := net.Listen("tcp", fmt.Sprintf("%s:%d", this.Ip, this.Port))
	if err != nil {
		fmt.Println("net.Listen err:", err)
		return
	}
	defer listener.Close()

	// start a goroutine to listen broadcast message
	go this.ListenMessage()
	for {
		conn, err := listener.Accept()
		if err != nil {
			fmt.Println("listener.Accept err:", err)
			continue
		}

		go this.Handler(conn)
	}
}

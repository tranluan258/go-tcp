package server

import (
	"bufio"
	"fmt"
	"net"
	"os"
)

type client struct {
	conn net.Conn
}

type sever struct {
	host string
	port string
}

func New() *sever {
	return &sever{
		host: "localhost",
		port: "8080",
	}
}

func (s *sever) Run() {
	listener, err := net.Listen("tcp", fmt.Sprintf("%s:%s", s.host, s.port))
	if err != nil {
		os.Exit(100)
	}

	fmt.Printf("Server started on %s:%s\n", s.host, s.port)

	defer listener.Close()

	for {
		conn, _ := listener.Accept()
		fmt.Printf("New connection from %s", conn.RemoteAddr().String())
		client := &client{
			conn: conn,
		}
		go client.handleRequest()
	}
}

func (client *client) handleRequest() {
	reader := bufio.NewReader(client.conn)

	for {
		msg, err := reader.ReadString('\n')
		if err != nil {
			client.conn.Close()
			return
		}
		fmt.Printf("Incoming message: %s", string(msg))

		client.conn.Write([]byte("Message received.\n"))
	}
}

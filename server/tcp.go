package server

import (
	"bufio"
	"context"
	"fmt"
	"log"
	"net"
	"os"
	"os/signal"
	"sync"
	"time"
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
		log.Fatal("Error starting server: ", err)
	}

	fmt.Printf("Server started on %s:%s\n", s.host, s.port)

	ctx, stop := signal.NotifyContext(context.Background(), os.Interrupt)
	defer stop()
	go func() {
		<-ctx.Done()
		_ = listener.Close()
	}()

	var wg sync.WaitGroup
	for {
		conn, err := listener.Accept()
		if err != nil {
			if ctx.Err() != nil {
				fmt.Println("Server shutting down.")
				break
			}
			continue
		}

		fmt.Printf("New connection from %s \n", conn.RemoteAddr().String())

		client := &client{
			conn: conn,
		}

		wg.Add(1)
		go client.handleRequest(&wg, ctx)
	}

	wg.Wait()
	fmt.Println("Server stopped.")
}

func (client *client) handleRequest(wg *sync.WaitGroup, ctx context.Context) {
	defer wg.Done()
	defer client.conn.Close()

	reader := bufio.NewReader(client.conn)

	_ = client.conn.SetReadDeadline(time.Now().Add(5 * time.Minute))

	for {
		msg, err := reader.ReadString('\n')
		if err != nil {
			client.conn.Close()
			return
		}

		select {
		case <-ctx.Done():
			fmt.Println("Server shutting down, closing connection.")
			return
		default:
		}

		fmt.Printf("Incoming message: %s", string(msg))

		client.conn.Write([]byte("Message received.\n"))
	}
}

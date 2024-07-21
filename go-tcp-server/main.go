package main

import (
	"context"
	"fmt"
	"net"
	"time"
)

type message struct {
	from    string
	payload []byte
}
type Server struct {
	listernAddr string
	listener    net.Listener
	quitch      chan struct{}
	msgCh       chan message
	connections []net.Conn
}

func NewServer(addr string) *Server {
	return &Server{
		listernAddr: addr,
		quitch:      make(chan struct{}),
		msgCh:       make(chan message, 10),
	}
}

func (s *Server) closeAllConnections() {
	for _, conn := range s.connections {
		conn.Close()
	}
}
func (s *Server) Start() error {
	ln, err := net.Listen("tcp", s.listernAddr)
	if err != nil {
		return err
	}
	s.listener = ln
	fmt.Println("Server started listening on tcp ", s.listener.Addr().String())
	ctx, cancel := context.WithCancel(context.Background())
	defer func() {
		cancel()
		s.closeAllConnections()
		ln.Close()
	}()
	go s.AcceptLoop(ctx)
	<-s.quitch
	fmt.Println("SERVER CLOSED")
	return nil
}

func (s *Server) AcceptLoop(ctx context.Context) {
	for {
		select {
		case <-ctx.Done():
			fmt.Println("Stopping Accept. Server Terminated.")
			return
		default:
			conn, err := s.listener.Accept()
			if err != nil {
				fmt.Println("accept error: ", err)
				continue
			}
			s.connections = append(s.connections, conn)
			fmt.Println("Recieved a New Connection")
			go s.ReadLoop(ctx, conn)
		}
	}
}

func (s *Server) ReadLoop(ctx context.Context, conn net.Conn) {

	defer conn.Close()
	buf := make([]byte, 2048)
	for {

		select {
		case <-ctx.Done():
			fmt.Println("Stopping Reading from Connection. Server Terminated.")
			return
		default:
			fmt.Println("READING...")
			//conn.SetDeadline(time.Now().Add(time.Second * 5))
			n, err := conn.Read(buf)
			if err != nil {
				fmt.Println("Read Errpr: ", err)
				continue
			}
			msg := message{
				from:    conn.RemoteAddr().String(),
				payload: buf[:n],
			}
			s.msgCh <- msg
		}

	}
}
func main() {
	server := NewServer(":3000")
	go func() {
		for msg := range server.msgCh {
			fmt.Printf("Received message %s from connection %s\n", string(msg.payload), msg.from)
		}
	}()
	go server.Start()
	time.Sleep(time.Second * 20)
	server.quitch <- struct{}{}
	time.Sleep(time.Second * 3600)
}

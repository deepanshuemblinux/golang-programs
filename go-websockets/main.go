package main

import (
	"fmt"
	"io"
	"net/http"
	"sync"

	"github.com/sirupsen/logrus"
	"golang.org/x/net/websocket"
)

type Server struct {
	connections map[*websocket.Conn]bool
	mu          sync.Mutex
}

func NewServer() *Server {
	return &Server{
		connections: make(map[*websocket.Conn]bool),
		mu:          sync.Mutex{},
	}
}

func (s *Server) handleWS(conn *websocket.Conn) {
	fmt.Println("new incoming connection from websocket client: ", conn.RemoteAddr())
	s.mu.Lock()
	defer s.mu.Unlock()
	s.connections[conn] = true
	s.readLoop(conn)
}

func (s *Server) readLoop(conn *websocket.Conn) {
	buf := make([]byte, 1024)
	remote_addr := conn.RemoteAddr().String()
	for {
		n, err := conn.Read(buf)
		if err != nil {
			if err == io.EOF {
				logrus.Info("Connection to client closed", remote_addr)
			}
			logrus.Error(err)
			continue
		}
		msg := buf[:n]
		logrus.WithFields(logrus.Fields{
			"message": string(msg),
			"from":    remote_addr,
		}).Info("Recieved a message ")
		conn.Write([]byte("Thank you for the message!"))
	}

}
func main() {
	server := NewServer()
	http.Handle("/ws", websocket.Handler(server.handleWS))
	http.ListenAndServe(":8080", nil)
}

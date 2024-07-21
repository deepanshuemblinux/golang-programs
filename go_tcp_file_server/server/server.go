package server

import (
	"context"
	"errors"
	"fmt"
	"log"
	"net"
	"strings"
	"time"
)

const (
	BUFF_SIZE int = 1024
)

type fileServer struct {
	address     string
	listener    net.Listener
	connections []net.Conn
	serv_close  chan bool
}

func NewFileServer(address string) *fileServer {
	return &fileServer{
		address:    address,
		serv_close: make(chan bool),
	}
}

func (fs *fileServer) Run() {
	listener, err := net.Listen("tcp", fs.address)
	defer func() {
		for _, conn := range fs.connections {
			conn.Close()
		}
		fs.listener.Close()
	}()
	if err != nil {
		panic(fmt.Errorf("Error: %s", err))
	}
	fs.listener = listener
	ctx, cancel := context.WithCancel(context.Background())
	log.Printf("Runnig server on %s", listener.Addr().String())
	go fs.acceptLoop(ctx)
	<-fs.serv_close
	log.Println("Closing File Server")
	cancel()

}

func (fs *fileServer) acceptLoop(ctx context.Context) {
	//log.Println("Accept loop running")
	for {
		select {
		case <-ctx.Done():
			log.Println("Exiting acceptLoop")
			return
		default:
			conn, err := fs.listener.Accept()
			if err != nil {
				if !errors.Is(err, net.ErrClosed) {
					log.Println(err)
				} else {
					return
				}
			}
			conn.SetDeadline(time.Now().Add(time.Minute * 2))
			fs.connections = append(fs.connections, conn)
			go handleConnection(ctx, conn)
		}

	}
}

func handleConnection(ctx context.Context, conn net.Conn) {
	defer conn.Close()
	for {
		select {
		case <-ctx.Done():
			log.Println("Exiting handle connections loop")
			return
		default:
			buf := make([]byte, BUFF_SIZE)
			for {
				n, err := conn.Read(buf)
				if err != nil {
					log.Println(err)
					return
				}

				fmt.Printf("Client %s sent %s\n", conn.RemoteAddr().String(), strings.TrimRight(string(buf[:n]), "\n"))
			}
		}
	}

}
func (fs *fileServer) Close() {
	fs.serv_close <- true
}

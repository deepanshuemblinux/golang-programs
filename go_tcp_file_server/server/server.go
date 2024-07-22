package server

import (
	"context"
	"errors"
	"fmt"
	"io"
	"log"
	"net"
	"os"
	"strconv"
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
			file_info := make([]byte, BUFF_SIZE)
			n, _ := conn.Read(file_info)
			fmt.Printf("num bytes read is %d\n", n)
			fmt.Println(string(file_info[:n]))
			conn.Write([]byte("Read"))
			file_info_slice := strings.Split(string(file_info[:n]), " ")
			file_name := file_info_slice[0]
			file_size, err := strconv.Atoi(file_info_slice[1])
			fmt.Printf("File Name is %s, file size is %d\n", file_name, file_size)
			if err != nil {
				return
			}
			file, err := os.Create(file_name)
			if err != nil {
				log.Printf("Error opening the file %s to write\n", file_name)
				return
			}
			numLoops := file_size / BUFF_SIZE
			numRemainingBytes := file_size % BUFF_SIZE
			for i := 0; i < numLoops; i++ {
				_, err := io.ReadFull(conn, buf)
				file.Write(buf)
				if err != nil {
					log.Println(err)
					return
				}
			}
			_, err = io.ReadAtLeast(conn, buf, numRemainingBytes)
			fmt.Printf("numLoops is %d, numRemBytes is %d\n", numLoops, numRemainingBytes)
			file.Write(buf[:numRemainingBytes])
			file.Sync()
		}
	}

}
func (fs *fileServer) Close() {
	fs.serv_close <- true
}

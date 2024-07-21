package client

import (
	"errors"
	"fmt"
	"io"
	"log"
	"net"
	"os"
	"sync"
)

const (
	BUFF_SIZE int = 1024
)

type PoolableConnection interface {
	getId() int
}

type connection struct {
	id   int
	conn net.Conn
}

func (c *connection) getId() int {
	return c.id
}

type fileClient struct {
	serverAddr        string
	numConn           int
	activeConnections []connection
	idleConnections   []connection
	mu                sync.Mutex
}

func Init(serverAddr string, numConn int) *fileClient {
	fc := fileClient{
		serverAddr: serverAddr,
		numConn:    numConn,
	}
	for i := 0; i < fc.numConn; i++ {
		c, err := net.Dial("tcp", fc.serverAddr)
		if err != nil {
			log.Fatal(err)
		}
		conStruct := connection{
			id:   i,
			conn: c,
		}
		fc.idleConnections = append(fc.idleConnections, conStruct)
	}
	return &fc
}

func (c *fileClient) loan() *connection {
	c.mu.Lock()
	defer c.mu.Unlock()
	if len(c.idleConnections) == 0 {
		return nil
	}
	ret := c.idleConnections[0]
	c.idleConnections = c.idleConnections[1:]
	c.activeConnections = append(c.activeConnections, ret)
	fmt.Println("connection loaned")
	return &ret
}

func (c *fileClient) receive(target *connection) error {
	c.mu.Lock()
	defer c.mu.Unlock()
	if len(c.activeConnections) == 0 {
		return fmt.Errorf("No active connections to be received by Pool")
	}
	for idx, con := range c.activeConnections {
		if con.getId() == target.getId() {
			c.activeConnections = append(c.activeConnections[:idx], c.activeConnections[idx+1:]...)
			c.idleConnections = append(c.idleConnections, con)
			break
		}
	}
	fmt.Println("connection received back")
	return nil
}
func (c *fileClient) SendFile(file *os.File) error {
	con := c.loan()
	defer c.receive(con)
	buff := make([]byte, BUFF_SIZE)
	for {
		numBytes, err := io.ReadFull(file, buff)

		if err != nil {
			if errors.Is(err, io.ErrUnexpectedEOF) {
				con.conn.Write(buff[:numBytes])
				break
			} else if errors.Is(err, io.EOF) {
				break
			} else {
				return err
			}

		}
		con.conn.Write(buff[:numBytes])
	}
	return nil
}

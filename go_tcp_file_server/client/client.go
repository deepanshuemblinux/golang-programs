package client

import (
	"errors"
	"fmt"
	"io"
	"log"
	"net"
	"os"
	"strconv"
	"sync"
	"time"
)

const (
	BUFF_SIZE int = 1024
)

func retryLoan(numRetries, waitTimeSec int, f func() (*connection, error)) (*connection, error) {
	for i := 0; i < numRetries; i++ {
		ret, err := f()
		if err == nil {
			return ret, nil
		}
		timer := time.NewTimer(time.Second * time.Duration(waitTimeSec))
		<-timer.C
	}
	return nil, fmt.Errorf("Timed out")
}

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

func (c *fileClient) loan() (*connection, error) {
	c.mu.Lock()
	defer c.mu.Unlock()
	if len(c.idleConnections) == 0 {
		return nil, fmt.Errorf("Could not loan connection")
	}
	ret := c.idleConnections[0]
	c.idleConnections = c.idleConnections[1:]
	c.activeConnections = append(c.activeConnections, ret)
	fmt.Println("connection loaned")
	return &ret, nil
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
	file_info, err := file.Stat()
	if err != nil {
		return err
	}
	con, err := retryLoan(5, 2, c.loan)
	if err != nil {
		return err
	}
	defer c.receive(con)
	buff := make([]byte, BUFF_SIZE)
	file_size := strconv.FormatInt(file_info.Size(), 10)
	if err != nil {
		return err
	}
	file_info_buff := []byte(file_info.Name() + " " + file_size)
	fmt.Println(string(file_info_buff))
	con.conn.Write(file_info_buff)
	con.conn.Read(file_info_buff)

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

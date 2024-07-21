package main

import (
	"log"
	"os"
	"strings"
	"time"

	"github.com/deepanshuemblinux/go-tcp-file-server/client"
	"github.com/deepanshuemblinux/go-tcp-file-server/server"
)

func main() {
	ch := make(chan any)
	cmd := strings.ToLower(os.Args[1])
	switch cmd {
	case "server":
		file_server := server.NewFileServer(":8081")

		//fmt.Printf("%+v\n", file_server)
		go file_server.Run()
		time.Sleep(time.Second * 5)
	case "client":
		file_client := client.Init(":8081", 2)
		f1, err := os.Open("client/client.go")
		if err != nil {
			log.Fatal(err)
		}
		f2, err := os.Open("server/server.go")
		if err != nil {
			log.Fatal(err)
		}
		go file_client.SendFile(f1)
		go file_client.SendFile(f2)

	}
	<-ch
}

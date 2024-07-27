package main

import (
	"github.com/deepanshuemblinux/go-rate-limiter/api"
	"github.com/deepanshuemblinux/go-rate-limiter/service"
)

func main() {
	srvc := service.NewTextMessageService()
	server := api.NewAPIServer(":8080", srvc)
	server.Run()
	ch := make(chan bool)
	<-ch
}

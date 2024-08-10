package main

import (
	"log"

	"github.com/deepanshuemblinux/go-chat-websockets/db"
	"github.com/deepanshuemblinux/go-chat-websockets/router"
	"github.com/deepanshuemblinux/go-chat-websockets/user"
)

func main() {
	database, err := db.NewDatabase()
	if err != nil {
		log.Fatal(err)
	}
	userRep := user.NewRepository(database.GetDB())
	userSvc := user.NewService(userRep)
	handler := user.NewHandler(userSvc)
	router.InitRouter(handler)
	router.Start(":8080")
}

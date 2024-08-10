package router

import (
	"net/http"

	"github.com/deepanshuemblinux/go-chat-websockets/user"
)

func InitRouter(userHandler *user.Handler) {
	http.HandleFunc("/signup", userHandler.CreateUser)
	http.HandleFunc("/login", userHandler.Login)
	http.HandleFunc("/logout", userHandler.Logout)
}

func Start(addr string) {
	http.ListenAndServe(addr, nil)
}

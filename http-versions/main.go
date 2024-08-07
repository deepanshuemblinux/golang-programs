package main

import (
	"fmt"
	"log"
	"net/http"

	"golang.org/x/net/http2"
)

func main() {
	http.Handle("/", http.FileServer(http.Dir("./static")))

	server := &http.Server{
		Addr: ":8080",
	}

	http2.ConfigureServer(server, &http2.Server{})
	fmt.Println("Starting web server on ::8080")

	log.Fatal(server.ListenAndServeTLS("server.crt", "server.key"))
}

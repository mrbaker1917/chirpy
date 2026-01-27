package main

import (
	"fmt"
	"log"
	"net/http"
)

func main() {
	serveMux := http.NewServeMux()
	var httpServer http.Server
	httpServer.Addr = ":8080"
	fs := http.FileServer(http.Dir("."))
	serveMux.Handle("/", fs)
	httpServer.Handler = serveMux

	fmt.Printf("Server is running on localhost%s\n", httpServer.Addr)
	err := httpServer.ListenAndServe()

	if err != nil {
		log.Println("server could not be set up: ", err)
	}

}

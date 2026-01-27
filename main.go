package main

import (
	"fmt"
	"net/http"
)

func main() {
	serveMux := http.NewServeMux()
	var httpServer http.Server
	httpServer.Addr = ":8080"
	httpServer.Handler = serveMux
	err := httpServer.ListenAndServe()

	if err != nil {
		fmt.Println("server could not be set up: ", err)
	}

}
